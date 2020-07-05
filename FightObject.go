package GameFight

import (
	"errors"
	"math/rand"
)

type FightObject struct {
	FightDBData     FightDBData
	bAttacker       bool               //是否为攻击方
	pFightCell      *FightCell         //战斗单元
	skillList       [MaxSkillNum]Skill //主动技能列表
	commonKill      Skill              //普通攻击
	equipSkillCount int
	equipSkillList  [MaxEquipNumPerHero]Skill          //装备技能
	impactList      [MaxImpactNumber]Impact            //buff列表
	impactEffect    [EmAttributeNumber]AttributeEffect //技能影响
	fightDistance   int                                //战斗条长度
	attackInfo      *AttackInfo                        //出手信息
}

func (f *FightObject) CleanUp() {
	f.FightDBData.CleanUp()
	for i := 0; i < MaxSkillNum; i++ {
		f.skillList[i].CleanUp()
	}
	for i := 0; i < MaxImpactNumber; i++ {
		f.impactList[i].CleanUp()
	}
	for i := 0; i < int(EmAttributeNumber); i++ {
		f.impactEffect[i].AttrValue = 0
	}
	f.commonKill.CleanUp()
	f.fightDistance = 0
}

func (f *FightObject) HeartBeat(uTime int) bool {
	//清空上回合的战斗信息
	f.attackInfo.CleanUp()
	//该回合自己的英雄信息
	pFightCell := f.GetFightCell()
	pRoundInfo := pFightCell.GetRoundInfo()
	pFObjectInfo := pRoundInfo.GetFObjectInfoByGuid(f.GetGuid())

	if pFObjectInfo != nil {
		pFObjectInfo.AttackSpeed = f.GetAttackSpeed()
	}

	if !f.IsActive() {
		return true
	}
	if uTime == 1 {
		//被动技能
		f.CastPassiveSkill(uTime)
	}
	pEnemyList := f.GetEnemyList()
	if pEnemyList.GetActiveCount() == 0 {
		return true
	}

	f.fightDistance += f.GetAttackSpeed()
	nDistance := Distance
	if f.fightDistance >= nDistance {
		f.fightDistance = 0

		bRet := f.SkillHeartBeat(uTime)
		//fmt.Println("释放技能", bRet)
		if !bRet {
			bRet = f.CastCommonSkill(uTime)
		}
	}
	//fmt.Printf("英雄: %+v 战斗条:%+v 是否触发攻击:%+v \n", f.GetGuid(), f.fightDistance, f.fightDistance >= nDistance)

	if pFObjectInfo != nil {
		pFObjectInfo.EndDistance = f.fightDistance
		if pFObjectInfo.MaxHP < f.GetMaxHP() {
			pFObjectInfo.MaxHP = f.GetMaxHP()
		}
		if pFObjectInfo.MaxMP < f.GetMaxMP() {
			pFObjectInfo.MaxMP = f.GetMaxMP()
		}
		//fmt.Println(pFObjectInfo.String())
	}
	return true
}

func (f FightObject) GetGuid() X_GUID {
	return f.FightDBData.Guid
}
func (f FightObject) GetMatrixID() int {
	return f.FightDBData.MatrixID
}
func (f *FightObject) SetMatrixID(index int) {
	f.FightDBData.MatrixID = index
}
func (f FightObject) GetTableID() int {
	return f.FightDBData.TableID
}
func (f FightObject) GetQuality() int {
	return f.FightDBData.Quality
}
func (f FightObject) GetColor() int {
	return f.FightDBData.Color
}

func (f FightObject) GetProfession() int {
	return f.FightDBData.Profession
}
func (f FightObject) GetFightDistance() int {
	return f.fightDistance
}
func (f FightObject) GetLevel() int {
	return f.FightDBData.Level
}
func (f *FightObject) InitFightDBData(fightDBData FightDBData) bool {
	f.FightDBData = fightDBData
	return true
}

func (f *FightObject) InitSkill() bool {
	for i := 0; i < f.FightDBData.SkillCount; i++ {
		if f.FightDBData.Skill[i] > 0 {
			bret := f.skillList[i].Init(f.FightDBData.Skill[i], f)
			if bret == false {
				return bret
			}
		}
	}
	f.equipSkillCount = f.FightDBData.EquipSkillCount
	for i := 0; i < f.equipSkillCount; i++ {
		if f.FightDBData.EquipSkill[i] > 0 {
			bret := f.equipSkillList[i].Init(f.FightDBData.EquipSkill[i], f)
			if bret == false {
				return bret
			}
		}
	}
	bret := f.commonKill.Init(f.FightDBData.Profession, f)
	if bret == false {
		return bret
	}
	return true
}

