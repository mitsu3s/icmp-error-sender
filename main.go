// main.go
package main

import (
	"fmt"
	"os"

	"github.com/mitsu3s/icmp-error-sender/internal"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <error_type> <srcIP> <dstIP>")
		os.Exit(1)
	}

	errorType := os.Args[1]
	srcIP := os.Args[2]
	dstIP := os.Args[3]

	switch errorType {
	case "redirect":
		internal.Redirect(srcIP, dstIP)
	default:
		fmt.Printf("Unknown error type: %s\n", errorType)
		os.Exit(1)
	}
}
