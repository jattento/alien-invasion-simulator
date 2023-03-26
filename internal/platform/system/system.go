package system

import (
	"io"
	"os"
)

type Manager struct {
	OpenFunc func(string) (io.ReadCloser, error)
}

func NewManager() *Manager {
	return &Manager{
		OpenFunc: func(s string) (io.ReadCloser, error) { return os.Open(s) },
	}
}
