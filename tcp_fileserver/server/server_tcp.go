package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"tcp_fileserver/utils"
)

const BUFFER_SIZE = 1024

var channelFiles map[int]ChannelFile

const dbFilename = "channels.gob"

func Start(port int) error {
	ip := utils.GetLocalIP()
	server, errStart := net.Listen("tcp", ip+":"+strconv.Itoa(port))
	if errStart != nil {
		return errors.New("There was an error starting the server" + errStart.Error())
	}

	//load ChannelFiles saved in db
	var errChannel error
	channelFiles, errChannel = LoadChannels(dbFilename)

	if errChannel != nil {
		errDetail := errChannel.(*utils.ErrorF)
		if errDetail.CodeError != 3 { //if error is different from "file doesnt exist", show error message
			return errors.New("There was an error starting the server" + errChannel.Error())
		} else {
			//initialize channels map
			channelFiles = make(map[int]ChannelFile)
		}
	}

	//TODO: loop until disconnect
	for {
		connection, errConn := server.Accept()
		if errConn != nil {
			return errors.New("There was a error with the connection" + errConn.Error())
		}
		fmt.Println("connected")
		go connectionHandler(connection)
	}
}

var validSubscribe = regexp.MustCompile(`(subscribe)[ |\t]+(\d+)`)
var validSend = regexp.MustCompile(`(send)[ |\t]+([\d]+)+[ |\t]+["]([a-zA-Z0-9\s_\\.\-\(\):]+)["]`)

func connectionHandler(connection net.Conn) error {
	buffer := make([]byte, BUFFER_SIZE)
	_, errReading := connection.Read(buffer)
	if errReading != nil {
		return errors.New("There is an error reading from connection" + errReading.Error())
	}

	strBuffer := string(buffer)
	command := strings.Split(strBuffer, " ")
	switch command[0] {
	case "subscribe":
		groups := validSubscribe.FindStringSubmatch(strBuffer)
		if groups == nil {
			utils.SendError(connection, "Bad format of subscribe command. Usage: \n"+
				"'subscribe <channel>'\n"+
				"Where channel is a number to receive the file")
			return nil
		}

		channel, _ := strconv.Atoi(groups[2])

		//check if the channel exists in DB file channel
		channelFile, ok := channelFiles[channel]
		if !ok {
			utils.SendError(connection, "Channel does not exist in DB file")
			return nil
		}
		fmt.Println("Sending the filename info from the db channel")
		utils.SendMessage(connection, fmt.Sprintf(`{"filename": "%s"}`, channelFiles[channel].Filename), false)

		fmt.Printf("Sending file to client on channel %d...\n", channel)
		//Send file to client on channel %d
		err := sendFiletoClient(connection, channelFiles[channel].FullPath)
		if err != nil {
			fmt.Printf("Error sending file: %v\n", err)
			//utils.SendError(connection, err.Error())
			connection.Close()
			return nil
		}
		//if send file to client on channel is successfully, increment Downloaded and save to database
		channelFile.Downloaded += 1
		channelFiles[channel] = channelFile

		var errChannel error
		errChannel = SaveChannels(channelFiles, dbFilename)
		if errChannel != nil {
			fmt.Printf("Error saving channels file: %v\n", errChannel)
			return nil
		}

		connection.Close()
	case "send":
		//fmt.Println(strBuffer)
		groups := validSend.FindStringSubmatch(strBuffer)
		if groups == nil {
			utils.SendError(connection, `Bad format of send command. Usage: 
										'send <channel> "<filename>"'
										Where channel is a number and "filename" the filename to send (not full path, only filename) (quoted enclosed) on the server`)
			return nil
		}
		channel, _ := strconv.Atoi(groups[2])
		filename := groups[3]
		currentPath, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		//verify if the channel doesnt exists
		if _, ok := channelFiles[channel]; !ok {
			//add new channel to map of channels
			channelFiles[channel] = ChannelFile{
				Channel:    channel,
				Filename:   filename,
				FullPath:   currentPath + "/filesServed/" + filename,
				Downloaded: 0,
			}
		} else {
			//channel file exists already, send error
			utils.SendError(connection, fmt.Sprintf("The channel file %d is already in use", channel))
			return nil
		}

		//save in DB file channel
		var errChannel error
		errChannel = SaveChannels(channelFiles, dbFilename)
		if errChannel != nil {
			utils.SendError(connection, errChannel.Error())
			return nil
		}

		fmt.Printf("Saving file from client to channel %d...\n", channel)
		//Saving file from client to channel
		err = saveFilefromClient(connection, channelFiles[channel].FullPath)
		if err != nil {
			utils.SendError(connection, err.Error())
		}
		connection.Close()
	default:
		utils.SendError(connection, "Bad command\n")
	}
	//}
	connection.Close()
	return nil
}

func sendFiletoClient(connection net.Conn, filepath string) error {
	file, err2 := os.Open(strings.TrimSpace(filepath)) // For read access.
	if err2 != nil {
		return utils.NewError(10, "Error reading file from Server: "+err2.Error())
	}
	defer file.Close()

	n, err2 := io.Copy(connection, file)
	if err2 != nil {
		return utils.NewError(11, "Error copying file from Server: "+err2.Error())
	}
	fmt.Printf("Bytes send to the client: %v\n", n)
	return nil
}

func saveFilefromClient(connection net.Conn, filepath string) error {
	file, err := os.Create(strings.TrimSpace(filepath)) // For read access.
	if err != nil {
		return errors.New("Error creating file from Client: " + err.Error())
	}
	defer file.Close()

	utils.SendMessage(connection, "Senddataplease", false)

	n, err := io.Copy(file, connection) //waiting reading and saving bytes from client
	if err != nil {
		return errors.New("Error copying file from Client: " + err.Error())
	}
	fmt.Println(n, "bytes sent to server")
	return nil
}
