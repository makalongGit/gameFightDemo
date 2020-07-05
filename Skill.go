package GameFight

import (
	"errors"
	"fmt"
	"math/rand"
)

type Skill struct {
	skillID      int          //技能id
	skillType    EmSkillType  //技能类型
	skillTime    int          //技能时间
	coolDownTime int          //技能冷却时间
	pCaster      *FightObject //施放者
	pTarget      *FightObject //技能的目标
}

func (s Skill) String() string {
	return fmt.Sprintf("施放者 %+v 使用技能id %+v [技能类型 %+v 技能时间 %+v 技能冷却时间 %+v]",
		s.pCaster.GetGuid(), s.skillID, s.skillType, s.skillTime, s.coolDownTime)
}

func (s Skill) GetSkillRow() (TableRowSkill, error) {
	id := s.GetSkillID()
	if skill, ok := G_Skill[id]; ok {
		return skill, nil
	}
	return TableRowSkill{}, errors.New("config err")
}

func (s Skill) IsValid() bool {
	if s.skillID <= 0 {
		return false
	}
	return true
}

func (s Skill) CheckCondition(nRound int) bool {
	if s.skillType == EmSkillTypeHeroPassive || s.skillType == EmSkillTypeEquipPassive {
		return false
	}
	if nRound < s.skillTime {
		return false
	}
	if s.pCaster == nil {
		return false
	}
	pRow, err := s.GetSkillRow()
	if err != nil {
		return false
	}
	if s.pCaster.GetMP() < pRow.NeedMP {
		return false
	}
	nRand := rand.Intn(10000)
	if nRand > pRow.SkillRate {
		return false
	}
	return true

}

/*
   3	4	5

   0	1	2
  ------------
   0	1	2

   3	4	5

*/
func (s *Skill) SelectTarget() bool {
	var (
		pEnemyList   *FightObjList
		pFightObject *FightObject
		nCount       int
		nIndexList   [MaxMatrixCellCount]int
	)
	s.pTarget = nil
	pRow, err := s.GetSkillRow()
	if err != nil {
		panic(err)
	}
	pEnemyList = s.pCaster.GetEnemyList()
	switch pRow.SelectTargetOpt {
	case EmSkillTargetOptAuto: //自动选择
		s.pTarget = s.pCaster
		break
	case EmSkillTargetOptOrder: //0 顺序选择
		for i := 0; i < MaxMatrixCellCount/2; i++ {
			frontIndex := (s.pCaster.GetMatrixID() + i) % (MaxMatrixCellCount / 2)
			backIndex := frontIndex + MaxMatrixCellCount/2
			pFightObject = pEnemyList.GetFightObject(frontIndex)
			if pFightObject != nil && pFightObject.IsActive() {
				s.pTarget = pFightObject
				break
			}
			pFightObject = pEnemyList.GetFightObject(backIndex)
			if pFightObject != nil && pFightObject.IsActive() {
				s.pTarget = pFightObject
				break
			}
		}
		break
	case EmSkillTargetOptRand: //随机选择
		nCount = 0
		nIndexList = [MaxMatrixCellCount]int{}
		for i := 0; i < MaxMatrixCellCount; i++ {
			pFightObject = pEnemyList.GetFightObject(i)
			if pFightObject != nil && pFightObject.IsActive() {
				nIndexList[nCount] = i
				nCount++
			}
		}
		if nCount <= 0 {
			return false
		}
		rankIndex := rand.Intn(nCount)
		s.pTarget = pEnemyList.GetFightObject(nIndexList[rankIndex])
		break
	case EmSkillTargetOptSlow: //后排优先
		frontIndex := s.pCaster.GetMatrixID() % (MaxMatrixCellCount / 2)
		//后排 对线
		for i := 0; i < MaxMatrixCellCount/2; i++ {
			index := (frontIndex+i)%(MaxMatrixCellCount/2) + MaxMatrixCellCount/2
			pFightObject = pEnemyList.GetFightObject(index)
			if pFightObject != nil && pFightObject.IsActive() {
				s.pTarget = pFightObject
				return true
			}
		}
		//前排 对线
		for i := 0; i < MaxMatrixCellCount/2; i++ {
			index := (frontIndex + i) % (MaxMatrixCellCount / 2)
			pFightObject = pEnemyList.GetFightObject(index)
			if pFightObject != nil && pFightObject.IsActive() {
				s.pTarget = pFightObject
				return true
			}
		}
		break
	default:
		break
	}

	if s.pTarget != nil {
		return true
	}
	return false
}

