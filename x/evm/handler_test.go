package evm

import (
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/ethermint/crypto"
	eminttypes "github.com/cosmos/ethermint/types"
	"github.com/cosmos/ethermint/utils"
	"github.com/cosmos/ethermint/x/evm/types"
	//ethcmn "github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
)

// pragma solidity ^0.5.1;

// contract Test {
//     event Hello(uint256 indexed world);

//     constructor() public {
//         emit Hello(17);
//     }
// }

// {
// 	"linkReferences": {},
// 	"object": "6080604052348015600f57600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a2603580604b6000396000f3fe6080604052600080fdfea165627a7a723058206cab665f0f557620554bb45adf266708d2bd349b8a4314bdff205ee8440e3c240029",
// 	"opcodes": "PUSH1 0x80 PUSH1 0x40 MSTORE CALLVALUE DUP1 ISZERO PUSH1 0xF JUMPI PUSH1 0x0 DUP1 REVERT JUMPDEST POP PUSH1 0x11 PUSH32 0x775A94827B8FD9B519D36CD827093C664F93347070A554F65E4A6F56CD738898 PUSH1 0x40 MLOAD PUSH1 0x40 MLOAD DUP1 SWAP2 SUB SWAP1 LOG2 PUSH1 0x35 DUP1 PUSH1 0x4B PUSH1 0x0 CODECOPY PUSH1 0x0 RETURN INVALID PUSH1 0x80 PUSH1 0x40 MSTORE PUSH1 0x0 DUP1 REVERT INVALID LOG1 PUSH6 0x627A7A723058 KECCAK256 PUSH13 0xAB665F0F557620554BB45ADF26 PUSH8 0x8D2BD349B8A4314 0xbd SELFDESTRUCT KECCAK256 0x5e 0xe8 DIFFICULTY 0xe EXTCODECOPY 0x24 STOP 0x29 ",
// 	"sourceMap": "25:119:0:-;;;90:52;8:9:-1;5:2;;;30:1;27;20:12;5:2;90:52:0;132:2;126:9;;;;;;;;;;25:119;;;;;;"
// }

func TestHandler_Logs(t *testing.T) {
	// create logger, codec and root multi-store
	cdc := newTestCodec()

	// The ParamsKeeper handles parameter storage for the application
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	paramsKeeper := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	// Set specific supspaces
	authSubspace := paramsKeeper.Subspace(auth.DefaultParamspace)
	ak := auth.NewAccountKeeper(cdc, accKey, authSubspace, eminttypes.ProtoBaseAccount)
	ek := NewKeeper(ak, storageKey, codeKey, blockKey, cdc)

	gasLimit := uint64(100000)
	gasPrice := big.NewInt(1000000)

	priv1, _ := crypto.GenerateKey()

	bytecode, err := utils.HexToBytes("0x6080604052348015600f57600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a2603580604b6000396000f3fe6080604052600080fdfea165627a7a723058206cab665f0f557620554bb45adf266708d2bd349b8a4314bdff205ee8440e3c240029")
	if err != nil {
		t.Fatal(err)
	}

	tx := types.NewEthereumTxMsg(1, nil, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(1), priv1.ToECDSA())

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	// mount stores
	keys := []*sdk.KVStoreKey{accKey, storageKey, codeKey, blockKey}
	for _, key := range keys {
		cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, nil)
	}

	err = cms.LoadLatestVersion()
	if err != nil {
		t.Fatal(err)
	}

	ms := cms.CacheMultiStore()
	ctx := sdk.NewContext(ms, abci.Header{}, false, logger)
	ctx = ctx.WithBlockHeight(1).WithChainID("1")

	result := handleETHTxMsg(ctx, ek, tx)
	t.Log(result.Data)
	t.Log(ek.bloom)
	t.Log(ek.logs)

	// TODO: assert that we can set/get logs from the keeper
}