func (f *FightObject) GetMaxHP() int {
	nEffectValue := f.impactEffect[EmAttributeMaxHp].AttrValue
	nEffectValue += int(f.FightDBData.MaxHP * f.impactEffect[EmAttributePercentMaxHp].AttrValue / 100)
	nEndValue := f.FightDBData.MaxHP + nEffectValue
	nEndValue = ChkMin(nEndValue, 0)
	return nEndValue
}

func (f *FightObject) GetMaxMP() int {
	nEffectValue := f.impactEffect[EmAttributeMaxMp].AttrValue
	nEffectValue += f.FightDBData.MaxMP * f.impactEffect[EmAttributeMaxMp].AttrValue
	nEndValue := f.FightDBData.MaxMP + nEffectValue
	nEndValue = ChkMin(nEndValue, 0)
	return nEndValue
}

func (f *FightObject) GetMP() int {
	nEffectValue := f.impactEffect[EmAttributeMp].AttrValue
	nEndValue := f.FightDBData.Mp + nEffectValue
	nMaxMP := f.GetMaxMP()
	if nEndValue > nMaxMP {
		nEndValue = nMaxMP
	}
	nEndValue = ChkMin(nEndValue, 0)
	f.impactEffect[EmAttributeMp].AttrValue = 0
	f.FightDBData.Mp = nEndValue
	return nEndValue
}

//技能条
func (f *FightObject) GetAttackSpeed() int {
	nEffectValue := f.impactEffect[EmAttributeAttackSpeed].AttrValue
	nEffectValue += int(f.FightDBData.AttackSpeed * f.impactEffect[EmAttributePercentAttackSpeed].AttrValue / 100)
	nEndValue := f.FightDBData.AttackSpeed + nEffectValue
	nEndValue = ChkMin(nEndValue, 0)
	return nEndValue
}

func (f *FightObject) SetMP(nMP int) {
	nMaxMP := f.GetMaxMP()
	if nMP > nMaxMP {
		nMP = nMaxMP
	}
	if nMP < 0 {
		nMP = 0
	}
	f.FightDBData.Mp = nMP
}

func (f *FightObject) GetAttackInfo() *AttackInfo {
	if f.attackInfo == nil {
		f.attackInfo = new(AttackInfo)
	}
	return f.attackInfo
}

func (f *FightObject) GetFightCell() *FightCell {
	return f.pFightCell
}

func (f *FightObject) SetFightCell(cell *FightCell) {
	f.pFightCell = cell
}

func (f *FightObject) SetAttacker(bAttacker bool) {
	f.bAttacker = bAttacker
}
func (f *FightObject) IsAttacker() bool {
	return f.bAttacker
}
func (f *FightObject) GetImpactList() [MaxImpactNumber]Impact {
	return f.impactList
}

//物理攻击
func (f *FightObject) GetPhysicAttack() int {

	nEffectValue := f.impactEffect[EmAttributePhysicAttack].AttrValue
	nEffectValue += int(f.FightDBData.PhysicAttack * f.impactEffect[EmAttributePercentPhysicAttack].AttrValue / 100)

	nEndValue := f.FightDBData.PhysicAttack + nEffectValue
	nEndValue = ChkMin(nEndValue, 0)
	return nEndValue
}

//魔法攻击
func (f *FightObject) GetMagicAttack() int {
	nEffectValue := f.impactEffect[EmAttributeMagicAttack].AttrValue
	nEffectValue += int(f.FightDBData.MagicAttack * f.impactEffect[EmAttributePercentMagicAttack].AttrValue / 100)
	nEndValue := f.FightDBData.MagicAttack + nEffectValue
	nEndValue = ChkMin(nEndValue, 0)
	return nEndValue
}

//物理防御
func (f *FightObject) GetPhysicDefend() int {
	nEffectValue := f.impactEffect[EmAttributePhysicDefence].AttrValue
	nEffectValue += int(f.FightDBData.PhysicDefence * f.impactEffect[EmAttributePercentPhysicDefence].AttrValue / 100)
	nEndValue := f.FightDBData.PhysicDefence + nEffectValue
	nEndValue = ChkMin(nEndValue, 0)
	return nEndValue
}

