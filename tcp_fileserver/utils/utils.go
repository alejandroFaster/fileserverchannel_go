package utils

import (
	"errors"
	"fmt"
	"net"
	"strconv"
)

func Validate_port(strPort string) (int, error) {
	var port, err = strconv.Atoi(strPort)
	if err != nil {
		return 0, errors.New(" Error to convert port number. The port must be a number between 1 and 65535")
	}

	if port > 65535 || port < 1 {
		return 0, errors.New(" The port must be a number between 1 and 65535")
	}

	return port, nil
}

func ValidateIPAddress(ip string) error {
	if net.ParseIP(ip) == nil {
		return errors.New(ip + " is not a valid IP address")
	}
	return nil
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func SendError(connection net.Conn, message string) {
	SendMessage(connection, message, true)
}

func SendMessage(connection net.Conn, message string, error bool) {
	_, err := connection.Write([]byte(message + "\r\n"))
	if err != nil {
		fmt.Print(err.Error())
	}
	if error {
		//Close the connection if it's a error message
		connection.Close()
	}

}
