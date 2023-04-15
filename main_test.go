package main

import (
	"fmt"
	"os"
	"testing"
)

func TestHello(t *testing.T) {
	dataDir := "./test_data"
	dir, err := os.ReadDir(dataDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, d := range dir {
		filepath := fmt.Sprintf("%s/%s", dataDir, d.Name())
		t.Log(filepath)
		b, err := os.ReadFile(filepath)
		if err != nil {
			t.Fatal(err)
		}
		unknownTxBytes(b)
		return // only one file
	}

}
