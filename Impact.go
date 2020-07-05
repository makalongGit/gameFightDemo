package GameFight

import (
	"errors"
)

type Impact struct {
	LogicFuncArray [EmImpactLogicCount]func()
	impactID       int          //ID
	pHolder        *FightObject //拥有者
	pCaster        *FightObject //施放者
	skillID        int          //技能ID
	startTime      int          //开始时间
	life           int          //开始周期
	logicTime      int          //开始作用时间
	term           int          //开始间隔
	conAttTimes    int          //连击次数
}

func (impact Impact) GetImpactRow() (TableRowImpact, error) {
	id := impact.GetImpactID()
	if impact, ok := G_Impact[id]; ok {
		return impact, nil
	} else {
		return TableRowImpact{}, errors.New("config err")
	}
}

func (impact Impact) GetImpactID() int {
	return impact.impactID
}

func (impact Impact) GetHolder() *FightObject {
	return impact.pHolder
}

func (impact Impact) GetCaster() *FightObject {
	return impact.pCaster
}

func (impact Impact) IsValid() bool {
	if impact.impactID > 0 {
		return true
	}
	return false
}

func (impact *Impact) Init(nImpactID, conAttTimes, nRound, nSkillID int, pCaster, pHolder *FightObject) bool {
	impact.impactID = nImpactID
	impact.pHolder = pHolder
	impact.pCaster = pCaster
	impact.skillID = nSkillID
	impact.startTime = nRound
	impact.life = nRound
	impact.logicTime = nRound
	impact.conAttTimes = 0
	impact.term = 1

	pRow, err := impact.GetImpactRow()
	if err != nil {
		return false
	}
	switch pRow.LogicID {
	case EmImpactLogic0, EmImpactLogic1, EmImpactLogic6:
		if conAttTimes > 0 {
			impact.conAttTimes = conAttTimes
		}
		break
	case EmImpactLogic2,EmImpactLogic3:
		//2=物理持续攻击
		//逻辑参数3：持续时间，单位回合（10）
		//逻辑参数4：生效间隔，单位回合（2）
		impact.life += pRow.Param[2] - pRow.Param[3]
		impact.term = pRow.Param[3]
		break
	case EmImpactLogic4,EmImpactLogic5:
		//逻辑参数3：持续时间，单位回合（10）
		impact.life += pRow.Param[2] - 1
		break
	default:
		break
	}
	impact.LogicFuncArray[0] = impact.ImpactLogic0
	impact.LogicFuncArray[1] = impact.ImpactLogic1
	impact.LogicFuncArray[2] = impact.ImpactLogic2
	impact.LogicFuncArray[3] = impact.ImpactLogic3
	impact.LogicFuncArray[4] = impact.ImpactLogic4
	impact.LogicFuncArray[5] = impact.ImpactLogic5
	impact.LogicFuncArray[6] = impact.ImpactLogic6
	return true
}

