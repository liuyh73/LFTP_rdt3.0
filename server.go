package main

import (
	"github.com/liuyh73/LFTP/LFTP/models"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"strconv"
)

const (
	server_ip       = "127.0.0.1"
	server_port     = "8808"
	server_send_len = 1993
	server_recv_len = 2000
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
		buf := make([]byte, server_recv_len)
		_, clientUDPAddr, err := serverSocket.ReadFromUDP(buf)

		checkErr(err)
		packet := &models.Packet{}
		packet.FromBytes(buf)
		fmt.Println(packet)
		dataStr := string(packet.Data)
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
	packet := models.NewPacket(byte(0), byte(0), byte(0), []byte("Connected!"))
	_, err := serverSocket.WriteToUDP(packet.ToBytes(), clientUDPAddr)
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
		var err1, err2 error
		buf := make([]byte, server_send_len)
		_, err1 = file.Read(buf)
		if err1 == io.EOF {
			packet := models.NewPacket(byte(0), byte(0), byte(1), buf)
			_, err2 = serverSocket.WriteToUDP(packet.ToBytes(), clientUDPAddr)
			fmt.Println("Write Length:"+strconv.Itoa(int(packet.Length)))
			if err2 != nil {
				fmt.Println(err2)
			}
			fmt.Printf("Finished to download the file %s.\n", file.Name())
			break
		}
		packet := models.NewPacket(byte(0), byte(0), byte(0), buf)
		_, err2 = serverSocket.WriteToUDP(packet.ToBytes(), clientUDPAddr)
		fmt.Println("Write Length:"+strconv.Itoa(int(packet.Length)))
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
