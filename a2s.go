package main

import (
	"errors"
	"fmt"
	"net"
	"time"
)

func queryA2SInfo(addr string, retries int, timeout time.Duration) (int, int, error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return 0, 0, err
	}
	defer conn.Close()

	baseReq := append([]byte{0xFF, 0xFF, 0xFF, 0xFF, 'T'}, []byte("Source Engine Query\x00")...)
	buf := make([]byte, 1400)

	var lastErr error
	for attempt := 0; attempt < retries; attempt++ {
		req := baseReq

		conn.SetDeadline(time.Now().Add(timeout))
		if _, err := conn.Write(req); err != nil {
			lastErr = err
			continue
		}

		n, err := conn.Read(buf)
		if err != nil {
			lastErr = err
			continue
		}

		resp := buf[:n]
		if len(resp) < 5 {
			lastErr = errors.New("response too short")
			continue
		}

		if resp[4] == 'A' && len(resp) >= 9 {
			req = append(baseReq, resp[5:9]...)
			conn.SetDeadline(time.Now().Add(timeout))
			if _, err := conn.Write(req); err != nil {
				lastErr = err
				continue
			}
			n, err = conn.Read(buf)
			if err != nil {
				lastErr = err
				continue
			}
			resp = buf[:n]
		}

		if len(resp) < 6 || resp[4] != 'I' {
			lastErr = errors.New("invalid response")
			continue
		}

		return parseA2SInfo(resp)
	}

	return 0, 0, fmt.Errorf("after %d attempts: %w", retries, lastErr)
}

func parseA2SInfo(buf []byte) (int, int, error) {
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