//0=单次物理攻击；
//逻辑参数1：从英雄物理攻击中取得的倍率（例：150）
//逻辑参数2：额外增加的物理伤害（例：50）
//如英雄物理攻击为100，则最终的技能物理攻击=英雄物理攻击（100）*逻辑参数1（150/100）+逻辑参数2（50）
func (impact Impact) ImpactLogic0() {
	var (
		pImpactInfo   *ImpactInfo
		pAttackInfo   *AttackInfo
		pSkillAttack  *SkillAttack
		pFightCell    *FightCell
		index         int
		pRow          TableRowImpact
		nPhysicAttack int
		nAttack       int
		nDefend       int
		nDecay        int
		nDamage       int
		fConAttHurt   int
	)
	index = InvalidId
	pAttackInfo = impact.pCaster.GetAttackInfo()
	if pAttackInfo != nil && pAttackInfo.IsValid() {
		pSkillAttack = pAttackInfo.GetSkillAttack(impact.skillID)
		if pSkillAttack != nil {
			pImpactInfo = pSkillAttack.GetImpactInfo(impact.impactID)
			if pImpactInfo != nil {
				index = pImpactInfo.GetTargetIndex(impact.pHolder.GetGuid())
			}
		}
	}

	pFightCell = impact.pHolder.GetFightCell()
	nPhysicAttack = impact.pCaster.GetPhysicAttack()

	pRow, err := impact.GetImpactRow()
	if err != nil {
		panic(err)
	}
	nAttack = int(nPhysicAttack*pRow.Param[0]/100 + pRow.Param[1])
	if pFightCell.GetFightType() == EmTypeFightStair && impact.pCaster.IsAttacker() {
		nAttack += nAttack * pFightCell.GetPlusAtt() / 100
	}

	nDefend = impact.pHolder.GetPhysicDefend()
	nDecay = impact.pHolder.GetPhysicHurtDecay()
	//本次物理攻击伤害=(自身当前经过连击计算后的物理攻击-目标物理防御)*(1-目标物理伤害减免/100)
	nDamage = CalcDamage(nAttack, nDefend, nDecay)
	nDamage = ChkMin(nDamage, 0)
	impact.pHolder.SetHP(impact.pHolder.GetHP() - nDamage)

	if pImpactInfo != nil && index >= 0 {
		pImpactInfo.Hurts[index][0] = nDamage
	}
	//连击
	fConAttHurt = impact.pCaster.GetConAttHurt() / 100
	for i := 1; i < impact.conAttTimes; i++ {
		nAttack = int(nAttack * fConAttHurt)
		nDamage = CalcDamage(nAttack, nDefend, nDecay)
		nDamage = ChkMin(nDamage, 0)
		impact.pHolder.SetHP(impact.pHolder.GetHP() - nDamage)

		if pImpactInfo != nil && index >= 0 {
			pImpactInfo.Hurts[index][i] = nDamage
		}

	}
}

func (impact Impact) ImpactLogic1() {
	var (
		pImpactInfo   *ImpactInfo
		pAttackInfo   *AttackInfo
		pSkillAttack  *SkillAttack
		pFightCell    *FightCell
		index         int
		pRow          TableRowImpact
		nPhysicAttack int
		nAttack       int
		nDefend       int
		nDecay        int
		nDamage       int
		fConAttHurt   int
	)
	index = InvalidId
	pAttackInfo = impact.pCaster.GetAttackInfo()
	if pAttackInfo != nil && pAttackInfo.IsValid() {
		pSkillAttack = pAttackInfo.GetSkillAttack(impact.skillID)
		if pSkillAttack != nil {
			pImpactInfo = pSkillAttack.GetImpactInfo(impact.impactID)
			if pImpactInfo != nil {
				index = pImpactInfo.GetTargetIndex(impact.pHolder.GetGuid())
			}
		}
	}

	pFightCell = impact.pHolder.GetFightCell()
	nPhysicAttack = impact.pCaster.GetPhysicAttack()
	pRow, err := impact.GetImpactRow()
	if err != nil {
		panic(err)
	}
	nAttack = int(nPhysicAttack*pRow.Param[0]/100 + pRow.Param[1])
	if pFightCell.GetFightType() == EmTypeFightStair && impact.pCaster.IsAttacker() {
		nAttack += int(nAttack * pFightCell.GetPlusAtt() / 100)
	}

	nDefend = impact.pHolder.GetMagicDefend()
	nDecay = impact.pHolder.GetMagicHurtDecay()
	nDamage = CalcDamage(nAttack, nDefend, nDecay)
	nDamage = ChkMin(nDamage, 0)
	//本次魔法攻击伤害=(自身当前经过连击计算后的魔法攻击-目标魔法防御)*(1-目标魔法伤害减免/100)

	impact.pHolder.SetHP(impact.pHolder.GetHP() - nDamage)
	if pImpactInfo != nil && index >= 0 {
		pImpactInfo.Hurts[index][0] = nDamage
	}
	//连击
	fConAttHurt = impact.pCaster.GetConAttHurt() / 100
	for i := 1; i < impact.conAttTimes; i++ {
		nAttack = int(nAttack * fConAttHurt)
		nDamage = CalcDamage(nAttack, nDefend, nDecay)
		nDamage = ChkMin(nDamage, 0)
		impact.pHolder.SetHP(impact.pHolder.GetHP() - nDamage)

		if pImpactInfo != nil && index >= 0 {
			pImpactInfo.Hurts[index][i] = nDamage
		}
	}
}

