package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileService(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	svc := NewFileService(tmpDir)

	// Test write and read
	testPath := "test.txt"
	testContent := "Hello, World!"

	if err := svc.WriteFile(testPath, testContent); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	content, err := svc.ReadFile(testPath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if content != testContent {
		t.Errorf("Content mismatch: got %q, want %q", content, testContent)
	}

	// Test GetFileInfo
	info, err := svc.GetFileInfo(testPath)
	if err != nil {
		t.Fatalf("GetFileInfo failed: %v", err)
	}

	if info.Name != "test.txt" {
		t.Errorf("Wrong filename: got %q, want test.txt", info.Name)
	}

	if info.Size != int64(len(testContent)) {
		t.Errorf("Wrong file size: got %d, want %d", info.Size, len(testContent))
	}

	// Test ListDir
	files, err := svc.ListDir(".")
	if err != nil {
		t.Fatalf("ListDir failed: %v", err)
	}

	if len(files) < 1 {
		t.Error("ListDir returned empty list")
	}

	// Test DeleteFile
	if err := svc.DeleteFile(testPath); err != nil {
		t.Fatalf("DeleteFile failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, testPath)); !os.IsNotExist(err) {
		t.Error("File was not deleted")
	}
}

func TestFileServiceCreateDir(t *testing.T) {
	tmpDir := t.TempDir()
	svc := NewFileService(tmpDir)

	dirPath := "newdir/subdir"
	if err := svc.CreateDir(dirPath); err != nil {
		t.Fatalf("CreateDir failed: %v", err)
	}

	fullPath := filepath.Join(tmpDir, dirPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Error("Directory was not created")
	}
}
