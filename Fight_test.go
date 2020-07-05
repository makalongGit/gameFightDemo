package GameFight

import (
	"fmt"
	"testing"
)

func TestFightCell_Fight(t *testing.T) {
	fightCell := FightCell{
		mAttackList:  nil,
		mDefenceList: nil,
		mFightInfo:   new(FightInfo),
		mRoundInfo:   new(FightRoundInfo),
		mFightType:   0,
		mPlusAtt:     0,
	}

	attackObjList := FightObjList{
		mOwnerGuid:  1,
		mObjectList: [6]*FightObject{},
	}
	attackerHeroList := [MaxMatrixCellCount]int{9}
	fightCell.mAttackList = &attackObjList
	defendObjList := FightObjList{
		mOwnerGuid:  1,
		mObjectList: [6]*FightObject{},
	}
	fightCell.mDefenceList = &defendObjList
	defendHeroList := [MaxMatrixCellCount]int{10}
	FillFightObjList(attackerHeroList, fightCell.GetAttackList())
	FillFightObjList(defendHeroList, fightCell.GetDefenceList())
	fightCell.InitAttackList()
	fightCell.InitDefendList(0)
	fightCell.Fight()


}

func FillFightObjList(heroList [MaxMatrixCellCount]int, pFightObjList *FightObjList) bool {
	for i := 0; i < MaxMatrixCellCount; i++ {
		fightObject := new(FightObject)
		if FillFightObjByCfg(heroList[i], fightObject) {
			pFightObjList.FillObject(i, fightObject)
		}
	}
	return true
}

//heroId 英雄id
func FillFightObjByCfg(heroId int, fightObject *FightObject) bool {
	pHeroRow, ok := G_HeroAttr[heroId]
	if !ok {
		panic("cfg error")
	}
	fightDBData := FightDBData{
		Guid:            X_GUID(heroId),
		TableID:         heroId,
		Type:            0,
		Quality:         pHeroRow.InitQuality,
		MatrixID:        0,
		Profession:      pHeroRow.Profession,
		Level:           pHeroRow.InitLevel,
		HP:              pHeroRow.InitHP,
		Mp:              pHeroRow.InitMP,
		SkillCount:      MaxSkillNum,
		Skill:           [2]int{pHeroRow.MagicSkillID1, pHeroRow.MagicSkillID2},
		EquipSkillCount: 0,
		EquipSkill:      [6]int{},
		PhysicAttack:    pHeroRow.InitPhysicAttack,
		MagicAttack:     pHeroRow.InitMagicAttack,
		PhysicDefence:   pHeroRow.InitPhysicDefence,
		MagicDefence:    pHeroRow.InitMagicDefence,
		MaxHP:           pHeroRow.InitHP,
		MaxMP:           pHeroRow.InitMP,
		Hit:             pHeroRow.InitHit,
		Dodge:           pHeroRow.InitDodge,
		Strike:          pHeroRow.InitStrike,
		StrikeHurt:      pHeroRow.InitStrikeHurt,
		Continuous:      pHeroRow.InitContinuous,
		ConAttHurt:      pHeroRow.InitConAttHurt,
		ConAttTimes:     pHeroRow.InitConAttTimes,
		BackAttack:      pHeroRow.InitBackAttack,
		BackAttHurt:     pHeroRow.InitBackAttHurt,
		AttackSpeed:     pHeroRow.InitAttackSpeed,
		PhysicHurtDecay: pHeroRow.InitPhysicHurtDecay,
		MagicHurtDecay:  pHeroRow.InitMagicHurtDecay,
		Exp:             0,
		GrowRate:        0,
		BearPoint:       0,
		Equip:           [6]int{InvalidId, InvalidId, InvalidId, InvalidId, InvalidId, InvalidId},
		Color:           1,
	}
	fightObject.InitFightDBData(fightDBData)
	return true

}

func TestX_GUID_IsValid(t *testing.T) {
	guid := X_GUID(2)
	fmt.Println(guid.IsValid())
}

func TestFuncInArr(t *testing.T){
	arr := [2]func(){func(){ fmt.Println(1)}, func(){ fmt.Println(2)}}
	arr[1]()

}