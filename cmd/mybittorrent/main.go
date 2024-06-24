package main

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"

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

		peers := getPeers(fileName);
		if peers == nil {
			return;
		}
		for i := 0; i < len(peers); i++ {
			fmt.Println(peers[i]);
		}

	} else if command == "handshake" {
		fileName := os.Args[2];
		hostWithPort := os.Args[3];

		peerID, conn := doHandshake(fileName, hostWithPort);
		defer conn.Close();

		fmt.Println("Peer ID:", peerID);

	} else if command == "download_piece" {
		outputPath := os.Args[3];
		fileName := os.Args[4];
		pieceNumber, _ := strconv.Atoi(os.Args[5]);

		content, _ := os.ReadFile(fileName);
		totalLength := torrentfileparser.GetLength(string(content)).(int);
		pieceLength := torrentfileparser.GetPieceLength(string(content)).(int);
		piecesHash := torrentfileparser.GetPiecesHash(string(content));

		pieceSize := pieceLength
		if pieceNumber == len(piecesHash)-1 {
			remaining := totalLength % pieceLength
			if remaining > 0 {
				pieceSize = remaining
			}
		}

		peers := getPeers(fileName)

		_, conn := doHandshake(fileName, peers[1])
		defer conn.Close()

		bitfieldBuffer := make([]byte, 2048)
		_, err := conn.Read(bitfieldBuffer)
		if err != nil {
			fmt.Println("Error reading bitfield:", err)
			return
		}
		bitfieldMessageId := bitfieldBuffer[4]
		if bitfieldMessageId != 5 {
			fmt.Println("Expected bitfield message ID 5, but got:", bitfieldMessageId)
			return
		}

		interestedMessageId := byte(2)
		_, err = conn.Write([]byte{0, 0, 0, 1, interestedMessageId})
		if err != nil {
			fmt.Println("Error sending interested message:", err)
			return
		}

		unChokeBuffer := make([]byte, 5)
		_, err = conn.Read(unChokeBuffer)
		if err != nil {
			fmt.Println("Error reading unchoke message:", err)
			return
		}
		unchokeMessageId := unChokeBuffer[4]
		if unchokeMessageId != 1 {
			fmt.Println("Expected unchoke message ID 1, but got:", unchokeMessageId)
			return
		}

		blockSize := 16 * 1024
		totalBlocks := (pieceSize + blockSize - 1) / blockSize

		var pieceData []byte
		for i := 0; i < totalBlocks; i++ {
			requestMessageId := byte(6)
			begin := i * blockSize
			length := blockSize
			if i == totalBlocks-1 && pieceSize%blockSize != 0 {
				length = pieceSize % blockSize
			}

			payload := make([]byte, 12)
			binary.BigEndian.PutUint32(payload[0:4], uint32(pieceNumber))
			binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
			binary.BigEndian.PutUint32(payload[8:12], uint32(length))

			message := append([]byte{0, 0, 0, 13, requestMessageId}, payload...)

			_, err = conn.Write(message)
			if err != nil {
				fmt.Println("Error sending request message:", err)
				return
			}

			responseBuffer := make([]byte, 4+1+8+length)
			_, err = io.ReadFull(conn, responseBuffer)
			if err != nil {
				fmt.Println("Error reading response message:", err)
				return
			}

			responseMessageId := responseBuffer[4]
			if responseMessageId != 7 {
				fmt.Println("Expected piece message ID 7, but got:", responseMessageId)
				return
			}

			pieceData = append(pieceData, responseBuffer[13:]...)
		}

		err = os.WriteFile(outputPath, pieceData, 0644)
		if err != nil {
			fmt.Println("Error writing piece data to file:", err)
			return
		}

		fmt.Println("Piece", pieceNumber, "downloaded to", outputPath)
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}

func doHandshake(fileName string, hostWithPort string) (string, net.Conn) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	conn, err := net.Dial("tcp", hostWithPort)
	if err != nil {
		fmt.Println("Error connecting to peer:", err)
		return "", nil
	}

	infoHash := torrentfileparser.GetInfoHash(string(content))

	protocol := "BitTorrent protocol"
	protocolLength := byte(len(protocol))
	reservedBytes := make([]byte, 8)
	infoHashBytes, _ := hex.DecodeString(infoHash)
	peerId := []byte("00112233445566778899")

	handshake := append([]byte{protocolLength}, protocol...)
	handshake = append(handshake, reservedBytes...)
	handshake = append(handshake, infoHashBytes...)
	handshake = append(handshake, peerId...)

	_, err = conn.Write(handshake)
	if err != nil {
		fmt.Println("Error sending handshake:", err)
		return "", nil
	}

	responseBuffer := make([]byte, 68)
	_, err = conn.Read(responseBuffer)
	if err != nil {
		fmt.Println("Error reading handshake response:", err)
		return "", nil
	}

	receivedPeerID := responseBuffer[48:]

	return hex.EncodeToString(receivedPeerID), conn
}

func getPeers(fileName string) []string {
	content, err := os.ReadFile(fileName);
	if err != nil {
		fmt.Println(err);
		return nil;
	}

	trackerUrl := torrentfileparser.GetTrackerURL(string(content));

	parsedUrl, err := url.Parse(trackerUrl.(string));
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return nil;
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
		return nil;
	}

	defer response.Body.Close();

	body, err := io.ReadAll(response.Body);
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil;
	}

	decodedBody, err := decodebencode.DecodeBencode(string(body));
	if err != nil {
		fmt.Println("Error decoding response body:", err)
		return nil;
	}

	peersString := decodedBody.(map[string]interface{})["peers"].(string);

	var peers []string;

	for i := 0; i < len(peersString); i += 6 {
		peer := peersString[i : i+6]
	
		ipBytes := peer[:4]
		portBytes := peer[4:]
	
		ip := fmt.Sprintf("%d.%d.%d.%d", ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
	
		port := int(portBytes[0])<<8 + int(portBytes[1])
	
		peers = append(peers, fmt.Sprintf("%s:%d", ip, port));
	}

	return peers;
}
