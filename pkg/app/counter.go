package app

import "sort"

// Counter ...TODO
// Note: Do not track top values. Expected case, that would be worse runtime
// than O(N*T) runtime of on-demand solution.
type Counter struct {
	Values map[string]uint32
}

// NewCounter ...TODO
func NewCounter() *Counter {
	return &Counter{Values: map[string]uint32{}}
}

// Count ...TODO
func (hc *Counter) Count(hex string) {
	hc.Values[hex] = hc.Values[hex] + 1
}

// Slice ...TODO
func (hc *Counter) Slice() []HexPair {
	hps := []HexPair{}
	for hex, count := range hc.Values {
		hps = append(hps, HexPair{Hex: hex, Count: count})
	}
	return hps
}

// Top ...TODO
func (hc *Counter) Top(n int) []string {
	pairs := hc.Slice()
	sort.Sort(sort.Reverse(ByCount(pairs)))
	strs := []string{}
	for _, hp := range pairs {
		strs = append(strs, hp.Hex)
	}
	return strs[0:n]
}

// HexPair ...TODO
type HexPair struct {
	Hex   string
	Count uint32
}

// ByCount implements sort.Interface for HexPairs based on HexPair.Count
type ByCount []HexPair

func (hps ByCount) Len() int           { return len(hps) }
func (hps ByCount) Swap(i, j int)      { hps[i], hps[j] = hps[j], hps[i] }
func (hps ByCount) Less(i, j int) bool { return hps[i].Count < hps[j].Count }
