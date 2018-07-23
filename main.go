package main

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"strconv"

	c3 "github.com/c3systems/c3-sdk-go"
)

const (
	dataKey      string  = "data"
	unknownImage string  = "unknown"
	threshold    float32 = 0.2
)

var client = c3.NewC3()

type jsonResults map[string]uint64

// App ...
type App struct {
}

func (a *App) setItem(key string, value []byte) error {
	client.State().Set([]byte(key), value)
	return nil
}

func (a *App) getItem(key string) ([]byte, bool) {
	return client.State().Get([]byte(key))
}

func (a *App) getAllResults() string {
	v, found := client.State().Get([]byte(dataKey))
	if !found {
		return "{}"
	}
	return string(v)
}

func (a *App) getResultsForType(t string) (string, error) {
	var data jsonResults
	dataBytes, found := a.getItem(dataKey)
	if found {
		if err := json.Unmarshal(dataBytes, &data); err != nil {
			return "", err
		}
	}
	if data == nil {
		data = make(jsonResults)
	}

	return strconv.FormatUint(data[t], 10), nil
}

func (a *App) processImage(imageHex, format string) (string, error) {
	log.Println("process image called")
	var resultType = ""
	imageBytes, err := hex.DecodeString(imageHex)
	if err != nil {
		log.Printf("error decoding; %v", err)
		return "", err
	}
	results, err := recognizeHandler(imageBytes, format)
	if err != nil {
		log.Printf("error with recognize handler; %v", err)
		return "", err
	}
	if results == nil || len(results) == 0 {
		resultType = unknownImage
	} else {
		if results[0].Probability >= threshold {
			resultType = results[0].Label
		} else {
			resultType = unknownImage
		}
	}

	var data jsonResults
	dataBytes, found := a.getItem(dataKey)
	if found {
		if err := json.Unmarshal(dataBytes, &data); err != nil {
			log.Printf("error unmarshalling; %v", err)
			return "", err
		}
	}
	if data == nil {
		data = make(jsonResults)
	}

	data[resultType]++
	dataBytes, err = json.Marshal(data)
	if err != nil {
		log.Printf("error marshalling; %v", err)
		return "", err
	}

	if err = a.setItem(dataKey, dataBytes); err != nil {
		log.Printf("error setting item; %v", err)
		return "", err
	}

	log.Printf("processing image done; result type; %s", resultType)

	return resultType, nil
}

func main() {
	if err := loadModel(); err != nil {
		log.Fatalf("err loading model\n%v", err)
	}

	log.Println("running")
	app := &App{}
	if err := client.RegisterMethod("processImage", []string{"string", "string"}, app.processImage); err != nil {
		log.Fatal(err)
	}
	if err := client.RegisterMethod("getAllResults", nil, app.getAllResults); err != nil {
		log.Fatal(err)
	}
	if err := client.RegisterMethod("getResultsForType", []string{"string"}, app.getResultsForType); err != nil {
		log.Fatal(err)
	}

	client.Serve()
}
