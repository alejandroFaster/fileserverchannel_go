package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"tcp_fileserver/utils"
)

const BUFFER_SIZE = 256

var validJSONFilename = regexp.MustCompile(`{"(filename)["]:[ |\t]*["]([a-zA-Z0-9\s_\\.\-\(\):]+)"}`)

func SendFile(fullpath string, ip string, port int, channel int) error {
	conn, err := connect(ip, port)
	if err != nil {
		return err
	}
	filename := filepath.Base(fullpath)

	//send message using custom protocol
	///Ex. send 10001 "filename.ext"
	utils.SendMessage(conn, fmt.Sprintf(`send %d "%s"`, channel, filename), false)

	//wait for message error or OK
	buffer := make([]byte, BUFFER_SIZE)
	_, error := conn.Read(buffer)
	if error != nil {
		conn.Close()
		return utils.NewError(15, "There is an error reading from connection: "+error.Error())
	}
	strBuffer := strings.TrimSpace(string(bytes.Trim(buffer, "\x00")))

	if strBuffer != "Senddataplease" {
		conn.Close()
		return utils.NewError(15, "\nThere is an protocol error custom, the error message is: "+strBuffer)
	}

	//copy the file from client to server if doesnt errors received by the custom protocol
	file, err2 := os.Open(strings.TrimSpace(fullpath)) // For read access.
	if err2 != nil {
		conn.Close()
		return utils.NewError(10, "Error reading file from Client: "+err2.Error())
	}
	defer file.Close()

	n, err2 := io.Copy(conn, file)
	if err2 != nil {
		conn.Close()
		return utils.NewError(11, "Error copying file from Client: "+err2.Error())
	}
	fmt.Printf("Bytes send to the server: %v\n", n)

	conn.Close()
	return nil
}

func DownloadFile(ip string, port int, channel int) error {
	//connect to the server
	conn, err := connect(ip, port)
	if err != nil {
		return err
	}
	//send subscribe request using custom protocol
	///Ex. subscribe 24
	utils.SendMessage(conn, fmt.Sprintf(`subscribe %d`, channel), false)
	//wait for message error or OK
	//get information about the filename by the channel number
	buffer := make([]byte, BUFFER_SIZE)
	_, error := conn.Read(buffer)
	if error != nil {
		conn.Close()
		return utils.NewError(15, "There is an error reading from connection: "+error.Error())
	}
	strBuffer := strings.TrimSpace(string(bytes.Trim(buffer, "\x00")))

	groups := validJSONFilename.FindStringSubmatch(strBuffer)
	if groups == nil {
		conn.Close()
		return utils.NewError(20, "\""+strBuffer+"\"\nBad format of filename JSON Answer from server. Format Usage: \n"+
			"{\"filename\": \"<filename.ext>\"}'\n")
	}

	filename := groups[2]
	//download from server to client
	///saving all the bytes received from server to the filename received by the channel
	file, errDownload := os.Create(strings.TrimSpace(filename)) // For read access.
	if errDownload != nil {
		return utils.NewError(16, "Error creating file from Server to Client: "+err.Error())
	}
	defer file.Close()

	n, errDownload := io.Copy(file, conn) //waiting reading and saving bytes from client
	if errDownload != nil {
		return errors.New("Error downloading file from Server to client: " + err.Error())
	}
	fmt.Println(n, "bytes received from server")

	conn.Close()
	return nil
}

func connect(ip string, port int) (net.Conn, *utils.ErrorF) {
	connection, err := net.Dial("tcp", ip+":"+fmt.Sprint(port))
	if err != nil {
		return nil, utils.NewError(5, "There was an error making a connection, the error is: "+err.Error())
	}
	return connection, nil
}
