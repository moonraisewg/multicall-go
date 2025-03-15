<br />
<div align="center">
<h3 align="center">Multicall Go</h3>
  <p align="center">
    Multicall used on EVM farmily chains to optimize onchain index experience
    <br />
    <a href="https://github.com/moonraisewg/multicall-go/issues">Report Bug</a>
    ·
    <a href="https://github.com/moonraisewg/multicall-go/issues">Request Feature</a>
  </p>
</div>

----

## Usage

#### Download and install it:

```sh
go get github.com/moonraisewg/multicall-go
```

## Todo

- [x] Multicall 
- [x] Multicall V2
- [x] Struct return type
- [x] Tuple return type
- [ ] Struct as param(s)

#### Canonical example:

```go
package main

import (
	"fmt"
	"github.com/moonraisewg/multicall-go/call"
	"log"
	"math/big"
)

var ContractExamples = map[string]call.Chain{
	"0xB8c77482e45F1F44dE1745F52C74426C631bDD52": call.Ethereum, // map[BNB:chain]
	"0x0000000000000000000000000000000000000000": call.Ethereum, // map[null_address:chain]
}

var Contracts = []string{
	"0xB8c77482e45F1F44dE1745F52C74426C631bDD52", // BNB on Ethereum
	"0xdAC17F958D2ee523a2206206994597C13D831ec7", // USDT on Ethereum
	"0xc0a47dFe034B400B47bDaD5FecDa2621de6c4d95", // UniswapFactory on Ethereum
}

func main() {
	ExampleMultiCallV1()
	ExampleMultiCallV2()
	ExampleMultiCallWithBatch()
}

func ExampleMultiCallV1() {
	for address, chain := range ContractExamples {
		caller := call.NewContractBuilder().
			WithChainConfig(call.DefaultChainConfigs[chain]).
			AddMethod("totalSupply()(uint256)")
		_, result, err := caller.AddCall("ts", address, "totalSupply").Call(nil)
		if err != nil {
			fmt.Printf("Error to call %s contract on %s\n", address, chain)
		} else {
			fmt.Printf("Call %s success with total supply is %d in decimals\n", address, result["ts"][0].(*big.Int))
		}
	}
}

func ExampleMultiCallV2() {
	for address, chain := range ContractExamples {
		caller := call.NewContractBuilder().
			WithChainConfig(call.DefaultChainConfigs[chain]).
			AddMethod("totalSupply()(uint256)")
		result, err := caller.AddCall("ts", address, "totalSupply").FlexibleCall(true)
		if err != nil {
			fmt.Printf("Error to call %s contract on %s\n", address, chain)
		} else {
			success := result["ts"].Success
			if success {
				fmt.Printf("Call %s success with total supply is %d\n", address, result["ts"].ReturnData[0].(*big.Int))
			} else {
				fmt.Printf("Call %s failed\n", address)
			}
		}
	}
}

func ExampleMultiCallWithBatch() {
	caller := call.NewContractBuilder().
		WithChainConfig(call.DefaultChainConfigs[call.Ethereum]).
		AddMethod("totalSupply()(uint256)")
	for _, address := range Contracts {
		caller.AddCall(address, address, "totalSupply")
	}

	results, err := caller.FlexibleCall(false)

	if err != nil {
		log.Fatal(err)
	} else {
		for key, result := range results {
			success := result.Success
			if success {
				fmt.Printf("Call %s success with total supply is %d\n", key, result.ReturnData[0].(*big.Int))
			} else {
				fmt.Printf("Call %s failed\n", key)
			}
		}
	}
}

```

