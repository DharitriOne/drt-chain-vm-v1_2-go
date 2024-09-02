package scenarioexec

import (
	"fmt"
	"sort"
	"strings"

	"github.com/DharitriOne/drt-chain-core-go/core"
	worldmock "github.com/DharitriOne/drt-chain-vm-v1_2-go/mock/world"
	er "github.com/DharitriOne/drt-chain-vm-v1_2-go/scenarios/expression/reconstructor"
	mj "github.com/DharitriOne/drt-chain-vm-v1_2-go/scenarios/json/model"
	mjwrite "github.com/DharitriOne/drt-chain-vm-v1_2-go/scenarios/json/write"
	oj "github.com/DharitriOne/drt-chain-vm-v1_2-go/scenarios/orderedjson"
)

const includeProtectedStorage = false

func (ae *VMTestExecutor) convertMockAccountToScenarioFormat(account *worldmock.Account) (*mj.Account, error) {
	var storageKeys []string
	for storageKey := range account.Storage {
		storageKeys = append(storageKeys, storageKey)
	}

	sort.Strings(storageKeys)
	var storageKvps []*mj.StorageKeyValuePair
	for _, storageKey := range storageKeys {
		storageValue := account.Storage[storageKey]
		includeKey := includeProtectedStorage || !strings.HasPrefix(storageKey, core.ProtectedKeyPrefix)
		if includeKey && len(storageValue) > 0 {
			storageKvps = append(storageKvps, &mj.StorageKeyValuePair{
				Key: mj.JSONBytesFromString{
					Value:    []byte(storageKey),
					Original: ae.exprReconstructor.Reconstruct([]byte(storageKey), er.NoHint),
				},
				Value: mj.JSONBytesFromTree{
					Value:    storageValue,
					Original: &oj.OJsonString{Value: ae.exprReconstructor.Reconstruct(storageValue, er.NoHint)},
				},
			})
		}
	}

	tokenData, err := account.GetFullMockDCDTData()
	if err != nil {
		return nil, err
	}
	var dcdtNames []string
	for dcdtName := range tokenData {
		dcdtNames = append(dcdtNames, dcdtName)
	}
	sort.Strings(dcdtNames)
	var scenDCDT []*mj.DCDTData
	for _, dcdtName := range dcdtNames {
		dcdtObj := tokenData[dcdtName]

		var scenRoles []string
		for _, mockRoles := range dcdtObj.Roles {
			scenRoles = append(scenRoles, string(mockRoles))
		}

		var scenInstances []*mj.DCDTInstance
		for _, mockInstance := range dcdtObj.Instances {
			var uri mj.JSONBytesFromTree
			if len(mockInstance.TokenMetaData.URIs) > 0 {
				uri = mj.JSONBytesFromTree{
					Value:    mockInstance.TokenMetaData.URIs[0],
					Original: &oj.OJsonString{Value: ae.exprReconstructor.Reconstruct(mockInstance.TokenMetaData.URIs[0], er.NoHint)},
				}
			}

			scenInstances = append(scenInstances, &mj.DCDTInstance{
				Nonce: mj.JSONUint64{
					Value:    mockInstance.TokenMetaData.Nonce,
					Original: ae.exprReconstructor.ReconstructFromUint64(mockInstance.TokenMetaData.Nonce),
				},
				Balance: mj.JSONBigInt{
					Value:    mockInstance.Value,
					Original: ae.exprReconstructor.ReconstructFromBigInt(mockInstance.Value),
				},
				Uri: uri,
			})
		}

		scenDCDT = append(scenDCDT, &mj.DCDTData{
			TokenIdentifier: mj.JSONBytesFromString{
				Value:    dcdtObj.TokenIdentifier,
				Original: ae.exprReconstructor.Reconstruct(dcdtObj.TokenIdentifier, er.StrHint),
			},
			Instances: scenInstances,
			LastNonce: mj.JSONUint64{
				Value:    dcdtObj.LastNonce,
				Original: ae.exprReconstructor.ReconstructFromUint64(dcdtObj.LastNonce),
			},
			Roles: scenRoles,
		})
	}

	return &mj.Account{
		Address: mj.JSONBytesFromString{
			Value:    account.Address,
			Original: ae.exprReconstructor.Reconstruct([]byte(account.Address), er.AddressHint),
		},
		Nonce: mj.JSONUint64{
			Value:    account.Nonce,
			Original: ae.exprReconstructor.ReconstructFromUint64(account.Nonce),
		},
		Balance: mj.JSONBigInt{
			Value:    account.Balance,
			Original: ae.exprReconstructor.ReconstructFromBigInt(account.Balance),
		},
		Storage:  storageKvps,
		DCDTData: scenDCDT,
	}, nil
}

// DumpWorld prints the state of the MockWorld to stdout.
func (ae *VMTestExecutor) DumpWorld() error {
	fmt.Print("world state dump:\n")
	var scenAccounts []*mj.Account

	for _, account := range ae.World.AcctMap {
		scenAccount, err := ae.convertMockAccountToScenarioFormat(account)
		if err != nil {
			return err
		}
		scenAccounts = append(scenAccounts, scenAccount)
	}

	ojAccount := mjwrite.AccountsToOJ(scenAccounts)
	s := oj.JSONString(ojAccount)
	fmt.Println(s)

	return nil
}
