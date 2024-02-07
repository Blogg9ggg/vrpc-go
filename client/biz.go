package client

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"
	"vrpc/codec"
)

func parseOptions(opts ...*codec.Option) (*codec.Option, error) {
	// if opts is nil or pass nil as parameter
	if len(opts) == 0 || opts[0] == nil {
		return codec.DefaultOption, nil
	}

	if len(opts) != 1 {
		return nil, errors.New("number of options is more than 1")
	}

	opt := opts[0]
	opt.MagicNumber = codec.DefaultOption.MagicNumber
	if opt.CodecType == "" {
		opt.CodecType = codec.DefaultOption.CodecType
	}

	log.Println("parseOptions:", opt)

	return opt, nil
}

// Dial connects to an RPC server at the specified network address
func Dial(network, address string, opts ...*codec.Option) (client *Client, err error) {
	return dialTimeout(NewClient, network, address, opts...)
}

type newClientFunc func(conn net.Conn, opt *codec.Option) (client *Client, err error)
type clientResult struct {
	client *Client
	err    error
}

func dialTimeout(f newClientFunc, network, address string, opts ...*codec.Option) (client *Client, err error) {
	opt, err := parseOptions(opts...)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTimeout(network, address, opt.ConnectTimeout)
	if err != nil {
		return nil, err
	}

	// close the connection if client is nil
	defer func() {
		if client == nil {
			_ = conn.Close()
		}
	}()

	ch := make(chan clientResult)
	go func() {
		client, err := f(conn, opt)
		ch <- clientResult{client: client, err: err}
	}()

	log.Println("opt.ConnectTimeout:", opt.HandleTimeout)
	if opt.HandleTimeout == 0 {
		result := <-ch
		return result.client, result.err
	}

	select {
	case <-time.After(opt.HandleTimeout):
		return nil, fmt.Errorf("rpc client: connect timeout: expect within %s", opt.ConnectTimeout)
	case result := <-ch:
		return result.client, result.err
	}
}
