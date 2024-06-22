package decodebencode

import (
	"fmt"
	"unicode"
)

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