package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/jinzhu/now"
)

const apiUrl string = "https://www.ifirma.pl/iapi/fakturakraj.json"

func main() {

	flagCfg := flag.String("config", fmt.Sprint(os.Getenv("HOME"), "/.config/ifirma.hcl"), "Path to config file")
	flagInvoice := flag.String("invoice", "", "Invoice ID")
	flagPrice := flag.Float64("net_price", 0, "Net Price")

	flag.Parse()

	iToken, found := os.LookupEnv("IFIRMA_FV_TOKEN")
	if !found {
		log.Fatalf("Environment variable IFIRMA_FV_TOKEN not defined")
		return
	}
	iEmail, found := os.LookupEnv("IFIRMA_EMAIL")
	if !found {
		log.Fatalf("Environment variable IFIRMA_EMAIL not defined")
		return
	}

	fmt.Println("Using config:", *flagCfg)
	fmt.Println("Selected invoice:", *flagInvoice)
	fmt.Println("Price:", *flagPrice)

	var root IFRoot
	err := hclsimple.DecodeFile(*flagCfg, nil, &root)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
		return
	}

	var invoice *IFInvoice
	for _, inv := range root.Invoices {
		if *flagInvoice == inv.ID {
			invoice = &inv
			break
		}
	}

	if invoice == nil {
		log.Fatalf("Invoice not found in config: %s", *flagInvoice)
		return
	}

	req := IFRequest{
		Paid:          0,
		PaidDoc:       0,
		IssueType:     "NET",
		BankAccountNo: strings.ReplaceAll(root.Payment.Bank, " ", ""),
		IssuedAt:      extractDate(invoice.IssuedAt),
		SoldAt:        extractDate(invoice.SoldAt),
		SoldAtFormat:  "DZN",
		PaymentMethod: "PRZ",
		SignatureType: "BPO",
		ContractorID:  invoice.To,
		Positions: []IFRequestPos{
			IFRequestPos{
				Vat:      0.23,
				Quantity: invoice.Positions[0].Quantity,
				Price:    *flagPrice,
				Name:     invoice.Positions[0].Name,
				Unit:     invoice.Positions[0].Unit,
				VatType:  "PRC",
			},
		},
	}

	out, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("Json marshalling failed: %v", err)
		return
	}
	// https://api.ifirma.pl/naglowek-autoryzacji/
	hashToken, err := hex.DecodeString(iToken)
	if err != nil {
		log.Fatalf("Cannot decode token %v", hashToken)
	}
	hash := hmac.New(sha1.New, hashToken)
	io.WriteString(hash, apiUrl)
	io.WriteString(hash, iEmail)
	io.WriteString(hash, "faktura")
	hash.Write(out)
	sum := hex.EncodeToString(hash.Sum(nil))
	auth := "IAPIS user=" + iEmail + ", hmac-sha1=" + sum

	client := resty.New()
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json; charset=UTF-8").
		SetHeader("Authentication", auth).
		SetBody(out).
		Post(apiUrl)

	fmt.Println(resp)
}

func extractDate(d string) string {
	if d == "" {
		lastMonth := now.BeginningOfMonth().AddDate(0, 0, -1)
		lastMonth = now.With(lastMonth).EndOfMonth()
		return lastMonth.Format("2006-01-02")
	}

	return d
}
