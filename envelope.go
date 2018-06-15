package addrvrf

const InitialSequence = 1

type Envelope struct {
	Sequence int
	Input    AddressInput
	Output   AddressOutput
}

type AddressInput struct {
	Street  string
	City    string
	State   string
	ZIPCode string
}

type AddressOutput struct {
	Status        string
	DeliveryLine1 string
	LastLine      string
	Street        string
	City          string
	State         string
	ZIPCode       string
}
