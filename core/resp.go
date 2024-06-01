package core

import (
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

func DecodeArrayString(msg []byte) ([]string, error) {
	value, err := decode(msg)

	if err != nil {
		return nil, err
	}

	value2 := value.([]interface{})

	res := make([]string, len(value2))

	for index, val := range value2 {
		res[index] = val.(string)
	}

	return res, nil
}

func decode(msg []byte) (interface{}, error) {

	if len(msg) == 0 {
		return nil, errors.New("no data")
	}

	value, _, err := decodeOne(msg)

	return value, err
}

func Encode(value interface{}, isSimple bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimple {
			return []byte(fmt.Sprintf("+%s\r\n", v))
		} else {
			return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v))
		}
	}
	return []byte("")
}
