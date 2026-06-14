package main

import (
	"flag"
	"net"
	"strconv"
	"log"
)

type Config struct {
	host string
	port int
}

var config Config

func setUpFlags(){
	flag.StringVar(&config.host, "host", "0.0.0.0", "Host for db server")
	flag.IntVar(&config.port, "port", 7379, "Port number")
	flag.Parse()
}

func runTcpServer(host string, port int){
	var totalClients int = 0

	ln, err := net.Listen("tcp", host + ":" + strconv.Itoa(port))
	if err != nil {
		log.Fatal("Error starting TCP server: ", err)
		panic(err)
	}

	log.Println("TCP server listening on " + host + ":" + strconv.Itoa(port))
	
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
			panic(err)
		}
		totalClients+=1
		log.Println("New client connected: ", conn.RemoteAddr(), "Total Clients", totalClients)

		for {
			cmd, err := readCommand(conn)
			if err != nil {
				conn.Close()
				totalClients -= 1
				log.Println("Client disconnected: ", conn.RemoteAddr(), "Total Clients", totalClients)
				break
			}

			if err := respond(conn, cmd); err != nil {
				log.Println("Error responding to client: ", err)
			}
		}
	}
}

func readCommand(conn net.Conn) (string, error){
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func respond(conn net.Conn, response string) error {
	_, err := conn.Write([]byte(response))
	if(err != nil){
		log.Println("Error writing response to client: ", err)
		return err
	}
	return nil
}


func main() {
	setUpFlags()
	runTcpServer(config.host, config.port)
}