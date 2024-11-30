package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *telnetClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *telnetClient) Send() error {
	reader := bufio.NewReader(c.in)
	for {
		input, err := reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}

		_, err = c.conn.Write([]byte(input))
		if err != nil {
			return err
		}
	}
}

func (c *telnetClient) Receive() error {
	reader := bufio.NewReader(c.conn)
	for {
		response, err := reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return err
		}
		_, err = fmt.Fprint(c.out, response)
		if err != nil {
			return err
		}
	}
}
