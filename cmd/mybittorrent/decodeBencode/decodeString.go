package decodebencode

import "strconv"

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