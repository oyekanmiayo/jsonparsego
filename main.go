package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("files/empty.json")
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	var data []byte
	for {
		buffer := make([]byte, 1)
		n, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				os.Exit(1)
			}
			break
		}
		data = append(data, buffer[:n]...)
	}

	// Empty and invalid json
	if len(data) == 0 {
		os.Exit(1)
	}

	tokenList, err := scanTokens(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	valid, err := parseTokensV2(tokenList)
	if valid {
		os.Exit(0)
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
}
