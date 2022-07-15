package main

import (
	"fmt"
	cmd "github.com/go-play/pkg/cmd/server"
	"os"
)

/*
./server -grpc-port=9090 -http-port=8080 -db-host=localhost:3306 -db-user=root -db-password=root@123 -db-schema=todo -log-level=-1 -log-time-format=2006-01-02T15:04:05.999999999Z07:00
*/
func main() {
	if err := cmd.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
