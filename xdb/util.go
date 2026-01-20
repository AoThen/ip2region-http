// Copyright 2022 The Ip2Region Authors. All rights reserved.
// Use of this source code is governed by a Apache2.0-style
// license that can be found in the LICENSE file.

package xdb

import (
	"bytes"
	"fmt"
	"net"
)

func ParseIP(ip string) ([]byte, error) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return nil, fmt.Errorf("invalid ip address: %s", ip)
	}

	v4 := parsedIP.To4()
	if v4 != nil {
		return v4, nil
	}

	return nil, fmt.Errorf("invalid ip address: %s", ip)
}

func IP2String(ip []byte) string {
	return net.IP(ip[:]).String()
}

func IPCompare(ip1, ip2 []byte) int {
	return bytes.Compare(ip1, ip2)
}
