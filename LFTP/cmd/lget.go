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
	"bytes"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/liuyh73/ftp/LFTP/config"
	"github.com/spf13/cobra"
)

var lgetFile string

// lgetCmd represents the lget command
var lgetCmd = &cobra.Command{
	Use:   "lget",
	Short: "lget command helps us to get a file from server.",
	Long:  `We can use LFTP lget <file> to get a file from server.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("lget called")
		if !connectToServer() {
			return
		}
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
		_, err = clientSocket.Write([]byte("lget: " + lgetFile))
		checkErr(err)
		// 创建文件句柄
		outputFile, err := os.OpenFile(lgetFile, os.O_CREATE|os.O_TRUNC, 0600)
		checkErr(err)
		for {
			res := make([]byte, config.SERVER_RECV_LEN)
			// lenth, err = clientSocket.Read(res)
			// length, remoteAddr *UDPAddr, err = clientSocket.ReadFromUDP(res)
			_, err = clientSocket.Read(res)
			checkErr(err)
			var resStr string
			resStr = string(bytes.TrimRight(res[:], "\x00"))
			//fmt.Println(resStr)
			if resStr == "end" {
				break
			}
			outputFile.Write(res)
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
