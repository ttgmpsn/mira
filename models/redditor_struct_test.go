package models

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func BenchmarkCreateRedditor(b *testing.B) {
	data, _ := ioutil.ReadFile("./tests/redditor.json")
	redditorExampleJSON := string(data)
	for i := 0; i < b.N; i++ {
		sub := Redditor{}
		json.Unmarshal([]byte(redditorExampleJSON), &sub)
	}
}
