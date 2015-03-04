package tournament

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

const (
	LOSS = iota
	DRAW
	WIN
)

type inputEntry struct {
	teams    [2]string
	outcomes [2]int
}

type teamResult struct {
	team   string
	played int
	wins   int
	draws  int
	losses int
	points int
}

type TeamResultSlice []teamResult

// sort.Interface implementation, sorts on points, descending
func (s TeamResultSlice) Len() int {
	return len(s)
}

func (s TeamResultSlice) Less(i, j int) bool {
	return s[i].points > s[j].points
}

func (s TeamResultSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func readInput(reader io.Reader) ([]inputEntry, error) {
	scanner := bufio.NewScanner(reader)
	var entries []inputEntry
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ";")
		if len(parts) == 3 {
			t1, t2 := parts[0], parts[1]
			teams := [2]string{t1, t2}
			var outcomes [2]int
			if parts[2] == "win" {
				outcomes = [2]int{WIN, LOSS}
			} else if parts[2] == "loss" {
				outcomes = [2]int{LOSS, WIN}
			} else if parts[2] == "draw" {
				outcomes = [2]int{DRAW, DRAW}
			} else {
				continue
			}
			entries = append(entries, inputEntry{teams: teams, outcomes: outcomes})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

func tallyEntries(entries []inputEntry) map[string]teamResult {
	var results = make(map[string]teamResult)
	for _, entry := range entries {
		for i := 0; i < 2; i++ {
			team := entry.teams[i]
			outcome := entry.outcomes[i]
			result, present := results[team]
			if !present {
				result.team = team
				// The rest is 0, which is correct
			}
			switch outcome {
			case WIN:
				result.wins += 1
				result.points += 3
			case DRAW:
				result.draws += 1
				result.points += 1
			case LOSS:
				result.losses += 1
			}
			result.played += 1
			results[team] = result
		}
	}
	return results
}

func report(writer io.Writer, resultMap map[string]teamResult) {
	var entries TeamResultSlice = make([]teamResult, 0, len(resultMap))
	for _, entry := range resultMap {
		entries = append(entries, entry)
	}
	sort.Sort(entries)
	fmt.Fprintf(writer, "Team                           | MP |  W |  D |  L |  P\n")
	for _, entry := range entries {
		fmt.Fprintf(writer, "%-30s | %2d | %2d | %2d | %2d | %2d\n",
			entry.team, entry.played, entry.wins, entry.draws, entry.losses, entry.points)
	}
}

func Tally(reader io.Reader, writer io.Writer) error {
	entries, err := readInput(reader)
	if err != nil {
		return err
	}
	report(writer, tallyEntries(entries))
	return nil
}
