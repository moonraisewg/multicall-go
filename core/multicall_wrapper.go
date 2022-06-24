package core

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
)

type Call struct {
	Name     string         `json:"name"`
	Method   string         `json:"method"`
	Target   common.Address `json:"target"`
	CallData []byte         `json:"call_data"`
}

type CallResponse struct {
	Method     string
	ReturnData []byte `json:"returnData"`
}

func (call Call) GetMultiCall() MultiCall {
	return MultiCall{Target: call.Target, CallData: call.CallData}
}

type Caller struct {
	Client          *ethclient.Client
	Abi             abi.ABI
	ContractAddress common.Address
}

func NewCaller(client *ethclient.Client, contractAddress common.Address) (*Caller, error) {
	mcAbi, err := abi.JSON(strings.NewReader(MultiMetaData.ABI))
	if err != nil {
		return nil, err
	}

	return &Caller{
		Client:          client,
		Abi:             mcAbi,
		ContractAddress: contractAddress,
	}, nil
}

func (caller *Caller) Execute(calls []Call, blockNumber *big.Int) (*big.Int, map[string]CallResponse, error) {
	var multiCalls = make([]MultiCall, 0, len(calls))
	for _, call := range calls {
		multiCalls = append(multiCalls, call.GetMultiCall())
	}
	callData, err := caller.Abi.Pack("aggregate", multiCalls)
	if err != nil {
		return nil, nil, err
	}
	resp, err := caller.Client.CallContract(context.Background(), ethereum.CallMsg{To: &caller.ContractAddress, Data: callData}, blockNumber)
	if err != nil {
		return nil, nil, err
	}

	responses, err := caller.Abi.Unpack("aggregate", resp)

	if err != nil {
		return nil, nil, err
	}

	results := make(map[string]CallResponse)
	for i, response := range responses[1].([][]byte) {
		results[calls[i].Name] = CallResponse{
			Method:     calls[i].Method,
			ReturnData: response,
		}
	}
	return responses[0].(*big.Int), results, nil
}
