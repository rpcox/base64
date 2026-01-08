// A simple tool to generate base64 encoded strings using standard and url alphabets
package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
)

const (
	tool             = "b64gen"
	version          = "0.1.0"
	ErrSelectEncoder = iota
	ErrSelectOperation
	ErrFileOpen
)

// Setting the encoding to use is required. The Standard or URL alphabet must be used.
func SetEncoding(std, url bool) *base64.Encoding {
	var encoding *base64.Encoding

	if std {
		encoding = base64.StdEncoding
	} else if url {
		encoding = base64.URLEncoding
	} else {
		fmt.Println("select an encoder: -std or -url required")
		os.Exit(ErrSelectEncoder)
	}

	return encoding
}

// Encode the string with the specified encoding
func Encode(str []string, encoding *base64.Encoding) [][]byte {
	var bs [][]byte

	for _, s := range str {
		dst := make([]byte, encoding.EncodedLen(len(s)))
		encoding.Encode(dst, []byte(s))
		bs = append(bs, dst)
	}

	return bs
}

// Decode the string with the specified encoding
func Decode(str []string, encoding *base64.Encoding) [][]byte {
	var bs [][]byte

	for _, s := range str {
		dst := make([]byte, encoding.DecodedLen(len(s)))
		encoding.Decode(dst, []byte(s))
		bs = append(bs, dst)
	}

	return bs
}

// Display result to console
func Spew(s []string, b [][]byte, enumerate bool) {
	count := 1
	for i, v := range s {
		if enumerate {
			fmt.Fprintf(os.Stdout, "%8d: %-32s  %s\n", count, v, string(b[i]))
			count++
		} else {
			fmt.Fprintf(os.Stdout, "%-32s  %s\n", v, string(b[i]))
		}
	}
}

func DecodeFromFile(fileName string, encoding *base64.Encoding, enumerate bool) {
	fh, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(ErrFileOpen)
	}
	defer fh.Close()

	count := 1
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		text := scanner.Text()
		dst := make([]byte, encoding.DecodedLen(len(text)))
		encoding.Decode(dst, []byte(text))
		if enumerate {
			fmt.Fprintf(os.Stdout, "%8d: %-32s %s\n", count, text, string(dst))
			count++
		} else {
			fmt.Fprintf(os.Stdout, "%-32s  %s\n", text, string(dst))
		}
	}
}

func EncodeFromFile(fileName string, encoding *base64.Encoding, enumerate bool) {
	fh, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(ErrFileOpen)
	}
	defer fh.Close()

	count := 1
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		text := scanner.Text()
		dst := make([]byte, encoding.EncodedLen(len(text)))
		encoding.Encode(dst, []byte(text))
		if enumerate {
			fmt.Fprintf(os.Stdout, "%8d: %-32s %s\n", count, text, string(dst))
			count++
		} else {
			fmt.Fprintf(os.Stdout, "%-32s  %s\n", text, string(dst))
		}
	}
}

var (
	encode = flag.NewFlagSet("enc", flag.ExitOnError)
	decode = flag.NewFlagSet("dec", flag.ExitOnError)
)

func Version(b bool) {
	if b {
		fmt.Fprintf(os.Stdout, "%s v%s\n", tool, version)
		os.Exit(0)
	}
}

func main() {
	// encode flags
	_estd := encode.Bool("std", false, "Use the RFC 4648 'Standard' Base 64 alphabet")
	_eurl := encode.Bool("url", false, "Use the RFC 4648 'URL and Filename Safe' Base 64 alphabet")
	_eenum := encode.Bool("enum", false, "Enumerate the results by line")
	_efile := encode.String("f", "", "Specify a file with Base 64 encoded lines")
	// decode flags
	_dstd := decode.Bool("std", false, "Use the RFC 4648 'Standard' Base 64 alphabet")
	_durl := decode.Bool("url", false, "Use the RFC 4648 'URL and Filename Safe' Base 64 alphabet")
	_denum := decode.Bool("enum", false, "Enumerate the results by line")
	_dfile := decode.String("f", "", "Specify a file with Base 64 encoded lines")

	_version := flag.Bool("version", false, "Display version and exit")
	flag.Usage = Usage
	flag.Parse()
	Version(*_version)

	var encoding *base64.Encoding

	switch os.Args[1] {
	case "enc", "encode":
		encode.Parse(os.Args[2:])
		encoding = SetEncoding(*_estd, *_eurl)
		if *_efile != "" {
			EncodeFromFile(*_efile, encoding, *_eenum)
			os.Exit(0)
		} else {
			remaining := encode.Args()
			bs := Encode(remaining, encoding)
			Spew(remaining, bs, *_eenum)
		}
	case "dec", "decode":
		decode.Parse(os.Args[2:])
		encoding = SetEncoding(*_dstd, *_durl)
		if *_dfile != "" {
			DecodeFromFile(*_dfile, encoding, *_denum)
			os.Exit(0)
		}
		remaining := decode.Args()
		bs := Decode(remaining, encoding)
		Spew(remaining, bs, *_denum)
	default:
		fmt.Println("'encode' or 'decode' required. nothing to do")
		os.Exit(ErrSelectOperation)
	}

}
