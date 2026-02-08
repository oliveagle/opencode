package core

import (
	"io"
	"os"
	"path/filepath"

	"github.com/anomalyco/opencode/pkg/types"
)

type FileService struct {
	baseDir string
}

func NewFileService(baseDir string) *FileService {
	return &FileService{baseDir: baseDir}
}

func (f *FileService) ReadFile(path string) (string, error) {
	fullPath := filepath.Join(f.baseDir, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (f *FileService) WriteFile(path, content string) error {
	fullPath := filepath.Join(f.baseDir, path)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(fullPath, []byte(content), 0644)
}

func (f *FileService) ListDir(path string) ([]types.FileInfo, error) {
	fullPath := filepath.Join(f.baseDir, path)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	var files []types.FileInfo
	for _, e := range entries {
		info, _ := e.Info()
		files = append(files, types.FileInfo{
			Path:    filepath.Join(path, e.Name()),
			Name:    e.Name(),
			IsDir:   e.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime().Unix(),
		})
	}
	return files, nil
}

func (f *FileService) GetFileInfo(path string) (*types.FileInfo, error) {
	fullPath := filepath.Join(f.baseDir, path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	return &types.FileInfo{
		Path:    path,
		Name:    info.Name(),
		Size:    info.Size(),
		IsDir:   info.IsDir(),
		ModTime: info.ModTime().Unix(),
	}, nil
}

func (f *FileService) CreateDir(path string) error {
	fullPath := filepath.Join(f.baseDir, path)
	return os.MkdirAll(fullPath, 0755)
}

func (f *FileService) DeleteFile(path string) error {
	fullPath := filepath.Join(f.baseDir, path)
	return os.RemoveAll(fullPath)
}

func (f *FileService) CopyFile(src, dst string) error {
	srcFile, err := os.Open(filepath.Join(f.baseDir, src))
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstPath := filepath.Join(f.baseDir, dst)
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
