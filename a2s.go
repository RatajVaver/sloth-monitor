package main

import (
	"errors"
	"net"
	"time"
)

func queryA2SInfo(addr string) (int, int, error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return 0, 0, err
	}
	defer conn.Close()

	req := append([]byte{0xFF, 0xFF, 0xFF, 0xFF, 'T'}, []byte("Source Engine Query\x00")...)
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	if _, err := conn.Write(req); err != nil {
		return 0, 0, err
	}

	buf := make([]byte, 1400)
	n, err := conn.Read(buf)
	if err != nil {
		return 0, 0, err
	}

	buf = buf[:n]
	if len(buf) < 6 || buf[4] != 'I' {
		return 0, 0, errors.New("invalid response")
	}

	off := 6
	for i := 0; i < 4; i++ {
		noff := skipNullTerm(buf, off)
		if noff < 0 {
			return 0, 0, errors.New("malformed response")
		}
		off = noff
	}

	if off+4 > len(buf) {
		return 0, 0, errors.New("response truncated")
	}

	return int(buf[off+2]), int(buf[off+3]), nil
}

func skipNullTerm(b []byte, off int) int {
	for i := off; i < len(b); i++ {
		if b[i] == 0 {
			return i + 1
		}
	}
	return -1
}
