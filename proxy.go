package main

import (
	"net"
	"strconv"
)

type Address struct {
	ip string
	port int
}

func resolveAddress(addr Address) *net.TCPAddr {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr.ip + ":" + strconv.Itoa(addr.port))
	if (err != nil) {
		return nil
	}
	return tcpAddr
}

func connect(addr *net.TCPAddr) *net.TCPConn {
	conn, err := net.DialTCP("tcp", nil, addr)
	if (err != nil) {
		return nil
	}

	conn.SetReadBuffer(1024)
	return conn
}

func relay(input *net.TCPConn, output *net.TCPConn) {
	defer output.Close()

	buff := make([]byte, 1024)
	for {
		read, err := input.Read(buff)
		if (read == 0) { return }
		if (err != nil) {
			return
		}

		output.Write(buff[:read])
	}
}

func StartConnection(input Address, output Address) {
	inputAddr := resolveAddress(input)
	outputAddr := resolveAddress(output)
	if (inputAddr == nil || outputAddr == nil) { return }

	inputConn := connect(inputAddr)
	if (inputConn == nil) { return }
	outputConn := connect(outputAddr)
	if (outputConn == nil) {
		inputConn.Close()
		return
	}

	go relay(inputConn, outputConn)
	go relay(outputConn, inputConn)
}
