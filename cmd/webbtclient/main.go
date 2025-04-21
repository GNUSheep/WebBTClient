package main

import (
	"WebBTClient/internal/bedecoder"
)

func main() {
	decoder := bedecoder.NewDecoder("file.torrent")
	torrent_file := decoder.Decode()
}
