package main

import (
	"WebBTClient/internal/bedecoder"
	"WebBTClient/internal/trackers"
	"fmt"
)

func main() {
	decoder := bedecoder.NewDecoder("announce")
	torrent_file := decoder.Decode()

	fmt.Println(torrent_file.(map[string]any)["peers"])

	tracker := trackers.NewTracker(torrent_file.(map[string]any))
	tracker.GetPeers()
}
