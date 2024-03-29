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
	"os"
	"time"

	"github.com/liuyh73/LFTP/LFTP/models"

	"github.com/liuyh73/LFTP/LFTP/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var host string
var port string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "LFTP",
	Short: "A network application, LFTP, to support large file transfer.",
	Long: `LFTP users User Datagram Protocol(UDP) as the transport layer protocol, but it transfers files reliably like tcp.
		LFTP also implements flow control function and congestion control function similar as TCP.
		LFTP server sid supports multiple clients at the same time and we can check the debug information when programs are exected.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.LFTP.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".LFTP" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".LFTP")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func connectToServer() bool {
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

	packetConn := models.NewPacket(byte(0), byte(0), byte(1), byte(0), []byte("conn: "+host))

	_, err = clientSocket.Write(packetConn.ToBytes())
	checkErr(err)
	// 读取服务器传回的数据
	res := make([]byte, config.CLIENT_RECV_LEN)
	_, err = clientSocket.Read(res)
	checkErr(err)

	packet := &models.Packet{}
	packet.FromBytes(res)
	fmt.Println(string(packet.Data))
	if string(packet.Data) == "Connected!" {
		return true
	}
	return false
}
