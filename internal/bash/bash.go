package bash

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/anomalyco/opencode/internal/util/log"
)

// Executor handles bash command execution
type Executor struct {
	workDir string
	timeout time.Duration
	env     []string
	mu      sync.Mutex
}

// ExecutorOption configures the executor
type ExecutorOption func(*Executor)

// WithWorkDir sets the working directory
func WithWorkDir(dir string) ExecutorOption {
	return func(e *Executor) {
		e.workDir = dir
	}
}

// WithTimeout sets the default timeout
func WithTimeout(d time.Duration) ExecutorOption {
	return func(e *Executor) {
		e.timeout = d
	}
}

// WithEnv sets additional environment variables
func WithEnv(env map[string]string) ExecutorOption {
	return func(e *Executor) {
		for k, v := range env {
			e.env = append(e.env, k+"="+v)
		}
	}
}

// NewExecutor creates a new bash executor
func NewExecutor(opts ...ExecutorOption) *Executor {
	e := &Executor{
		timeout: 30 * time.Second,
		env:     os.Environ(),
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Result represents the result of a command execution
type Result struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	ExitCode   int    `json:"exitCode"`
	Command    string `json:"command"`
	Duration   int64  `json:"duration"` // milliseconds
	WorkingDir string `json:"workingDir,omitempty"`
}

// Execute runs a bash command
func (e *Executor) Execute(ctx context.Context, command string) (*Result, error) {
	return e.ExecuteWithTimeout(ctx, command, e.timeout)
}

// ExecuteWithTimeout runs a bash command with a specific timeout
func (e *Executor) ExecuteWithTimeout(ctx context.Context, command string, timeout time.Duration) (*Result, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	start := time.Now()

	// Create timeout context
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Prepare command
	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	cmd.Env = e.env
	cmd.Dir = e.workDir

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Default.Debug("Executing command", map[string]interface{}{
		"command": command,
		"timeout": timeout.String(),
	})

	err := cmd.Run()

	result := &Result{
		Stdout:     stdout.String(),
		Stderr:     stderr.String(),
		Command:    command,
		Duration:   time.Since(start).Milliseconds(),
		WorkingDir: e.workDir,
	}

	// Get exit code
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				result.ExitCode = status.ExitStatus()
			} else {
				result.ExitCode = 1
			}
		} else if ctx.Err() == context.DeadlineExceeded {
			result.ExitCode = -1
			result.Stderr = fmt.Sprintf("Command timed out after %s", timeout)
		} else {
			result.ExitCode = 1
		}
	} else {
		result.ExitCode = 0
	}

	log.Default.Debug("Command completed", map[string]interface{}{
		"command":   command,
		"exitCode":  result.ExitCode,
		"duration":  result.Duration,
	})

	return result, nil
}

// ExecuteInteractive runs a command with interactive output
func (e *Executor) ExecuteInteractive(ctx context.Context, command string, stdoutHandler, stderrHandler func(string)) (*Result, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	start := time.Now()

	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	cmd.Env = e.env
	cmd.Dir = e.workDir

	// Create pipes for streaming output
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// Stream output
	var stdoutBuf, stderrBuf bytes.Buffer
	done := make(chan struct{})

	go func() {
		scanner := bytes.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			stdoutBuf.WriteString(line + "\n")
			if stdoutHandler != nil {
				stdoutHandler(line)
			}
		}
		done <- struct{}{}
	}()

	go func() {
		scanner := bytes.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			stderrBuf.WriteString(line + "\n")
			if stderrHandler != nil {
				stderrHandler(line)
			}
		}
		done <- struct{}{}
	}()

	// Wait for both streams
	<-done
	<-done

	err = cmd.Wait()

	result := &Result{
		Stdout:     stdoutBuf.String(),
		Stderr:     stderrBuf.String(),
		Command:    command,
		Duration:   time.Since(start).Milliseconds(),
		WorkingDir: e.workDir,
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				result.ExitCode = status.ExitStatus()
			}
		}
	}

	return result, nil
}

// SetWorkDir changes the working directory
func (e *Executor) SetWorkDir(dir string) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	if _, err := os.Stat(abs); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", abs)
	}
	e.workDir = abs
	return nil
}

// SetEnv sets an environment variable
func (e *Executor) SetEnv(key, value string) {
	// Remove existing key
	var newEnv []string
	prefix := key + "="
	for _, env := range e.env {
		if !strings.HasPrefix(env, prefix) {
			newEnv = append(newEnv, env)
		}
	}
	newEnv = append(newEnv, key+"="+value)
	e.env = newEnv
}

// GetEnv gets an environment variable
func (e *Executor) GetEnv(key string) string {
	prefix := key + "="
	for _, env := range e.env {
		if strings.HasPrefix(env, prefix) {
			return strings.TrimPrefix(env, prefix)
		}
	}
	return ""
}

// IsCommandAvailable checks if a command is available
func (e *Executor) IsCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// Which finds the path to a command
func (e *Executor) Which(command string) string {
	path, err := exec.LookPath(command)
	if err != nil {
		return ""
	}
	return path
}

// DefaultExecutor is the default executor
var DefaultExecutor = NewExecutor()

// Execute is a convenience function using the default executor
func Execute(ctx context.Context, command string) (*Result, error) {
	return DefaultExecutor.Execute(ctx, command)
}

// ExecuteWithTimeout is a convenience function using the default executor
func ExecuteWithTimeout(ctx context.Context, command string, timeout time.Duration) (*Result, error) {
	return DefaultExecutor.ExecuteWithTimeout(ctx, command, timeout)
}