func (s *Skill) GetConAttTimes() (nConAttTimes int) {
	nContinuous := s.pCaster.GetContinuous()
	if nContinuous > 100 {
		nConAttTimes = s.pCaster.GetConAttTimes()
		return
	}
	nRand := rand.Intn(100) + 1
	if nRand < nContinuous {
		nConAttTimes = s.pCaster.GetConAttTimes()
		return
	}
	return InvalidId
}

func (s *Skill) Init(SkillID int, Caster *FightObject) bool {
	s.skillID = SkillID
	s.pCaster = Caster

	SkillRow, err := s.GetSkillRow()
	if err != nil {
		return false
	}
	s.skillType = EmSkillType(SkillRow.SkillType)
	s.skillTime = SkillRow.StartRound
	s.coolDownTime = SkillRow.CoolDownTime
	return true
}

func (s Skill) GetSkillID() int {
	return s.skillID
}

func (s Skill) GetSkillType() EmSkillType {
	return s.skillType
}

func (s *Skill) CleanUp() {
	s.skillID = InvalidId
	s.skillType = EmSkillTypeInvalid
	s.pCaster = nil
	s.pTarget = nil
}

//获取施放者
func (s *Skill) GetCaster() *FightObject {
	return s.pCaster
}

//设置技能目标列表
func (s *Skill) SetSkillTarget(pTarget *FightObject) {
	s.pTarget = pTarget
}

