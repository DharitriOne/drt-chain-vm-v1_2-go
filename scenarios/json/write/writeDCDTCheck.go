package scenjsonwrite

import (
	mj "github.com/DharitriOne/drt-chain-vm-v1_2-go/scenarios/json/model"
	oj "github.com/DharitriOne/drt-chain-vm-v1_2-go/scenarios/orderedjson"
)

func checkDCDTDataToOJ(dcdtItems []*mj.CheckDCDTData) *oj.OJsonMap {
	dcdtItemsOJ := oj.NewMap()
	for _, dcdtItem := range dcdtItems {
		dcdtItemsOJ.Put(dcdtItem.TokenIdentifier.Original, checkDCDTItemToOJ(dcdtItem))
	}
	return dcdtItemsOJ
}

func checkDCDTItemToOJ(dcdtItem *mj.CheckDCDTData) oj.OJsonObject {
	if isCompactCheckDCDT(dcdtItem) {
		return checkBigIntToOJ(dcdtItem.Instances[0].Balance)
	}

	dcdtItemOJ := oj.NewMap()

	// instances
	if len(dcdtItem.Instances) == 1 {
		appendCheckDCDTInstanceToOJ(dcdtItem.Instances[0], dcdtItemOJ)
	} else {
		var convertedList []oj.OJsonObject
		for _, dcdtInstance := range dcdtItem.Instances {
			dcdtInstanceOJ := oj.NewMap()
			appendCheckDCDTInstanceToOJ(dcdtInstance, dcdtInstanceOJ)
			convertedList = append(convertedList, dcdtInstanceOJ)
		}
		instancesOJList := oj.OJsonList(convertedList)
		dcdtItemOJ.Put("instances", &instancesOJList)
	}

	if len(dcdtItem.LastNonce.Original) > 0 {
		dcdtItemOJ.Put("lastNonce", checkUint64ToOJ(dcdtItem.LastNonce))
	}

	// roles
	if len(dcdtItem.Roles) > 0 {
		var convertedList []oj.OJsonObject
		for _, roleStr := range dcdtItem.Roles {
			convertedList = append(convertedList, &oj.OJsonString{Value: roleStr})
		}
		rolesOJList := oj.OJsonList(convertedList)
		dcdtItemOJ.Put("roles", &rolesOJList)
	}
	if len(dcdtItem.Frozen.Original) > 0 {
		dcdtItemOJ.Put("frozen", checkUint64ToOJ(dcdtItem.Frozen))
	}

	return dcdtItemOJ
}

func appendCheckDCDTInstanceToOJ(dcdtInstance *mj.CheckDCDTInstance, targetOj *oj.OJsonMap) {
	if len(dcdtInstance.Nonce.Original) > 0 {
		targetOj.Put("nonce", checkUint64ToOJ(dcdtInstance.Nonce))
	}
	if len(dcdtInstance.Balance.Original) > 0 {
		targetOj.Put("balance", checkBigIntToOJ(dcdtInstance.Balance))
	}
}

func isCompactCheckDCDT(dcdtItem *mj.CheckDCDTData) bool {
	if len(dcdtItem.Instances) != 1 {
		return false
	}
	if len(dcdtItem.Instances[0].Nonce.Original) > 0 {
		return false
	}
	if len(dcdtItem.Roles) > 0 {
		return false
	}
	if len(dcdtItem.Frozen.Original) > 0 {
		return false
	}
	return true
}
