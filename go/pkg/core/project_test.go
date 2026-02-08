package core

import (
	"testing"

	"github.com/anomalyco/opencode/pkg/types"
)

func TestProjectService(t *testing.T) {
	tmpDir := t.TempDir()
	svc, err := NewProjectService(tmpDir)
	if err != nil {
		t.Fatalf("NewProjectService failed: %v", err)
	}

	// Test CreateProject
	proj := &types.ProjectConfig{
		Name:  "test-project",
		Path:  "/tmp/test",
		Root:  "/tmp/test",
		Alias: "test",
	}

	if err := svc.CreateProject(proj); err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}

	// Test GetProject
	retrieved := svc.GetProject("test-project")
	if retrieved == nil {
		t.Fatal("GetProject returned nil")
	}

	if retrieved.Name != "test-project" {
		t.Errorf("Project name mismatch: got %q, want test-project", retrieved.Name)
	}

	// Test ListProjects
	projects := svc.ListProjects()
	if len(projects) != 1 {
		t.Errorf("ListProjects returned %d projects, want 1", len(projects))
	}

	// Test DeleteProject
	if err := svc.DeleteProject("test-project"); err != nil {
		t.Fatalf("DeleteProject failed: %v", err)
	}

	if svc.GetProject("test-project") != nil {
		t.Error("Project was not deleted")
	}

	if len(svc.ListProjects()) != 0 {
		t.Error("ListProjects should be empty")
	}
}

func TestProjectServicePersistence(t *testing.T) {
	tmpDir := t.TempDir()

	// Create and save
	svc1, _ := NewProjectService(tmpDir)
	proj := &types.ProjectConfig{
		Name:  "persistent",
		Path:  "/tmp/p",
		Root:  "/tmp/p",
		Alias: "p",
	}
	svc1.CreateProject(proj)

	// Load from disk
	svc2, err := NewProjectService(tmpDir)
	if err != nil {
		t.Fatalf("Failed to reload service: %v", err)
	}

	retrieved := svc2.GetProject("persistent")
	if retrieved == nil {
		t.Fatal("Persistence failed - project not found after reload")
	}

	if retrieved.Alias != "p" {
		t.Errorf("Persistent data mismatch: got %q, want p", retrieved.Alias)
	}
}
