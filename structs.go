package main

import (
	"fmt"
	"log"
	"net"
)

type App struct {
	Name    string   `json:"Name"`
	Ports   []int    `json:"Ports"`
	Targets []string `json:"Targets"`
}
type Config struct {
	Apps []App `json:"Apps"`
}

type Server struct {
	Port    int      `json:"Ports"`
	Targets []string `json:"Targets"`
}

func (s Server) Listen(channel chan net.Conn) {
	listenAddr := fmt.Sprintf("localhost:%d", s.Port)
	listener, err := net.Listen("tcp", listenAddr)
	handleErrors(fmt.Sprintf("Error occured while listening on address %s", listenAddr), err)

	defer listener.Close()
	for {
		conn, err := listener.Accept()
		log.Println("Recieved a connection on ", conn.LocalAddr().String())
		handleErrors(fmt.Sprintf("Error occured while accepting connections on address %s", listenAddr), err)
		go func() { channel <- conn }()
	}

}
