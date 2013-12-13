package gobtce

import (
	"testing"
)

func TestRefreshTicker(t *testing.T) {
	x := PairTicker{PairBTCUSD, Ticker{}}
	if err := x.Refresh(); err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v\n", x)
}

func TestRefreshTickers(t *testing.T) {
	x := PairTickers{
		PairBTCUSD: Ticker{},
		PairLTCUSD: Ticker{},
		PairBTCRUR: Ticker{},
		PairLTCRUR: Ticker{},
	}
	n := len(x)
	if err := x.Refresh(); err != nil {
		t.Error(err)
		return
	}
	if len(x) != n {
		t.Errorf("Expected %v pairs, got %v", n, len(x))
		return
	}
	t.Logf("%#v\n", x)
}