//魔法防御
func (f *FightObject) GetMagicDefend() int {
	nEffectValue := f.impactEffect[EmAttributeMagicDefence].AttrValue
	nEffectValue += int(f.FightDBData.MagicDefence * f.impactEffect[EmAttributePercentMagicDefence].AttrValue / 100)
	nEndValue := f.FightDBData.MagicDefence + nEffectValue
	nEndValue = ChkMin(nEndValue, 0)
	return nEndValue
}

//物理减免
func (f *FightObject) GetPhysicHurtDecay() int {

	nEffectValue := f.impactEffect[EmAttributePhysicHurtDecay].AttrValue
	nEffectValue += int(f.FightDBData.PhysicHurtDecay * f.impactEffect[EmAttributePhysicHurtDecay].AttrValue / 100)
	nEndValue := f.FightDBData.PhysicHurtDecay + nEffectValue
	nEndValue = ChkMin(nEndValue, 0)
	//TODO 读最大物理减免值
	return nEndValue
}

//物理减免
func (f *FightObject) GetMagicHurtDecay() int {

	nEffectValue := f.impactEffect[EmAttributeMagicHurtDecay].AttrValue
	nEffectValue += int(f.FightDBData.MagicHurtDecay * f.impactEffect[EmAttributeMagicHurtDecay].AttrValue / 100)
	nEndValue := f.FightDBData.MagicHurtDecay + nEffectValue
	nEndValue = ChkMin(nEndValue, 0)
	//TODO 读最大魔法减免值
	return nEndValue
}

//命中率
func (f *FightObject) GetHit() int {
	nEffectValue := f.impactEffect[EmAttributeHit].AttrValue
	nEndValue := f.FightDBData.Hit + nEffectValue
	nEndValue = ChkMin(nEndValue, 0)
	return nEndValue
}

//闪避点数
func (f *FightObject) GetDodge() int {
	nEffectValue := f.impactEffect[EmAttributeDodge].AttrValue
	nEndValue := f.FightDBData.Dodge + nEffectValue
	nEndValue = ChkMin(nEndValue, 0)
	return nEndValue
}

//暴击率
func (f *FightObject) GetStrike() int {
	nEffectValue := f.impactEffect[EmAttributeStrike].AttrValue
	nEndValue := f.FightDBData.Strike + nEffectValue
	nEndValue = ChkMin(nEndValue, 0)
	//TODO 读最大暴击率
	return nEndValue
}

//暴击伤害
func (f *FightObject) GetStrikeHurt() int {
	nEffectValue := f.impactEffect[EmAttributeHurtStrike].AttrValue
	nEndValue := f.FightDBData.StrikeHurt + nEffectValue
	nEndValue = ChkMin(nEndValue, 0)
	return nEndValue
}

//连击
func (f *FightObject) GetContinuous() int {
	nEffectValue := f.impactEffect[EmAttributeContinuous].AttrValue
	nEndValue := f.FightDBData.Continuous + nEffectValue
	nEndValue = ChkMin(nEndValue, 0)
	return nEndValue
}

//连击次数
func (f *FightObject) GetConAttTimes() int {
	nEffectValue := f.impactEffect[EmAttributeContinuousTimes].AttrValue
	nEndValue := f.FightDBData.ConAttTimes + nEffectValue
	nEndValue = ChkMin(nEndValue, 0)
	return nEndValue
}

//连击伤害
func (f *FightObject) GetConAttHurt() int {
	nEffectValue := f.impactEffect[EmAttributeHurtContinuous].AttrValue
	nEndValue := f.FightDBData.ConAttHurt + nEffectValue
	nEndValue = ChkMin(nEndValue, 0)
	return nEndValue
}

//反击
func (f *FightObject) GetAttackBack() int {
	nEffectValue := f.impactEffect[EmAttributeBackAttack].AttrValue
	nEndValue := f.FightDBData.BackAttack + nEffectValue
	nEndValue = ChkMin(nEndValue, 0)
	return nEndValue
}

//反击伤害
func (f *FightObject) GetBackAttHurt() int {
	nEffectValue := f.impactEffect[EmAttributeHurtBackAttack].AttrValue
	nEndValue := f.FightDBData.BackAttHurt + nEffectValue
	nEndValue = ChkMin(nEndValue, 0)
	return nEndValue
}

