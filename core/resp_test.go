package core

import (
	"fmt"
	"log"
	"testing"
)

func TestLen(t *testing.T) {
	if l, _ := getLen([]byte("$100\r\n"), 1); l != 100 {
		t.Fail()
	}
}

func TestDecode(t *testing.T) {

	//cases := map[string]interface{}{
	//	"+OK\r\n":            "OK",
	//	"-Error message\r\n": "Error message",
	//	":1000\r\n":          "1000",
	//	"$5\r\nhello\r\n":    "hello",
	//}
	//
	//for k, v := range cases {
	//	if val, _, err := decodeOne([]byte(k)); val != v || err != nil {
	//		t.Fail()
	//	}
	//}

	cases2 := map[string][]interface{}{
		"*3\r\n:1000\r\n:2000\r\n:3000\r\n":                        {1000, 2000, 3000},
		"*3\r\n+Hello\r\n+World\r\n+Mumbai\r\n":                    {"Hello", "World", "Mumbai"},
		"*0\r\n":                                                   {},
		"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n":                     {"hello", "world"},
		"*3\r\n:1\r\n:2\r\n:3\r\n":                                 {int64(1), int64(2), int64(3)},
		"*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n":            {int64(1), int64(2), int64(3), int64(4), "hello"},
		"*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Hello\r\n-World\r\n": {[]int64{int64(1), int64(2), int64(3)}, []interface{}{"Hello", "World"}},
	}

	for k, v := range cases2 {

		val, _, err := decodeOne([]byte(k))

		if err != nil {
			log.Fatal(err.Error())
		}

		array := val.([]interface{})

		for i := 0; i < len(array); i++ {

			println(fmt.Sprintf("%v", array[i]))
			println(fmt.Sprintf("%v", v[i]))

			if fmt.Sprintf("%v", array[i]) != fmt.Sprintf("%v", v[i]) {
				t.Fail()
			}
		}
	}
}

func TestDecodeArrayString(t *testing.T) {
	given, err := DecodeArrayString([]byte("*3\r\n+Hello\r\n+World\r\n+Mumbai\r\n"))

	expected := []string{
		"Hello",
		"World",
		"Mumbai",
	}

	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}

	for i := range given {
		if expected[i] != given[i] {
			t.Fail()
		}
	}
}
