package raiden

import (
	"fmt"
	"testing"
)

func TestGetTokenList(t *testing.T) {
	url := "http://localhost:5001/api/v1"
	r := NewRaiden(url)
	out, err := r.GetTokenList()
	if err != nil {
		t.Error(err)
	}
	data := string(out)
	t.Error(data)
}

func TestJoinNetwork(t *testing.T) {
	url := "http://localhost:5001/api/v1"
	token := "0x95B2d84De40a0121061b105E6B54016a49621B44"
	funds := "100"
	r := NewRaiden(url)
	out, err := r.JoinNetwork(token, funds)
	if err != nil {
		t.Error(err)
	}
	data := string(out)
	t.Error(data)
}

func TestOpenChannel(t *testing.T) {
	url := "http://localhost:5001/api/v1"
	token := "0x95B2d84De40a0121061b105E6B54016a49621B44"
	raidenhub := "0x1F916ab5cf1B30B22f24Ebf435f53Ee665344Acf"
	deposit := "10"
	timeout := "1000"
	r := NewRaiden(url)
	open, err := r.OpenChannel(raidenhub, token, deposit, timeout)
	if err != nil {
		t.Error(err)
	}
	t.Error(open.Balance)
}

func TestPayToken(t *testing.T) {
	url := "http://localhost:5001/api/v1"
	token := "0x95B2d84De40a0121061b105E6B54016a49621B44"
	raidenhub := "0x1F916ab5cf1B30B22f24Ebf435f53Ee665344Acf"
	r := NewRaiden(url)
	msg, err := r.PayToken(token, raidenhub, "1")
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal(msg)
}

func TestQueryChannel(t *testing.T) {
	url := "http://localhost:5001/api/v1"
	token := "0x95B2d84De40a0121061b105E6B54016a49621B44"
	raidenhub := "0x1F916ab5cf1B30B22f24Ebf435f53Ee665344Acf"
	r := NewRaiden(url)
	open, err := r.QueryChannel(token, raidenhub)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(open)
}

func TestPaymentHistory(t *testing.T) {
	url := "http://localhost:5001/api/v1"
	token := "0x95B2d84De40a0121061b105E6B54016a49621B44"
	raidenhub := "0x1F916ab5cf1B30B22f24Ebf435f53Ee665344Acf"
	r := NewRaiden(url)
	open, err := r.PaymentHistory(token, raidenhub)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(open)
}
