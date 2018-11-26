package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/liuyh73/LFTP/LFTP/models"
)

const (
	server_ip       = "127.0.0.1"
	server_port     = "8808"
	server_send_len = 1992
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
	packet := models.NewPacket(byte(0), byte(0), byte(0), byte(0), []byte("Connected!"))
	_, err := serverSocket.WriteToUDP(packet.ToBytes(), clientUDPAddr)
	checkErr(err)
	fmt.Println("Connected to " + clientUDPAddr.String())
}

func handleGetFile(serverSocket *net.UDPConn, clientUDPAddr *net.UDPAddr, pathname string) {
	_, err := os.Stat(pathname)
	// serverSocket.SetDeadline(time.Now().Add(10 * time.Second))
	// lget file不存在
	if os.IsNotExist(err) {
		fmt.Printf("The file %s doesn't exist", pathname)
		packetSnd := models.NewPacket(byte(0), byte(0), byte(0), byte(0), []byte(fmt.Sprintf("The file %s doesn't exist", pathname)))
		serverSocket.WriteToUDP(packetSnd.ToBytes(), clientUDPAddr)
		return
	}
	// 打开该文件
	file, err := os.Open(pathname)
	defer file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred on opening the inputfile: %s\nDoes the file exist?\n", pathname)
		return
	}
	for {
		var err1, err2 error
		packetRcv := &models.Packet{}
		var packetSnd *models.Packet
		// 发送packet 0
		buf := make([]byte, server_send_len)
		_, err1 = file.Read(buf)
		if err1 == io.EOF {
			packetSnd = models.NewPacket(byte(0), byte(0), byte(1), byte(1), buf)
			_, err2 = serverSocket.WriteToUDP(packetSnd.ToBytes(), clientUDPAddr)
			fmt.Println("Write Length:" + strconv.Itoa(int(packetSnd.Length)))
			checkErr(err2)
			fmt.Printf("Finished to download the file %s.\n", file.Name())
			// 设置定时器
			break
		}

		packetSnd = models.NewPacket(byte(0), byte(0), byte(1), byte(1), buf)
		_, err2 = serverSocket.WriteToUDP(packetSnd.ToBytes(), clientUDPAddr)
		fmt.Println("Write Length:" + strconv.Itoa(int(packetSnd.Length)))
		checkErr(err2)
		// 设置定时器

		// 等待ACK 0
		for {
			buf := make([]byte, server_recv_len)
			_, err1 = serverSocket.Read(buf)
			checkErr(err1)
			packetRcv.FromBytes(buf)
			if packetRcv.ACK == 0 {
				break
			}
			// TODO 如果超时，重新发送数据包packetSnd, 设置定时器
		}

		// ACK == 0
		// 取消定时器

		// 等待来自上层的调用1
		// 发送packet 1
		buf = make([]byte, server_send_len)
		_, err1 = file.Read(buf)
		if err1 == io.EOF {
			packetSnd = models.NewPacket(byte(1), byte(0), byte(1), byte(1), buf)
			_, err2 = serverSocket.WriteToUDP(packetSnd.ToBytes(), clientUDPAddr)
			fmt.Println("Write Length:" + strconv.Itoa(int(packetSnd.Length)))
			checkErr(err2)
			fmt.Printf("Finished to download the file %s.\n", file.Name())
			// 设置定时器
			break
		}
		packetSnd = models.NewPacket(byte(1), byte(0), byte(1), byte(0), buf)
		_, err2 = serverSocket.WriteToUDP(packetSnd.ToBytes(), clientUDPAddr)
		fmt.Println("Write Length:" + strconv.Itoa(int(packetSnd.Length)))
		checkErr(err2)
		// 设置定时器

		// 等待ACK 1
		for {
			buf := make([]byte, server_recv_len)
			_, err1 = serverSocket.Read(buf)
			checkErr(err1)
			packetRcv.FromBytes(buf)
			if packetRcv.ACK == 1 {
				break
			}
			// TODO 如果超时，重新发送数据包packetSnd, 设置定时器
		}

		// ACK == 1
		// 取消定时器
	}
}

func handlePutFile(serverSocket *net.UDPConn, clientUDPAddr *net.UDPAddr) {

}

func handleList(serverSocket *net.UDPConn, clientUDPAddr *net.UDPAddr) {

}
