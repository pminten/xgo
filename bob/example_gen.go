//+build gentests

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	d, err := ioutil.ReadFile("example_gen.json")
	if err != nil {
		log.Fatal(err)
	}
	var j []struct {
		Case string
		In   string
		Rep  int `json:"random-repeat"`
		Want string
	}
	err = json.Unmarshal(d, &j)
	if len(j) == 0 {
		log.Print("no data found with expected structure")
		log.Fatal(err)
	}
	f, err := os.Create("cases_test.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fmt.Fprint(f, `package bob

var testCases = []struct {
	desc string
	in   string
	rep  int
	want string
}{`)
	for _, tc := range j {
		fmt.Fprintf(f, `
	{
		%q,
		%q,
		%d,
		%q,
	},`, tc.Case, tc.In, tc.Rep, tc.Want)
	}
	fmt.Fprint(f, "\n}\n")
}
