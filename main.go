package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/conallob/mcp-redfish/internal/redfish"
	"github.com/conallob/mcp-redfish/internal/tools"
	"github.com/mark3labs/mcp-go/server"
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

	client := redfish.NewClient(redfish.Config{
		Host:     *host,
		Username: *username,
		Password: *password,
		Insecure: *insecure,
	})

	s := server.NewMCPServer("mcp-redfish", version)
	tools.Register(s, client, *readOnly)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envBool(key string) bool {
	v := os.Getenv(key)
	return v == "1" || v == "true" || v == "yes"
}
