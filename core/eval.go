package core

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"time"
)

var RespNil = []byte("$-1\r\n")
var RespOK = []byte("+OK\r\n")
var RespZero = []byte(":0\r\n")
var RespOne = []byte(":1\r\n")
var RespMinusOne = []byte(":-1\r\n")
var RespMinusTwo = []byte(":-2\r\n")

func EvalPING(args []string) []byte {
	if len(args) >= 2 {
		return Encode(errors.New("ERR wrong no of arguments for 'ping' command"), false)
	}
	if len(args) == 0 {
		return Encode("PONG", true)
	} else {
		return Encode(args[0], false)
	}
}

func EvalGet(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong no of arguments for `GET` command"), false)
	}

	key := args[0]
	// check key in map

	obj := Get(key)
	if obj == nil {
		return RespNil
	}
	// check key is not expired in map

	if obj.ExpiredAt != -1 && obj.ExpiredAt <= time.Now().UnixMilli() {
		return RespNil
	}

	// return resp
	return Encode(obj.Value, false)
}

func EvalTTL(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong no of arguments for `TTL` command"), false)
	}

	key := args[0]
	// check key in map

	obj := Get(key)

	// if obj not exists then return error code -2
	if obj == nil {
		return RespMinusTwo
	}

	// if no expiration set
	if obj.ExpiredAt == -1 {
		return RespMinusOne
	}

	// calculation time duration remaining

	durationMS := obj.ExpiredAt - time.Now().UnixMilli()
	return Encode(durationMS/1000, false)
}

func EvalSet(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong no of arguments for `set` command"), false)
	}

	var key, value = args[0], args[1]

	var exDurationMs int64 = -1

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "px", "PX":
			i++
			if i == len(args) {
				return Encode(errors.New("(error) syntax error"), false)
			}
			exDurationMs2, err := strconv.ParseInt(args[i], 10, 64)

			if err != nil {
				return Encode(errors.New("(Error) ERR value is not integer of out of range"), false)
			}
			exDurationMs = exDurationMs2

		case "ex", "EX":
			i++
			if i == len(args) {
				return Encode(errors.New("(error) syntax error"), false)
			}
			exDurationSec, err := strconv.ParseInt(args[i], 10, 64)

			if err != nil {
				return Encode(errors.New("(error) ERR syntax error"), false)
			}

			exDurationMs = exDurationSec * 1000
		default:
			return Encode(errors.New("(error) ERR syntax error"), false)
		}
	}

	Put(key, NewObj(value, exDurationMs))
	return Encode("OK", true)
}

func EvalExp(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong no of arguments for `EXP` command"), false)
	}

	key := args[0]
	// check key in map

	exDurationSec, err := strconv.ParseInt(args[1], 10, 64)

	if err != nil {
		return Encode(errors.New("(Error) ERR value is not integer of out of range"), false)
	}

	obj := Get(key)

	// if obj not exists then return code 0
	if obj == nil {
		return RespZero
	}

	obj.ExpiredAt = time.Now().UnixMilli() + exDurationSec*1000
	return RespOne
}

func EvalDel(args []string) []byte {
	deletedCnt := 0

	for _, key := range args {
		if ok := Del(key); ok {
			deletedCnt++
		}
	}
	return Encode(deletedCnt, false)
}

func EvalAndRespond(cmds RedisCmds, conn io.ReadWriter) {

	var response []byte
	buf := bytes.NewBuffer(response)

	for _, cmd := range cmds {
		switch cmd.Cmd {
		case "PING":
			buf.Write(EvalPING(cmd.Args))
		case "GET":
			buf.Write(EvalGet(cmd.Args))
		case "SET":
			buf.Write(EvalSet(cmd.Args))
		case "TTL":
			buf.Write(EvalTTL(cmd.Args))
		case "DEL":
			buf.Write(EvalDel(cmd.Args))
		case "EXPIRE":
			buf.Write(EvalExp(cmd.Args))
		default:
			buf.Write(EvalPING(cmd.Args))
		}
	}
	conn.Write(buf.Bytes())
}
