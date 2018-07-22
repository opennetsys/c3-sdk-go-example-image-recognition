package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"sort"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

//// ClassifyResult ...
//type ClassifyResult struct {
//Filename string        `json:"filename"`
//Labels   []LabelResult `json:"labels"`
//}

// LabelResult ...
type LabelResult struct {
	Label       string  `json:"label"`
	Probability float32 `json:"probability"`
}

var (
	graphModel   *tf.Graph
	sessionModel *tf.Session
	labels       []string
)

func loadModel() error {
	// Load inception model
	model, err := ioutil.ReadFile("model/tensorflow_inception_graph.pb")
	if err != nil {
		return err
	}

	graphModel = tf.NewGraph()
	if err := graphModel.Import(model, ""); err != nil {
		return err
	}

	sessionModel, err = tf.NewSession(graphModel, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Load labels
	labelsFile, err := os.Open("model/imagenet_comp_graph_label_strings.txt")
	if err != nil {
		return err
	}

	defer labelsFile.Close()
	scanner := bufio.NewScanner(labelsFile)

	// Labels are separated by newlines
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}

	return err
}

func recognizeHandler(image []byte, imageFormat string) ([]LabelResult, error) {
	imageBuffer := bytes.NewBuffer(image)

	// Make tensor
	tensor, err := makeTensorFromImage(imageBuffer, imageFormat)
	if err != nil {
		return nil, err
	}

	// Run inference
	output, err := sessionModel.Run(
		map[tf.Output]*tf.Tensor{
			graphModel.Operation("input").Output(0): tensor,
		},
		[]tf.Output{
			graphModel.Operation("output").Output(0),
		},
		nil)
	if err != nil {
		return nil, err
	}

	// Return best labels
	return findBestLabels(output[0].Value().([][]float32)[0]), nil
}

// ByProbability ...
type ByProbability []LabelResult

// Len ...
func (a ByProbability) Len() int { return len(a) }

// Swap ...
func (a ByProbability) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less ...
func (a ByProbability) Less(i, j int) bool { return a[i].Probability > a[j].Probability }

func findBestLabels(probabilities []float32) []LabelResult {
	// Make a list of label/probability pairs
	var resultLabels []LabelResult
	for i, p := range probabilities {
		if i >= len(labels) {
			break
		}
		resultLabels = append(resultLabels, LabelResult{Label: labels[i], Probability: p})
	}
	// Sort by probability
	sort.Sort(ByProbability(resultLabels))
	// Return top 5 labels
	return resultLabels[:5]
}
