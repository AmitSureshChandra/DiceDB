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
		cmds, err := readCommand(conn)
		if err != nil {
			*conClients--
			log.Println("client disconnected with add ", conn.LocalAddr(), "total concurrent req ", conClients)
			if err == io.EOF {
				break
			}
			log.Println(err.Error())
		}

		respond(cmds, conn)
	}
}

func respond(cmds core.RedisCmds, conn io.ReadWriter) {
	core.EvalAndRespond(cmds, conn)
}

func respondError(err error, conn io.ReadWriter) {
	_, err = conn.Write([]byte(fmt.Sprintf("-%s\r\n", err.Error())))
	if err != nil {
		log.Println(err.Error())
	}
}

func readCommand(conn io.ReadWriter) (core.RedisCmds, error) {

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)

	if err != nil {
		return nil, err
	}

	cmds := make(core.RedisCmds, 0)

	values, err := core.Decode(buffer[:n])

	if err != nil {
		return nil, err
	}

	for _, value := range values {
		tokens, err := toArrayString(value.([]interface{}))

		if err != nil {
			return nil, err
		}
		cmds = append(cmds, &core.RedisCmd{
			Cmd:  strings.ToUpper(tokens[0]),
			Args: tokens[1:],
		})
	}

	return cmds, err
}

func toArrayString(value []interface{}) ([]string, error) {
	arr := make([]string, len(value))
	pos := 0
	for _, item := range value {
		arr[pos] = item.(string)
		pos++
	}
	return arr, nil
}
