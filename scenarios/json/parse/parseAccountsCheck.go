package scenjsonparse

import (
	"errors"
	"fmt"

	mj "github.com/DharitriOne/drt-chain-vm-v1_2-go/scenarios/json/model"
	oj "github.com/DharitriOne/drt-chain-vm-v1_2-go/scenarios/orderedjson"
)

func (p *Parser) processCheckAccount(acctRaw oj.OJsonObject) (*mj.CheckAccount, error) {
	acctMap, isMap := acctRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled account object is not a map")
	}

	acct := mj.CheckAccount{
		Nonce:         mj.JSONCheckUint64Unspecified(),
		Balance:       mj.JSONCheckBigIntUnspecified(),
		Username:      mj.JSONCheckBytesUnspecified(),
		IgnoreStorage: true,
		Code:          mj.JSONCheckBytesUnspecified(),
		Owner:         mj.JSONCheckBytesUnspecified(),
		AsyncCallData: mj.JSONCheckBytesUnspecified(),
	}
	var err error

	for _, kvp := range acctMap.OrderedKV {
		switch kvp.Key {
		case "comment":
			acct.Comment, err = p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid check account comment: %w", err)
			}
		case "nonce":
			acct.Nonce, err = p.processCheckUint64(kvp.Value)
			if err != nil {
				return nil, errors.New("invalid account nonce")
			}
		case "balance":
			acct.Balance, err = p.processCheckBigInt(kvp.Value, bigIntUnsignedBytes)
			if err != nil {
				return nil, errors.New("invalid account balance")
			}
		case "dcdt":
			acct.IgnoreDCDT = IsStar(kvp.Value)
			if !acct.IgnoreDCDT {
				dcdtMap, dcdtOk := kvp.Value.(*oj.OJsonMap)
				if !dcdtOk {
					return nil, errors.New("invalid DCDT map")
				}
				for _, dcdtKvp := range dcdtMap.OrderedKV {
					tokenNameStr, err := p.ExprInterpreter.InterpretString(dcdtKvp.Key)
					if err != nil {
						return nil, fmt.Errorf("invalid dcdt token identifer: %w", err)
					}
					tokenName := mj.NewJSONBytesFromString(tokenNameStr, dcdtKvp.Key)
					dcdtItem, err := p.processCheckDCDTData(tokenName, dcdtKvp.Value)
					if err != nil {
						return nil, fmt.Errorf("invalid dcdt value: %w", err)
					}
					acct.CheckDCDTData = append(acct.CheckDCDTData, dcdtItem)
				}
			}
		case "username":
			acct.Username, err = p.parseCheckBytes(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid account username: %w", err)
			}
		case "storage":
			acct.IgnoreStorage = IsStar(kvp.Value)
			if !acct.IgnoreStorage {
				// TODO: convert to a more permissive format
				storageMap, storageOk := kvp.Value.(*oj.OJsonMap)
				if !storageOk {
					return nil, errors.New("invalid account storage")
				}
				for _, storageKvp := range storageMap.OrderedKV {
					byteKey, err := p.ExprInterpreter.InterpretString(storageKvp.Key)
					if err != nil {
						return nil, fmt.Errorf("invalid account storage key: %w", err)
					}
					byteVal, err := p.processSubTreeAsByteArray(storageKvp.Value)
					if err != nil {
						return nil, fmt.Errorf("invalid account storage value: %w", err)
					}
					stElem := mj.StorageKeyValuePair{
						Key:   mj.NewJSONBytesFromString(byteKey, storageKvp.Key),
						Value: byteVal,
					}
					acct.CheckStorage = append(acct.CheckStorage, &stElem)
				}
			}
		case "code":
			acct.Code, err = p.parseCheckBytes(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid account code: %w", err)
			}
		case "owner":
			acct.Owner, err = p.parseCheckBytes(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid account owner: %w", err)
			}
		case "asyncCallData":
			acct.AsyncCallData, err = p.parseCheckBytes(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid asyncCallData: %w", err)
			}

		default:
			return nil, fmt.Errorf("unknown account field: %s", kvp.Key)
		}
	}

	return &acct, nil
}

func (p *Parser) processCheckAccountMap(acctMapRaw oj.OJsonObject) (*mj.CheckAccounts, error) {
	var checkAccounts = &mj.CheckAccounts{
		OtherAccountsAllowed: false,
		Accounts:             nil,
	}

	preMap, isPreMap := acctMapRaw.(*oj.OJsonMap)
	if !isPreMap {
		return nil, errors.New("unmarshalled check account map object is not a map")
	}
	for _, acctKVP := range preMap.OrderedKV {
		if acctKVP.Key == "+" {
			checkAccounts.OtherAccountsAllowed = true
		} else {
			acct, acctErr := p.processCheckAccount(acctKVP.Value)
			if acctErr != nil {
				return nil, acctErr
			}
			acctAddr, hexErr := p.parseAccountAddress(acctKVP.Key)
			if hexErr != nil {
				return nil, hexErr
			}
			acct.Address = acctAddr
			checkAccounts.Accounts = append(checkAccounts.Accounts, acct)
		}
	}
	return checkAccounts, nil
}
