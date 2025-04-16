package main

import (
	"WebBTClient/internal/bedecoder"
)

func main() {
	decoder := bedecoder.NewDecoder("paw")
	decoder.Decode()
}
