package transaction

import (
	"encoding/json"
	"fmt"
	"github.com/incognito-core-libs/gozilliqa-sdk/provider"
	"github.com/incognito-core-libs/gozilliqa-sdk/util"
	"strconv"
	"strings"
	"time"
)

type State int

const (
	Initialised State = iota
	Pending
	Confirmed
	Rejected
)

type Transaction struct {
	ID           string
	Version      string
	Nonce        string
	Amount       string
	GasPrice     string
	GasLimit     string
	Signature    string
	Receipt      TransactionReceipt
	SenderPubKey string
	ToAddr       string
	Code         string
	Data         interface{}

	Status          State
	ContractAddress string
	Priority        bool
}

func (t *Transaction) toTransactionParam() TxParams {
	data, _ := json.Marshal(t.Data)
	param := TxParams{
		ID:           t.ID,
		Version:      t.Version,
		Nonce:        t.Nonce,
		Amount:       t.Amount,
		GasPrice:     t.GasPrice,
		GasLimit:     t.GasLimit,
		Signature:    t.Signature,
		Receipt:      t.Receipt,
		SenderPubKey: t.SenderPubKey,
		Code:         t.Code,
		Data:         string(data),
	}

	if t.ToAddr == "" {
		param.ToAddr = "0000000000000000000000000000000000000000"
	} else {
		param.ToAddr = t.ToAddr
	}
	return param
}

func (t *Transaction) ToTransactionPayload() provider.TransactionPayload {
	version, _ := strconv.ParseInt(t.Version, 10, 32)
	nonce, _ := strconv.ParseInt(t.Nonce, 10, 32)
	data, _ := json.Marshal(t.Data)

	p := provider.TransactionPayload{
		Version:   int(version),
		Nonce:     int(nonce),
		ToAddr:    util.ToCheckSumAddress(t.ToAddr)[2:],
		Amount:    t.Amount,
		PubKey:    strings.ToLower(t.SenderPubKey),
		GasPrice:  t.GasPrice,
		GasLimit:  t.GasLimit,
		Code:      t.Code,
		Signature: strings.ToLower(t.Signature),
		Priority:  t.Priority,
	}

	if string(data) != "\"\"" {
		p.Data = string(data)
	}

	if p.ToAddr == "0000000000000000000000000000000000000000" {
		p.ToAddr = "0x0000000000000000000000000000000000000000"
	}
	return p
}

func (t *Transaction) TrackTx(hash string, provider *provider.Provider) bool {
	response := provider.GetTransaction(hash)
	if response == nil || response.Error != nil {
		return false
	}

	result := response.Result.(map[string]interface{})
	t.ID = result["ID"].(string)

	receipt, ok := result["receipt"].(map[string]interface{})
	if !ok {
		return false
	}

	t.Receipt.CumulativeGas = receipt["cumulative_gas"].(string)
	t.Receipt.EpochNum = receipt["epoch_num"].(string)
	t.Receipt.Success = receipt["success"].(bool)

	if !t.Receipt.Success {
		t.Status = Rejected
	} else {
		t.Status = Confirmed
	}
	return true
}

func (t *Transaction) Confirm(hash string, maxAttempts, interval int, provider *provider.Provider) {
	t.Status = Pending
	for i := 0; i < maxAttempts; i++ {
		fmt.Println("track " + hash)
		tracked := t.TrackTx(hash, provider)
		time.Sleep(time.Duration(interval) * time.Second)
		if tracked {
			fmt.Println("confirmed! " + hash)
			return
		}
	}
	t.Status = Rejected
}

func (t *Transaction) Bytes() ([]byte, error) {
	txParams := t.toTransactionParam()
	bytes, err := EncodeTransactionProto(txParams)

	if err != nil {
		return nil, err
	} else {
		return bytes, nil
	}
}

func (t *Transaction) isPending() bool {
	return t.Status == Pending
}

func (t *Transaction) isInitialised() bool {
	return t.Status == Initialised
}

func (t *Transaction) isConfirmed() bool {
	return t.Status == Confirmed
}

func (t *Transaction) isRejected() bool {
	return t.Status == Rejected
}
