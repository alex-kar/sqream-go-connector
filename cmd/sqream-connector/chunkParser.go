package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
)

type ChunkParser struct {
}

type ResultChunkHeader struct {
	ColSzs []int32
	Rows   int32
}

func (parser ChunkParser) Header(bytes []byte) (ResultChunkHeader, error) {
	if bytes == nil {
		return ResultChunkHeader{}, errors.New("Byte array is nil")
	}
	if len(bytes) < 8 {
		return ResultChunkHeader{}, errors.New("Byte arrays has less then 8 bytes")
	}
	headerLength := binary.LittleEndian.Uint64(bytes[:8])
	var result ResultChunkHeader
	json.Unmarshal(bytes[8:8+headerLength], &result)
	return result, nil
}
