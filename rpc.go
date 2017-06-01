package main

import (
	"net"
	"net/rpc"
)

func ListenRpc() error {
	if cfg.RpcPort == "" {
		log.Infoln("no rpc port specified, rpc disabled")
		return nil
	}

	ln, err := net.Listen("TCP", cfg.RpcPort)
	if err != nil {
		return err
	}

	rpc.Accept(ln)

	return nil
}
