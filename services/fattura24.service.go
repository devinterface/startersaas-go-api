package services

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"devinterface.com/goaas-api-starter/models"
	"github.com/stripe/stripe-go/v71"
)

type Fattura24Service struct{ BaseService }

type Row struct {
	XMLName        xml.Name `xml:"Row"`
	Code           string   `xml:"Date"`
	Description    string   `xml:"Description"`
	Price          string   `xml:"Price"`
	VatCode        int      `xml:"VatCode"`
	VatDescription string   `xml:"VatDescription"`
	Qty            int      `xml:"Qty"`
}

type Payment struct {
	XMLName xml.Name `xml:"Payment"`
	Date    string   `xml:"Date"`
	Paid    bool     `xml:"Paid"`
	Amount  string   `xml:"Amount"`
}

type Document struct {
	XMLName           xml.Name  `xml:"Document"`
	TotalWithoutTax   string    `xml:"TotalWithoutTax"`
	VatAmount         string    `xml:"VatAmount"`
	DocumentType      string    `xml:"DocumentType"`
	SendEmail         bool      `xml:"SendEmail"`
	FePaymentCode     string    `xml:"FePaymentCode"`
	Object            string    `xml:"Object"`
	Total             string    `xml:"Total"`
	PaymentMethodName string    `xml:"PaymentMethodName"`
	Payments          []Payment `xml:"Payments>Payment"`
	CustomerName      string    `xml:"CustomerName"`
	CustomerAddress   string    `xml:"CustomerAddress"`
	CustomerVatCode   string    `xml:"CustomerVatCode"`
	CustomerCellPhone string    `xml:"CustomerCellPhone"`
	CustomerEmail     string    `xml:"CustomerEmail"`
	FeCustomerPec     string    `xml:"FeCustomerPec"`
	FeDestinationCode string    `xml:"FeDestinationCode"`
	FootNotes         string    `xml:"FootNotes"`
	Rows              []Row     `xml:"Rows>Row"`
}

type Fattura24 struct {
	XMLName  xml.Name `xml:"Fattura24"`
	Document Document
}

// GenerateInvoice func
func (fattura24Service *Fattura24Service) GenerateInvoice(accountID interface{}, event stripe.Event) (done bool, err error) {
	account, err := accountService.ByID(accountID)
	amountPaid := event.Data.Object["amount_paid"].(float64)
	if amountPaid > 0 {
		periodEnd := event.Data.Object["period_end"].(float64)
		periodStart := event.Data.Object["period_start"].(float64)
		paidAt := event.Data.Object["created"].(float64)
		document, _ := fattura24Service.makeInvoiceDocument(account, paidAt, amountPaid, periodStart, periodEnd)
		out, _ := xml.MarshalIndent(document, " ", "  ")
		documentStr := fmt.Sprintf("<?xml version=\"1.0\" encoding=\"UTF-8\"?>%s", string(out))
		endpoint := os.Getenv("FATTURA24_URL")
		data := url.Values{}
		data.Set("apiKey", os.Getenv("FATTURA24_KEY"))
		data.Set("xml", documentStr)

		client := &http.Client{}
		r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
		if err != nil {
			log.Fatal(err)
		}
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		res, err := client.Do(r)
		if err != nil {
			return false, err
		}
		defer res.Body.Close()
	}
	return true, err
}

func (fattura24Service *Fattura24Service) DummyGenerateInvoice(accountID interface{}) (done bool, err error) {
	account, err := accountService.ByID(accountID)
	amountPaid := float64(990)
	if amountPaid > 0 {
		paidAt := float64(1602743028)
		periodEnd := float64(1605421412)
		periodStart := float64(1602743012)
		document, _ := fattura24Service.makeInvoiceDocument(account, paidAt, amountPaid, periodStart, periodEnd)
		out, _ := xml.MarshalIndent(document, " ", "  ")
		documentStr := fmt.Sprintf("<?xml version=\"1.0\" encoding=\"UTF-8\"?>%s", string(out))
		jsonData := make(map[string]string)
		jsonData["apiKey"] = os.Getenv("FATTURA24_KEY")
		jsonData["xml"] = documentStr
		jsonValue, _ := json.Marshal(jsonData)
		_, err := http.Post(os.Getenv("FATTURA24_URL"), "application/x-www-form-urlencoded", bytes.NewBuffer(jsonValue))
		if err != nil {
			return false, err
		}
	}
	return true, err
}

func (fattura24Service *Fattura24Service) makeInvoiceDocument(account *models.Account, paidAt float64, amountPaid float64, periodStart float64, periodEnd float64) (document Fattura24, err error) {
	finalCents := amountPaid / 100
	netPriceCents := (100 * finalCents) / 122
	vatAmountCents := finalCents - netPriceCents

	final := fmt.Sprintf("%.2f", finalCents)
	netPrice := fmt.Sprintf("%.2f", netPriceCents)
	vatAmount := fmt.Sprintf("%.2f", vatAmountCents)

	paidAtDate := time.Unix(int64(paidAt), 0)
	periodEndDate := time.Unix(int64(periodEnd), 0)
	periodStartDate := time.Unix(int64(periodStart), 0)

	var fattura24Document Fattura24 = Fattura24{
		Document: Document{
			TotalWithoutTax:   netPrice,
			VatAmount:         vatAmount,
			DocumentType:      "FE",
			SendEmail:         false,
			FePaymentCode:     "MP08",
			Object:            "Minimarket24 rata piano abbonamento",
			Total:             final,
			PaymentMethodName: "Carta di credito",
			Payments: []Payment{{
				Paid:   true,
				Amount: final,
				Date:   paidAtDate.Format("2006-01-02"),
			}},
			CustomerName:      account.CompanyName,
			CustomerAddress:   account.CompanyBillingAddress,
			CustomerVatCode:   account.CompanyVat,
			CustomerCellPhone: account.CompanyPhone,
			CustomerEmail:     account.CompanyEmail,
			FeCustomerPec:     account.CompanyPec,
			FeDestinationCode: account.CompanySdi,
			FootNotes:         "Grazie per aver utilizzato MiniMarket24",
			Rows: []Row{{
				Code:           "Abbonamento MiniMarket24",
				Description:    fmt.Sprintf("Rinnovo abbonamento MiniMarket24 - da %s a %s", periodStartDate.Format("2006-01-02"), periodEndDate.Format("2006-01-02")),
				Price:          netPrice,
				VatCode:        22,
				VatDescription: "22%%",
				Qty:            1,
			}},
		},
	}
	return fattura24Document, err
}
