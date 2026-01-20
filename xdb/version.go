// Copyright 2022 The Ip2Region Authors. All rights reserved.
// Use of this source code is governed by a Apache2.0-style
// license that can be found in the LICENSE file.

package xdb

import (
	"bytes"
	"fmt"
	"strings"
)

type Version struct {
	Id               int
	Name             string
	Bytes            int
	SegmentIndexSize int
	IPCompare        func([]byte, []byte) int
}

func (v *Version) String() string {
	return fmt.Sprintf(
		"{id:%d, name:%s, bytes:%d, segment_index_size:%d}",
		v.Id, v.Name, v.Bytes, v.SegmentIndexSize,
	)
}

const (
	IPv4VersionNo = 4
)

var IPv4 = &Version{
	Id:               IPv4VersionNo,
	Name:             "IPv4",
	Bytes:            4,
	SegmentIndexSize: 14,
	IPCompare: func(ip1, ip2 []byte) int {
		ip2[0], ip2[3] = ip2[3], ip2[0]
		ip2[1], ip2[2] = ip2[2], ip2[1]
		return bytes.Compare(ip1, ip2)
	},
}

func VersionFromName(name string) (*Version, error) {
	switch strings.ToUpper(name) {
	case "V4", "IPV4":
		return IPv4, nil
	default:
		return nil, fmt.Errorf("invalid version name `%s`", name)
	}
}