//2=物理持续攻击
//逻辑参数1：从英雄物理攻击中取得的倍率（例：150）
//逻辑参数2：额外增加的物理伤害（例：50）
//逻辑参数3：持续时间，单位回合（10）
//逻辑参数4：生效间隔，单位回合（2）
//此效果的最终效果为：每2回合对目标造成一次物理伤害，持续10回合（生效5次），每次造成的物理攻击具体数值=英雄物理攻击（100）*逻辑参数1（150/100）+逻辑参数2（50）
func (impact Impact) ImpactLogic2() {
	var (
		pImpactInfo   *ImpactInfo
		pAttackInfo   *AttackInfo
		pSkillAttack  *SkillAttack
		pFightCell    *FightCell
		index         int
		nPhysicAttack int
		nAttack       int
		nDefend       int
		nDecay        int
		nDamage       int
		pRoundInfo    *FightRoundInfo
		pFObjectInfo  *FObjectInfo
	)
	index = InvalidId
	pAttackInfo = impact.pCaster.GetAttackInfo()
	if pAttackInfo != nil && pAttackInfo.IsValid() {
		pSkillAttack = pAttackInfo.GetSkillAttack(impact.skillID)
		if pSkillAttack != nil {
			pImpactInfo = pSkillAttack.GetImpactInfo(impact.impactID)
			if pImpactInfo != nil {
				index = pImpactInfo.GetTargetIndex(impact.pHolder.GetGuid())
			}
		}
	}
	pFightCell = impact.pHolder.GetFightCell()
	if impact.logicTime <= impact.life {
		nPhysicAttack = impact.pCaster.GetPhysicAttack()
		pRow, err := impact.GetImpactRow()
		if err != nil {
			panic(err)
		}
		nAttack = int(nPhysicAttack*pRow.Param[0]/100 + pRow.Param[1])
		if pFightCell.GetFightType() == EmTypeFightStair && impact.pCaster.IsAttacker() {
			nAttack += nAttack * pFightCell.GetPlusAtt() / 100
		}
		nDefend = impact.pHolder.GetPhysicDefend()
		nDecay = impact.pHolder.GetPhysicHurtDecay()
		//本次物理攻击伤害=(自身当前经过连击计算后的物理攻击-目标物理防御)*(1-目标物理伤害减免/100)
		nDamage = CalcDamage(nAttack, nDefend, nDecay)
		nDamage = ChkMin(nDamage, 0)
		impact.pHolder.SetHP(impact.pHolder.GetHP() - nDamage)

		if impact.startTime == impact.logicTime {
			if pImpactInfo != nil && index >= 0 {
				pImpactInfo.Hurts[index][0] = nDamage
			}
		} else if impact.startTime < impact.logicTime {
			pRoundInfo = pFightCell.GetRoundInfo()
			pFObjectInfo = pRoundInfo.GetFObjectInfoByGuid(impact.pHolder.GetGuid())
			if pFObjectInfo != nil {
				for i := 0; i < pFObjectInfo.ImpactCount; i++ {
					if pFObjectInfo.ImpactList[i] == impact.impactID {
						pFObjectInfo.ImpactHurt[i] = nDamage
						break
					}
				}
			}
		}
	}
	impact.logicTime += impact.term
}

