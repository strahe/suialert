package benchmark

import (
	"bytes"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

var (
	changeTypes []string
)

func GetChangeType() string {
	return changeTypes[rand.Intn(len(changeTypes))]
}

//type CoinBalanceChange struct {
//	PackageId         string       `json:"packageId"`
//	TransactionModule string       `json:"transactionModule"`
//	Sender            string       `json:"sender"`
//	ChangeType        string       `json:"changeType"`
//	Owner             *ObjectOwner `json:"owner"`
//	CoinType          string       `json:"coinType"`
//	CoinObjectId      string       `json:"coinObjectId"`
//	Version           int64        `json:"version"`
//	Amount            int64        `json:"amount"`
//}

// MakeRule make a single dummy rule
func MakeRule(name string, seq int) string {
	rname := name + strconv.Itoa(seq)
	buff := &bytes.Buffer{}
	buff.WriteString("rule ")
	buff.WriteString(rname)
	buff.WriteString(" \"")
	buff.WriteString(strconv.Itoa(seq))
	buff.WriteString(" ")
	buff.WriteString(strings.ToTitle(lo.RandomString(8, lo.LettersCharset)))
	buff.WriteString(" ")
	buff.WriteString(strings.ToTitle(lo.RandomString(8, lo.LettersCharset)))
	buff.WriteString("\"")
	buff.WriteString(" salience ")
	buff.WriteString(strconv.Itoa(rand.Intn(100) + 10))
	buff.WriteString(" {\n\t")
	buff.WriteString("when\n\t\t")
	buff.WriteString("Event.ChangeType")
	buff.WriteString(" == \"")
	buff.WriteString(lo.Sample[string]([]string{"Pay", "Gas", "Receive"}))
	buff.WriteString("\"")
	buff.WriteString(" && \n\t\t")
	buff.WriteString("Event.Amount")
	buff.WriteString(lo.Sample[string]([]string{" > ", " >= ", " == ", " < ", " <= "}))
	buff.WriteString(strconv.Itoa(rand.Intn(100000000)))
	buff.WriteString("\n\tthen\n\t\t")
	buff.WriteString("Event.Matched(")
	buff.WriteString(strconv.Itoa(rand.Intn(10000)))
	buff.WriteString(");\n")
	buff.WriteString("}\n\n")

	return buff.String()
}

// GenRandomRule simply generate count number of simple parse-able rule into a file
func GenRandomRule(fileName string, count int) error {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close() // nolint: errcheck
	for i := 1; i <= count; i++ {
		_, err := f.WriteString(MakeRule(lo.RandomString(40, lo.AllCharset), i))
		if err != nil {
			return err
		}
	}
	return nil
}