func (f *FightObject) GetHeroAttr() (TableRowHeroAttr, error) {
	id := f.FightDBData.TableID
	if heroAttr, ok := G_HeroAttr[id]; ok {
		return heroAttr, nil
	} else {
		return TableRowHeroAttr{}, errors.New("config err")
	}
}

func (f *FightObject) GetFloatingHurt() int {
	pRowHero, err := f.GetHeroAttr()
	if err != nil {
		panic(err)
	}
	nFloatingHurt := pRowHero.FloatingHurt
	if nFloatingHurt > 0 {
		nFloatingHurt = rand.Intn(nFloatingHurt)
	}
	return nFloatingHurt
}

func (f *FightObject) CoolDownHeartBeat(uTime int) bool {
	return true
}

func (f *FightObject) ClearImpactEffect() {
	for i := 0; i < int(EmAttributeNumber); i++ {
		f.impactEffect[i].AttrValue = 0
	}
}

//impact攻击逻辑
func (f *FightObject) ImpactHeartBeat(uTime int) {
	for i := 0; i < MaxImpactNumber; i++ {
		if f.impactList[i].IsValid() {
			f.impactList[i].HeartBeat(uTime)
		}
	}
}

func (f *FightObject) CompareAttackOrder(pFightObj *FightObject) int {
	matrixID := f.GetMatrixID()
	objMatrixID := pFightObj.GetMatrixID()
	var fAttack, objAttack int
	if f.IsAttacker() {
		fAttack = 1
	} else {
		fAttack = 0
	}
	if pFightObj.IsAttacker() {
		objAttack = 1
	} else {
		objAttack = 0
	}
	nOrder := (1-fAttack)*MaxMatrixCellCount + matrixID
	objOrder := (1-objAttack)*MaxMatrixCellCount + objMatrixID
	return nOrder - objOrder
}

func (f *FightObject) GetImpactLogicType(nImpactID int) EmTypeImpactLogic {
	pRowImpact, err := ImpactTableRow(nImpactID)
	if err != nil {
		panic(err)
	}
	switch pRowImpact.LogicID {
	case EmImpactLogic0, EmImpactLogic1, EmImpactLogic6:
		return EmTypeImpactLogicSingle
	case EmImpactLogic2, EmImpactLogic3, EmImpactLogic5:
		return EmTypeImpactLogicDeBuff
	case EmImpactLogic4:
		return EmTypeImpactLogicBuff
	default:
		return EmTypeImpactLogicInvalid
	}
	return EmTypeImpactLogicInvalid
}

func (f *FightObject) AddImpact(nImpactID, conAttTimes, nRound, nSkillID int, pCaster *FightObject) EmImpactResult {
	pRowImpactNew, err := ImpactTableRow(nImpactID)
	if err != nil {
		panic(err)
	}
	logicType := f.GetImpactLogicType(nImpactID)
	if logicType == EmTypeImpactLogicSingle {
		newImpact := new(Impact)
		newImpact.Init(nImpactID, conAttTimes, nRound, nSkillID, pCaster, f)
		newImpact.HeartBeat(nRound)
		return EmImpactResultNormal
	} else {
		//mutex
		if pRowImpactNew.ImpactMutexID >= 0 {
			for i := 0; i < MaxImpactNumber; i++ {
				if f.impactList[i].IsValid() {
					pRowImpact, err := ImpactTableRow(f.impactList[i].GetImpactID())
					if err != nil {
						panic(err)
					}
					if pRowImpactNew.ImpactMutexID == pRowImpact.ImpactMutexID {
						f.impactList[i].Init(nImpactID, conAttTimes, nRound, nSkillID, pCaster, f)
						f.impactList[i].HeartBeat(nRound)
						return EmImpactResultNormal
					} else {
						return EmImpactResultFail
					}
				}
			}
		}

		if logicType == EmTypeImpactLogicBuff {
			for i := 0; i < MaxBuffNumber; i++ {
				if !f.impactList[i].IsValid() {
					f.impactList[i].Init(nImpactID, conAttTimes, nRound, nSkillID, pCaster, f)
					f.impactList[i].HeartBeat(nRound)
					return EmImpactResultNormal
				}
			}
		}

		if logicType == EmTypeImpactLogicDeBuff {
			for i := MaxBuffNumber; i < MaxImpactNumber; i++ {
				if !f.impactList[i].IsValid() {
					f.impactList[i].Init(nImpactID, conAttTimes, nRound, nSkillID, pCaster, f)
					f.impactList[i].HeartBeat(nRound)
					return EmImpactResultNormal
				}
			}
		}
	}
	return EmImpactResultFail
}

