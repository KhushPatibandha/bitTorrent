package decodebencode

import (
	"fmt"
	"strconv"
	"unicode"
)

func DecodeBencode(bencodedString string) (interface{}, error) {
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