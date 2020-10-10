package raiden

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func Req(method, url string, data []byte) (*http.Response, error) {
	reader := bytes.NewReader(data)
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}

func (r *Raiden) GetTokenList() ([]byte, error) {
	url := fmt.Sprintf("%v/%v", r.url, "tokens")
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func (r *Raiden) JoinNetwork(token, funds string) ([]byte, error) {
	type request struct {
		Funds string `json:"funds"`
	}
	req := request{
		Funds: funds,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%v/%v/%v", r.url, "connections", token)
	resp, err := Req("PUT", url, data)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

type Opening struct {
	TokenAddress        string `json:"token_address"`
	PartnerAddress      string `json:"partner_address"`
	SettleTimeout       string `json:"settle_timeout"`
	RevealTimeout       string `json:"reveal_timeout"`
	Balance             string `json:"balance"`
	TokenNetworkAddress string `json:"token_network_address"`
	TotalDeposit        string `json:"total_deposit"`
	State               string `json:"state"`
	ChannelIdentifier   string `json:"channel_identifier"`
	TotalWithdraw       string `json:"total_withdraw"`
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
	resp, err := Req("PUT", url, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	response := new(Opening)
	return response, json.NewDecoder(resp.Body).Decode(response)
}

func (r *Raiden) QueryChannel(token, other string) (*Opening, error) {
	url := fmt.Sprintf("%v/%v/%v/%v", r.url, "channels", token, other)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	response := new(Opening)
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
	_, err = Req("PATCH", url, data)
	return err
}

func (r *Raiden) PayToken(token, other, amount string) (string, error) {
	type request struct {
		Amount string `json:"amount"`
	}
	req := request{
		Amount: amount,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%v/%v/%v/%v", r.url, "payments", token, other)
	resp, err := Req("POST", url, data)
	type Message struct {
		Message string `json:"message"`
	}
	response := new(Message)
	return response.Message, json.NewDecoder(resp.Body).Decode(response)
}

type PaymentHistory []struct {
	Identifier   string `json:"identifier"`
	LogTime      string `json:"log_time"`
	Target       string `json:"target"`
	Amount       string `json:"amount"`
	Event        string `json:"event"`
	TokenAddress string `json:"token_address"`
}

func (r *Raiden) PaymentHistory(token, other string) (*PaymentHistory, error) {
	url := fmt.Sprintf("%v/%v/%v/%v", r.url, "payments", token, other)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	response := new(PaymentHistory)
	return response, json.NewDecoder(resp.Body).Decode(response)
}
