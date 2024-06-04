package server

import (
	"dicedb/config"
	"dicedb/core"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

func SetUpSyncServer() {
	listen, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	log.Println("server started on", config.Host+":"+strconv.Itoa(config.Port))

	if err != nil {
		return
	}
	conClients := 0

	for {
		conn, err := listen.Accept()

		conClients++

		log.Println("client connected with add ", conn.RemoteAddr(), "total concurrent req ", conClients)
		if err != nil {
			return
		}

		handleSyncConnection(conn, &conClients)
	}
}

func handleSyncConnection(conn net.Conn, conClients *int) {

	for {
		cmd, err, _ := readCommand(conn)

		if err != nil {
			*conClients--
			log.Println("client disconnected with add ", conn.LocalAddr(), "total concurrent req ", conClients)
			if err == io.EOF {
				break
			}
			log.Println(err.Error())
		}

		respond(cmd, conn)
	}
}

func respond(cmd *core.RedisCmd, conn io.ReadWriter) {
	log.Print("Command : ", cmd.Cmd)
	log.Print("Args : ", strings.Join(cmd.Args, ","))

	err := core.EvalAndRespond(cmd, conn)

	if err != nil {
		respondError(err, conn)
	}
}

func respondError(err error, conn io.ReadWriter) {
	_, err = conn.Write([]byte(fmt.Sprintf("-%s\r\n", err.Error())))
	if err != nil {
		log.Println(err.Error())
	}
}

func readCommand(conn io.ReadWriter) (*core.RedisCmd, error, int) {

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)

	if err != nil {
		return nil, err, 0
	}

	tokens, err := core.DecodeArrayString(buffer[:n])

	if err != nil {
		return nil, err, 0
	}

	return &core.RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil, n
}
