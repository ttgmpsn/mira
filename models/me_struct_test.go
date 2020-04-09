package models

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func BenchmarkCreateMe(b *testing.B) {
	data, _ := ioutil.ReadFile("./tests/me.json")
	meExampleJSON := string(data)
	for i := 0; i < b.N; i++ {
		sub := Me{}
		json.Unmarshal([]byte(meExampleJSON), &sub)
	}
}
