package app

import (
	"sort"
)

// Counter counts the frequency of string values.
// Note: do not track top values continuously. Expected case, that would result
// in worse runtime than O(N*T) runtime of the on-demand solution.
type Counter struct {
	Values map[string]uint32
}

// NewCounter returns a pointer to a new Counter instance
func NewCounter() *Counter {
	return &Counter{Values: map[string]uint32{}}
}

// Count adds to the frequency of the given hex value
func (hc *Counter) Count(hex string) {
	hc.Values[hex] = hc.Values[hex] + 1
}

// Slice converts a Counter into a slice of (string, count) pairs
func (hc *Counter) Slice() []HexPair {
	hps := []HexPair{}
	for hex, count := range hc.Values {
		hps = append(hps, HexPair{Hex: hex, Count: count})
	}
	return hps
}

// Top returns the n most frequent string values, sorted by frequency
func (hc *Counter) Top(n int) []string {
	pairs := hc.Slice()
	sort.Sort(sort.Reverse(ByCount(pairs)))
	strs := []string{}
	for _, hp := range pairs {
		strs = append(strs, hp.Hex)
	}
	// Ending index should be bounded by length of srts
	e := len(strs)
	if n <= e {
		e = n
	}
	return strs[0:e]
}

// HexPair represents a string and its frequency
type HexPair struct {
	Hex   string
	Count uint32
}

// ByCount implements sort.Interface for HexPairs based on HexPair.Count
type ByCount []HexPair

func (hps ByCount) Len() int      { return len(hps) }
func (hps ByCount) Swap(i, j int) { hps[i], hps[j] = hps[j], hps[i] }
func (hps ByCount) Less(i, j int) bool {
	if hps[i].Count == hps[j].Count {
		return hps[i].Hex < hps[j].Hex
	}
	return hps[i].Count < hps[j].Count
}
