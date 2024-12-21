package emvcoqrph

import (
	"fmt"
	"sort"

	"strings"

	"github.com/howeyc/crc16"
)

func formatCrc(value string) (string, error) {
	crc := crc16.ChecksumCCITT([]byte(value))
	crcHex := fmt.Sprintf("%04X", crc) // Format CRC as a 4-digit uppercase hexadecimal string
	return crcHex, nil
}

func appendCrc(value string, ID map[string]string) (string, error) {
	toCrc := value + ID["IDCRC"] + "04"
	crc, err := formatCrc(toCrc)
	if err != nil {
		return "", err
	}
	return toCrc + crc, nil
}

func encodeUnstructuredData(data UnstructuredEMVQRData) (string, error) {
	var builder strings.Builder

	sortedKeys := make([]string, 0, len(data))
	for k := range data {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		value := data[key]
		var valueStr string
		switch v := value.(type) {
		case string:
			valueStr = v
		case UnstructuredEMVQRData:
			embedded, err := encodeUnstructuredData(v)
			if err != nil {
				return "", err
			}
			valueStr = embedded
		default:
			return "", fmt.Errorf("unsupported data type for key %s", key)
		}
		builder.WriteString(key)
		builder.WriteString(fmt.Sprintf("%02d", len(valueStr)))
		builder.WriteString(valueStr)
	}

	return builder.String(), nil
}

func encodeEMVQR(data UnstructuredEMVQRData) (string, error) {
	dataCopy := make(UnstructuredEMVQRData)
	for k, v := range data {
		if k != ID["IDCRC"] {
			dataCopy[k] = v
		}
	}

	encodedData, err := encodeUnstructuredData(dataCopy)
	if err != nil {
		return "", err
	}

	return appendCrc(encodedData, ID)
}
