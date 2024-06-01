package core

import (
	"errors"
	"net"
)

func EvalPING(args []string, conn net.Conn) error {
	var b []byte
	if len(args) >= 2 {
		return errors.New("ERR wrong no of arguments for 'ping' command")
	}
	if len(args) == 0 {
		b = Encode("PONG", true)
	} else {
		b = Encode(args[0], false)
	}

	_, err := conn.Write(b)
	return err
}

func EvalAndRespond(cmd *RedisCmd, conn net.Conn) error {
	switch cmd.Cmd {
	case "PING":
		return EvalPING(cmd.Args, conn)
	default:
		return EvalPING(cmd.Args, conn)
	}
}
