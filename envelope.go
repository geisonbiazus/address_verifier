package addrvrf

type Envelope struct {
	Input  AddressInput
	Output AddressOutput
}

type AddressInput struct {
	Street  string
	City    string
	State   string
	ZIPCode string
}

type AddressOutput struct {
	DeliveryLine1 string
	LastLine      string
	Street        string
	City          string
	State         string
	ZIPCode       string
}
