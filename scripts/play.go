package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/c3systems/c3/core/sandbox"
)

func main() {
	if len(os.Args) != 5 {
		log.Fatal("expected args: image id, initial state, file name, file extension")
	}

	imageID := os.Args[1]
	initialState := os.Args[2]
	fileName := os.Args[3]
	fileExtension := os.Args[4]

	if imageID == "" || initialState == "" || fileName == "" || fileExtension == "" {
		log.Fatal("expected args: image id, initial state, file name, file extension")
	}

	sb := sandbox.New(nil)

	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("err reading file\n%v", err)
	}

	newState, err := sb.Play(&sandbox.PlayConfig{
		ImageID:      imageID,
		InitialState: []byte(initialState),
		Payload:      []byte(fmt.Sprintf(`[%q,%q,%q]`, "processImage", hex.EncodeToString(fileBytes), fileExtension)),
	})
	if err != nil {
		log.Fatalf("err playing in sandbox\n%v", err)
	}

	log.Printf("success! new state:\n%s", string(newState))
}
