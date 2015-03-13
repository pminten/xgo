// +build !example

package ledger

import (
	"errors"
	"fmt"
	"strings"
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
	var entriesCopy []Entry
	for _, e := range entries {
		entriesCopy = append(entriesCopy, e)
	}

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

	var s string
	if locale == "nl-NL" {
		s = "Datum" +
			strings.Repeat(" ", 10-len("Datum")) +
			" | " +
			"Omschrijving" +
			strings.Repeat(" ", 25-len("Omschrijving")) +
			" | " + "Verandering" + "\n"
	} else if locale == "en-US" {
		s = "Date" +
			strings.Repeat(" ", 10-len("Date")) +
			" | " +
			"Description" +
			strings.Repeat(" ", 25-len("Description")) +
			" | " + "Change" + "\n"
	} else {
		return "", errors.New("")
	}
	// Parallelism, always a great idea
	co := make(chan struct {
		s string
		e error
	})
	for _, et := range entriesCopy {
		go func(entry Entry) {
			if len(entry.Date) != 10 {
				co <- struct {
					s string
					e error
				}{e: errors.New("")}
			}
			d1, d2, d3, d4, d5 := entry.Date[0:4], entry.Date[4], entry.Date[5:7], entry.Date[7], entry.Date[8:10]
			if d2 != '-' {
				co <- struct {
					s string
					e error
				}{e: errors.New("")}
			}
			if d4 != '-' {
				co <- struct {
					s string
					e error
				}{e: errors.New("")}
			}
			de := entry.Description
			if len(de) > 27 {
				de = de[:24] + "..."
			} else {
				de = de + strings.Repeat(" ", 25-len(de))
			}
			var d string
			if locale == "nl-NL" {
				d = d5 + "-" + d3 + "-" + d1
			} else if locale == "en-US" {
				d = d3 + "/" + d5 + "/" + d1
			}
			a := locInfo.Currency(symbol, entry.Change)
			var al int
			for _ = range a {
				al++
			}
			co <- struct {
				s string
				e error
			}{s: d + strings.Repeat(" ", 10-len(d)) + " | " + de + " | " +
				strings.Repeat(" ", 13-al) + a + "\n"}
		}(et)
	}
	for _ = range entriesCopy {
		v := <-co
		if v.e != nil {
			return "", v.e
		} else {
			s += v.s
		}
	}
	return s, nil
}

var currencySymbols = map[string]string{
	"USD": "$",
	"EUR": "€",
}

type localeInfo struct {
	currency func(symbol string, cents int, negative bool) string
}

func (f localeInfo) Currency(symbol string, cents int) string {
	negative := false
	if cents < 0 {
		cents = cents * -1
		negative = true
	}
	return f.currency(symbol, cents, negative)
}

var locales = map[string]localeInfo{
	"nl-NL": {
		currency: dutchCurrencyFormat,
	},
	"en-US": {
		currency: americanCurrencyFormat,
	},
}

// The sign and amount are passed in separately to simplify some logic.
func dutchCurrencyFormat(symbol string, cents int, negative bool) string {
	var s string
	s += symbol
	s += " "
	s += moneyToString(cents, ".", ",")
	if negative {
		s += "-"
	} else {
		s += " "
	}
	return s
}

func americanCurrencyFormat(symbol string, cents int, negative bool) string {
	var s string
	if negative {
		s += "("
	}
	s += symbol
	s += moneyToString(cents, ",", ".")
	if negative {
		s += ")"
	} else {
		s += " "
	}
	return s
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
	s := strings.Join(revParts, thousandsSep)
	s += decimalSep
	s += centsPart
	return s
}
