package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"

	"golang.org/x/sys/unix"
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
			// log.Println("Command before decoding is:",cmd)
			decodedCommand,err := decode([]byte(cmd))

			// log.Println("Decoded Command is", decodedCommand)

			// log.Println("Received command from client", conn.RemoteAddr(), cmd)
			// if err := respond(conn, cmd); err != nil {
			// 	log.Println("Error responding to client: ", err)
			// }

			//here we will handle exact commands

			// cmdParams, err := findCmdParams(decodedCommand)

			res, err := execCommand(decodedCommand)

			if err != nil {
				fmt.Print("Failed to execute Command")
			}

			output,_ := encodeToSimpleString(res, err)

			respond(conn, output)
		}
	}
}

// func runAsyncTcpServe(host string, port int){
// 	ln, err := net.Listen("tcp", host + ":" + strconv.Itoa(port))
// 	if err != nil {
// 		println("Error in creating listner socket for server")
// 		panic(err)
// 	}
// 	tcpListener := ln.(*net.TCPListener)
// 	listnerFile, err := tcpListener.File()

// 	if err != nil {
// 		println("Error in getting file for the listner socket")
// 		panic(err)
// 	}

// 	listnerFd := int(listnerFile.Fd())
// 	defer unix.Close(listnerFd)

// 	if err := unix.SetNonblock(listnerFd, true); err != nil {
// 		log.Fatalf("Failed to set non-blocking: %v", err)
// 	}

// 	epollfd, err := unix.EpollCreate1(0)

// 	defer unix.Close(epollfd)

// 	// address of the listner socket event that we will register in epoll
// 	listnerSocketEvent := &unix.EpollEvent{
// 		Events: unix.EPOLLIN,
// 		Fd: int32(listnerFd),
// 	}

// 	err = unix.EpollCtl(epollfd, unix.EPOLL_CTL_ADD, listnerFd, listnerSocketEvent)

// 	if err != nil {
// 		println("Error registering listner socket in epoll")
// 		panic("err")
// 	}

// 	events := make([]unix.EpollEvent, 20000)

// 	for {
// 		numEvents, err := unix.EpollWait(epollfd, events, -1)
// 		if err != nil {
// 			if err == unix.EINTR {
// 				continue;
// 			}
// 			println("Error while checking sockets ready for io from epoll");
// 		}

// 		for i:=0; i<numEvents; i++ {
// 			currentFd := int(events[i].Fd)

// 			if(currentFd == listnerFd){
// 				clientFd, _, err := unix.Accept(listnerFd)
// 				if err != nil {
// 					fmt.Println("Error while accepting the connection and making a new socket for new client")
// 					continue;
// 				}

// 				unix.SetNonblock(clientFd, true)

// 				clientEvent := &unix.EpollEvent{
// 					Events: unix.EPOLLIN,
// 					Fd:     int32(clientFd),
// 				}

// 				// register the new client to epoll
// 				err = unix.EpollCtl(epollfd, unix.EPOLL_CTL_ADD, clientFd, clientEvent)

// 				if err != nil {
// 					fmt.Println("Error registering the client socket to epoll")
// 				}
// 			} else {
// 				buffer := make([]byte, 1024)
// 				n, err := unix.Read(currentFd, buffer)
// 				if n<=0 || (err != nil) {
// 					fmt.Printf("Client on FD %d disconnected or errored. Cleaning up.\n", currentFd)
// 					unix.EpollCtl(epollfd, unix.EPOLL_CTL_DEL, currentFd, nil)
// 					unix.Close(currentFd)
// 					continue 
// 				}

// 				bufferString := string(buffer[:n])


// 				// log.Println("Command before decoding is:",cmd)
// 				decodedCommand,err := decode([]byte(bufferString))

// 				// log.Println("Decoded Command is", decodedCommand)

// 				// log.Println("Received command from client", conn.RemoteAddr(), cmd)
// 				// if err := respond(conn, cmd); err != nil {
// 				// 	log.Println("Error responding to client: ", err)
// 				// }

// 				//here we will handle exact commands

// 				// cmdParams, err := findCmdParams(decodedCommand)

// 				res, err := execCommand(decodedCommand)

// 				if(err != nil){
// 					fmt.Print("Failed to execute Command")
// 					continue;
// 				}

// 				output,_ := encodeToSimpleString(res, err)

// 				err = unix.Write(currentFd, []byte(output))

// 				if(err != nil){
// 					log.Println("Error writing response to client: ", err)
// 					unix.EpollCtl(epollfd, unix.EPOLL_CTL_DEL, currentFd, nil)
// 					unix.Close(currentFd)
// 				}
// 			}
// 		}
// 	}
// }

