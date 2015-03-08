package quorum

import (
	"bytes"
	"fmt"
	"net"
	"testing"
	"time"
)

// Define a function RunQuorum(port int, clients int) error.
//
// Also define an exported TestVersion with a value that matches
// the internal testVersion here.

const testVersion = 1

// Use a port in the Dynamic Ports range (see RFC6335).
const PORT = 54321

func TestRunQuorumZeroFail(t *testing.T) {
	if RunQuorum(PORT, 0) == nil {
		t.Fatalf("No failure when RunQuorum was given invalid argument '0'")
	}
}

func TestRunQuorumNegativeFail(t *testing.T) {
	if RunQuorum(PORT, -3) == nil {
		t.Fatalf("No failure when RunQuorum was given invalid argument '-3'")
	}
}

// interact runs a client which votes for 'vote'. Returns the quorum vote or 'error'.
func interact(vote string) (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", PORT))
	if err != nil {
		return "", err
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	_, err = conn.Write([]byte(vote))
	if err != nil {
		return "", err
	}
	conn.CloseWrite()
	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(conn)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// Either outcome or error is used, never both.
type clientResult struct {
	outcome string
	err     error
}

// runClients runs multiple clients and collects the results or
// returns the first error from a worker.
func runClients(votes []string) ([]string, error) {
	resultChan := make(chan clientResult, len(votes))
	for _, vote := range votes {
		go func(v string) {
			outcome, err := interact(v)
			if err == nil {
				resultChan <- clientResult{outcome: outcome}
			} else {
				resultChan <- clientResult{err: err}
			}
		}(vote)
	}
	results := make([]string, 0, len(votes))
	for i := 0; i < len(votes); i++ {
		select {
		case result := <-resultChan:
			if result.err == nil {
				results = append(results, result.outcome)
			} else {
				return nil, result.err
			}
		case <-time.After(10 * time.Second):
			return nil, fmt.Errorf("resultChan read timed out")
		}
	}
	return results, nil
}

// Either outcomes or error is used, never both.
type clientsResult struct {
	outcomes []string
	err      error
}

func quorumTest(t *testing.T, votes []string, expected string) {
	quorumChan := make(chan error)
	clientsChan := make(chan clientsResult)
	go func() {
		quorumChan <- RunQuorum(PORT, len(votes))
	}()
	go func() {
		outcomes, err := runClients(votes)
		if err == nil {
			clientsChan <- clientsResult{outcomes: outcomes}
		} else {
			clientsChan <- clientsResult{err: err}
		}
	}()

	select {
	case result := <-clientsChan:
		if result.err != nil {
			t.Fatalf("Error from clients for votes %v: %v", votes, result.err)
		} else {
			first := result.outcomes[0]
			for _, outcome := range result.outcomes {
				if outcome != first {
					t.Fatalf("Different outcomes for votes %v, no quorum: %q and %q",
						votes, outcome, first)
				}
			}
			if first != expected {
				t.Fatalf("Wrong quorum reached for votes %v: expected %q, got %q",
					votes, expected, first)
			}
		}
	}
	select {
	case err := <-quorumChan:
		if err != nil {
			t.Fatalf("Unexpected error from RunQuorum for votes %v: %v", votes, err)
		}
	case <-time.After(15 * time.Second):
		t.Fatalf("RunQuorum didn't properly terminate for votes %v", votes)
	}
}

var testCases = []struct {
	votes    []string
	expected string
}{
	{
		votes:    []string{"a", "b", "a"},
		expected: "a",
	},
	{
		votes:    []string{"a"},
		expected: "a",
	},
	{
		votes:    []string{"a", "a", "a", "b"},
		expected: "a",
	},
	{
		votes:    []string{"a", "a", "a", "a", "a", "b", "b", "b", "c"},
		expected: "a",
	},
	{
		votes:    []string{"a", "b"},
		expected: "$no_quorum",
	},
	{
		votes:    []string{"a", "a", "b", "b"},
		expected: "$no_quorum",
	},
	{
		votes:    []string{"a", "b", "c"},
		expected: "$no_quorum",
	},
}

func TestRunQuorum(t *testing.T) {
	for _, tt := range testCases {
		quorumTest(t, tt.votes, tt.expected)
	}
}

func BenchRunQuorum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tt := range testCases {
			finishChan := make(chan bool)
			go func() {
				RunQuorum(PORT, len(tt.votes))
				finishChan <- true
			}()
			go func() {
				runClients(tt.votes)
				finishChan <- true
			}()
			<-finishChan
			<-finishChan
		}
	}
}