func (s *Skill) GetTargetList(nType EmImpactTarget, targetList *TargetList) bool {
	targetList.CleanUp()
	switch nType {
	case EmImpactTargetOptSelf: //自身
		targetList.Add(s.pCaster)
		break
	case EmImpactTargetOwnerSignle: //己方个体
		pOwnerList := s.pCaster.GetOwnerList()
		if pOwnerList != nil {
			nIndexList := [MaxMatrixCellCount]int{}
			nCount := 0
			for i := 0; i < MaxMatrixCellCount; i++ {
				pFightObject := pOwnerList.GetFightObject(i)
				if pFightObject != nil && pFightObject.IsActive() {
					nIndexList[nCount] = i
					nCount++
				}
			}
			if nCount <= 0 {
				return false
			}
			randIndex := rand.Intn(nCount)
			targetList.Add(pOwnerList.GetFightObject(randIndex))
		}
		break
	case EmImpactTargetOwnerAll: //己方全体
		pOwnerList := s.pCaster.GetOwnerList()
		if pOwnerList != nil {
			for i := 0; i < MaxMatrixCellCount; i++ {
				pFightObject := pOwnerList.GetFightObject(i)
				if pFightObject != nil && pFightObject.IsActive() {
					targetList.Add(pFightObject)
				}
			}
		}
		break
	case EmImpactTargetEnemySignle: //敌方个体
		targetList.Add(s.pTarget)
		break
	case EmImpactTargetEnemyFront: //敌方横排
		matrixIndex := s.pTarget.GetMatrixID()
		pEnemyList := s.pCaster.GetEnemyList()
		if pEnemyList != nil {
			frontTargetList := TargetList{}
			backTargetList := TargetList{}
			for i := 0; i < MaxMatrixCellCount/2; i++ {
				pFightObject := pEnemyList.GetFightObject(i)
				if pFightObject != nil && pFightObject.IsActive() {
					frontTargetList.Add(pFightObject)
				}
			}
			for i := MaxMatrixCellCount / 2; i < MaxMatrixCellCount; i++ {
				pFightObject := pEnemyList.GetFightObject(i)
				if pFightObject != nil && pFightObject.IsActive() {
					backTargetList.Add(pFightObject)
				}
			}
			pTarget := &frontTargetList
			if matrixIndex < MaxMatrixCellCount/2 {
				if frontTargetList.GetCount() <= 0 {
					pTarget = &backTargetList
				}
			} else {
				if backTargetList.GetCount() > 0 {
					pTarget = &backTargetList
				}
			}
			targetList = pTarget

		}
		break
	case EmImpactTargetEnemyBehind: //敌方后排
		pEnemyList := s.pCaster.GetEnemyList()
		if pEnemyList != nil {
			for i := MaxMatrixCellCount / 2; i < MaxMatrixCellCount; i++ {
				pFightObject := pEnemyList.GetFightObject(i)
				if pFightObject != nil && pFightObject.IsActive() {
					targetList.Add(pFightObject)
				}
			}
			if targetList.GetCount() == 0 {
				for i := 0; i < MaxMatrixCellCount/2; i++ {
					pFightObject := pEnemyList.GetFightObject(i)
					if pFightObject != nil && pFightObject.IsActive() {
						targetList.Add(pFightObject)
					}
				}
			}
		}
		break
	case EmImpactTargetEnemyAll: //敌方全体
		pEnemyList := s.pCaster.GetEnemyList()
		if pEnemyList != nil {
			for i := 0; i < MaxMatrixCellCount; i++ {
				pFightObject := pEnemyList.GetFightObject(i)
				if pFightObject != nil && pFightObject.IsActive() {
					targetList.Add(pFightObject)
				}
			}
		}
		break
	case EmImpactTargetEnemyLine: //敌方目标竖排
		matrixIndex := s.pTarget.GetMatrixID()
		lineIndex := (matrixIndex + MaxMatrixCellCount/2) % MaxMatrixCellCount
		targetList.Add(s.pTarget)
		pEnemyList := s.pCaster.GetEnemyList()
		if pEnemyList != nil {
			pFightObject := pEnemyList.GetFightObject(lineIndex)
			if pFightObject != nil && pFightObject.IsActive() {
				targetList.Add(pFightObject)
			}
		}
		break
	case EmImpactTargetEnemyAround: //敌方目标及周围
		matrixIndex := s.pTarget.GetMatrixID()
		lineIndex := (matrixIndex + MaxMatrixCellCount/2) % MaxMatrixCellCount
		targetList.Add(s.pTarget)
		pEnemyList := s.pCaster.GetEnemyList()
		if pEnemyList != nil {
			pFightObject := pEnemyList.GetFightObject(lineIndex)
			if pFightObject != nil && pFightObject.IsActive() {
				targetList.Add(pFightObject)
			}
			if (matrixIndex+1)*2%MaxMatrixCellCount > 0 {
				pFightObject = pEnemyList.GetFightObject(matrixIndex + 1)
				if pFightObject != nil && pFightObject.IsActive() {
					targetList.Add(pFightObject)
				}
			}
			if (matrixIndex-1) >= 0 && matrixIndex*2%MaxMatrixCellCount > 0 {
				pFightObject = pEnemyList.GetFightObject(matrixIndex - 1)
				if pFightObject != nil && pFightObject.IsActive() {
					targetList.Add(pFightObject)
				}
			}
		}
		break
	case EmImpactTargetEnemyBehinDone: //后排个体
		pEnemyList := s.pCaster.GetEnemyList()
		if pEnemyList != nil {
			frontIndex := s.pTarget.GetMatrixID() % (MaxMatrixCellCount / 2)
			//后排对线
			for i := 0; i < MaxMatrixCellCount/2; i++ {
				index := (frontIndex+i)&(MaxMatrixCellCount/2) + MaxMatrixCellCount/2
				pFightObject := pEnemyList.GetFightObject(index)
				if pFightObject != nil && pFightObject.IsActive() {
					targetList.Add(pFightObject)
					return true
				}
			}
			//前排对线
			for i := 0; i < MaxMatrixCellCount/2; i++ {
				index := (frontIndex + i) % (MaxMatrixCellCount / 2)
				pFightObject := pEnemyList.GetFightObject(index)
				if pFightObject != nil && pFightObject.IsActive() {
					targetList.Add(pFightObject)
					return true
				}
			}
		}
		break
	case EmImpactTargetOwnerMinHp: //己方血最少
		index := InvalidId
		nMinHP := InvalidValue
		pOwnerList := s.pCaster.GetOwnerList()
		if pOwnerList != nil {
			for i := 0; i < MaxMatrixCellCount; i++ {
				pFightObject := pOwnerList.GetFightObject(i)
				if pFightObject != nil && pFightObject.IsActive() {
					if nMinHP < 0 {
						nMinHP = pFightObject.GetHP()
						index = i
					}
					if nMinHP > pFightObject.GetHP() {
						nMinHP = pFightObject.GetHP()
						index = i
					}
				}
			}
			if index >= 0 {
				targetList.Add(pOwnerList.GetFightObject(index))
			}
		}
		break
	case EmImpactTargetOwnerMinMp: //己方蓝最少
		index := InvalidId
		nMinMP := InvalidValue
		pOwnerList := s.pCaster.GetOwnerList()
		if pOwnerList != nil {
			for i := 0; i < MaxMatrixCellCount; i++ {
				pFightObject := pOwnerList.GetFightObject(i)
				if pFightObject != nil && pFightObject.IsActive() {
					if nMinMP < 0 {
						nMinMP = pFightObject.GetMP()
						index = i
					}
					if nMinMP > pFightObject.GetMP() {
						nMinMP = pFightObject.GetMP()
						index = i
					}
				}
			}
			if index >= 0 {
				targetList.Add(pOwnerList.GetFightObject(index))
			}
		}
		break
	default:
		break
	}
	return true
}

