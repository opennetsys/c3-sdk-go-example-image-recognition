package main

import (
	"encoding/json"
	"log"
	"strconv"

	c3 "github.com/c3systems/sdk-go"
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

func (a *App) setItem(key, value string) error {
	client.State().Set(key, value)
	return nil
}

func (a *App) getItem(key string) string {
	return client.State().Get(key)
}

func (a *App) getAllResults() string {
	return client.State().Get(dataKey)
}

func (a *App) getResultsForType(t string) (string, error) {
	var data jsonResults
	dataStr := a.getItem(dataKey)
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return "", err
	}
	if data == nil {
		data = make(jsonResults)
	}

	return strconv.FormatUint(data[t], 10), nil
}

func (a *App) processImage(image, format string) (string, error) {
	var resultType = ""

	results, err := recognizeHandler(image, format)
	if err != nil {
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
	dataStr := a.getItem(dataKey)
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return "", err
	}
	if data == nil {
		data = make(jsonResults)
	}

	data[resultType]++
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	dataStr = string(dataBytes)

	if err = a.setItem(dataKey, dataStr); err != nil {
		return "", err
	}

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
