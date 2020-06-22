package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/incognito-core-libs/jsonrpc"
)

type Provider struct {
	host      string
	rpcClient jsonrpc.RPCClient
}

func NewProvider(host string) *Provider {

	rpcClient := jsonrpc.NewClient(host)
	return &Provider{host: host, rpcClient: rpcClient}
}

func (provider *Provider) GetNetworkId() *jsonrpc.RPCResponse {
	return provider.call("GetNetworkId")
}

func (provider *Provider) GetBlockchainInfo() *jsonrpc.RPCResponse {
	return provider.call("GetBlockchainInfo")
}

func (provider *Provider) GetShardingStructure() *jsonrpc.RPCResponse {
	return provider.call("GetShardingStructure")
}

func (provider *Provider) GetDsBlock(block_number string) *jsonrpc.RPCResponse {
	return provider.call("GetDsBlock", block_number)
}

func (provider *Provider) GetLatestDsBlock() *jsonrpc.RPCResponse {
	return provider.call("GetLatestDsBlock")
}

func (provider *Provider) GetNumDSBlocks() *jsonrpc.RPCResponse {
	return provider.call("GetNumDSBlocks")
}

func (provider *Provider) GetDSBlockRate() *jsonrpc.RPCResponse {
	return provider.call("GetDSBlockRate")
}

func (provider *Provider) DSBlockListing(ds_block_listing int) *jsonrpc.RPCResponse {
	return provider.call("DSBlockListing", ds_block_listing)
}

func (provider *Provider) GetTxBlock(tx_block string) *jsonrpc.RPCResponse {
	return provider.call("GetTxBlock", tx_block)
}

func (provider *Provider) GetLatestTxBlock() *jsonrpc.RPCResponse {
	return provider.call("GetLatestTxBlock")
}

func (provider *Provider) GetNumTxBlocks() *jsonrpc.RPCResponse {
	return provider.call("GetNumTxBlocks")
}

func (provider *Provider) GetTxBlockRate() *jsonrpc.RPCResponse {
	return provider.call("GetTxBlockRate")
}

func (provider *Provider) TxBlockListing(page int) *jsonrpc.RPCResponse {
	return provider.call("TxBlockListing", page)
}

func (provider *Provider) GetNumTransactions() *jsonrpc.RPCResponse {
	return provider.call("GetNumTransactions")
}

func (provider *Provider) GetTransactionRate() *jsonrpc.RPCResponse {
	return provider.call("GetTransactionRate")
}

func (provider *Provider) GetCurrentMiniEpoch() *jsonrpc.RPCResponse {
	return provider.call("GetCurrentMiniEpoch")
}

func (provider *Provider) GetCurrentDSEpoch() *jsonrpc.RPCResponse {
	return provider.call("GetCurrentDSEpoch")
}

func (provider *Provider) GetPrevDifficulty() *jsonrpc.RPCResponse {
	return provider.call("GetPrevDifficulty")
}

func (provider *Provider) GetPrevDSDifficulty() *jsonrpc.RPCResponse {
	return provider.call("GetPrevDSDifficulty")
}

func (provider *Provider) CreateTransaction(payload TransactionPayload) *jsonrpc.RPCResponse {
	//r, _ := json.Marshal(payload)
	//fmt.Println(string(r))
	fmt.Println(payload)
	return provider.call("CreateTransaction", &payload)
}

func (provider *Provider) GetTransaction(transaction_hash string) *jsonrpc.RPCResponse {
	return provider.call("GetTransaction", transaction_hash)
}

func (provider *Provider) GetRecentTransactions() *jsonrpc.RPCResponse {
	return provider.call("GetRecentTransactions")
}

func (provider *Provider) GetTransactionsForTxBlock(tx_block_number string) *jsonrpc.RPCResponse {
	return provider.call("GetTransactionsForTxBlock", tx_block_number)
}

func (provider *Provider) GetNumTxnsTxEpoch() *jsonrpc.RPCResponse {
	return provider.call("GetNumTxnsTxEpoch")
}

func (provider *Provider) GetNumTxnsDSEpoch() *jsonrpc.RPCResponse {
	return provider.call("GetNumTxnsDSEpoch")
}

func (provider *Provider) GetMinimumGasPrice() *jsonrpc.RPCResponse {
	return provider.call("GetMinimumGasPrice")
}

func (provider *Provider) GetSmartContractCode(contract_address string) *jsonrpc.RPCResponse {
	return provider.call("GetSmartContractCode", contract_address)
}

func (provider *Provider) GetSmartContractInit(contract_address string) *jsonrpc.RPCResponse {
	return provider.call("GetSmartContractInit", contract_address)
}

func (provider *Provider) GetSmartContractState(contract_address string) *jsonrpc.RPCResponse {
	return provider.call("GetSmartContractState", contract_address)
}

func (provider *Provider) GetSmartContractSubState(contractAddress string, params ...interface{}) (string, error) {
	//we should hack here for now
	type req struct {
		Id      string      `json:"id"`
		Jsonrpc string      `json:"jsonrpc"`
		Method  string      `json:"method"`
		Params  interface{} `json:"params"`
	}

	p := []interface{}{
		contractAddress,
	}

	for _, v := range params {
		p = append(p, v)
	}

	r := &req{
		Id:      "1",
		Jsonrpc: "2.0",
		Method:  "GetSmartContractSubState",
		Params:  p,
	}

	b, _ := json.Marshal(r)
	reader := bytes.NewReader(b)
	request, err := http.NewRequest("POST", provider.host, reader)

	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	return string(result), nil

}

func (provider *Provider) GetSmartContracts(user_address string) *jsonrpc.RPCResponse {
	return provider.call("GetSmartContracts", user_address)
}

func (provider *Provider) GetContractAddressFromTransactionID(transaction_id string) *jsonrpc.RPCResponse {
	return provider.call("GetContractAddressFromTransactionID", transaction_id)
}

func (provider *Provider) GetBalance(user_address string) *jsonrpc.RPCResponse {
	return provider.call("GetBalance", user_address)
}

func (provider *Provider) call(method_name string, params ...interface{}) *jsonrpc.RPCResponse {
	response, err := provider.rpcClient.Call(method_name, params)
	if err != nil {
		return nil
	} else {
		return response
	}
}
