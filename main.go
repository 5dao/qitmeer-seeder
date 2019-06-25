package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
)

var (
	amgr *Manager
	wg   sync.WaitGroup
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "loadConfig: %v\n", err)
		os.Exit(1)
	}
	amgr, err = NewManager(filepath.Join(defaultHomeDir,
		activeNetParams.Name))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "NewManager: %v\n", err)
		os.Exit(1)
	}

	amgr.AddAddresses([]net.IP{net.ParseIP(cfg.Seeder)})

	wg.Add(1)
	go creep()

	dnsServer := NewDNSServer(cfg.Host, cfg.Nameserver, cfg.Listen)
	go dnsServer.StartTCP()

	wg.Wait()
}
