package evm

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ethermint/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	"math/big"
)

type (
	// GenesisState defines the application's genesis state. It contains all the
	// information required and accounts to initialize the blockchain.
	GenesisState struct {
		Accounts []GenesisAccount `json:"accounts"`
	}

	// GenesisAccount defines an account to be initialized in the genesis state.
	GenesisAccount struct {
		Address ethcmn.Address `json:"address"`
		Balance *big.Int        `json:"balance"`
		Code    []byte         `json:"code,omitempty"`
		Storage types.Storage  `json:"storage,omitempty"`
	}
)

func ValidateGenesis(data GenesisState) error {
	for _, acct := range data.Accounts {
		if len(acct.Address.Bytes()) == 0 {
			return fmt.Errorf("Invalid GenesisAccount Error: Missing Address")
		}
		if acct.Balance == nil {
			return fmt.Errorf("Invalid GenesisAccount Error: Missing Balance")
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Accounts: []GenesisAccount{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.Accounts {
		keeper.SetCode(ctx, record.Address, record.Code)
		keeper.CreateGenesisAccount(ctx, record)
	}
	return []abci.ValidatorUpdate{}
}

// TODO: Implement
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{Accounts: nil}
}
