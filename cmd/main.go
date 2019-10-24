package main

import (
	"fmt"
	"fwu-server/internal/config"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

const BUFFERSIZE = 1024

func firmwareName(filename string) {
	fwimg, err := filepath.Glob("core-image-minimal-*.ubifs")
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	if len(fwimg) != 1 {
		log.Error().Err(err).Msg("[  ERR ] file not found")
		return
	}

	fmt.Println("[  OK  ] found", fwimg)

	buf, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatal().Err(err).Caller().Msg("")
	}

	c := config.Conf{}

	if err = yaml.Unmarshal(buf, &c); err != nil {
		log.Fatal().Err(err).Caller().Msg("")
	}

	c.FwName = fwimg[0]

	b, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatal().Err(err).Caller().Msg("")
	}

	f, err := os.OpenFile("config.yml", os.O_RDWR, 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("could not open file")
	}
	defer f.Close()

	if _, err := f.Write(b); err != nil {
		log.Error().Err(err).Msg("could not write file")
	}
}

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

func sendFileToClient(connection net.Conn, conf config.Conf) {

	fmt.Println("[  OK  ] A client has connected!")

	defer connection.Close()
	file, err := os.Open(conf.FwName)
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
	fmt.Println("[  OK  ] File has been sent, closing connection!")
	return
}

func main() {

	// search new firmware and add to task
	firmwareName("config.yml")

	// Read config yaml
	conf := config.NewConfig("config.yml")

	fmt.Println("[  OK  ] Starting server ...")
	fmt.Println("[  OK  ] Firmware ready:", conf.FwName)

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
		go sendFileToClient(connection, conf)
	}
}