func (s *Skill) CastSkill(nRound int) bool {
	//fmt.Println("技能攻击流程:")
	var (
		nRand        int
		nConAttTimes int
		pRowImpact   TableRowImpact
		pAttackInfo  *AttackInfo
		pSkillAttack *SkillAttack
		impactInfo   *ImpactInfo
		targetList   *TargetList
	)
	pRow, err := s.GetSkillRow()
	if err != nil {
		panic(err)
	}
	s.pCaster.SetMP(s.pCaster.GetMP() - pRow.NeedMP)
	s.skillTime = nRound + pRow.CoolDownTime

	pAttackInfo = s.pCaster.GetAttackInfo()
	pAttackInfo.CastGuid = s.pCaster.GetGuid()
	pAttackInfo.Skilled = true
	pSkillAttack = pAttackInfo.AllocSkillAttack()
	if pSkillAttack == nil {
		fmt.Println("技能攻击信息为空")
		return false
	}
	pSkillAttack.SkillID = s.skillID
	pSkillAttack.SkillTarget = s.pTarget.GetGuid()
	pSkillAttack.CostMp = pRow.NeedMP

	//连击
	nConAttTimes = s.GetConAttTimes()
	for i := 0; i < MaxSkillImpactCount; i++ {
		nRand = rand.Intn(10000)
		if nRand <= pRow.ImpactRate[i] && pRow.ImpactID[i] > 0 {
			pRowImpact, err = ImpactTableRow(pRow.ImpactID[i])
			if err != nil {
				panic(err)
			}
			if !(pRowImpact.LogicID == 0 || pRowImpact.LogicID == 1) {
				nConAttTimes = 0
			}
			impactInfo = new(ImpactInfo)
			impactInfo.SkillID = s.skillID
			impactInfo.ImpactID = pRow.ImpactID[i]
			impactInfo.ConAttTimes = nConAttTimes

			targetList = new(TargetList)
			s.GetTargetList(EmImpactTarget(pRow.ImpactTargetType[i]), targetList)
			if targetList.GetCount() == 0 {
				continue
			}
			for i := 0; i < targetList.GetCount(); i++ {
				pTarget := targetList.GetFightObject(i)
				if pTarget != nil && pTarget.IsActive() {
					impactInfo.TargetList[impactInfo.TargetCount] = pTarget.GetGuid()
					impactInfo.TargetCount++
				}
			}
			pSkillAttack.AddImpactInfo(impactInfo)
			for index := 0; index < targetList.GetCount(); index++ {
				pTarget := targetList.GetFightObject(index)
				if pTarget != nil && pTarget.IsActive() {
					fmt.Println(pRow.ImpactID[i])
					pTarget.AddImpact(pRow.ImpactID[i], nConAttTimes, nRound, s.skillID, s.pCaster)
				}
			}
		}
	}
	return true
}

//魔法技能流程
func (s *Skill) SkillLogic(nRound int) bool {
	bRet := s.IsValid()
	if bRet == false {
		fmt.Println("技能id不合法")
		return false
	}
	bRet = s.CheckCondition(nRound)
	if bRet == false {
		//fmt.Println("技能达不到释放要求")
		return false
	}
	bRet = s.SelectTarget()
	if bRet == false {
		fmt.Println("目标选择错误")
		return false
	}
	//fmt.Printf("技能信息: %+v \n", s.String())
	bRet = s.CastSkill(nRound)
	if bRet == false {
		fmt.Println("技能逻辑错误")
		return false
	}
	return true
}

//被动技能
func (s *Skill) PassiveSkillLogic(nRound int) bool {
	bRet := s.IsValid()
	if bRet == false {
		return false
	}
	if nRound != 1 {
		return false
	}
	if s.skillType != EmSkillTypeHeroPassive && s.skillType != EmSkillTypeEquipPassive {
		return false
	}
	bRet = s.CheckCondition(nRound)
	if bRet == false {
		return false
	}
	bRet = s.SelectTarget()
	if bRet == false {
		return false
	}
	fmt.Println(s.String())
	bRet = s.CastSkill(nRound)
	if bRet == false {
		return false
	}
	return true
}

