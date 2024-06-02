package core

import (
	"errors"
	"io"
)

func EvalPING(args []string, conn io.ReadWriter) error {
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

func EvalAndRespond(cmd *RedisCmd, conn io.ReadWriter) error {
	switch cmd.Cmd {
	case "PING":
		return EvalPING(cmd.Args, conn)
	default:
		return EvalPING(cmd.Args, conn)
	}
}
