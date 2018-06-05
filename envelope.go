package addrvrf

type Envelope struct {
	Input  AddressInput
	Output AddressOutput
}

type AddressInput struct {
	City string
}

type AddressOutput struct {
	City string
}