//3=持续魔法攻击
//逻辑参数1：从英雄物理攻击（因为英雄默认没有魔法攻击）中取得的倍率（例：150）
//逻辑参数2：额外增加的魔法伤害（例：50）
//逻辑参数3：持续时间，单位回合（10）
//逻辑参数4：生效间隔，单位回合（2）
//此效果的最终效果为：每2回合对目标造成一次魔法伤害，持续10回合（生效5次），每次造成的魔法攻击具体数值=英雄物理攻击（100）*逻辑参数1（150/100）+逻辑参数2（50）
func (impact *Impact) ImpactLogic3() {
	var (
		pImpactInfo   *ImpactInfo
		pAttackInfo   *AttackInfo
		pSkillAttack  *SkillAttack
		pFightCell    *FightCell
		index         int
		nPhysicAttack int
		nAttack       int
		nDefend       int
		nDecay        int
		nDamage       int
		pRoundInfo    *FightRoundInfo
		pFObjectInfo  *FObjectInfo
	)
	index = InvalidId
	pAttackInfo = impact.pCaster.GetAttackInfo()
	if pAttackInfo != nil && pAttackInfo.IsValid() {
		pSkillAttack = pAttackInfo.GetSkillAttack(impact.skillID)
		if pSkillAttack != nil {
			pImpactInfo = pSkillAttack.GetImpactInfo(impact.impactID)
			if pImpactInfo != nil {
				index = pImpactInfo.GetTargetIndex(impact.pHolder.GetGuid())
			}
		}
	}
	pFightCell = impact.pHolder.GetFightCell()
	if impact.logicTime <= impact.life {
		nPhysicAttack = impact.pCaster.GetPhysicAttack()
		pRow, err := impact.GetImpactRow()
		if err != nil {
			panic(err)
		}
		nAttack = int(nPhysicAttack*pRow.Param[0]/100 + pRow.Param[1])
		if pFightCell.GetFightType() == EmTypeFightStair && impact.pCaster.IsAttacker() {
			nAttack += nAttack * pFightCell.GetPlusAtt() / 100
		}
		nDefend = impact.pHolder.GetMagicDefend()
		nDecay = impact.pHolder.GetMagicHurtDecay()
		//本次物理攻击伤害=(自身当前经过连击计算后的物理攻击-目标物理防御)*(1-目标物理伤害减免/100)
		nDamage = CalcDamage(nAttack, nDefend, nDecay)
		nDamage = ChkMin(nDamage, 0)
		impact.pHolder.SetHP(impact.pHolder.GetHP() - nDamage)
	}
	if impact.startTime == impact.logicTime {
		if pImpactInfo != nil && index >= 0 {
			pImpactInfo.Hurts[index][0] = nDamage
		}
	} else if impact.startTime < impact.logicTime {
		pRoundInfo = pFightCell.GetRoundInfo()
		pFObjectInfo = pRoundInfo.GetFObjectInfoByGuid(impact.pHolder.GetGuid())
		if pFObjectInfo != nil {
			for i := 0; i < pFObjectInfo.ImpactCount; i++ {
				if pFObjectInfo.ImpactList[i] == impact.impactID {
					pFObjectInfo.ImpactHurt[i] = nDamage
					break
				}
			}
		}
	}
	impact.logicTime += impact.term
}

