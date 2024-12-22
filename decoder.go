package emvcoqrph

import (
	"errors"
	"strconv"
)

func decodeInner(text string) (map[string]string, error) {
	result := make(map[string]string)
	remaining := text
	for len(remaining) > 0 {
		if len(remaining) < 4 {
			return nil, errors.New("malformed: still " + remaining)
		}
		rawKey := remaining[:2]
		rawLength := remaining[2:4]
		length, err := strconv.Atoi(rawLength)
		if err != nil {
			return nil, errors.New("malformed length for field " + rawKey + ": " + rawLength)
		}
		if len(remaining) < 4+length {
			return nil, errors.New("unexpected end of file for field " + rawKey + ": " + strconv.Itoa(length) + " " + text)
		}
		value := remaining[4 : 4+length]
		remaining = remaining[4+length:]
		result[rawKey] = value
	}
	return result, nil
}

func decodeEMVQR(text string, ID map[string]string) (UnstructuredEMVQRData, error) {
	result, err := decodeInner(text)
	if err != nil {
		return nil, err
	}

	decodedResult := make(UnstructuredEMVQRData)
	for key, value := range result {
		if key == ID["IDAdditionalDataFieldTemplate"] ||
			key == ID["IDMerchantInformationLanguageTemplate"] ||
			(key >= ID["IDMerchantAccountInformationRangeStart"] && key <= ID["IDMerchantAccountInformationRangeEnd"]) ||
			(key >= ID["IDRFUForEMVCoRangeStart"] && key <= ID["IDRFUForEMVCoRangeEnd"]) ||
			(key >= ID["IDUnreservedTemplatesRangeStart"] && key <= ID["IDUnreservedTemplatesRangeEnd"]) {
			innerDecoded, err := decodeInner(value)
			if err != nil {
				return nil, err
			}
			decodedResult[key] = innerDecoded
		} else {
			decodedResult[key] = value
		}
	}
	return decodedResult, nil
}
