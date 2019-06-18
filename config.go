package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"net"
	"os"
	"path/filepath"
	"qitmeer/params"
	"strings"
)

const (
	defaultListenPort = "18130"
)

var (
	// Default network parameters
	activeNetParams = &params.MainNetParams
	//get current path
	defaultHomeDir, _ = os.Getwd()
)

// See loadConfig for details on the configuration load process.
type config struct {
	Host       string `short:"H" long:"host" description:"Seed DNS address"`
	Listen     string `long:"listen" short:"l" description:"Listen on address:port"`
	Nameserver string `short:"n" long:"nameserver" description:"hostname of nameserver"`
	Seeder     string `short:"s" long:"default seeder" description:"IP address of a working node"`
	TestNet    bool   `long:"testnet" description:"Use testnet"`
}

func loadConfig() (*config, error) {
	err := os.MkdirAll(defaultHomeDir, 0755)
	if err != nil {
		// Show a nicer error message if it's because a symlink is
		// linked to a directory that does not exist (probably because
		// it's not mounted).
		if e, ok := err.(*os.PathError); ok && os.IsExist(err) {
			if link, lerr := os.Readlink(e.Path); lerr == nil {
				str := "is symlink %s -> %s mounted?"
				err = fmt.Errorf(str, e.Path, link)
			}
		}

		str := "failed to create home directory: %v"
		err := fmt.Errorf(str, err)
		_, _ = fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	// Default config.
	cfg := config{
		//Listen: normalizeAddress("localhost", defaultListenPort),
		Listen: normalizeAddress("0.0.0.0", defaultListenPort),
	}

	preCfg := cfg
	preParser := flags.NewParser(&preCfg, flags.Default)
	_, err = preParser.Parse()
	if err != nil {
		e, ok := err.(*flags.Error)
		if ok && e.Type == flags.ErrHelp {
			os.Exit(0)
		}
		preParser.WriteHelp(os.Stderr)
		return nil, err
	}

	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	// Load additional config from file.
	parser := flags.NewParser(&cfg, flags.Default)
	// Parse command line options again to ensure they take precedence.
	_, err = parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			parser.WriteHelp(os.Stderr)
		}
		return nil, err
	}

	if len(cfg.Host) == 0 {
		str := "Please specify a hostname"
		err := fmt.Errorf(str)
		_, _ = fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	if len(cfg.Nameserver) == 0 {
		str := "Please specify a nameserver"
		err := fmt.Errorf(str)
		_, _ = fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	if len(cfg.Seeder) == 0 {
		str := "Please specify a seeder"
		err := fmt.Errorf(str)
		_, _ = fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	cfg.Listen = normalizeAddress(cfg.Listen, defaultListenPort)

	if cfg.TestNet {
		activeNetParams = &params.TestNetParams
	}

	return &cfg, nil
}

// normalizeAddress returns addr with the passed default port appended if
// there is not already a port specified.
func normalizeAddress(addr, defaultPort string) string {
	host, port, err := net.SplitHostPort(addr)
	log.Println(host, port)
	if err != nil {
		return net.JoinHostPort(addr, defaultPort)
	}
	return addr
}
