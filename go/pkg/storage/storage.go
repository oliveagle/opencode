package storage

import (
	"log"
)

type Storage struct {
	// TODO: add storage fields (database connection, cache, etc.)
}

func New() *Storage {
	log.Println("Initializing storage layer")
	return &Storage{}
}

// TODO: implement storage interface methods
// ReadFile, WriteFile, ListDir, GetFileInfo, etc.
