package core

import (
	"dicedb/config"
	"fmt"
	"log"
	"os"
)

func DumpAllAOF() {

	fp, err := os.OpenFile(config.AOFFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 666)
	defer fp.Close()

	if err != nil {
		log.Println(err.Error())
		return
	}

	for key, obj := range store {
		Dump(fp, key, obj)
	}
}

func Dump(fp *os.File, key string, obj *Obj) {
	write, err := fp.Write(Encode([]string{"SET", key, obj.Value.(string)}, false))
	fmt.Println("written ", write)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