func runAsyncTcpServe(host string, port int) {
	ln, err := net.Listen("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		println("Error in creating listner socket for server")
		panic(err)
	}
	tcpListener := ln.(*net.TCPListener)
	listnerFile, err := tcpListener.File()

	if err != nil {
		println("Error in getting file for the listner socket")
		panic(err)
	}

	listnerFd := int(listnerFile.Fd())
	defer unix.Close(listnerFd)

	if err := unix.SetNonblock(listnerFd, true); err != nil {
		log.Fatalf("Failed to set non-blocking: %v", err)
	}
	log.Println("TCP server listening on " + host + ":" + strconv.Itoa(port))

	kqfd, err := unix.Kqueue()
	if err != nil {
		log.Fatalf("Failed to create kqueue: %v", err)
	}
	defer unix.Close(kqfd)

	// address of the listner socket event that we will register in kqueue
	listnerSocketEvent := unix.Kevent_t{
		Ident:  uint64(listnerFd),
		Filter: unix.EVFILT_READ,
		Flags:  unix.EV_ADD,
	}

	_, err = unix.Kevent(kqfd, []unix.Kevent_t{listnerSocketEvent}, nil, nil)

	if err != nil {
		println("Error registering listner socket in kqueue")
		panic("err")
	}

	events := make([]unix.Kevent_t, 20000)

	for {
		numEvents, err := unix.Kevent(kqfd, nil, events, nil)
		if err != nil {
			if err == unix.EINTR {
				continue
			}
			println("Error while checking sockets ready for io from kqueue")
		}

		for i := 0; i < numEvents; i++ {
			currentFd := int(events[i].Ident)

			if currentFd == listnerFd {
				clientFd, _, err := unix.Accept(listnerFd)
				if err != nil {
					fmt.Println("Error while accepting the connection and making a new socket for new client")
					continue
				}

				unix.SetNonblock(clientFd, true)

				clientEvent := unix.Kevent_t{
					Ident:  uint64(clientFd),
					Filter: unix.EVFILT_READ,
					Flags:  unix.EV_ADD,
				}

				// register the new client to kqueue
				_, err = unix.Kevent(kqfd, []unix.Kevent_t{clientEvent}, nil, nil)

				if err != nil {
					fmt.Println("Error registering the client socket to kqueue")
				}
			} else {
				buffer := make([]byte, 1024)
				n, err := unix.Read(currentFd, buffer)
				if n <= 0 || (err != nil) {
					// fmt.Printf("Client on FD %d disconnected or errored. Cleaning up.\n", currentFd)
					delEvent := unix.Kevent_t{
						Ident:  uint64(currentFd),
						Filter: unix.EVFILT_READ,
						Flags:  unix.EV_DELETE,
					}
					unix.Kevent(kqfd, []unix.Kevent_t{delEvent}, nil, nil)
					unix.Close(currentFd)
					continue
				}

				bufferString := string(buffer[:n])

				// log.Println("Command before decoding is:",cmd)
				decodedCommand, err := decode([]byte(bufferString))

				// log.Println("Decoded Command is", decodedCommand)

				// log.Println("Received command from client", conn.RemoteAddr(), cmd)
				// if err := respond(conn, cmd); err != nil {
				// 	log.Println("Error responding to client: ", err)
				// }

				//here we will handle exact commands

				// cmdParams, err := findCmdParams(decodedCommand)

				res, err := execCommand(decodedCommand)

				// if err != nil {
				// 	fmt.Print("Failed to execute Command")
					// continue
				// }

				output, _ := encodeToSimpleString(res, err)

				_, err = unix.Write(currentFd, []byte(output))

				if err != nil {
					log.Println("Error writing response to client: ", err)
					delEvent := unix.Kevent_t{
						Ident:  uint64(currentFd),
						Filter: unix.EVFILT_READ,
						Flags:  unix.EV_DELETE,
					}
					unix.Kevent(kqfd, []unix.Kevent_t{delEvent}, nil, nil)
					unix.Close(currentFd)
				}
			}
		}
	}
}



func findCmdParams(cmd string) ([]string, error){

	arrayCmd := []string{}

	currString := ""

	for i:=0 ; i<len(cmd);i++{
		if cmd[i] == ' ' {
			if currString != "" {
				arrayCmd = append(arrayCmd, currString)
			}
			currString = ""
			continue;
		}
		currString = currString + string(cmd[i])
	}
	if currString != "" {
		arrayCmd = append(arrayCmd, currString)
	}
	return arrayCmd, nil
}

func readCommand(conn net.Conn) (string, error){
	buf := make([]byte, 1024)
	// copy the data from the 
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
	// runTcpServer(config.host, config.port)
	runAsyncTcpServe(config.host, config.port);
}