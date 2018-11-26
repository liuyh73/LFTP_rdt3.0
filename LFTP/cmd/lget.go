// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/liuyh73/LFTP/LFTP/models"
	"strconv"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/liuyh73/LFTP/LFTP/config"
	"github.com/spf13/cobra"
)

var lgetFile string

// lgetCmd represents the lget command
var lgetCmd = &cobra.Command{
	Use:   "lget",
	Short: "lget command helps us to get a file from server.",
	Long:  `We can use LFTP lget <file> to get a file from server.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !connectToServer() {
			return
		}
		lgetPacket := models.NewPacket(byte(0), byte(0), byte(0), []byte("lget: " + lgetFile))
		fmt.Println(lgetPacket)
		fmt.Println(lgetFile)
		// 获取raddr
		serverAddr := host + ":" + port
		raddr, err := net.ResolveUDPAddr("udp", serverAddr)
		checkErr(err)
		// 获取客户端套接字
		// net.DialUDP("udp", localAddr *UDPAddr, remoteAddr *UDPAddr)
		clientSocket, err := net.DialUDP("udp", nil, raddr)
		checkErr(err)
		defer clientSocket.Close()
		// 设置等待响应时间
		clientSocket.SetDeadline(time.Now().Add(10 * time.Second))
		// 向服务器发送请求
		_, err = clientSocket.Write(lgetPacket.ToBytes())
		checkErr(err)
		// 创建文件句柄
		outputFile, err := os.OpenFile(lgetFile, os.O_CREATE|os.O_TRUNC, 0600)
		checkErr(err)
		for {
			buf := make([]byte, config.CLIENT_RECV_LEN)
			packet := &models.Packet{}
			// lenth, err = clientSocket.Read(res)
			// length, remoteAddr *UDPAddr, err = clientSocket.ReadFromUDP(res)
			_, err = clientSocket.Read(buf)
			checkErr(err)
			packet.FromBytes(buf)
			//fmt.Println(resStr)
			if packet.Finished == byte(1) {
				fmt.Println("end")
				break
			}
			length, err := outputFile.Write(packet.Data)
			fmt.Println("Read lenth: "+strconv.Itoa(length))
			checkErr(err)
		}
		fmt.Printf("Finished to download the file %s.\n", lgetFile)
	},
}

func init() {
	rootCmd.AddCommand(lgetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lgetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lgetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	lgetCmd.Flags().StringVarP(&lgetFile, "file", "f", "", "lgetfile filename")
	lgetCmd.Flags().StringVarP(&host, "host", "H", config.SERVER_IP, "Server host")
	lgetCmd.Flags().StringVarP(&port, "port", "P", config.SERVER_PORT, "Server port")
}
