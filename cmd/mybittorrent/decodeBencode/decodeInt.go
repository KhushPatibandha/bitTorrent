package decodebencode

import "strconv"

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