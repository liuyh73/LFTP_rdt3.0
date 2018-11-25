package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const (
	server_ip       = "127.0.0.1"
	server_port     = "8808"
	server_recv_len = 200
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {
	serverAddr := server_ip + ":" + server_port
	serverUDPAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	checkErr(err)

	serverSocket, err := net.ListenUDP("udp", serverUDPAddr)
	checkErr(err)
	defer serverSocket.Close()

	for {
		data := make([]byte, server_recv_len)
		_, clientUDPAddr, err := serverSocket.ReadFromUDP(data)

		checkErr(err)

		dataStr := string(data[:bytes.IndexByte(data, 0)])
		fmt.Println("Received:", dataStr)
		if strings.Split(dataStr, ": ")[0] == "conn" {
			handleConn(serverSocket, clientUDPAddr)
		} else if strings.Split(dataStr, ": ")[0] == "lget" {
			handleGetFile(serverSocket, clientUDPAddr, strings.Split(dataStr, ": ")[1])
		} else if strings.Split(dataStr, ": ")[0] == "lsend" {
			handlePutFile(serverSocket, clientUDPAddr)
		} else if strings.Split(dataStr, ": ")[0] == "list" {
			handleList(serverSocket, clientUDPAddr)
		}
	}
}

func handleConn(serverSocket *net.UDPConn, clientUDPAddr *net.UDPAddr) {
	_, err := serverSocket.WriteToUDP([]byte("Connected!"), clientUDPAddr)
	checkErr(err)
	fmt.Println("Connected to " + clientUDPAddr.String())
}

func handleGetFile(serverSocket *net.UDPConn, clientUDPAddr *net.UDPAddr, pathname string) {
	_, err := os.Stat(pathname)

	if os.IsNotExist(err) {
		fmt.Printf("The file %s doesn't exist", pathname)
		serverSocket.WriteToUDP([]byte(fmt.Sprintf("The file %s doesn't exist", pathname)), clientUDPAddr)
		return
	}
	file, err := os.Open(pathname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred on opening the inputfile: %s\nDoes the file exist?\n", pathname)
	}
	defer file.Close()
	for {
		data := make([]byte, 200)
		_, err1 := file.Read(data)
		_, err2 := serverSocket.WriteToUDP(data, clientUDPAddr)
		if err1 == io.EOF {
			//serverSocket.WriteToUDP([]byte("end"), clientUDPAddr)
			fmt.Printf("Finished to download the file %s.\n", file.Name())
			break
		}

		if err1 != nil {
			fmt.Println(err1)
		}

		if err2 != nil {
			fmt.Println(err2)
		}
	}
}

func handlePutFile(serverSocket *net.UDPConn, clientUDPAddr *net.UDPAddr) {

}

func handleList(serverSocket *net.UDPConn, clientUDPAddr *net.UDPAddr) {

}
