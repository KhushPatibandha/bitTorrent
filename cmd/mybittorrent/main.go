package main

import (
	"encoding/json"
	"fmt"
	"os"

	decodebencode "github.com/codecrafters-io/bittorrent-starter-go/cmd/mybittorrent/decodeBencode"
	torrentfileparser "github.com/codecrafters-io/bittorrent-starter-go/cmd/mybittorrent/torrentFileParser"
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

		content, err := os.ReadFile(fileName);
		if err != nil {
			fmt.Println(err);
			return;
		}

		fmt.Println("Tracker URL:", torrentfileparser.GetTrackerURL(string(content)));
		fmt.Println("Length:", torrentfileparser.GetLength(string(content)));
		fmt.Println("Info Hash:", torrentfileparser.GetInfoHash(string(content)));
		fmt.Println("Piece Length:", torrentfileparser.GetPieceLength(string(content)));
		fmt.Println("Pieces Hash:");

		piecesHash := torrentfileparser.GetPiecesHash(string(content));
		
		for i := 0; i < len(piecesHash); i++ {
			fmt.Println(piecesHash[i]);
		}
		
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
