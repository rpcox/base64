package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"
)

var offset int64

func NextByte(r *bufio.Reader, em *EncodingMap) func(yield func(byte) bool) {
	return func(yield func(byte) bool) {
		for {
			b, err := r.ReadByte()
			offset++

			if err != nil {
				if err == io.EOF {
					// end of file
					break
				}

				fmt.Printf("error reading byte at offset %d: %v\n", offset, err)
				break
			}

			if !em.InAlphabet(b) {
				// Get another byte
				continue
			}

			if !yield(b) {
				return
			}
		}
	}
}

type Encoding interface {
	InAlphabet(b byte) bool
}

type EncodingMap struct {
	abc map[byte]rune
}

func NewEncodingMap(s string) *EncodingMap {
	var em EncodingMap
	em.abc = make(map[byte]rune)
	for _, v := range s {
		em.abc[byte(v)] = v
	}

	return &em
}

func (em *EncodingMap) InAlphabet(b byte) bool {
	_, ok := em.abc[b]
	return ok
}

//stdAlphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
//urlAlphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

func main() {
	_file := flag.String("f", "", "identify the file")
	_bufferLen := flag.Int("buffer-len", 20, "specify buffer length")
	flag.Parse()

	fh, err := os.Open(*_file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	defer fh.Close()
	r := bufio.NewReader(fh)

	em := NewEncodingMap("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	//                    666667777777777888888888899991
	//                    567890123456789012345678907890
	//em := NewEncodingMap("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")
	re := regexp.MustCompile(`^\w{3}\d{9}\w{3}`)
	dst := make([]byte, base64.StdEncoding.DecodedLen(*_bufferLen))

	var b []byte
	start := time.Now()

	for x := range NextByte(r, em) {
		b = append(b, x)
		if len(b) != *_bufferLen {
			continue
		}
		//fmt.Println(string(b))

		_, err := base64.StdEncoding.Decode(dst, b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v : offset %d: char = %x buf = %v\n", err, offset-int64(*_bufferLen), b[0], b)
			continue
		}

		if re.Match(dst) {
			fmt.Fprintf(os.Stdout, "   offset: %-d %s %s \n", offset-int64(*_bufferLen), string(b), string(dst))
		}

		b = b[1:] // drop first byte in slice
	}

	//fmt.Println(offset)
	fmt.Fprintf(os.Stdout, "  elapsed: %v\n", time.Since(start))
}
