package decodebencode

import (
	"fmt"
	"unicode"
)

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