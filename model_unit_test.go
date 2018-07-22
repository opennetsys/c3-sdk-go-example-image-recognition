// +built unit

package main

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestRecognizeHandler(t *testing.T) {
	type test struct {
		inputFile string
		expected  []LabelResult
		err       error
	}

	tests := []test{
		test{
			inputFile: "./images/cat/cat.jpg",
			expected: []LabelResult{
				LabelResult{
					Label:       "tiger cat",
					Probability: 0.4221901,
				},
				LabelResult{
					Label:       "tabby",
					Probability: 0.25731775,
				},
				LabelResult{
					Label:       "Egyptian cat",
					Probability: 0.24682136,
				},
				LabelResult{
					Label:       "lynx",
					Probability: 0.06381659,
				},
				LabelResult{
					Label:       "Persian cat",
					Probability: 0.004809883,
				},
			},
			err: nil,
		},
		test{
			inputFile: "./images/dog/dog.jpg",
			expected: []LabelResult{
				LabelResult{
					Label:       "English foxhound",
					Probability: 0.4735971,
				},
				LabelResult{
					Label:       "Walker hound",
					Probability: 0.45291826,
				},
				LabelResult{
					Label:       "beagle",
					Probability: 0.060902625,
				},
				LabelResult{
					Label:       "basset",
					Probability: 0.002712648,
				},
				LabelResult{
					Label:       "Saluki",
					Probability: 0.0024361422,
				},
			},
			err: nil,
		},
		test{
			inputFile: "./images/cow/cow.jpg",
			expected: []LabelResult{
				LabelResult{
					Label:       "Staffordshire bullterrier",
					Probability: 0.20728138,
				},
				LabelResult{
					Label:       "Boston bull",
					Probability: 0.16594343,
				},
				LabelResult{
					Label:       "American Staffordshire terrier",
					Probability: 0.1386188,
				},
				LabelResult{
					Label:       "French bulldog",
					Probability: 0.11638059,
				},
				LabelResult{
					Label:       "Great Dane",
					Probability: 0.099053636,
				},
			},
			err: nil,
		},
	}

	if err := loadModel(); err != nil {
		t.Fatalf("err loading model\n%v", err)
	}
	for idx, tt := range tests {
		fileBytes, err := ioutil.ReadFile(tt.inputFile)
		if err != nil {
			t.Errorf("test %d failed\nerr reading file\n%v", idx+1, err)
			continue
		}

		result, err := recognizeHandler(string(fileBytes), "jpg")
		if err != tt.err {
			t.Errorf("test %d failed\nexpected err %v\nreceived err %v", idx+1, tt.err, err)
			continue
		}

		if !reflect.DeepEqual(tt.expected, result) {
			t.Errorf("test %d failed\n expected %v\nreceived %v", idx+1, tt.expected, result)
		}
	}
}
