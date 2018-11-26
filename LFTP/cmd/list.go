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
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/liuyh73/LFTP/LFTP/config"
	"github.com/spf13/cobra"
)

var limit int

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list command helps us to get the file list from server.",
	Long:  `We can use LFTP list to get the file list from server.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !connectToServer() {
			fmt.Println("You don't connect to the server.")
			return
		}
		// 方法一：
		// 获取clientSocket可以使用net.Dial("udp", address string)
		// serverAddr := host + ":" + port
		// clientSocket, err := net.Dial("udp", serverAddr)
		// 方法二：
		serverAddr := host + ":" + port
		raddr, err := net.ResolveUDPAddr("udp", serverAddr)
		checkErr(err)
		// net.DialUDP("udp", localAddr *UDPAddr, remoteAddr *UDPAddr)
		clientSocket, err := net.DialUDP("udp", nil, raddr)
		checkErr(err)
		defer clientSocket.Close()
		clientSocket.SetDeadline(time.Now().Add(5 * time.Second))
		_, err = clientSocket.Write([]byte("list: " + strconv.Itoa(limit) + " files"))
		checkErr(err)
		data := make([]byte, 1024)
		// lenth, err = clientSocket.Read(data)
		// length, remoteAddr *UDPAddr, err = clientSocket.ReadFromUDP(data)
		_, _, err = clientSocket.ReadFromUDP(data)
		checkErr(err)
		fmt.Println(string(data))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listCmd.Flags().IntVarP(&limit, "limit", "l", 10, "get the list with a number of `limit` files.")
	listCmd.Flags().StringVarP(&host, "host", "H", config.SERVER_IP, "Server host")
	listCmd.Flags().StringVarP(&port, "port", "P", config.SERVER_PORT, "Server port")
}
