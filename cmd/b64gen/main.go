package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
)

const (
	ErrSelectEncoder = iota
	ErrSelectOperation
	ErrFileOpen
)

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

func Encode(str []string, encoding *base64.Encoding) [][]byte {
	var bs [][]byte

	for _, s := range str {
		dst := make([]byte, encoding.EncodedLen(len(s)))
		encoding.Encode(dst, []byte(s))
		bs = append(bs, dst)
	}

	return bs
}

func Decode(str []string, encoding *base64.Encoding) [][]byte {
	var bs [][]byte

	for _, s := range str {
		dst := make([]byte, encoding.DecodedLen(len(s)))
		encoding.Decode(dst, []byte(s))
		bs = append(bs, dst)
	}

	return bs
}

func Spew(s []string, b [][]byte, enumerate bool) {
	count := 1
	for i, v := range s {
		if enumerate {
			fmt.Fprintf(os.Stdout, "%8d:\t%s\t%s\n", count, v, string(b[i]))
			count++
		} else {
			fmt.Fprintf(os.Stdout, "%s\t%s\n", v, string(b[i]))
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
			fmt.Fprintf(os.Stdout, "%8d:\t%s\t%s\n", count, text, string(dst))
			count++
		} else {
			fmt.Fprintf(os.Stdout, "%s\t%s\n", text, string(dst))
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
			fmt.Fprintf(os.Stdout, "%8d:\t%s\t%s\n", count, text, string(dst))
			count++
		} else {
			fmt.Fprintf(os.Stdout, "%s\t%s\n", text, string(dst))
		}
	}
}

func main() {
	encode := flag.NewFlagSet("enc", flag.ExitOnError)
	_estd := encode.Bool("std", false, "dafadk")
	_eurl := encode.Bool("url", false, "dafadk")
	_eenum := encode.Bool("enum", false, "dafadk")
	_efile := encode.String("f", "", "dladj")
	decode := flag.NewFlagSet("dec", flag.ExitOnError)
	_dstd := decode.Bool("std", false, "dafadk")
	_durl := decode.Bool("url", false, "dafadk")
	_denum := decode.Bool("enum", false, "dafadk")
	_dfile := decode.String("f", "", "dladj")
	flag.Usage = Usage
	flag.Parse()

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
