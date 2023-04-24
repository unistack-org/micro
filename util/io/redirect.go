package io

import (
	"os"
	"sync"
)

var osStderrMu sync.Mutex

var OrigStderr = func() *os.File {
	fd, err := dupFD(os.Stderr.Fd())
	if err != nil {
		panic(err)
	}

	return os.NewFile(fd, os.Stderr.Name())
}()
