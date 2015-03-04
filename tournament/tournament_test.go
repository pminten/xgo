package tournament

import (
	"bytes"
	"strings"
	"testing"
)

var testCases = []struct {
	description string
	input       string
	expected    string
}{
	{
		description: "good",
		input: `
Allegoric Alaskians;Blithering Badgers;win
Devastating Donkeys;Courageous Californians;draw
Devastating Donkeys;Allegoric Alaskians;win
Courageous Californians;Blithering Badgers;loss
Blithering Badgers;Devastating Donkeys;loss
Allegoric Alaskians;Courageous Californians;win
`,
		expected: `
Team                           | MP |  W |  D |  L |  P
Devastating Donkeys            |  3 |  2 |  1 |  0 |  7
Allegoric Alaskians            |  3 |  2 |  0 |  1 |  6
Blithering Badgers             |  3 |  1 |  0 |  2 |  3
Courageous Californians        |  3 |  0 |  1 |  2 |  1
`[1:], // [1:] = strip initial readability newline
	},
	{
		description: "ignore bad lines",
		input: `
Allegoric Alaskians;Blithering Badgers;win
Devastating Donkeys_Courageous Californians;draw
Devastating Donkeys;Allegoric Alaskians;win

Courageous Californians;Blithering Badgers;loss
Bla;Bla;Bla
Blithering Badgers;Devastating Donkeys;loss
# Yackity yackity yack
Allegoric Alaskians;Courageous Californians;win
Devastating Donkeys;Courageous Californians;draw
Devastating Donkeys@Courageous Californians;draw
Devastating Donkeys;Allegoric Alaskians;dra
`,
		expected: `
Team                           | MP |  W |  D |  L |  P
Devastating Donkeys            |  3 |  2 |  1 |  0 |  7
Allegoric Alaskians            |  3 |  2 |  0 |  1 |  6
Blithering Badgers             |  3 |  1 |  0 |  2 |  3
Courageous Californians        |  3 |  0 |  1 |  2 |  1
`[1:],
	},
	{
		description: "incomplete competition",
		input: `
Allegoric Alaskians;Blithering Badgers;win
Devastating Donkeys;Allegoric Alaskians;win
Courageous Californians;Blithering Badgers;loss
Allegoric Alaskians;Courageous Californians;win
`,
		expected: `
Team                           | MP |  W |  D |  L |  P
Allegoric Alaskians            |  3 |  2 |  0 |  1 |  6
Blithering Badgers             |  2 |  1 |  0 |  1 |  3
Devastating Donkeys            |  1 |  1 |  0 |  0 |  3
Courageous Californians        |  2 |  0 |  0 |  2 |  0
`[1:],
	},
}

// Simply strip the spaces of all the strings to get a canonical
// input. The spaces are only for readability of the tests.
func prepare(lines []string) []string {
	newLines := make([]string, len(lines))
	for i, l := range lines {
		newLines[i] = strings.Replace(l, " ", "", -1)
	}
	return newLines
}

func TestTally(t *testing.T) {
	for _, tt := range testCases {
		reader := strings.NewReader(tt.input)
		var buffer bytes.Buffer
		err := Tally(reader, &buffer)
		actual := buffer.String()
		// We don't expect errors for any of the test cases
		if err != nil {
			t.Fatalf("Tally for input named %q returned error %q. Error not expected.",
				tt.description, err)
		}
		if actual != tt.expected {
			t.Fatalf("Tally for input named %q was expected to return...\n%s\n...but returned...\n%s",
				tt.description, tt.expected, actual)
		}
	}
}

func BenchmarkTally(b *testing.B) {

	b.StopTimer()

	for _, tt := range testCases {
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			var buffer bytes.Buffer
			Tally(strings.NewReader(tt.input), &buffer)
		}

		b.StopTimer()
	}

}
