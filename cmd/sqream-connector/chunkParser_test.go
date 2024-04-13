package main

import (
	"os"
	"reflect"
	"testing"
)

func TestHeader(t *testing.T) {
	parser := new(ChunkParser)
	file, err := os.ReadFile("chunkRawData.txt")
	if err != nil {
		t.Fail()
	}

	header, err := parser.Header(file)
	if err != nil {
		t.Fail()
	}

	if !reflect.DeepEqual(header.ColSzs, []int32{1, 4}) {
		t.Errorf("Wrong column sizes %v", header.ColSzs)
	}

	if header.Rows != 1 {
		t.Errorf("Wrong amount of rows %d", header.Rows)
	}
}
