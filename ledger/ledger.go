// +build !example

package ledger

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

const TestVersion = 1

type Entry struct {
	Date        string // "Y-m-d"
	Description string
	Change      int // in cents
}

func FormatLedger(currency string, locale string, entries []Entry) (string, error) {
	symbol, found := currencySymbols[currency]
	if !found {
		return "", fmt.Errorf("Invalid or unknown currency %q", currency)
	}
	locInfo, found := locales[locale]
	if !found {
		return "", fmt.Errorf("Invalid or unknown locale %q", locale)
	}
	entriesCopy := make([]Entry, len(entries))
	copy(entriesCopy, entries)

	m1 := map[bool]int{true: 0, false: 1}
	m2 := map[bool]int{true: -1, false: 1}
	es := entriesCopy
	for len(es) > 1 {
		first, rest := es[0], es[1:]
		success := false
		for !success {
			success = true
			for i, e := range rest {
				if (m1[e.Date == first.Date]*m2[e.Date < first.Date]*4 +
					m1[e.Description == first.Description]*m2[e.Description < first.Description]*2 +
					m1[e.Change == first.Change]*m2[e.Change < first.Change]*1) < 0 {
					es[0], es[i+1] = es[i+1], es[0]
					success = false
				}
			}
		}
		es = es[1:]
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%-10s | %-25s | %s\n",
		locInfo.translations["date"],
		locInfo.translations["descr"],
		locInfo.translations["change"]))
	// Parallelism, always a great idea
	co := make(chan struct {
		s string
		e error
	})
	for _, et := range entriesCopy {
		go func(entry Entry) {
			date, err := time.Parse("2006-01-02", entry.Date)
			if err != nil {
				co <- struct {
					s string
					e error
				}{e: err}
			}
			description := entry.Description
			if len(description) > 27 {
				description = description[:24] + "..."
			}
			co <- struct {
				s string
				e error
			}{s: fmt.Sprintf("%-10s | %-25s | %13s\n",
				locInfo.Date(date),
				description,
				locInfo.Currency(symbol, entry.Change))}
		}(et)
	}
	for _ = range entriesCopy {
		v := <-co
		if v.e != nil {
			return "", v.e
		} else {
			buf.WriteString(v.s)
		}
	}
	return buf.String(), nil
}

var currencySymbols = map[string]string{
	"USD": "$",
	"EUR": "€",
}

type localeInfo struct {
	currency     func(symbol string, cents int, negative bool) string
	dateFormat   string
	translations map[string]string
}

func (f localeInfo) Currency(symbol string, cents int) string {
	negative := false
	if cents < 0 {
		cents = cents * -1
		negative = true
	}
	return f.currency(symbol, cents, negative)
}

func (f localeInfo) Date(t time.Time) string {
	return t.Format(f.dateFormat)
}

//	Currency(symbol string, cents int) string
//	Date(symbol string, cents int) string

var locales = map[string]localeInfo{
	"nl-NL": {
		currency:   dutchCurrencyFormat,
		dateFormat: "02-01-2006",
		translations: map[string]string{
			"date":   "Datum",
			"descr":  "Omschrijving",
			"change": "Verandering",
		},
	},
	"en-US": {
		currency:   americanCurrencyFormat,
		dateFormat: "01/02/2006",
		translations: map[string]string{
			"date":   "Date",
			"descr":  "Description",
			"change": "Change",
		},
	},
}

// The sign and amount are passed in separately to simplify some logic.
func dutchCurrencyFormat(symbol string, cents int, negative bool) string {
	var buf bytes.Buffer
	buf.WriteString(symbol)
	buf.WriteRune(' ')
	buf.WriteString(moneyToString(cents, ".", ","))
	if negative {
		buf.WriteRune('-')
	} else {
		buf.WriteRune(' ') // To keep alignment in report
	}
	return buf.String()
}

func americanCurrencyFormat(symbol string, cents int, negative bool) string {
	var buf bytes.Buffer
	if negative {
		buf.WriteRune('(')
	}
	buf.WriteString(symbol)
	buf.WriteString(moneyToString(cents, ",", "."))
	if negative {
		buf.WriteRune(')')
	} else {
		buf.WriteRune(' ') // To keep alignment in report
	}
	return buf.String()
}

// Precondition: cents is not negative
func moneyToString(cents int, thousandsSep, decimalSep string) string {
	centsStr := fmt.Sprintf("%03d", cents) // Pad to 3 digits
	centsPart := centsStr[len(centsStr)-2:]
	rest := centsStr[:len(centsStr)-2]
	var parts []string
	for len(rest) > 3 {
		parts = append(parts, rest[len(rest)-3:])
		rest = rest[:len(rest)-3]
	}
	if len(rest) > 0 {
		parts = append(parts, rest)
	}
	revParts := make([]string, 0, len(parts))
	for i := len(parts) - 1; i >= 0; i-- {
		revParts = append(revParts, parts[i])
	}
	var buf bytes.Buffer
	buf.WriteString(strings.Join(revParts, thousandsSep))
	buf.WriteString(decimalSep)
	buf.WriteString(centsPart)
	return buf.String()
}

func sortEntries(entries []Entry) {
}
