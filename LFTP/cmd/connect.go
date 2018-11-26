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
	"time"

	"github.com/liuyh73/LFTP/LFTP/config"
	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connect command helps us to connect to server.",
	Long:  `We can use LFTP connect to connect to server.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Connecting...")
		// 获取raddr
		serverAddr := host + ":" + port
		raddr, err := net.ResolveUDPAddr("udp", serverAddr)
		checkErr(err)
		// 获取客户端套接字
		clientSocket, err := net.DialUDP("udp", nil, raddr)
		checkErr(err)
		defer clientSocket.Close()
		// 设置等待响应时间
		clientSocket.SetDeadline(time.Now().Add(5 * time.Second))
		// 向服务器发送请求
		_, err = clientSocket.Write([]byte("conn: "))
		checkErr(err)
		// 读取服务器传回的数据
		res := make([]byte, config.CLIENT_RECV_LEN)
		_, err = clientSocket.Read(res)
		checkErr(err)
		resStr := string(res[:bytes.IndexByte(res, 0)])
		fmt.Println(resStr)
		if resStr == "Connected!" {
		}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// connectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// connectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	connectCmd.Flags().StringVarP(&host, "host", "H", config.SERVER_IP, "Server host")
	connectCmd.Flags().StringVarP(&port, "port", "P", config.SERVER_PORT, "Server port")
}