func (impact *Impact) ImpactLogic4() {
	var (
		pImpactInfo  *ImpactInfo
		pAttackInfo  *AttackInfo
		pSkillAttack *SkillAttack
		pFightCell   *FightCell
		index        int
		pRoundInfo   *FightRoundInfo
		pFObjectInfo *FObjectInfo
	)
	index = InvalidId
	pAttackInfo = impact.pCaster.GetAttackInfo()
	if pAttackInfo != nil && pAttackInfo.IsValid() {
		pSkillAttack = pAttackInfo.GetSkillAttack(impact.skillID)
		if pSkillAttack != nil {
			pImpactInfo = pSkillAttack.GetImpactInfo(impact.impactID)
			if pImpactInfo != nil {
				index = pImpactInfo.GetTargetIndex(impact.pHolder.GetGuid())
			}
		}
	}
	if impact.logicTime <= impact.life {
		pRow, err := impact.GetImpactRow()
		if err != nil {
			panic(err)
		}
		impact.pHolder.ChangeEffect(EmAttribute(pRow.Param[0]), pRow.Param[1], false)
		if impact.startTime == impact.logicTime {
			if pImpactInfo != nil && index >= 0 {
				if EmAttribute(pRow.Param[0]) == EmAttributeHp {
					impact.pHolder.SetHP(impact.pHolder.GetHP())
					pImpactInfo.Hurts[index][0] = pRow.Param[1] * -1
				}
				if EmAttribute(pRow.Param[0]) == EmAttributeMp {
					pImpactInfo.Mp[index] = pRow.Param[1] * -1
				}
			}
		} else if impact.startTime == impact.logicTime {
			pFightCell = impact.pHolder.GetFightCell()
			pRoundInfo = pFightCell.GetRoundInfo()
			pFObjectInfo = pRoundInfo.GetFObjectInfoByGuid(impact.pHolder.GetGuid())
			if pFObjectInfo != nil {
				for i := 0; i < pFObjectInfo.ImpactCount; i++ {
					if pFObjectInfo.ImpactList[i] == impact.impactID {
						if EmAttribute(pRow.Param[0]) == EmAttributeHp {
							pFObjectInfo.ImpactHurt[i] = pRow.Param[1] * -1
						}
						if EmAttribute(pRow.Param[0]) == EmAttributeMp {
							pFObjectInfo.ImpactMP[i] = pRow.Param[1] * -1
						}
						break
					}
				}
			}
		}
	}
	impact.logicTime += 1
}

//5=debuff削弱类
//逻辑参数1：改变的英雄属性id，读取AttributeData.tab表。
//逻辑参数2：改变的具体数值
//逻辑参数3：持续时间，单位回合（10）
//最终可实现的效果如敌人攻击减少X点持续10回合。
func (impact *Impact) ImpactLogic5() {
	var (
		pImpactInfo  *ImpactInfo
		pAttackInfo  *AttackInfo
		pSkillAttack *SkillAttack
		pFightCell   *FightCell
		index        int
		pRoundInfo   *FightRoundInfo
		pFObjectInfo *FObjectInfo
	)
	index = InvalidId
	pAttackInfo = impact.pCaster.GetAttackInfo()
	if pAttackInfo != nil && pAttackInfo.IsValid() {
		pSkillAttack = pAttackInfo.GetSkillAttack(impact.skillID)
		if pSkillAttack != nil {
			pImpactInfo = pSkillAttack.GetImpactInfo(impact.impactID)
			if pImpactInfo != nil {
				index = pImpactInfo.GetTargetIndex(impact.pHolder.GetGuid())
			}
		}
	}
	if impact.logicTime <= impact.life {
		pRow, err := impact.GetImpactRow()
		if err != nil {
			panic(err)
		}
		impact.pHolder.ChangeEffect(EmAttribute(pRow.Param[0]), pRow.Param[1]*(-1), false)
		if impact.startTime == impact.logicTime {
			if pImpactInfo != nil && index >= 0 {
				if EmAttribute(pRow.Param[0]) == EmAttributeHp {
					impact.pHolder.SetHP(impact.pHolder.GetHP())
					pImpactInfo.Hurts[index][0] = pRow.Param[1]
				}
				if EmAttribute(pRow.Param[0]) == EmAttributeMp {
					pImpactInfo.Mp[index] = pRow.Param[1]
				}
			}
		} else if impact.startTime == impact.logicTime {
			pFightCell = impact.pHolder.GetFightCell()
			pRoundInfo = pFightCell.GetRoundInfo()
			pFObjectInfo = pRoundInfo.GetFObjectInfoByGuid(impact.pHolder.GetGuid())
			if pFObjectInfo != nil {
				for i := 0; i < pFObjectInfo.ImpactCount; i++ {
					if pFObjectInfo.ImpactList[i] == impact.impactID {
						if EmAttribute(pRow.Param[0]) == EmAttributeHp {
							pFObjectInfo.ImpactHurt[i] = pRow.Param[1]
						}
						if EmAttribute(pRow.Param[0]) == EmAttributeMp {
							pFObjectInfo.ImpactMP[i] = pRow.Param[1]
						}
						break
					}
				}
			}
		}
	}
	impact.logicTime += 1
}

