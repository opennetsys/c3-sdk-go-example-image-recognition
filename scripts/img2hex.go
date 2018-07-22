package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	fileBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("err reading file\n%v", err)
	}

	fmt.Println(hex.EncodeToString(fileBytes))
}
