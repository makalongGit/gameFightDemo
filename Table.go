package GameFight

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var G_Impact map[int]TableRowImpact
var G_Skill map[int]TableRowSkill
var G_HeroAttr map[int]TableRowHeroAttr

func init() {
	G_Impact = make(map[int]TableRowImpact)
	ReadImpactCsv("Impact.tab.csv")
	G_Skill = make(map[int]TableRowSkill)
	ReadSkillCsv("Skill.csv")
	G_HeroAttr = make(map[int]TableRowHeroAttr)
	ReadHeroAttrCsv("HeroAttr.csv")
}
func ReadImpactCsv(filename string) bool {
	//获取数据
	fileName := "./Config/" + filename
	cntb, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return false
	}

	r := csv.NewReader(strings.NewReader(string(cntb)))
	ss, _ := r.ReadAll()
	sz := len(ss)
	for i := 2; i < sz; i++ {
		tmp := TableRowImpact{}
		tmp.ImpactID, _ = strconv.Atoi(ss[i][0])
		tmp.Description = ss[i][1]
		tmp.LogicID, _ = strconv.Atoi(ss[i][2])
		for j := 0; j < MaxImpactLogicParamCount; j++ {
			tmp.Param[j], _ = strconv.Atoi(ss[i][3+j])
		}
		tmp.Icon = ss[i][MaxImpactLogicParamCount+1]
		tmp.ImpactMutexID, _ = strconv.Atoi(ss[i][MaxImpactLogicParamCount+2])
		tmp.ReplaceLevel, _ = strconv.Atoi(ss[i][MaxImpactLogicParamCount+3])
		tmp.ReplaceLevel, _ = strconv.Atoi(ss[i][MaxImpactLogicParamCount+3])
		tmp.DeadDisappear, _ = strconv.Atoi(ss[i][MaxImpactLogicParamCount+4])
		tmp.OfflineTimeGO, _ = strconv.Atoi(ss[i][MaxImpactLogicParamCount+5])
		tmp.ScriptID, _ = strconv.Atoi(ss[i][MaxImpactLogicParamCount+6])
		tmp.Effect = ss[i][MaxImpactLogicParamCount+7]
		tmp.Name = ss[i][MaxImpactLogicParamCount+8]
		tmp.Desc = ss[i][MaxImpactLogicParamCount+9]
		tmp.SkillEffect = ss[i][MaxImpactLogicParamCount+10]
		G_Impact[tmp.ImpactID] = tmp
	}
	return true
}
func ReadSkillCsv(filename string) bool {
	//获取数据
	fileName := "./Config/" + filename
	cntb, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return false
	}

	r := csv.NewReader(strings.NewReader(string(cntb)))
	ss, _ := r.ReadAll()
	sz := len(ss)
	for i := 2; i < sz; i++ {
		tmp := TableRowSkill{}
		tmp.SkillId, _ = strconv.Atoi(ss[i][0])
		tmp.Description = ss[i][1]
		tmp.Name = ss[i][2]
		tmp.SkillLevel, _ = strconv.Atoi(ss[i][3])
		tmp.Tips = ss[i][4]
		tmp.SelectTargetOpt, _ = strconv.Atoi(ss[i][5])
		tmp.SkillRate, _ = strconv.Atoi(ss[i][6])
		tmp.NeedMP, _ = strconv.Atoi(ss[i][7])
		tmp.CoolDownTime, _ = strconv.Atoi(ss[i][8])
		tmp.CoolDownID, _ = strconv.Atoi(ss[i][9])
		tmp.StartRound, _ = strconv.Atoi(ss[i][10])
		tmp.SkillType, _ = strconv.Atoi(ss[i][11])
		index := 9
		for j := 0; j < MaxSkillImpactCount; j++ {
			index = index + 3
			tmp.ImpactID[j], _ = strconv.Atoi(ss[i][index])
			tmp.ImpactRate[j], _ = strconv.Atoi(ss[i][index+1])
			tmp.ImpactTargetType[j], _ = strconv.Atoi(ss[i][index+2])
		}
		tmp.ScriptID, _ = strconv.Atoi(ss[i][24])
		tmp.HeroLevel, _ = strconv.Atoi(ss[i][25])
		tmp.ConsumeMoney, _ = strconv.Atoi(ss[i][26])
		tmp.LearnType, _ = strconv.Atoi(ss[i][27])

		G_Skill[tmp.SkillId] = tmp
	}
	return true
}
func ReadHeroAttrCsv(filename string) bool {
	//获取数据
	fileName := "./Config/" + filename
	cntb, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return false
	}

	r := csv.NewReader(strings.NewReader(string(cntb)))
	ss, _ := r.ReadAll()
	sz := len(ss)
	for i := 2; i < sz; i++ {
		tmp := TableRowHeroAttr{}
		tmp.ID, _ = strconv.Atoi(ss[i][0])
		tmp.SpiritID, _ = strconv.Atoi(ss[i][1])
		tmp.Name = ss[i][2]
		tmp.Profession, _ = strconv.Atoi(ss[i][3])
		tmp.InitQuality, _ = strconv.Atoi(ss[i][4])
		tmp.InitLevel, _ = strconv.Atoi(ss[i][5])
		tmp.LevelLimit, _ = strconv.Atoi(ss[i][6])
		tmp.LevelCrossRole, _ = strconv.Atoi(ss[i][7])
		tmp.InitExp, _ = strconv.Atoi(ss[i][8])
		tmp.TakeLevel, _ = strconv.Atoi(ss[i][9])
		tmp.EffectAttackByGrowRate, _ = strconv.Atoi(ss[i][10])
		tmp.EffectDefendByGrowRate, _ = strconv.Atoi(ss[i][11])
		tmp.EffectHPByGrowRate, _ = strconv.Atoi(ss[i][12])
		tmp.EffectMPByGrowRate, _ = strconv.Atoi(ss[i][13])
		tmp.EffectPhysicAttackByLevel, _ = strconv.Atoi(ss[i][14])
		tmp.EffectPhysicDefendByLevel, _ = strconv.Atoi(ss[i][15])
		tmp.EffectHpByLevel, _ = strconv.Atoi(ss[i][16])
		tmp.EffectMpByLevel, _ = strconv.Atoi(ss[i][17])
		tmp.InitAttackSpeed, _ = strconv.Atoi(ss[i][18])
		tmp.InitPhysicAttack, _ = strconv.Atoi(ss[i][19])
		tmp.InitMagicAttack, _ = strconv.Atoi(ss[i][20])
		tmp.InitPhysicDefence, _ = strconv.Atoi(ss[i][21])
		tmp.InitMagicDefence, _ = strconv.Atoi(ss[i][22])
		tmp.InitHP, _ = strconv.Atoi(ss[i][23])
		tmp.InitMP, _ = strconv.Atoi(ss[i][24])
		tmp.InitHit, _ = strconv.Atoi(ss[i][25])
		tmp.InitDodge, _ = strconv.Atoi(ss[i][26])
		tmp.InitStrike, _ = strconv.Atoi(ss[i][27])
		tmp.InitContinuous, _ = strconv.Atoi(ss[i][28])
		tmp.InitBackAttack, _ = strconv.Atoi(ss[i][29])
		tmp.InitStrikeHurt, _ = strconv.Atoi(ss[i][30])
		tmp.InitConAttTimes, _ = strconv.Atoi(ss[i][31])
		tmp.InitConAttHurt, _ = strconv.Atoi(ss[i][32])
		tmp.InitBackAttHurt, _ = strconv.Atoi(ss[i][33])
		tmp.InitPhysicHurtDecay, _ = strconv.Atoi(ss[i][34])
		tmp.InitMagicHurtDecay, _ = strconv.Atoi(ss[i][35])
		tmp.FloatingHurt, _ = strconv.Atoi(ss[i][36])
		tmp.PhysicSkillID, _ = strconv.Atoi(ss[i][37])
		tmp.MagicSkillID1, _ = strconv.Atoi(ss[i][38])
		tmp.LearnedLevel1, _ = strconv.Atoi(ss[i][39])
		tmp.MagicSkillID2, _ = strconv.Atoi(ss[i][40])
		tmp.LearnedLevel2, _ = strconv.Atoi(ss[i][41])
		tmp.MaxGrowPoint, _ = strconv.Atoi(ss[i][42])
		tmp.InitBearPoint, _ = strconv.Atoi(ss[i][43])
		tmp.BearParam, _ = strconv.Atoi(ss[i][44])
		tmp.RequireHumanLevel, _ = strconv.Atoi(ss[i][47])

		G_HeroAttr[tmp.ID] = tmp
	}
	return true
}

func ImpactTableRow(impactId int) (TableRowImpact, error) {
	if impact, ok := G_Impact[impactId]; ok {
		return impact, nil
	} else {
		return TableRowImpact{}, errors.New("config err")
	}
}
func SkillTableRow(skillId int) (TableRowSkill, error) {
	if skill, ok := G_Skill[skillId]; ok {
		return skill, nil
	} else {
		return TableRowSkill{}, errors.New("config err")
	}
}
func HeroTableRow(id int) (TableRowHeroAttr, error) {
	if heroAttr, ok := G_HeroAttr[id]; ok {
		return heroAttr, nil
	} else {
		return TableRowHeroAttr{}, errors.New("config err")
	}
}
