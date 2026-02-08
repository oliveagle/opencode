package installation

import (
	"os"
	"path/filepath"
	"runtime"
)

const Version = "1.1.53"

func IsLocal() bool {
	// Check if running from source (development mode)
	execPath, _ := os.Executable()
	return filepath.Base(filepath.Dir(execPath)) == "bin" || 
		filepath.Base(filepath.Dir(execPath)) == ".bin" ||
		os.Getenv("OPENCODE_DEV") == "1"
}

func ConfigDir() string {
	if dir := os.Getenv("OPENCODE_CONFIG_DIR"); dir != "" {
		return dir
	}
	
	homeDir, _ := os.HomeDir()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(homeDir, ".config", "opencode")
	case "windows":
		return filepath.Join(homeDir, "AppData", "Roaming", "opencode")
	default:
		if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
			return filepath.Join(xdgConfig, "opencode")
		}
		return filepath.Join(homeDir, ".config", "opencode")
	}
}

func DataDir() string {
	if dir := os.Getenv("OPENCODE_DATA_DIR"); dir != "" {
		return dir
	}
	
	homeDir, _ := os.HomeDir()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(homeDir, ".local", "share", "opencode")
	case "windows":
		return filepath.Join(homeDir, "AppData", "Local", "opencode")
	default:
		if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData != "" {
			return filepath.Join(xdgData, "opencode")
		}
		return filepath.Join(homeDir, ".local", "share", "opencode")
	}
}

func CacheDir() string {
	if dir := os.Getenv("OPENCODE_CACHE_DIR"); dir != "" {
		return dir
	}
	
	homeDir, _ := os.HomeDir()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(homeDir, ".cache", "opencode")
	case "windows":
		return filepath.Join(homeDir, "AppData", "Local", "opencode", "Cache")
	default:
		if xdgCache := os.Getenv("XDG_CACHE_HOME"); xdgCache != "" {
			return filepath.Join(xdgCache, "opencode")
		}
		return filepath.Join(homeDir, ".cache", "opencode")
	}
}
