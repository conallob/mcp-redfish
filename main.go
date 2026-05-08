package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"

	"github.com/conallob/mcp-redfish/internal/redfish"
	"github.com/conallob/mcp-redfish/internal/tools"
)

// version is set at build time via -ldflags "-X main.version=<ver>".
var version = "dev"

func main() {
	var (
		host        = flag.String("host", envOr("REDFISH_HOST", ""), "BMC hostname or IP (required)")
		username    = flag.String("username", envOr("REDFISH_USERNAME", ""), "BMC username")
		password    = flag.String("password", envOr("REDFISH_PASSWORD", ""), "BMC password")
		insecure    = flag.Bool("insecure", envBool("REDFISH_INSECURE"), "Skip TLS certificate verification")
		readOnly    = flag.Bool("read-only", envBool("REDFISH_READ_ONLY"), "Restrict to read-only tools only")
		solURL      = flag.String("sol-url", envOr("REDFISH_SOL_URL", ""), "WebSocket URL for Serial-over-LAN console (e.g. wss://bmc-host/console)")
		showVersion = flag.Bool("version", false, "Print version and exit")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("mcp-redfish %s\n", version)
		os.Exit(0)
	}

	if *host == "" {
		fmt.Fprintln(os.Stderr, "error: BMC host is required (--host or REDFISH_HOST)")
		flag.Usage()
		os.Exit(1)
	}

	cfg := redfish.Config{
		Host:     *host,
		Username: *username,
		Password: *password,
		Insecure: *insecure,
	}
	client := redfish.NewClient(cfg)

	var sol *redfish.SOLClient
	if *solURL != "" {
		sol = redfish.NewSOLClient(cfg, *solURL)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		sol.Start(ctx)
	}

	s := server.NewMCPServer("mcp-redfish", version)
	tools.Register(s, client, sol, *readOnly)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func envOr(key, def string) string { //nolint:unparam
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envBool(key string) bool {
	v := os.Getenv(key)
	return v == "1" || v == "true" || v == "yes"
}
