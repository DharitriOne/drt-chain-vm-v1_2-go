package scenarioexec

import (
	"errors"
	"math/big"

	"github.com/DharitriOne/drt-chain-core-go/core"
	"github.com/DharitriOne/drt-chain-core-go/data/dcdt"
	vmcommon "github.com/DharitriOne/drt-chain-vm-common-go"
	"github.com/DharitriOne/drt-chain-vm-common-go/builtInFunctions"
	worldmock "github.com/DharitriOne/drt-chain-vm-v1_2-go/mock/world"
	mj "github.com/DharitriOne/drt-chain-vm-v1_2-go/scenarios/json/model"
)

func convertAccount(testAcct *mj.Account) (*worldmock.Account, error) {
	storage := make(map[string][]byte)
	for _, stkvp := range testAcct.Storage {
		key := string(stkvp.Key.Value)
		storage[key] = stkvp.Value.Value
	}

	if len(testAcct.Address.Value) != 32 {
		return nil, errors.New("bad test: account address should be 32 bytes long")
	}

	account := &worldmock.Account{
		Address:         testAcct.Address.Value,
		Nonce:           testAcct.Nonce.Value,
		Balance:         big.NewInt(0).Set(testAcct.Balance.Value),
		BalanceDelta:    big.NewInt(0),
		DeveloperReward: big.NewInt(0),
		Username:        testAcct.Username.Value,
		Storage:         storage,
		Code:            testAcct.Code.Value,
		OwnerAddress:    testAcct.Owner.Value,
		AsyncCallData:   testAcct.AsyncCallData,
		ShardID:         uint32(testAcct.Shard.Value),
		IsSmartContract: len(testAcct.Code.Value) > 0,
		CodeMetadata: (&vmcommon.CodeMetadata{
			Payable:     true,
			Upgradeable: true,
			Readable:    true,
		}).ToBytes(), // TODO: add explicit fields in scenario JSON
	}

	for _, scenDCDTData := range testAcct.DCDTData {
		tokenName := scenDCDTData.TokenIdentifier.Value
		isFrozen := scenDCDTData.Frozen.Value > 0
		for _, instance := range scenDCDTData.Instances {
			tokenNonce := instance.Nonce.Value
			tokenKey := worldmock.MakeTokenKey(tokenName, tokenNonce)
			tokenBalance := instance.Balance.Value
			tokenData := &dcdt.DCDigitalToken{
				Value:      tokenBalance,
				Type:       uint32(core.Fungible),
				Properties: makeDCDTUserMetadataBytes(isFrozen),
				TokenMetaData: &dcdt.MetaData{
					Name:  tokenName,
					Nonce: tokenNonce,
				},
			}
			err := account.SetTokenData(tokenKey, tokenData)
			if err != nil {
				return nil, err
			}
			err = account.SetLastNonce(tokenName, scenDCDTData.LastNonce.Value)
			if err != nil {
				return nil, err
			}
		}
		err := account.SetTokenRolesAsStrings(tokenName, scenDCDTData.Roles)
		if err != nil {
			return nil, err
		}
	}

	return account, nil
}

func makeDCDTUserMetadataBytes(frozen bool) []byte {
	metadata := &builtInFunctions.DCDTUserMetadata{
		Frozen: frozen,
	}

	return metadata.ToBytes()
}

func convertNewAddressMocks(testNAMs []*mj.NewAddressMock) []*worldmock.NewAddressMock {
	var result []*worldmock.NewAddressMock
	for _, testNAM := range testNAMs {
		result = append(result, &worldmock.NewAddressMock{
			CreatorAddress: testNAM.CreatorAddress.Value,
			CreatorNonce:   testNAM.CreatorNonce.Value,
			NewAddress:     testNAM.NewAddress.Value,
		})
	}
	return result
}

func convertBlockInfo(testBlockInfo *mj.BlockInfo) *worldmock.BlockInfo {
	if testBlockInfo == nil {
		return nil
	}

	var randomsSeed [48]byte
	if testBlockInfo.BlockRandomSeed != nil {
		copy(randomsSeed[:], testBlockInfo.BlockRandomSeed.Value)
	}

	result := &worldmock.BlockInfo{
		BlockTimestamp: testBlockInfo.BlockTimestamp.Value,
		BlockNonce:     testBlockInfo.BlockNonce.Value,
		BlockRound:     testBlockInfo.BlockRound.Value,
		BlockEpoch:     uint32(testBlockInfo.BlockEpoch.Value),
		RandomSeed:     &randomsSeed,
	}

	return result
}

// this is a small hack, so we can reuse JSON printing in error messages
func convertLogToTestFormat(outputLog *vmcommon.LogEntry) *mj.LogEntry {
	testLog := mj.LogEntry{
		Address:    mj.JSONCheckBytesReconstructed(outputLog.Address),
		Identifier: mj.JSONCheckBytesReconstructed(outputLog.Identifier),
		Data:       mj.JSONCheckBytesReconstructed(outputLog.GetFirstDataItem()),
		Topics:     make([]mj.JSONCheckBytes, len(outputLog.Topics)),
	}
	for i, topic := range outputLog.Topics {
		testLog.Topics[i] = mj.JSONCheckBytesReconstructed(topic)
	}

	return &testLog
}

func generateTxHash(txIndex string) []byte {
	txIndexBytes := []byte(txIndex)
	if len(txIndexBytes) > 32 {
		return txIndexBytes[:32]
	}
	for i := len(txIndexBytes); i < 32; i++ {
		txIndexBytes = append(txIndexBytes, '.')
	}
	return txIndexBytes
}

func addDCDTToVMInput(dcdtData *mj.DCDTTxData, vmInput *vmcommon.VMInput) {
	if dcdtData != nil {
		vmInput.DCDTTransfers = make([]*vmcommon.DCDTTransfer, 1)
		vmInput.DCDTTransfers[0] = &vmcommon.DCDTTransfer{}
		vmInput.DCDTTransfers[0].DCDTTokenName = dcdtData.TokenIdentifier.Value
		vmInput.DCDTTransfers[0].DCDTValue = dcdtData.Value.Value
		vmInput.DCDTTransfers[0].DCDTTokenNonce = dcdtData.Nonce.Value
		if vmInput.DCDTTransfers[0].DCDTTokenNonce != 0 {
			vmInput.DCDTTransfers[0].DCDTTokenType = uint32(core.NonFungible)
		} else {
			vmInput.DCDTTransfers[0].DCDTTokenType = uint32(core.Fungible)
		}
	}
}
