package main

import (
	"encoding/json"
	"fmt"
	"os"

	decodebencode "github.com/codecrafters-io/bittorrent-starter-go/cmd/mybittorrent/decodeBencode"
)

func main() {
	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]
		
		decoded, err := decodebencode.DecodeBencode(bencodedValue);
		if err != nil {
			fmt.Println(err)
			return
		}
		
		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else if command == "info" {
		fileName := os.Args[2];

		contents, err := os.ReadFile(fileName);
		if err != nil {
			fmt.Println(err);
			return;
		}

		decoded, err := decodebencode.DecodeBencode(string(contents));
		if err != nil {
			fmt.Println(err);
			return;
		}

		fmt.Println("Tracker URL:", decoded.(map[string]interface{})["announce"]);
		fmt.Println("Length:", decoded.(map[string]interface{})["info"].(map[string]interface{})["length"]);

	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
