package core

import (
	"bytes"
	"errors"
	"fmt"
)

func decodeOne(msg []byte) (interface{}, int, error) {
	switch msg[0] {
	case '+':
		return decodeString(msg)
	case ':':
		return decodeInteger(msg)
	case '-':
		return decodeError(msg)
	case '$':
		return decodeBulkString(msg)
	case '*':
		return decodeArrays(msg)
	default:
		return "", 0, nil
	}
}

func decodeArrays(msg []byte) (interface{}, int, error) {
	l, pos := getLen(msg, 1)
	arr := make([]interface{}, l)
	i := 0

	for l > 0 {
		decoded, newPos, _ := decodeOne(msg[pos:])
		pos = newPos + pos
		arr[i] = decoded
		l--
		i++
	}

	return arr, pos, nil
}

func decodeBulkString(msg []byte) (interface{}, int, error) {
	length, start := getLen(msg, 1)
	end := start + length
	return string(msg[start:end]), end + 2, nil // actual message starts from 4th index
}

func getLen(bytes []byte, pos int) (int, int) {
	val := 0
	for bytes[pos] >= '0' && bytes[pos] <= '9' {
		val *= 10
		val += int(bytes[pos] - '0')
		pos++
	}
	return val, pos + 2 // skip \r\n in pos
}

func decodeError(msg []byte) (interface{}, int, error) {
	i := 1
	for ; msg[i] != '\r'; i++ {
	}
	return string(msg[1:i]), i + 2, nil
}

func decodeInteger(msg []byte) (interface{}, int, error) {
	i := 1
	for ; msg[i] != '\r'; i++ {
	}
	return string(msg[1:i]), i + 2, nil
}

func decodeString(msg []byte) (interface{}, int, error) {
	i := 1
	for ; msg[i] != '\r'; i++ {
	}
	return string(msg[1:i]), i + 2, nil
}

func Decode(msg []byte) ([]interface{}, error) {

	if len(msg) == 0 {
		return nil, errors.New("no data")
	}

	values := make([]interface{}, 0)
	index := 0
	for index < len(msg) {
		value, delta, err := decodeOne(msg)
		index += delta
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}

func Encode(value interface{}, isSimple bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimple {
			return []byte(fmt.Sprintf("+%s\r\n", v))
		} else {
			return EncodeString(v)
		}
	case []string:
		return encodeStringArray(v)
	case int, int8, int16, int32, int64:
		return []byte(fmt.Sprintf(":%d\r\n", v))
	case error:
		return []byte(fmt.Sprintf("-%s\r\n", v))
	default:
		return RespNil
	}
}

func encodeStringArray(arr []string) []byte {
	b := make([]byte, 0)
	buf := bytes.NewBuffer(b)

	for _, str := range arr {
		buf.Write(EncodeString(str))
	}

	return buf.Bytes()
}

func EncodeString(v string) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v))
}
