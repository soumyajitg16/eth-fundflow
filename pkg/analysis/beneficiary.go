package analysis

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
	"sync"

	"github.com/soumyajitg16/eth-fundflow.git/pkg/client" // Etherscan API client
)

// BeneficiaryHandler handles /beneficiary requests, fetching all tx types in parallel
func BeneficiaryHandler(es *client.EtherscanClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Validating query parameter
		addr := r.URL.Query().Get("address")
		if addr == "" {
			http.Error(w, "address query param required", http.StatusBadRequest)
			return
		}

		var wg sync.WaitGroup
		var normal, internal []client.Transaction
		var token []client.TokenTransfer
		var err1, err2, err3 error

		// Launch concurrent fetches
		wg.Add(3)
		go func() { defer wg.Done(); normal, err1 = es.FetchNormalTx(addr) }()
		go func() { defer wg.Done(); internal, err2 = es.FetchInternalTx(addr) }()
		go func() { defer wg.Done(); token, err3 = es.FetchTokenTransfers(addr) }()
		wg.Wait()

		// Handle fetch errors
		if err1 != nil || err2 != nil || err3 != nil {
			log.Printf("Error fetching tx: %v %v %v", err1, err2, err3)
			http.Error(w, "error fetching transactions", http.StatusInternalServerError)
			return
		}

		// Trace beneficiaries
		result := AnalyzeBeneficiaries(addr, normal, internal, token)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{"message": "success", "data": result}); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	}
}

// Beneficiary holds aggregated outflow data per address
type Beneficiary struct {
	BeneficiaryAddress string                   `json:"beneficiary_address"`
	Amount             float64                  `json:"amount"`      // in ETH
	Transactions       []map[string]interface{} `json:"transactions"` // ordered by time
}

// AnalyzeBeneficiaries aggregates and sorts outflows across tx types
func AnalyzeBeneficiaries(subject string, normal, internal []client.Transaction, token []client.TokenTransfer) []Beneficiary {
	// helper to parse wei string to ETH float
	parseWei := func(raw string) float64 {
		val, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			log.Printf("ParseWei error: %v", err)
			return 0
		}
		return val / 1e18
	}

	// map to accumulate beneficiaries
	acc := make(map[string]*Beneficiary)
	addOutflow := func(to string, amount float64, ts string, hash string) {
		if _, ok := acc[to]; !ok {
			acc[to] = &Beneficiary{BeneficiaryAddress: to, Amount: 0, Transactions: []map[string]interface{}{}}
		}
		acc[to].Amount += amount
		// convert timestamp to readable ISO
		t, _ := strconv.ParseInt(ts, 10, 64)
		timeStr := time.Unix(t, 0).UTC().Format(time.RFC3339)
		acc[to].Transactions = append(acc[to].Transactions, map[string]interface{}{ 
			"tx_amount": amount,
			"date_time": timeStr,
			"transaction_id": hash,
		})
	}

	// normal tx
	for _, tx := range normal {
		if tx.From == subject {
			addOutflow(tx.To, parseWei(tx.Value), tx.TimeStamp, tx.Hash)
		}
	}
	// internal tx
	for _, tx := range internal {
		if tx.From == subject {
			addOutflow(tx.To, parseWei(tx.Value), tx.TimeStamp, tx.Hash)
		}
	}
	// token tx
	for _, t := range token {
		if t.From == subject {
			addOutflow(t.To, parseWei(t.Value), t.TimeStamp, t.Hash)
		}
	}

	// move to slice & sort high to low
	out := make([]Beneficiary, 0, len(acc))
	for _, b := range acc {
		out = append(out, *b)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Amount > out[j].Amount })
	return out
}