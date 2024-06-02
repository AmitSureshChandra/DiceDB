package core

import "syscall"

type FDComm struct {
	Fd int
}

func (fd FDComm) Read(b []byte) (int, error) {
	return syscall.Read(fd.Fd, b)
}

func (fd FDComm) Write(b []byte) (int, error) {
	return syscall.Write(fd.Fd, b)
}
