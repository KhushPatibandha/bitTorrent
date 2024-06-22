package torrentfileparser

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"

	decodebencode "github.com/codecrafters-io/bittorrent-starter-go/cmd/mybittorrent/decodeBencode"
	bencode "github.com/jackpal/bencode-go"
)

func GetTrackerURL(content string) interface{} {
	decode := decode(content);
	return decode.(map[string]interface{})["announce"];
}

func GetLength(content string) interface{} {
	decode := decode(content);
	return decode.(map[string]interface{})["info"].(map[string]interface{})["length"]
}

func GetInfoHash(content string) string {
	decode := decode(content)
	info := decode.(map[string]interface{})["info"]
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, info)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hashRes := sha1.Sum(buf.Bytes());

	return fmt.Sprintf("%x", hashRes);
}

func decode(content string) interface{} {
	decoded, err := decodebencode.DecodeBencode(content);
	if err != nil {
		fmt.Println(err);
		os.Exit(1);
	}

	return decoded;
}