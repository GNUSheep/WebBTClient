package bedecoder

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"unicode"
)

type Decoder struct {
	reader *bufio.Reader
}

func NewDecoder(filepath string) *Decoder {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Failed to open a torrent file: ", err)
	}

	r := bufio.NewReader(file)

	d := &Decoder{
		reader: r,
	}

	return d
}

func (d *Decoder) decodeInt() (string, error) {
	var buffer []byte

	for {
		b, err := d.reader.ReadByte()
		if b == 'e' {
			break
		}

		if err == io.EOF {
			return "", errors.New("Missing 'e' terminator after integer")
		}
		if !unicode.IsDigit(rune(b)) && b != '-' {
			return "", errors.New("Integer contains not-digit characters")
		}
		if err != nil {
			return "", err
		}

		buffer = append(buffer, b)
	}

	return string(buffer), nil
}

func (d *Decoder) decodeByteString() (string, error) {
	var string_length_buffer []byte
	for {
		b, err := d.reader.ReadByte()
		if b == ':' {
			break
		}

		if err == io.EOF {
			return "", errors.New("Unexcepted EOF, excepted to have a length of string")
		}
		if !unicode.IsDigit(rune(b)) {
			return "", errors.New("Length contains not-digit characters or is negative")
		}
		if err != nil {
			return "", err
		}

		string_length_buffer = append(string_length_buffer, b)
	}
	string_length, err := strconv.ParseInt(string(string_length_buffer), 10, 64)
	if err != nil {
		return "", err
	}

	var string_buffer []byte
	for i := 0; i < int(string_length); i++ {
		b, err := d.reader.ReadByte()

		if err == io.EOF {
			return "", errors.New("Unexcepted EOF, excepted to have a byte of string")
		}
		if err != nil {
			return "", err
		}

		string_buffer = append(string_buffer, b)
	}

	return string(string_buffer), nil
}

func (d *Decoder) decodeList() ([]any, error) {
	var buffer []any

	for {
		b, err := d.reader.ReadByte()
		if b == 'e' {
			break
		}

		if err == io.EOF {
			return nil, errors.New("Missing 'e' terminator after list")
		}

		element, err := d.handleByte(b)
		if err != nil {
			return nil, err
		}
		buffer = append(buffer, element)
	}

	return buffer, nil
}

func (d *Decoder) decodeDict() (map[string]any, error) {
	buffer := make(map[string]any)

	for {
		b, err := d.reader.ReadByte()
		if b == 'e' {
			break
		}
		if err == io.EOF {
			return nil, errors.New("Missing 'e' terminator after dict")
		}
		if err != nil {
			return nil, err
		}

		if err := d.reader.UnreadByte(); err != nil {
			return nil, err
		}
		key, err := d.decodeByteString()
		if err != nil {
			return nil, err
		}

		b, err = d.reader.ReadByte()
		if err == io.EOF || b == 'e' {
			return nil, errors.New("Missing value after key in dict")
		}
		if err != nil {
			return nil, err
		}

		value, err := d.handleByte(b)
		buffer[key] = value
	}

	return buffer, nil
}

func (d *Decoder) handleByte(b byte) (any, error) {
	switch b {
	case 'i':
		i, err := d.decodeInt()
		if err != nil {
			return nil, err
		}

		return i, nil
	case 'l':
		l, err := d.decodeList()
		if err != nil {
			return nil, err
		}

		return l, nil
	case 'd':
		d, err := d.decodeDict()
		if err != nil {
			return nil, err
		}

		return d, nil
	default:
		if err := d.reader.UnreadByte(); err != nil {
			return nil, err
		}
		s, err := d.decodeByteString()
		if err != nil {
			log.Fatal(err)
		}

		return s, nil
	}
}

func (d *Decoder) Decode() {
	b, err := d.reader.ReadByte()
	if err != nil {
		log.Fatal("Error while reading a torrent file: ", err)
	}

	output, err := d.handleByte(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", output.(map[string]any))

}
