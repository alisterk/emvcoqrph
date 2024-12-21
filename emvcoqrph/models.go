package emvcoqrph

type (
	UnstructuredEMVQRData map[string]interface{}
)

type QRPHData struct {
	PointOfInitiationMethod string
	MerchantCategoryCode    string
	TransactionCurrency     string
	TransactionAmount       string
	CountryCode             string
	MerchantName            string
	MerchantCity            string
	PostalCode              string
	Recipient               QRPHRecipientInfo
	Raw                     map[string]interface{}
}

type QRPHRecipientInfo struct {
	Type          string
	BankCode      string
	AccountNumber string
}
