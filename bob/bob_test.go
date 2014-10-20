package bob

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

const testVersion = 1

// Retired testVersions
// (none) 79937f6d58e25ebafe12d1cb4a9f88f4de70cfd6

func TestHeyBob(t *testing.T) {
	if TestVersion != testVersion {
		t.Fatalf("Found TestVersion = %v, want %v", TestVersion, testVersion)
	}
	rand.Seed(time.Now().Unix())
	for _, tt := range testCases {
		in := tt.in
		if tt.rep > 0 {
			in = strings.Repeat(in, 1+rand.Intn(tt.rep))
		}
		t.Logf("Test case %q", in)
		if actual := Hey(in); actual != tt.want {
			t.Fatalf(msg, tt.desc, tt.in, actual, tt.want)
		}
	}
	t.Log("Tested", len(testCases), "cases.")
}

const msg = `
ALICE (%s): %s
BOB: %s

Expected Bob to respond: %s`

func BenchmarkBob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tt := range testCases {
			Hey(tt.in)
		}
	}
}
