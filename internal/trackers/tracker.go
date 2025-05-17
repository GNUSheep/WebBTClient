package trackers

import (
	"WebBTClient/internal/bedecoder"
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Tracker struct {
	Announce string
	InfoHash [20]byte
	Length   int64
}

func NewTracker(torrent_file map[string]any) *Tracker {
	announce := torrent_file["announce"].(string)

	var buf bytes.Buffer
	err := bedecoder.Encode(&buf, torrent_file["info"])
	if err != nil {
		log.Fatal(err)
	}
	info_hash := sha1.Sum(buf.Bytes())

	length := torrent_file["info"].(map[string]any)["length"].(int64)

	return &Tracker{
		Announce: announce,
		InfoHash: info_hash,
		Length:   length,
	}
}

func (t *Tracker) GetPeers() (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}

	base.RawQuery = url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{"-WBTC0001-f5kw9ysxz1rs"},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"left":       []string{strconv.FormatInt(t.Length, 10)},
		"port":       []string{string("6881")},
		"compact":    []string{string("1")},
	}.Encode()

	resp, err := http.Get(base.String())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
	return "", nil
}