func (impact Impact) ImpactLogic6() {
	var (
		pImpactInfo  *ImpactInfo
		pAttackInfo  *AttackInfo
		pSkillAttack *SkillAttack
		nMagicAttack int
		nAttack      int
		nDefend      int
		nDecay       int
		nDamage      int
		fConAttHurt  int
		pFightCell   *FightCell
		index        int
	)
	index = InvalidId
	pAttackInfo = impact.pCaster.GetAttackInfo()
	if pAttackInfo != nil && pAttackInfo.IsValid() {
		pSkillAttack = pAttackInfo.GetSkillAttack(impact.skillID)
		if pSkillAttack != nil {
			pImpactInfo = pSkillAttack.GetImpactInfo(impact.impactID)
			if pImpactInfo != nil {
				index = pImpactInfo.GetTargetIndex(impact.pHolder.GetGuid())
			}
		}
	}

	pFightCell = impact.pHolder.GetFightCell()
	nMagicAttack = impact.pCaster.GetMagicAttack()
	pRow, err := impact.GetImpactRow()
	if err != nil {
		panic(err)
	}
	nAttack = int(nMagicAttack*pRow.Param[0]/100 + pRow.Param[1])
	if pFightCell.GetFightType() == EmTypeFightStair && impact.pCaster.IsAttacker() {
		nAttack += int(nAttack * pFightCell.GetPlusAtt() / 100)
	}
	nDefend = impact.pHolder.GetMagicDefend()
	nDecay = impact.pHolder.GetMagicHurtDecay()
	nDamage = CalcDamage(nAttack, nDefend, nDecay)
	nDamage = ChkMin(nDamage, 0)
	//本次魔法攻击伤害=(自身当前经过连击计算后的魔法攻击-目标魔法防御)*(1-目标魔法伤害减免/100)

	impact.pHolder.SetHP(impact.pHolder.GetHP() - nDamage)
	if pImpactInfo != nil && index >= 0 {
		pImpactInfo.Hurts[index][0] = nDamage
	}

	//连击
	fConAttHurt = impact.pCaster.GetConAttHurt() / 100
	for i := 1; i < impact.conAttTimes; i++ {
		nAttack = int(nAttack * fConAttHurt)
		nDamage = CalcDamage(nAttack, nDefend, nDecay)
		nDamage = ChkMin(nDamage, 0)
		impact.pHolder.SetHP(impact.pHolder.GetHP() - nDamage)

		if pImpactInfo != nil && index >= 0 {
			pImpactInfo.Hurts[index][i] = nDamage
		}
	}

}

func (impact *Impact) CleanUp() {
	impact.impactID = InvalidId
	impact.pHolder = nil
	impact.pCaster = nil
	impact.startTime = 0
	impact.life = 0
	impact.term = 0
	impact.logicTime = 0
	impact.conAttTimes = 0
}

func (impact *Impact) HeartBeat(uTime int) bool {
	if !impact.IsValid() {
		return false
	}
	if uTime == impact.logicTime {
		pRow, err := impact.GetImpactRow()
		if err != nil {
			panic(err)
		}
		if pRow.LogicID >= 0 && pRow.LogicID < EmImpactLogicCount {
			logicFunc := impact.LogicFuncArray[pRow.LogicID]
			logicFunc()
		}
	}

	if uTime >= impact.life {
		impact.CleanUp()
	}
	return true
}