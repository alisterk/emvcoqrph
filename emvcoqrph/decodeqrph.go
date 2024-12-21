package emvcoqrph

import (
	"errors"
	"fmt"
	//"strconv"
)

func decodeRootQRInfo(root UnstructuredEMVQRData, constants map[string]string, pointOfInitiationMap map[string]string) (QRPHData, error) {
	pointOfInitiationMethod := pointOfInitiationMap[fmt.Sprint(root[constants["IDPointOfInitiationMethod"]])]
	return QRPHData{
		PointOfInitiationMethod: pointOfInitiationMethod,
		MerchantCategoryCode:    fmt.Sprint(root[constants["IDMerchantCategoryCode"]]),
		TransactionCurrency:     fmt.Sprint(root[constants["IDTransactionCurrency"]]),
		TransactionAmount:       fmt.Sprint(root[constants["IDTransactionAmount"]]),
		CountryCode:             fmt.Sprint(root[constants["IDCountryCode"]]),
		MerchantName:            fmt.Sprint(root[constants["IDMerchantName"]]),
		MerchantCity:            fmt.Sprint(root[constants["IDMerchantCity"]]),
		PostalCode:              fmt.Sprint(root[constants["IDPostalCode"]]),
	}, nil
}

func decodeRecipient(data UnstructuredEMVQRData) (QRPHRecipientInfo, error) {
	typeStr, ok := data["00"].(string)
	if !ok {
		return QRPHRecipientInfo{}, errors.New("missing or invalid recipient type")
	}
	bankCode, ok := data["01"].(string)
	if !ok {
		return QRPHRecipientInfo{}, errors.New("missing or invalid bank code")
	}
	accountNumber, ok := data["03"].(string)
	if !ok {
		return QRPHRecipientInfo{}, errors.New("missing or invalid account number")
	}

	return QRPHRecipientInfo{
		Type:          typeStr,
		BankCode:      bankCode,
		AccountNumber: accountNumber,
	}, nil
}

func decodeQRPH(data UnstructuredEMVQRData, constants map[string]string, pointOfInitiationMap map[string]string) (QRPHData, error) {
	root, err := decodeRootQRInfo(data, constants, pointOfInitiationMap)
	if err != nil {
		return QRPHData{}, err
	}

	recipientData, ok := data["27"].(UnstructuredEMVQRData)
	if !ok {
		recipientData, ok = data["28"].(UnstructuredEMVQRData)
		if !ok {
			return QRPHData{}, errors.New("QRPH data is missing recipient information")
		}
	}

	recipient, err := decodeRecipient(recipientData)
	if err != nil {
		return QRPHData{}, err
	}

	root.Recipient = recipient
	root.Raw = data
	return root, nil
}

func decodeQRPHFromText(text string, constants map[string]string, pointOfInitiationMap map[string]string) (QRPHData, error) {
	emvQRData, err := decodeEMVQR(text, constants)
	if err != nil {
		return QRPHData{}, err
	}
	return decodeQRPH(emvQRData, constants, pointOfInitiationMap)
}
