package main

type IFRoot struct {
	Payment  IFPayment   `hcl:"payment,block"`
	Invoices []IFInvoice `hcl:"invoice,block"`
}

type IFPayment struct {
	Bank string `hcl:"bank"`
}

type IFInvoice struct {
	ID        string  `hcl:"id,label"`
	To        string  `hcl:"to"`
	IssuedAt  string  `hcl:"issued_at"`
	SoldAt    string  `hcl:"sold_at"`
	Positions []IFPos `hcl:"pos,block"`
	Comment   string  `hcl:"comment,optional"`
}

type IFPos struct {
	Name     string  `hcl:"name"`
	Quantity int     `hcl:"quantity"`
	Unit     string  `hcl:"unit"`
	GTU      string  `hcl:"gtu,optional"`
	Vat      float64 `hcl:"vat"`
}

// https://api.ifirma.pl/wystawianie-faktury-sprzedaz%cc%87y-krajowej-towarow-i-uslug/
type IFRequest struct {
	Paid          float64        `json:"Zaplacono"`
	PaidDoc       float64        `json:"ZaplaconoNaDokumencie"`
	IssueType     string         `json:"LiczOd"`
	BankAccountNo string         `json:"NumerKontaBankowego"`
	IssuedAt      string         `json:"DataWystawienia"`
	SoldAt        string         `json:"DataSprzedazy"`
	SoldAtFormat  string         `json:"FormatDatySprzedazy"`
	PaymentMethod string         `json:"SposobZaplaty"`
	SignatureType string         `json:"RodzajPodpisuOdbiorcy"`
	ContractorID  string         `json:"IdentyfikatorKontrahenta"`
	Comment       string         `json:"Uwagi"`
	Positions     []IFRequestPos `json:"Pozycje"`
}

type IFRequestPos struct {
	Vat      float64 `json:"StawkaVat"`
	Quantity int     `json:"Ilosc"`
	Price    float64 `json:"CenaJednostkowa"`
	Name     string  `json:"NazwaPelna"`
	Unit     string  `json:"Jednostka"`
	VatType  string  `json:"TypStawkiVat"`
	GTU      string  `json:"GTU"`
}

type IFResponse struct {
	Response IFResponseDetails `json:"response"`
}

type IFResponseDetails struct {
	ID   uint64 `json:"Identyfikator"`
	Code int    `json:"Kod"`
	Msg  string `json:"Informacja"`
}