//装备技能逻辑
func (s *Skill) EquipSkillLogic(nRound int, pTarget *FightObject) bool {
	if s.skillType != EmSkillTypeEquipActive {
		return false
	}
	if nRound < s.skillTime {
		return false
	}
	pRow, err := s.GetSkillRow()
	if err != nil {
		panic(err)
	}
	if s.pCaster.GetMP() < pRow.NeedMP {
		return false
	}
	nRand := rand.Intn(10001)
	if nRand > pRow.SkillRate {
		return false
	}
	bRet := s.SelectTarget()
	if bRet == false {
		return false
	}
	bRet = s.CastSkill(nRound)
	if bRet == false {
		return false
	}
	return true
}

func (s Skill) CanHit() bool {
	nHit := s.pCaster.GetHit() - s.pTarget.GetDodge()
	if nHit >= 100 {
		return true
	}
	if nHit <= 0 {
		return false
	}
	nRand := rand.Intn(101)
	if nRand <= nHit {
		return true
	}
	return false
}

func (s Skill) CanStrike() bool {
	nStrike := s.pCaster.GetStrike()
	if nStrike >= 100 {
		return true
	}
	if nStrike <= 0 {
		return false
	}
	nRand := rand.Intn(101)
	if nRand <= nStrike {
		return true
	}
	return false
}

func (s Skill) CanBackAttack() bool {
	nAttackBack := s.pTarget.GetAttackBack()
	if nAttackBack >= 100 {
		return true
	}
	if nAttackBack <= 0 {
		return false
	}
	nRand := rand.Intn(101)
	if nRand <= nAttackBack {
		return true
	}
	return false
}

func (s *Skill) CommonSkillLogic(nRound int) bool {
	//fmt.Println("普通攻击流程: ")
	pAttackInfo := s.pCaster.GetAttackInfo()
	bRet := s.SelectTarget()

	if bRet == false {
		pAttackInfo.CastGuid = InvalidId
		return false
	}
	pAttackInfo.CastGuid = s.pCaster.GetGuid()

	//命中
	bRet = s.CanHit()
	pAttackInfo.BHit = bRet
	if bRet == false {
		return false
	}
	s.pCaster.CastEquipSkill(nRound, s.pTarget)
	if !s.pTarget.IsActive() {
		pAttackInfo.Skilled = false
		return true
	}

	pAttackInfo.Skilled = false
	nPhysicAttack := s.pCaster.GetPhysicAttack()
	bStrike := s.CanStrike()

	pAttackInfo.BStrike = bStrike
	pAttackInfo.SkillTarget = s.pTarget.GetGuid()
	if bStrike {
		//暴击攻击力=当前攻击力*（1+自身暴击伤害/100）
		nPhysicAttack = int(nPhysicAttack * (1 + s.pCaster.GetStrikeHurt()/100));
	}

	pFightCell := s.pCaster.GetFightCell()
	if pFightCell.GetFightType() == EmTypeFightStair && s.pCaster.IsAttacker() {
		nPhysicAttack += nPhysicAttack * pFightCell.GetPlusAtt() / 100
	}

	nDefend := s.pTarget.GetPhysicDefend()
	nDecay := s.pTarget.GetPhysicHurtDecay()
	//本次物理攻击伤害=(自身当前经过暴击计算后的物理攻击-目标物理防御)*(1-目标物理伤害减免/100)*(1+物理攻击伤害浮动)
	nDamage := CalcDamage(nPhysicAttack, nDefend, nDecay)
	nDamage = ChkMin(nDamage, 0)

	s.pTarget.SetHP(s.pTarget.GetHP() - nDamage)
	pAttackInfo.SkillTarget = s.pTarget.GetGuid()
	pAttackInfo.Hurt = nDamage

	//反击
	if s.pTarget.IsActive() && s.CanBackAttack() {
		//反击伤害=自身攻击*反击衰减
		backAttack := s.pTarget.GetPhysicAttack() * s.pTarget.GetBackAttHurt() / 100
		s.pCaster.SetHP(s.pCaster.GetHP() - backAttack)
		pAttackInfo.BBackAttack = true
		pAttackInfo.BackAttackHurt = backAttack
	}

	//fmt.Println(pAttackInfo.String())
	return true
}
