package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "Connection timeout")

	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Printf("usage: %s [--timeout=<duration>] <host> <port>", os.Args[0])
		return
	}

	host := args[0]
	port, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("invalid port: %v", err)
		return
	}

	client := NewTelnetClient(fmt.Sprintf("%s:%d", host, port), *timeout, os.Stdin, os.Stdout)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-signals:
			cancel()
			client.Close()
		case <-ctx.Done():
		}
	}()

	if err := client.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect:", err)
		return
	}

	go func() {
		if err := client.Send(); err != nil {
			fmt.Fprintln(os.Stderr, "Send error:", err)
			cancel()
		}
	}()

	if err := client.Receive(); err != nil {
		fmt.Fprintln(os.Stderr, "Receive error:", err)
		cancel()
	}
}
