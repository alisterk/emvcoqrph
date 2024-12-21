package emvcoqrph

import (
// "errors"
// "fmt"
// "strconv"
)

const (
	StaticPointOfInitiationMethod  = "11"
	DynamicPointOfInitiationMethod = "12"
)

func encodeQRPH(data QRPHData) (UnstructuredEMVQRData, error) {
	root := make(UnstructuredEMVQRData)
	pointOfInitiationMethod := StaticPointOfInitiationMethod
	if data.PointOfInitiationMethod == "dynamic" {
		pointOfInitiationMethod = DynamicPointOfInitiationMethod
	}

	root["IDPointOfInitiationMethod"] = pointOfInitiationMethod
	root["IDMerchantCategoryCode"] = data.MerchantCategoryCode
	root["IDTransactionCurrency"] = data.TransactionCurrency
	root["IDCountryCode"] = data.CountryCode
	root["IDMerchantName"] = data.MerchantName
	root["IDMerchantCity"] = data.MerchantCity

	if data.TransactionAmount != "" {
		root["IDTransactionAmount"] = data.TransactionAmount
	}

	if data.PostalCode != "" {
		root["IDPostalCode"] = data.PostalCode
	}

	recipients := map[string]string{
		"00": data.Recipient.Type,
		"01": data.Recipient.BankCode,
		"03": data.Recipient.AccountNumber,
	}
	recipientKey := "28"
	if data.Recipient.Type == "com.p2pqrpay" {
		recipientKey = "27"
	}

	return UnstructuredEMVQRData{
		"_raw":       data.Raw,
		"root":       root,
		recipientKey: recipients,
	}, nil
}

func encodeQRPHToText(data QRPHData) (string, error) {
	emvQRData, err := encodeQRPH(data)
	if err != nil {
		return "", err
	}

	// Assuming encodeEMVQR function exists for this purpose
	return encodeEMVQR(emvQRData)
}
