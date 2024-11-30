package main

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
}

var (
	mockAddr    = "127.0.0.1:9999"
	mockMsg     = "Hello, world!\n"
	mockTimeout = time.Second * 5
)

type mockReadWriteCloser struct {
	io.Reader
	io.Writer
}

func (m *mockReadWriteCloser) Close() error { return nil }

func startMockServer(t *testing.T, address string, response string) net.Listener {
	t.Helper()
	ln, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatalf("failed to create mock server: %v", err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				c.Write([]byte(response))
			}(conn)
		}
	}()

	return ln
}

func TestNewTelnetClient(t *testing.T) {
	in := &mockReadWriteCloser{strings.NewReader(mockMsg), &strings.Builder{}}
	out := &strings.Builder{}

	client := NewTelnetClient(mockAddr, mockTimeout, in, out)

	if client == nil {
		t.Fatal("NewTelnetClient returned nil")
	}

	if client.(*telnetClient).address != mockAddr {
		t.Errorf("expected address %v, got %v", mockAddr, client.(*telnetClient).address)
	}

	if client.(*telnetClient).timeout != mockTimeout {
		t.Errorf("expected timeout %v, got %v", mockTimeout, client.(*telnetClient).timeout)
	}
}

func TestTelnetClient_Connect(t *testing.T) {
	in := &mockReadWriteCloser{strings.NewReader(mockMsg), &strings.Builder{}}
	out := &strings.Builder{}

	ln := startMockServer(t, mockAddr, mockMsg)
	defer ln.Close()

	client := NewTelnetClient(mockAddr, mockTimeout, in, out)

	if err := client.Connect(); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	if client.(*telnetClient).conn == nil {
		t.Fatal("expected conn to be non-nil after Connect")
	}
}

func TestTelnetClient_Close(t *testing.T) {
	in := &mockReadWriteCloser{strings.NewReader(mockMsg), &strings.Builder{}}
	out := &strings.Builder{}

	client := NewTelnetClient(mockAddr, mockTimeout, in, out)

	connCh := make(chan net.Conn)
	ln, err := net.Listen("tcp", mockAddr)
	if err != nil {
		t.Fatalf("Listen failed: %v", err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Fatalf("Accept failed: %v", err)
			return
		}
		connCh <- conn
	}()

	err = client.Connect()
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	client.(*telnetClient).conn = <-connCh

	err = client.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

func TestTelnetClient_Send(t *testing.T) {
	in := &mockReadWriteCloser{strings.NewReader(mockMsg), &strings.Builder{}}
	out := &strings.Builder{}

	// Mock server to capture sent data
	listener, err := net.Listen("tcp", mockAddr)
	if err != nil {
		t.Fatalf("failed to create mock server: %v", err)
	}
	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		reader := bufio.NewReader(conn)
		received, _ := reader.ReadString('\n')
		if received != mockMsg {
			t.Errorf("expected %v, got %v", mockMsg, received)
		}
	}()

	client := NewTelnetClient(mockAddr, mockTimeout, in, out)
	err = client.Connect()
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	err = client.Send()
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
}

func TestTelnetClient_Receive(t *testing.T) {
	in := &mockReadWriteCloser{strings.NewReader(mockMsg), &strings.Builder{}}
	out := &strings.Builder{}

	ln := startMockServer(t, mockAddr, mockMsg)
	defer ln.Close()

	client := NewTelnetClient(mockAddr, mockTimeout, in, out)
	err := client.Connect()
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer client.Close()

	err = client.Receive()
	if err != nil {
		t.Fatalf("Receive failed: %v", err)
	}

	if out.String() != mockMsg {
		t.Errorf("expected %v, got %v", mockMsg, out.String())
	}
}
