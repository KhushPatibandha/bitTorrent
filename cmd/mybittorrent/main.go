package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
		
	} else if command == "peers" {
		fileName := os.Args[2];

		content, err := os.ReadFile(fileName);
		if err != nil {
			fmt.Println(err);
			return;
		}

		trackerUrl := torrentfileparser.GetTrackerURL(string(content));

		parsedUrl, err := url.Parse(trackerUrl.(string));
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return
		}

		params := url.Values{};

		infoHash := torrentfileparser.GetInfoHash(string(content));
		infoHashBytes, _ := hex.DecodeString(infoHash);

		params.Add("info_hash", string(infoHashBytes));
		params.Add("peer_id", "00112233445566778899");
		params.Add("port", "6881");
		params.Add("uploaded", "0");
		params.Add("downloaded", "0");
		
		leftLength := torrentfileparser.GetLength(string(content));
		
		params.Add("left", fmt.Sprintf("%v", leftLength));
		params.Add("compact", "1");

		parsedUrl.RawQuery = params.Encode();

		response, err := http.Get(parsedUrl.String());
		if err != nil {
			fmt.Println("Error making request to tracker:", err)
			return
		}

		defer response.Body.Close();

		body, err := io.ReadAll(response.Body);
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		decodedBody, err := decodebencode.DecodeBencode(string(body));
		if err != nil {
			fmt.Println("Error decoding response body:", err)
			return
		}

		peersString := decodedBody.(map[string]interface{})["peers"].(string);

		for i := 0; i < len(peersString); i += 6 {
			peer := peersString[i : i+6]
		
			ipBytes := peer[:4]
			portBytes := peer[4:]
		
			ip := fmt.Sprintf("%d.%d.%d.%d", ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
		
			port := int(portBytes[0])<<8 + int(portBytes[1])
		
			fmt.Printf("%s:%d\n", ip, port)
		}

	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
