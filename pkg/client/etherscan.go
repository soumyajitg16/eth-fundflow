package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Base URL of Etherscan API to fetch the transaction data
const baseURL = "https://api.etherscan.io/api"

// EtherscanClient holds API key and HTTP client for making requests
type EtherscanClient struct {
	apiKey string
	http   *http.Client
}

// NewEtherscanClient constructs an EtherscanClient with timeout
func NewEtherscanClient(key string) *EtherscanClient {
	return &EtherscanClient{
		apiKey: key,
		http:   &http.Client{Timeout: 10 * time.Second},
	}
}

// APIResponse models the common Etherscan JSON envelope
type APIResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

// Transaction represents a normal or internal Ethereum transaction
type Transaction struct {
	BlockNumber string `json:"blockNumber"`
	TimeStamp   string `json:"timeStamp"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	Hash        string `json:"hash"`
}

// TokenTransfer represents NFT transfer trasactions
type TokenTransfer struct {
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
	Hash            string `json:"hash"`
	TimeStamp       string `json:"timeStamp"`
	TokenName       string `json:"tokenName"`
	TokenSymbol     string `json:"tokenSymbol"`
	ContractAddress string `json:"contractAddress"`
}

// FetchNormalTx retrieves normal transactions (Etherscan "txlist" endpoint) for an address
func (c *EtherscanClient) FetchNormalTx(address string) ([]Transaction, error) {
	url := fmt.Sprintf("%s?module=account&action=txlist&address=%s&sort=asc&apikey=%s", baseURL, address, c.apiKey)
	return c.doFetch(url)
}

// FetchInternalTx retrieves internal transactions ( Etherscan "txlistinternal" endpoint) for an address
func (c *EtherscanClient) FetchInternalTx(address string) ([]Transaction, error) {
	url := fmt.Sprintf("%s?module=account&action=txlistinternal&address=%s&sort=asc&apikey=%s", baseURL, address, c.apiKey)
	return c.doFetch(url)
}

// FetchTokenTransfers retrieves ERC20/ERC721/ERC1155 transfers for an address
func (c *EtherscanClient) FetchTokenTransfers(address string) ([]TokenTransfer, error) {
	url := fmt.Sprintf("%s?module=account&action=tokentx&address=%s&sort=asc&apikey=%s", baseURL, address, c.apiKey)
	var resp struct {
		APIResponse
		Result []TokenTransfer `json:"result"`
	}
	if err := c.doRequest(url, &resp); err != nil {
		return nil, err
	}
	return resp.Result, nil
}

// doFetch is a helper for endpoints returning []Transaction
func (c *EtherscanClient) doFetch(url string) ([]Transaction, error) {
	var resp struct {
		APIResponse
		Result []Transaction `json:"result"`
	}
	if err := c.doRequest(url, &resp); err != nil {
		return nil, err
	}
	return resp.Result, nil
}

// doRequest performs the GET request and decodes JSON, with status check
func (c *EtherscanClient) doRequest(url string, out interface{}) error {
	resp, err := c.http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(" Etherscan API returned status %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}