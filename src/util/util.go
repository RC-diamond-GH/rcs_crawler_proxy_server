package util

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func HexDump(data []byte) {
	length := len(data)
	for i := 0; i < length; i += 16 {
		chunk := data[i:min(i+16, length)]
		var hexPart []string
		for _, b := range chunk {
			hexPart = append(hexPart, fmt.Sprintf("%02x", b))
		}
		hexLine := strings.Join(hexPart, " ")
		var asciiPart []rune
		for _, b := range chunk {
			asciiPart = append(asciiPart, toPrintable(b))
		}
		asciiLine := string(asciiPart)
		hexLine = fmt.Sprintf("%-47s", hexLine)
		fmt.Printf("%s  %s\n", hexLine, asciiLine)
	}
}

func toPrintable(b byte) rune {
	if b >= 0x20 && b <= 0x7E {
		return rune(b)
	}
	return '.'
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func HTTPDump(reader *http.Request) {
	method := reader.Method
	reqURL := reader.URL
	reqHeader := reader.Header

	data, err := io.ReadAll(reader.Body)
	if err != nil {
		logger.Error("Error while HTTP Dump getting Data: " + err.Error())
	}

	fmt.Printf("%s, %s, %s, %s\n", reqURL.Host, reqURL.Opaque, reqURL.Path, reqURL.Query().Encode())
	fmt.Printf("%s/%s?%s", reqURL.Host, reqURL.Path, reqURL.Query().Encode())
	//test
	fmt.Printf("\n%s %s HTTP/1.1\n", method, reqURL)
	for k, v := range reqHeader {
		print(k)
		print(": ")
		for _, content := range v {
			print(content)
			print(", ")
		}
		println()
	}
	if len(data) > 0 {
		println("\n")
		HexDump(data)
	}
}
