package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

func main() {
	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]
		
		decoded, err := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func decodeBencode(bencodedString string) (interface{}, error) {
	if unicode.IsDigit(rune(bencodedString[0])) {
		
		res, _, err := decodeString(bencodedString, 0);
		if err != nil {
			return "", err;
		}

		return res, nil;
	
	} else if rune(bencodedString[0]) == 'i' && rune(bencodedString[len(bencodedString)-1]) == 'e' {

		return strconv.Atoi(bencodedString[1 : len(bencodedString)-1]);

	} else if rune(bencodedString[0]) == 'l' && rune(bencodedString[len(bencodedString)-1]) == 'e' {

		res, _, err := decodeList(bencodedString, 1);
		if err != nil {
			return "", err;
		}

		return res, nil;

	} else if rune(bencodedString[0]) == 'd' && rune(bencodedString[len(bencodedString)-1]) == 'e' {

		if bencodedString == "de" {
			return map[string]interface{}{}, nil;
		}

		res, _, err := decodeDictionary(bencodedString, 1);
		if err != nil {
			return "", err;
		}

		return res, nil;

	} else {
		return "", fmt.Errorf("only strings are supported at the moment")
	}
}

func decodeString(bencodedString string, pointer int) (interface{}, int, error) {
	var firstColonIndex int

	for i := pointer; i < len(bencodedString); i++ {
		if bencodedString[i] == ':' {
			firstColonIndex = i
			break
		}
	}

	lengthStr := bencodedString[pointer : firstColonIndex]

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", -1, err;
	}

	end := firstColonIndex + 1 + length;

	return bencodedString[firstColonIndex+1 : end], end, nil
}

func decodeInt(bencodedString string, pointer int) (int, int, error) {
	var firstEIndex int;

	for i := pointer + 1; i < len(bencodedString); i++ {
		if bencodedString[i] == 'e' {
			firstEIndex = i;
			break;
		}
	}

	res := bencodedString[pointer + 1 :firstEIndex];
	resInt, err := strconv.Atoi(res);
	if err != nil {
		return -1, -1, err;
	}
	end := firstEIndex + 1;

	return resInt, end, nil;
}

func decodeList(bencodedString string, pointer int) (interface{}, int, error) {
    if len(bencodedString) == 2 {
		return []interface{}{}, pointer, nil;
	}
	
	var resSlice []interface{}

    for pointer < len(bencodedString) - 1 {
        var end int
        var err error
        var res interface{}

        switch {
			case unicode.IsDigit(rune(bencodedString[pointer])):
				res, end, err = decodeString(bencodedString, pointer)
			case bencodedString[pointer] == 'i':
				res, end, err = decodeInt(bencodedString, pointer)
			case bencodedString[pointer] == 'l':
				res, end, err = decodeList(bencodedString, pointer+1)
			case bencodedString[pointer] == 'e':
				return resSlice, pointer + 1, nil
			default:
				return nil, -1, fmt.Errorf("invalid bencoded string")
        }

        if err != nil {
            return nil, -1, err
        }

        resSlice = append(resSlice, res)
        pointer = end
    }

    return resSlice, pointer, nil
}

func decodeDictionary(bencodedString string, pointer int) (interface{}, int, error) {

	var resMap map[string]interface{} = make(map[string]interface{});

	// if odd means store the result as key and if even then store the result as value for the last key 
	count := 1;

	var key interface{};
	var value interface{};

	for pointer < len(bencodedString) - 1 {
		var end int;
		var err error;

		switch {
			case unicode.IsDigit(rune(bencodedString[pointer])):
				if count % 2 != 0 {
					key, end, err = decodeString(bencodedString, pointer);
				} else {
					value, end, err = decodeString(bencodedString, pointer);
				}
			case bencodedString[pointer] == 'i':
				if count % 2 != 0 {
					key, end, err = decodeInt(bencodedString, pointer);
				} else {
					value, end, err = decodeInt(bencodedString, pointer);
				}
			case bencodedString[pointer] == 'l':
				if count % 2 != 0 {
					key, end, err = decodeList(bencodedString, pointer+1);
				} else {
					value, end, err = decodeList(bencodedString, pointer+1);
				}
			case bencodedString[pointer] == 'd':
				if count % 2 != 0 {
					key, end, err = decodeDictionary(bencodedString, pointer+1);
				} else {
					value, end, err = decodeDictionary(bencodedString, pointer+1);
				}
			case bencodedString[pointer] == 'e':
				return resMap, pointer + 1, nil;
			default:
				return nil, -1, fmt.Errorf("invalid bencoded string");
		}

		if err != nil {
			return nil, -1, err;
		}

		if count % 2 == 0 {
			resMap[key.(string)] = value;
		}
		count++;
		pointer = end;

	}

	return resMap, pointer, nil;
}
