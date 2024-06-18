package core

import "golang.org/x/sys/unix"

func CreateEpoll() (int, error) {
	epollFD, err := unix.EpollCreate1(0)

	if err != nil {
		return -1, err
	}
	return epollFD, nil
}

func AddToPoll(epollFD int, fd int) error {
	event := unix.EpollEvent{Events: unix.EPOLLIN, Fd: int32(fd)}
	err := unix.EpollCtl(epollFD, unix.EPOLL_CTL_ADD, fd, &event)
	if err != nil {
		return err
	}
	return nil
}

func RemoveFromPoll(epollFD int, fd int) error {
	err := unix.EpollCtl(epollFD, unix.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}
	return nil
}

func WailForEvents(ePollFD int) ([]unix.EpollEvent, error) {
	events := make([]unix.EpollEvent, 10)

	n, err := unix.EpollWait(ePollFD, events, -1)

	if err != nil {
		return nil, err
	}
	return events[:n], nil
}
