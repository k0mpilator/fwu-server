package main

import (
	"fmt"
	"fwu-server/internal/config"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

const BUFFERSIZE = 1024

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

func sendFileToClient(connection net.Conn) {

	fmt.Println("[  OK  ] A client has connected!")

	defer connection.Close()
	file, err := os.Open("example.bin")
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	fmt.Println("[  OK  ] Sending filename and filesize!")
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("[  OK  ] Start sending file!")

	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Println("File has been sent, closing connection!")
	return
}

func main() {

	// Read config yaml
	conf := config.NewConfig("config.yml")

	fmt.Println("[  OK  ] Starting server ...")

	//listen interfaces
	server, err := net.Listen(conf.NetworkType, conf.NetworkPort)
	if err != nil {
		log.Error().Err(err).Msg("")
		os.Exit(1)
	}
	defer server.Close()

	for {
		//accept connection on port
		connection, err := server.Accept()
		if err != nil {
			log.Error().Err(err).Msg("")
			os.Exit(1)
		}

		fmt.Println("[  OK  ] Client connected")
		go sendFileToClient(connection)
	}
}
