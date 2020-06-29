package GameFight

import "fmt"

type FightObjList struct {
	mOwnerGuid  X_GUID
	mObjectList [MaxMatrixCellCount]*FightObject
}

func (list *FightObjList) CleanUp() {
	list.mOwnerGuid = InvalidId
	for i := 0; i < MaxMatrixCellCount; i++ {
		list.mObjectList[i].CleanUp()
	}
}

func (list *FightObjList) SetOwnerGuid(guid X_GUID) {
	list.mOwnerGuid = guid
}
func (list *FightObjList) GetOwnerGuid() X_GUID {
	return list.mOwnerGuid
}

func (list *FightObjList) GetActiveCount() int {
	nCount := 0
	for i := 0; i < MaxMatrixCellCount; i++ {
		if list.mObjectList[i].IsActive() {
			nCount += 1
		}
	}
	return nCount
}

func (list *FightObjList) GetInactiveCount() int {
	nCount := 0
	for i := 0; i < MaxMatrixCellCount; i++ {
		if list.mObjectList[i].IsValid() && !list.mObjectList[i].IsActive() {
			nCount += 1
		}
	}
	return nCount
}

func (list *FightObjList) GetFightObject(index int) *FightObject {
	if index >= 0 && index < MaxMatrixCellCount {
		return list.mObjectList[index]
	}
	return nil
}

func (list *FightObjList) FillObject(index int, object *FightObject) bool {
	if !object.IsValid() {
		return false
	}
	if index < 0 || index >= MaxMatrixCellCount {
		return false
	}
	if list.mObjectList[index] != nil && list.mObjectList[index].IsValid() {
		return false
	}
	object.SetMatrixID(index)
	list.mObjectList[index] = object
	return true
}

func (list *FightObjList) Init(fightDB *FightDB) bool {
	if !fightDB.IsValid() {
		return false
	}
	list.SetOwnerGuid(fightDB.HumanGuid)
	for i := 0; i < fightDB.FightCount; i++ {
		fightDBData := fightDB.FightDBData[i]
		if fightDBData.IsValid() {
			fightObj := new(FightObject)
			fightObj.InitFightDBData(fightDBData)
			index := fightObj.GetMatrixID()
			list.mObjectList[index] = fightObj
		}
	}
	return true
}

func (list *FightObjList) ImpactHeartBeat(uTime int) {
	for i := 0; i < MaxMatrixCellCount; i++ {
		if list.mObjectList[i].IsActive() {
			//清空英雄的impact
			list.mObjectList[i].ClearImpactEffect()
			//impact攻击逻辑
			list.mObjectList[i].ImpactHeartBeat(uTime)
		}
	}
}
//攻击
func (list *FightObjList) HeartBeat(uTime int) {
	for i := 0; i < MaxMatrixCellCount; i++ {
		if list.mObjectList[i].IsActive() {
			fmt.Printf("第%+v回合, 英雄id %+v \n",uTime, list.mObjectList[i].GetGuid())
			list.mObjectList[i].HeartBeat(uTime)
			pAttackInfo := list.mObjectList[i].GetAttackInfo()
			if pAttackInfo != nil && pAttackInfo.IsValid() {
				pFightCell := list.mObjectList[i].GetFightCell()
				pRoundInfo := pFightCell.GetRoundInfo()
				pRoundInfo.AddAttackInfo(*pAttackInfo)

			}
		}
	}
}
//清除所有buff
func (list *FightObjList) ClearImpactEffect() {
	for i := 0; i < MaxMatrixCellCount; i++ {
		if list.mObjectList[i].IsActive() {
			list.mObjectList[i].ClearImpactEffect()
		}
	}
}
