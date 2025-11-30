package main

import (
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := OpenDB(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	poll := time.Duration(cfg.Poll) * time.Second

	log.Printf("monitor: poll=%s, table=%s", poll, cfg.DBTable)

	for {
		s, err := GetNextServer(db, cfg.DBTable)
		if err != nil {
			log.Printf("error fetching server: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		host := stripPort(s.IP)
		addr := net.JoinHostPort(host, strconv.Itoa(s.QueryPort))

		log.Printf("checking server: id=%d addr=%s", s.ID, addr)
		players, maxp, err := queryA2SInfo(addr)
		if err != nil {
			log.Printf("query error: id=%d addr=%s: %v", s.ID, addr, err)
			if err := UpdateServerStatus(db, cfg.DBTable, s.ID, -1, 0, cfg.AvgRatio); err != nil {
				log.Printf("failed to update: id=%d: %v", s.ID, err)
			}
		} else {
			if err := UpdateServerStatus(db, cfg.DBTable, s.ID, players, maxp, cfg.AvgRatio); err != nil {
				log.Printf("failed to update: id=%d: %v", s.ID, err)
			} else {
				log.Printf("updated: id=%d players=%d/%d", s.ID, players, maxp)
			}
		}

		time.Sleep(poll)
	}
}

func stripPort(ip string) string {
	if h, _, err := net.SplitHostPort(ip); err == nil {
		return h
	}

	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		suf := ip[idx+1:]
		if _, err := strconv.Atoi(suf); err == nil {
			return ip[:idx]
		}
	}

	return ip
}
