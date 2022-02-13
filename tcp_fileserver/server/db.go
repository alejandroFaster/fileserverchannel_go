package server

import (
	"encoding/gob"
	"os"
	"tcp_fileserver/utils"
)

type ChannelFile struct {
	Channel    int
	Filename   string
	FullPath   string
	Downloaded int64
}

func SaveChannels(mapChannels map[int]ChannelFile, filename string) error {
	encodeFile, err := os.Create(filename)
	if err != nil {
		return utils.NewError(1, "Error trying to create '"+filename+"': "+err.Error())
	}

	encoder := gob.NewEncoder(encodeFile)

	if err := encoder.Encode(mapChannels); err != nil {
		return utils.NewError(2, "Error trying to encode mapChannels to gob format: "+err.Error())
	}
	encodeFile.Close()
	return nil
}

func LoadChannels(filename string) (map[int]ChannelFile, error) {
	decodeFile, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, utils.NewError(3, "File doesn't exists")
		}
		return nil, utils.NewError(4, "Error trying to open the filename: '"+filename+"'. The error is: "+err.Error())
	}
	decoder := gob.NewDecoder(decodeFile)
	channels := make(map[int]ChannelFile)
	decoder.Decode(&channels)

	return channels, nil
}
