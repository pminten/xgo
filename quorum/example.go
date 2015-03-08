package quorum

import (
	"bytes"
	"fmt"
	"net"
)

func RunQuorum(port int, clients int) error {
	if clients < 1 {
		return fmt.Errorf("Invalid number of clients: %d", clients)
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}
	defer listener.Close()
	voteChan := make(chan string, clients)
	outcomeChan := make(chan string, clients)
	errorChan := make(chan error)
	for i := 0; i < clients; i++ {
		go runClient(listener, voteChan, outcomeChan, errorChan)
	}
	votes := make([]string, 0, clients)
	for i := 0; i < clients; i++ {
		select {
		case vote := <-voteChan:
			votes = append(votes, vote)
		case err := <-errorChan:
			return err
		}
	}
	outcome := determineQuorum(votes)
	for i := 0; i < clients; i++ {
		outcomeChan <- outcome
	}
	return nil
}

func runClient(listener net.Listener, voteChan chan<- string, outcomeChan <-chan string, errorChan chan<- error) {
	var buffer bytes.Buffer
	conn, err := listener.Accept()
	if err != nil {
		errorChan <- err
		return
	}
	defer conn.Close()
	_, err = buffer.ReadFrom(conn)
	if err != nil {
		errorChan <- err
		return
	}
	voteChan <- buffer.String()
	outcome := <-outcomeChan
	_, err = conn.Write([]byte(outcome))
	if err != nil {
		errorChan <- err
		return
	}
}

func determineQuorum(votes []string) string {
	req := len(votes)/2 + 1
	voteMap := make(map[string]int)
	for _, vote := range votes {
		voteMap[vote] = voteMap[vote] + 1
	}
	for vote, amount := range voteMap {
		if amount >= req {
			return vote
		}
	}
	return "$no_quorum"
}
