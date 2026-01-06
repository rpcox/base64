package main

import (
	"fmt"
	"os"
)

func Usage() {
	s := `

	NAME
	   b64gen - encode/decode with Base 64

	SYNOPSIS
	   b64gen  encode | decode   [ -enum ]  -std | -url   STRING
	   b64gen  encode | decode   [ -enum ]  -std | -url   -f FILE

	DESCRIPTION
	   Convert a string to a Base 64 encoding for transmission over a medium that 
	   does not support other than simple ASCII, or decode a Base64 string.

	 decode | dec
		Set the tool mode to decode Base 64 strings
	
	 encode | enc
		Set the tool mode to encode strings to Base 64

	 -enum
		Enumerate the results by line
 	 -f string
		Specify a file with Base 64 encoded lines
 	 -std
 		Use the RFC 4648 'Standard' Base 64 alphabet
 	 -url
 		Use the RFC 4648 'URL and Filename Safe' Base 64 alphabet
	 -version
		Display tool version and exit
	   
	EXAMPLES

	Encode a string to Base 64 on the command line using the standard alphabet

	> b64gen enc -std 'ABC123'

	Decode a set of Base 64 strings from a file using the URL alphabet

	> b64gen dec -url -f mybase64strings.txt
`
	fmt.Fprintf(os.Stdout, "%s\n", s)
	os.Exit(0)
}
