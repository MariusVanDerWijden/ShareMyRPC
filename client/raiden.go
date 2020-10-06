package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Raiden struct {
	url string
}

// NewRaiden creates a new raiden connector with the given url.
// Ex. http://localhost:5001/api/v1
func NewRaiden(url string) *Raiden {
	return &Raiden{
		url: url,
	}
}

func putReq(url string, data []byte) (*http.Response, error) {
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// TODO set payload
	return http.DefaultClient.Do(req)
}

func (r *Raiden) GetTokenList() error {
	url := fmt.Sprintf("%v/%v", r.url, "tokens")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// TODO handle resp
	return nil
}

func (r *Raiden) JoinNetwork(token, funds string) error {
	type request struct {
		Funds string `json:"funds"`
	}
	req := request{
		Funds: funds,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%v/%v/%v", r.url, "connections", token)
	putReq(url, data)
	return nil
}

type Opening struct {
	/*
			"token_network_address": "0x3C158a20b47d9613DDb9409099Be186fC272421a",
		    "channel_identifier": "99",
		    "partner_address": "0x61C808D82A3Ac53231750daDc13c777b59310bD9",
		    "token_address": "0x9aBa529db3FF2D8409A1da4C9eB148879b046700",
		    "balance": "1337",
		    "total_deposit": "1337",
		    "total_withdraw": "0",
		    "state": "opened",
		    "settle_timeout": "500",
			"reveal_timeout": "50"
	*/
}

func (r *Raiden) OpenChannel(partner, token, deposit, timeout string) error {
	type request struct {
		PartnerAddress string `json:"partner_address"`
		TokenAddress   string `json:"token_address"`
		TotalDeposit   string `json:"total_deposit"`
		SettleTimeout  string `json:"settle_timeout"`
	}
	req := request{
		PartnerAddress: partner,
		SettleTimeout:  timeout,
		TokenAddress:   token,
		TotalDeposit:   deposit,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%v/%v", r.url, "channels")
	putReq(url, data)
	return nil
}

func (r *Raiden) QueryChannel(token, other string) error {
	url := fmt.Sprintf("%v/%v/%v/%v", r.url, "channels", token, other)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// TODO handle resp
	return nil
}

func (r *Raiden) DepositToken(token, other, amount string) error {
	type request struct {
		TotalDeposit string `json:"total_deposit"`
	}
	req := request{
		TotalDeposit: amount,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%v/%v/%v/%v", r.url, "channels", token, other)
	putReq(url, data)
	return nil
}

func (r *Raiden) PayToken(token, other, amount string) error {
	type request struct {
		Amount string `json:"amount"`
	}
	req := request{
		Amount: amount,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%v/%v/%v/%v", r.url, "channels", token, other)
	putReq(url, data)
	return nil
}
