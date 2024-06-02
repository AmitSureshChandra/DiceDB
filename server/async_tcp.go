package server

import (
	"dicedb/config"
	"dicedb/core"
	"io"
	"log"
	"net"
	"syscall"
)

var conClients = make([]int, 0)

func RunAsyncServer() error {
	log.Println("starting an asynchronous TCP server on", config.Host, config.Port)

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)

	if err != nil {
		return err
	}

	defer func(fd int) {
		err := syscall.Close(fd)
		if err != nil {
			panic(err.Error())
		}
	}(fd)

	if err := syscall.SetNonblock(fd, true); err != nil {
		return err
	}

	ip := net.ParseIP(config.Host).To4()

	add := syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ip[0], ip[1], ip[2], ip[3]},
	}
	// *2\r\n+PING\r\n+PONG\r\n

	if err := syscall.Bind(fd, &add); err != nil {
		log.Println(err.Error())
		return err
	}

	if err := syscall.Listen(fd, syscall.SOMAXCONN); err != nil {
		log.Println(err.Error())
		return err
	}

	for {

		err := acceptConn(fd)
		if err != nil {
			return err
		}

		err = readWrite()
		if err != nil {
			return err
		}
	}
}

func readWrite() error {
	for i := 0; i < len(conClients); i++ {
		err := handle(conClients[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func acceptConn(fd int) error {
	nfd, _, err := syscall.Accept(fd)

	if err != nil {
		if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
			//time.Sleep(100 * time.Millisecond)
			return nil
		}
		log.Println(err.Error())
		return err
	}

	// make non-blocking
	if err := syscall.SetNonblock(nfd, true); err != nil {
		return err
	}

	// store fd to further read & write
	conClients = append(conClients, nfd)

	println("new connection established ", conClients)
	return nil
}

func handle(nfd int) error {
	buffer := make([]byte, 1024)

	//n, err := syscall.Read(nfd, buffer)

	fdCmd := core.FDComm{Fd: nfd}

	cmd, err, n := readCommand(fdCmd)

	if err != nil {
		// if fd is not ready to read
		if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
			//time.Sleep(100 * time.Millisecond)
			return nil
		}

		// if connection is closed
		if err == io.EOF {
			closeFDConn(nfd)
		}
		log.Println(err.Error())
		return err
	} else {
		if n == 0 { // means connection broken from terminal
			closeFDConn(nfd)
		} else {
			println("reading on", nfd, "data :", string(buffer[:n]))
		}
	}

	// write if data available in buffer
	if n > 0 {
		_, err = syscall.Write(nfd, buffer[:n])

		if err != nil {
			if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
				//time.Sleep(100 * time.Millisecond)
				return nil
			}
			return err
		} else {
			println("writing on", nfd, "data :", string(buffer[:n]))
			respond(cmd, fdCmd)
		}
	}
	return nil
}

func closeFDConn(nfd int) {
	for i := 0; i < len(conClients); i++ {
		if conClients[i] == nfd {
			conClients = append(conClients[:i], conClients[i+1:]...)
			log.Print(nfd, " connection closed")
		}
	}
}
