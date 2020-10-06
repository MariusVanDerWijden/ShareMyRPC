package client

import (
	"bytes"
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
	reader := bytes.NewReader(data)
	req, err := http.NewRequest("PUT", url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
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
	TokenNetworkAddress string `json:"token_network_address"`
	ChannelIdentifier   string `json:"channel_identifier"`
	PartnerAddress      string `json:"partner_address"`
	TokenAddress        string `json:"token_address"`
	Balance             string `json:"balance"`
	TotalDeposit        string `json:"total_deposit"`
	TotalWithdraw       string `json:"total_withdraw"`
	State               string `json:"state"`
	SettleTimeout       string `json:"settle_timeout"`
	RevealTimeout       string `json:"reveal_timeout"`
}

func (r *Raiden) OpenChannel(partner, token, deposit, timeout string) (*Opening, error) {
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
		return nil, err
	}
	url := fmt.Sprintf("%v/%v", r.url, "channels")
	resp, err := putReq(url, data)
	if err != nil {
		return nil, err
	}
	var response *Opening
	return response, json.NewDecoder(resp.Body).Decode(response)
}

func (r *Raiden) QueryChannel(token, other string) (*Opening, error) {
	url := fmt.Sprintf("%v/%v/%v/%v", r.url, "channels", token, other)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var response *Opening
	return response, json.NewDecoder(resp.Body).Decode(response)
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
