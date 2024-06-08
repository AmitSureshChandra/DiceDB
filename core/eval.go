package core

import (
	"errors"
	"io"
	"strconv"
	"time"
)

var RespNil []byte = []byte("$-1\r\n")

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

func EvalGet(args []string, conn io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("(error) ERR wrong no of arguments for `GET` command")
	}

	key := args[0]
	// check key in map

	obj := Get(key)
	if obj == nil {
		conn.Write(RespNil)
		return nil
	}
	// check key is not expired in map

	if obj.ExpiredAt != -1 && obj.ExpiredAt <= time.Now().UnixMilli() {
		conn.Write(RespNil)
		return nil
	}

	// return resp
	conn.Write(Encode(obj.Value, false))
	return nil
}

func EvalTTL(args []string, conn io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("(error) ERR wrong no of arguments for `TTL` command")
	}

	key := args[0]
	// check key in map

	obj := Get(key)

	// if obj not exists then return error code -2
	if obj == nil {
		conn.Write([]byte(":-2\r\n"))
		return nil
	}

	// if no expiration set
	if obj.ExpiredAt == -1 {
		conn.Write([]byte(":-1\r\n"))
		return nil
	}

	// calculation time duration remaining

	durationMS := obj.ExpiredAt - time.Now().UnixMilli()

	conn.Write(Encode(durationMS/1000, false))

	return nil
}

func EvalSet(args []string, conn io.ReadWriter) error {
	if len(args) < 2 {
		return errors.New("(error) ERR wrong no of arguments for `set` command")
	}

	var key, value = args[0], args[1]

	var exDurationMs int64 = -1

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "px", "PX":
			i++
			if i == len(args) {
				return errors.New("(error) syntax error")
			}
			exDurationMs2, err := strconv.ParseInt(args[i], 10, 64)

			if err != nil {
				return errors.New("(Error) ERR value is not integer of out of range")
			}
			exDurationMs = exDurationMs2

		case "ex", "EX":
			i++
			if i == len(args) {
				return errors.New("(error) syntax error")
			}
			exDurationSec, err := strconv.ParseInt(args[i], 10, 64)

			if err != nil {
				return errors.New("(error) ERR syntax error")
			}

			exDurationMs = exDurationSec * 1000
		default:
			return errors.New("(error) ERR syntax error")
		}
	}

	Put(key, NewObj(value, exDurationMs))

	conn.Write(Encode("OK", true))
	return nil
}

func EvalAndRespond(cmd *RedisCmd, conn io.ReadWriter) error {
	switch cmd.Cmd {
	case "PING":
		return EvalPING(cmd.Args, conn)
	case "GET":
		return EvalGet(cmd.Args, conn)
	case "SET":
		return EvalSet(cmd.Args, conn)
	case "TTL":
		return EvalTTL(cmd.Args, conn)
	case "DEL":
		return EvalDel(cmd.Args, conn)
	case "EXPIRE":
		return EvalExp(cmd.Args, conn)
	default:
		return EvalPING(cmd.Args, conn)
	}
}

func EvalExp(args []string, conn io.ReadWriter) error {
	if len(args) != 2 {
		return errors.New("(error) ERR wrong no of arguments for `EXP` command")
	}

	key := args[0]
	// check key in map

	exDurationSec, err := strconv.ParseInt(args[1], 10, 64)

	if err != nil {
		return errors.New("(Error) ERR value is not integer of out of range")
	}

	obj := Get(key)

	// if obj not exists then return code 0
	if obj == nil {
		conn.Write([]byte(":0\r\n"))
		return nil
	}

	obj.ExpiredAt = time.Now().UnixMilli() + exDurationSec*1000

	// return 1 if timeout is set
	conn.Write([]byte(":1\r\n"))

	return nil
}

func EvalDel(args []string, conn io.ReadWriter) error {
	deletedCnt := 0

	for _, key := range args {
		if ok := Del(key); ok {
			deletedCnt++
		}
	}

	conn.Write(Encode(deletedCnt, false))
	return nil
}
