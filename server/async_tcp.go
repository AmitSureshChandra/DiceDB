package server

import (
	"dicedb/config"
	"dicedb/core"
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"net"
	"strconv"
	"syscall"
	"time"
)

var conClients = 0

var cronFreq = 1 * time.Second
var lastCronExcTime = time.Now()

func RunAsyncServer() error {

	listen, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		panic(err.Error())
	}

	defer listen.Close()

	epollFD, err := core.CreateEpoll()

	if err != nil {
		return err
	}

	defer syscall.Close(epollFD)

	f, err := listen.(*net.TCPListener).File()

	if err != nil {
		return err
	}

	listenerFD := f.Fd()

	if err := syscall.SetNonblock(int(listenerFD), true); err != nil {
		println(err.Error())
	}

	if err := core.AddToPoll(epollFD, int(listenerFD)); err != nil {
		return err
	}

	fmt.Println("Server is listening on 127.0.0.1:8080")

	for {

		// delete the expired keys
		if time.Now().After(lastCronExcTime.Add(cronFreq)) {
			core.DeleteExpireKeys()
			lastCronExcTime = time.Now()
		}

		events, err := core.WailForEvents(epollFD)

		if err != nil {
			return err
		}

		for _, event := range events {
			err := handleEvent(event, epollFD, listenerFD, listen)
			if err != nil {
				log.Println("error handling req")
				log.Println(err.Error())
			}
		}
	}
}

func handleEvent(event unix.EpollEvent, epollFD int, listenerFD uintptr, listen net.Listener) error {

	// listening for read events
	if event.Fd != int32(listenerFD) {
		err := handleConn(int(event.Fd))
		if err == nil {
			return nil
		}
		log.Println(err.Error())
		conClients--
		log.Println("client disconnected ", conClients)
		err = core.RemoveFromPoll(epollFD, int(event.Fd))
		if err != nil {
			log.Println(err.Error())
		}
		return nil
	}

	// listening for new connection
	conn, err := listen.Accept()

	if err != nil {
		return err
	}

	f, err := conn.(*net.TCPConn).File()

	if err != nil {
		return err
	}
	connFD := int(f.Fd())

	if err := syscall.SetNonblock(connFD, true); err != nil {
		println(err.Error())
	}

	if err := core.AddToPoll(epollFD, connFD); err != nil {
		conn.Close()
		return err
	}
	conClients++
	log.Println("client connected ", conClients)
	return nil
}

func handleConn(nfd int) error {
	fdCmd := core.FDComm{Fd: nfd}
	cmds, err := readCommand(fdCmd)
	if err != nil {
		return err
	}
	respond(cmds, fdCmd)
	return nil
}
