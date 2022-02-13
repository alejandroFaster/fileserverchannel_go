package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"tcp_fileserver/client"
	"tcp_fileserver/server"
	"tcp_fileserver/utils"
	"time"
) //When the package file is the same of this go file
//you can create folders with go files and add them to the code using the function name
//without add on this import zone

func showConsoleError() {
	fmt.Println("Operation invalid.")
	fmt.Println("Server mode must be started like this:")
	fmt.Println("./app start <port>")
	fmt.Println("where <port> is any port number to listen TCP connections")
	fmt.Println("\nClient mode must be started like this:")
	fmt.Println("./app receive <port> -channel <number> ")
	fmt.Println("where <number> is any number channel to receive files and <port> the port number to connect to the TCP Server")
	fmt.Println("Also, client mode must be started like this:")
	fmt.Println("./app send <filename.ext> <ip> <port> -channel <number>")
	fmt.Println("where <number> is any number channel to receive files and <port> the port number to connect to the TCP Server")
}

func main() {

	if len(os.Args) < 2 {
		showConsoleError()
		return
	}

	var mode = strings.ToLower(os.Args[1])

	if mode != "start" && mode != "receive" && mode != "send" {
		showConsoleError()
		return
	}

	var port int
	var channel int
	var err error

	switch mode {
	//server mode
	case "start":
		if len(os.Args) != 3 {
			log.Fatal("Server mode must be one parameter like this:")
			log.Fatal("./app start <port>")
			return
		}
		port, err = utils.Validate_port(os.Args[2]) //I dont need import the validateParams.go, because the package file is the same
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		ip := utils.GetLocalIP()
		now := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[%s] Starting server (%s) on port %d...\n", now, ip, port)
		//TODO: Starting server...
		server.Start(port)
	//client mode
	case "receive":
		if len(os.Args) != 6 {
			log.Fatal("Client mode must be 3 parameters to receive like this:")
			log.Fatal("./app receive <ip> <port> -channel <number> ")
			return
		}

		var ip = os.Args[2]
		err = utils.ValidateIPAddress(ip)
		if err != nil {
			log.Fatal(err)
			return
		}

		port, err = utils.Validate_port(os.Args[3])
		if err != nil {
			log.Fatal(err.Error())
			return
		}

		channel, err = strconv.Atoi(os.Args[5])
		if err != nil {
			log.Fatal("Error to read channel number. The channel must be a number.")
			return
		}
		fmt.Printf("Receive client file on IP Server, %s, port %d and channel %d...\n", ip, port, channel)
		//Receive client from channel...
		errDownload := client.DownloadFile(ip, port, channel)
		if errDownload != nil {
			log.Fatal("Error to send file: ", errDownload.Error())
			return
		}

	case "send":
		if len(os.Args) != 7 {
			log.Fatal("Client mode to send must have 4 parameters like this:")
			log.Fatal("./app send <filename.ext> <ip> <port> -channel <number>")
			return
		}
		var filename = os.Args[2]
		var ip = os.Args[3]
		err = utils.ValidateIPAddress(ip)
		if err != nil {
			log.Fatal(err)
			return
		}

		port, err = utils.Validate_port(os.Args[4])
		if err != nil {
			log.Fatal(err)
			return
		}

		channel, err = strconv.Atoi(os.Args[6])
		if err != nil {
			log.Fatal("Error to convert channel number. The channel must be a number.")
			return
		}
		fmt.Printf("Send client file called '%s', to IP address %s, TCP port %d on channel %d...\n", filename, ip, port, channel)

		errSend := client.SendFile(filename, ip, port, channel)
		if errSend != nil {
			log.Fatal("Error to send file: ", errSend.Error())
			return
		}

		fmt.Printf("File sended successfully called '%s', to IP address %s, TCP port %d on channel %d\n", filename, ip, port, channel)

	}
}
