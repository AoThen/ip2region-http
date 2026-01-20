// Copyright 2022 The Ip2Region Authors. All rights reserved.
// Use of this source code is governed by a Apache2.0-style
// license that can be found in the LICENSE file.

package xdb

import (
	"encoding/binary"
	"fmt"
	"os"
)

type Searcher struct {
	version     *Version
	vectorIndex []byte
	contentBuff []byte
}

func NewWithBuffer(version *Version, cBuff []byte) (*Searcher, error) {
	return &Searcher{
		version:     version,
		vectorIndex: nil,
		contentBuff: cBuff,
	}, nil
}

func (s *Searcher) IPVersion() *Version {
	return s.version
}

func (s *Searcher) Search(ip []byte) (string, error) {
	if len(ip) != s.version.Bytes {
		return "", fmt.Errorf("invalid ip address")
	}

	var il0, il1 = int(ip[0]), int(ip[1])
	var idx = il0*VectorIndexCols*VectorIndexSize + il1*VectorIndexSize
	var sPtr, ePtr = uint32(0), uint32(0)
	if s.vectorIndex != nil {
		sPtr = binary.LittleEndian.Uint32(s.vectorIndex[idx:])
		ePtr = binary.LittleEndian.Uint32(s.vectorIndex[idx+4:])
	} else if s.contentBuff != nil {
		sPtr = binary.LittleEndian.Uint32(s.contentBuff[HeaderInfoLength+idx:])
		ePtr = binary.LittleEndian.Uint32(s.contentBuff[HeaderInfoLength+idx+4:])
	}

	var bytes, dBytes = len(ip), len(ip) << 1
	var segIndexSize = uint32(s.version.SegmentIndexSize)
	var dataLen, dataPtr = 0, uint32(0)
	var buff = make([]byte, segIndexSize)
	var l, h = 0, int((ePtr - sPtr) / segIndexSize)
	for l <= h {
		m := (l + h) >> 1
		p := sPtr + uint32(m)*segIndexSize
		err := s.read(int64(p), buff)
		if err != nil {
			return "", fmt.Errorf("read segment index: %w", err)
		}

		if s.version.IPCompare(ip, buff[0:bytes]) < 0 {
			h = m - 1
		} else if s.version.IPCompare(ip, buff[bytes:dBytes]) > 0 {
			l = m + 1
		} else {
			dataLen = int(binary.LittleEndian.Uint16(buff[dBytes:]))
			dataPtr = binary.LittleEndian.Uint32(buff[dBytes+2:])
			break
		}
	}

	if dataLen == 0 {
		return "", fmt.Errorf("ip not found")
	}

	var regionBuff = make([]byte, dataLen)
	err := s.read(int64(dataPtr), regionBuff)
	if err != nil {
		return "", fmt.Errorf("read region: %w", err)
	}

	return string(regionBuff), nil
}

func (s *Searcher) SearchByStr(str string) (string, error) {
	ip, err := ParseIP(str)
	if err != nil {
		return "", err
	}

	return s.Search(ip)
}

func (s *Searcher) read(offset int64, buff []byte) error {
	if s.contentBuff != nil {
		cLen := copy(buff, s.contentBuff[offset:])
		if cLen != len(buff) {
			return fmt.Errorf("incomplete read")
		}
	} else {
		return fmt.Errorf("content buffer is required")
	}
	return nil
}

func LoadHeaderFromFile(dbFile string) (*Header, error) {
	handle, err := os.OpenFile(dbFile, os.O_RDONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("open xdb file: %w", err)
	}
	defer handle.Close()

	var buff = make([]byte, HeaderInfoLength)
	rLen, err := handle.Read(buff)
	if err != nil {
		return nil, err
	}
	if rLen != len(buff) {
		return nil, fmt.Errorf("incomplete read")
	}

	return NewHeader(buff)
}

func LoadContentFromFile(dbFile string) ([]byte, error) {
	handle, err := os.OpenFile(dbFile, os.O_RDONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("open xdb file: %w", err)
	}
	defer handle.Close()

	fi, err := handle.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}

	size := fi.Size()
	var buff = make([]byte, size)
	rLen, err := handle.Read(buff)
	if err != nil {
		return nil, err
	}
	if rLen != len(buff) {
		return nil, fmt.Errorf("incomplete read")
	}

	return buff, nil
}
