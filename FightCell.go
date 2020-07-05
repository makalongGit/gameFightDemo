package GameFight

import "fmt"

type FightCell struct {
	mAttackList  *FightObjList   //攻击方
	mDefenceList *FightObjList   //防御方
	mFightInfo   *FightInfo      //战斗信息提取
	mRoundInfo   *FightRoundInfo //回合信息

	mFightType EmTypeFight //战斗类型
	mPlusAtt   int         //攻加成
}

func (f *FightCell) CleanUp() {
	f.mAttackList.CleanUp()
	f.mDefenceList.CleanUp()
	f.mFightInfo.CleanUp()
	f.mRoundInfo.CleanUp()
	f.mFightType = EmTypeFightNormal
	f.mPlusAtt = 0
}

func (f *FightCell) GetFightType() EmTypeFight {
	return f.mFightType
}

func (f *FightCell) SetFightType(val EmTypeFight) {
	f.mFightType = val
}

func (f *FightCell) GetPlusAtt() int {
	return f.mPlusAtt
}

func (f *FightCell) SetPlusAtt(val int) {
	f.mPlusAtt = val
}

func (f *FightCell) GetRoundInfo() *FightRoundInfo {
	return f.mRoundInfo
}

//攻击方
func (f *FightCell) GetAttackList() *FightObjList {
	return f.mAttackList
}

//防守方
func (f *FightCell) GetDefenceList() *FightObjList {
	return f.mDefenceList
}

func (f *FightCell) IsOver() bool {
	if f.mAttackList.GetActiveCount() <= 0 {
		return true
	}
	if f.mDefenceList.GetActiveCount() <= 0 {
		return true
	}
	return false
}

func (f *FightCell) IsWin() bool {
	return f.mFightInfo.BWin
}

//初始化战斗信息
func (f *FightCell) initFightInfo() {
	//读表
	f.mFightInfo.MaxFightDistance = 10
	for i := 0; i < MaxMatrixCellCount; i++ {
		pFightObj := f.mAttackList.GetFightObject(i)
		if pFightObj != nil && pFightObj.IsValid() {
			objdata := FObjectData{}
			objdata.Guid = pFightObj.GetGuid()
			objdata.TableID = pFightObj.GetTableID()
			objdata.Quality = pFightObj.GetQuality()
			objdata.Color = pFightObj.GetColor()
			objdata.Profession = pFightObj.GetProfession()
			objdata.Level = pFightObj.GetLevel()
			objdata.MatrixID = pFightObj.GetMatrixID()
			f.mFightInfo.AddAttackObjectData(objdata)
		}

		pFightObj = f.mDefenceList.GetFightObject(i)
		if pFightObj != nil && pFightObj.IsValid() {
			objdata := FObjectData{}
			objdata.Guid = pFightObj.GetGuid()
			objdata.TableID = pFightObj.GetTableID()
			objdata.Quality = pFightObj.GetQuality()
			objdata.Color = pFightObj.GetColor()
			objdata.Profession = pFightObj.GetProfession()
			objdata.Level = pFightObj.GetLevel()
			objdata.MatrixID = pFightObj.GetMatrixID()
			f.mFightInfo.AddDefendObjectData(objdata)
		}
	}
}

//初始化每回合双方的初始数据
func (f *FightCell) initRoundInfo() {
	f.mRoundInfo.CleanUp()
	for i := 0; i < MaxMatrixCellCount; i++ {
		pFightObj := f.mAttackList.GetFightObject(i)
		if pFightObj != nil && pFightObj.IsValid() {
			objInfo := new(FObjectInfo)
			objInfo.Guid = pFightObj.GetGuid()
			objInfo.HP = pFightObj.GetHP()
			objInfo.MaxHP = pFightObj.GetMaxHP()
			objInfo.MP = pFightObj.GetMP()
			objInfo.MaxMP = pFightObj.GetMaxMP()
			objInfo.FightDistance = pFightObj.GetFightDistance()
			objInfo.AttackSpeed = pFightObj.GetAttackSpeed()
			pImpactList := pFightObj.GetImpactList()
			//每个英雄的impact列表
			for index := 0; index < MaxImpactNumber; index++ {
				if pImpactList[index].IsValid() {
					objInfo.AddImpact(pImpactList[index].GetImpactID(), 0, 0)
				}
			}
			f.mRoundInfo.AddAttackObjectInfo(objInfo)
		}

		pFightObj = f.mDefenceList.GetFightObject(i)
		if pFightObj != nil && pFightObj.IsValid() {
			objInfo := new(FObjectInfo)
			objInfo.Guid = pFightObj.GetGuid()
			objInfo.HP = pFightObj.GetHP()
			objInfo.MaxHP = pFightObj.GetMaxHP()
			objInfo.MP = pFightObj.GetMP()
			objInfo.MaxMP = pFightObj.GetMaxMP()
			objInfo.FightDistance = pFightObj.GetFightDistance()
			objInfo.AttackSpeed = pFightObj.GetAttackSpeed()
			pImpactList := pFightObj.GetImpactList()
			//每个英雄的impact列表
			for index := 0; index < MaxImpactNumber; index++ {
				if pImpactList[index].IsValid() {
					objInfo.AddImpact(pImpactList[index].GetImpactID(), 0, 0)
				}
			}
			f.mRoundInfo.AddDefendObjectInfo(objInfo)
		}
	}
}

func (f *FightCell) Fight() bool {
	f.initFightInfo()
	for nRound := 1; nRound <= MAX_FIGHT_ROUND; nRound++ {
		if f.IsOver() {
			if f.mAttackList.GetActiveCount() > 0 {
				f.mFightInfo.SetWin(true)
			} else {
				f.mFightInfo.SetWin(false)
			}
			return true
		}
		fmt.Printf("第%+v回合:\n", nRound)
		fmt.Printf("攻击方英雄存活数: %+v \n", f.mAttackList.GetActiveCount())
		fmt.Printf("防守方英雄存活数: %+v \n", f.mDefenceList.GetActiveCount())
		f.initRoundInfo()
		fmt.Println("战前信息",f.mRoundInfo.String())
		f.mAttackList.ImpactHeartBeat(nRound)
		f.mDefenceList.ImpactHeartBeat(nRound)
		f.mAttackList.HeartBeat(nRound)
		f.mDefenceList.HeartBeat(nRound)
		fmt.Println("战后信息",f.mRoundInfo.String())
		f.mFightInfo.AddRoundInfo(*f.mRoundInfo)
		//spew.Dump(f.mRoundInfo.AttackInfo)
	}
	f.mFightInfo.SetWin(false)
	return true
}

func (f *FightCell) InitAttackList() {
	for i := 0; i < MaxMatrixCellCount; i++ {
		pFightObj := f.mAttackList.GetFightObject(i)
		if pFightObj != nil {
			pFightObj.SetAttacker(true)
			pFightObj.SetFightCell(f)
			pFightObj.InitSkill()
		}
	}
}

func (f *FightCell) InitDefendList(nDefendType int) {
	for i := 0; i < MaxMatrixCellCount; i++ {
		pFightObj := f.mDefenceList.GetFightObject(i)
		if pFightObj != nil {
			pFightObj.SetAttacker(false)
			pFightObj.SetFightCell(f)
			pFightObj.InitSkill()
		}
	}
	f.mFightInfo.DefendType = nDefendType
}

func (f *FightCell) GetFightInfo() *FightInfo{
	return f.mFightInfo
}