func (f *FightObject) GetOwnerList() *FightObjList {
	if f.bAttacker {
		return f.pFightCell.GetAttackList()
	} else {
		return f.pFightCell.GetDefenceList()
	}
}

func (f *FightObject) GetHP() int {
	nEffectValue := f.impactEffect[EmAttributeHp].AttrValue
	nEndValue := f.FightDBData.HP + nEffectValue
	nMaxHP := f.GetMaxMP()

	if nEndValue > nMaxHP {
		nEndValue = nMaxHP
	}
	nEndValue = ChkMin(nEndValue, 0)
	f.FightDBData.HP = nEndValue
	f.impactEffect[EmAttributeHp].AttrValue = 0
	return nEndValue
}

func (f *FightObject) SetHP(nHP int) {
	nMaxHP := f.GetMaxMP()
	if nHP > nMaxHP {
		nHP = nMaxHP
	}
	if nHP < 0 {
		nHP = 0
	}
	f.FightDBData.HP = nHP
	if nHP == 0 {
		f.ClearImpact()
	}
}

func (f *FightObject) ClearImpact() bool {
	pOwnList := f.GetOwnerList()
	for i := 0; i < MaxMatrixCellCount; i++ {
		pFightObj := pOwnList.GetFightObject(i)
		if pFightObj != nil && pFightObj.IsActive() {
			pImpactList := pFightObj.GetImpactList()
			for j := 0; j < MaxImpactNumber; j++ {
				if pImpactList[j].IsValid() && pImpactList[j].GetCaster() == f {
					pImpactList[j].CleanUp()
				}
			}
		}
	}
	pOwnList = f.GetEnemyList()
	for i := 0; i < MaxMatrixCellCount; i++ {
		pFightObj := pOwnList.GetFightObject(i)
		if pFightObj != nil && pFightObj.IsActive() {
			pImpactList := pFightObj.GetImpactList()
			for j := 0; j < MaxImpactNumber; j++ {
				if pImpactList[j].IsValid() && pImpactList[j].GetCaster() == f {
					pImpactList[j].CleanUp()
				}
			}
		}
	}
	return true
}

func (f *FightObject) IsValid() bool {
	//fmt.Println(f.FightDBData.Guid)
	if f == nil {
		return false
	}
	return f.FightDBData.Guid.IsValid()
}

func (f *FightObject) IsActive() bool {
	if !f.IsValid() {
		return false
	}
	if f.GetHP() > 0 {
		return true
	}
	return false
}

func (f *FightObject) ChangeEffect(nAttrType EmAttribute, nValue int, bRemove bool) {
	if bRemove {
		f.impactEffect[nAttrType].AttrValue -= nValue
	} else {
		f.impactEffect[nAttrType].AttrValue += nValue
	}
}

//获取敌方对象列表
func (f *FightObject) GetEnemyList() *FightObjList {
	if f.bAttacker {
		return f.pFightCell.GetDefenceList()
	} else {
		return f.pFightCell.GetAttackList()
	}
}

//技能攻击
func (f *FightObject) SkillHeartBeat(uTime int) bool {

	for i := MaxSkillNum - 1; i >= 0; i-- {
		bLogic := f.skillList[i].SkillLogic(uTime)
		if bLogic {
			return true
		}
	}
	return false
}

//普通攻击
func (f *FightObject) CastCommonSkill(uTime int) bool {
	return f.commonKill.CommonSkillLogic(uTime)
}

//被动技能
func (f *FightObject) CastPassiveSkill(uTime int) {
	for i := 0; i < MaxSkillNum; i++ {
		//被动技能
		f.skillList[i].PassiveSkillLogic(uTime)
	}
	//武器被动技能
	for i := 0; i < f.equipSkillCount; i++ {
		f.equipSkillList[i].PassiveSkillLogic(uTime)
	}
}

//装备技能
func (f *FightObject) CastEquipSkill(uTime int, pTarget *FightObject) {
	for i := 0; i < f.equipSkillCount; i++ {
		f.equipSkillList[i].EquipSkillLogic(uTime, pTarget)
	}
}
