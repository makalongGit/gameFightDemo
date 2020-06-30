package GameFight

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type FightDBData struct {
	Guid            X_GUID                  //对象guid
	TableID         int                     //表格id
	Type            int                     //英雄出处
	Quality         int                     //品质
	MatrixID        int                     //阵法Id
	Profession      int                     //职业 0近战 1远程
	Level           int                     //等级
	HP              int                     //血量
	Mp              int                     //魔法值
	SkillCount      int                     //技能数量
	Skill           [MaxSkillNum]int        //技能列表
	EquipSkillCount int                     //装备技能数量
	EquipSkill      [MaxEquipNumPerHero]int //装备技能李彪

	//战斗属性
	PhysicAttack    int //物理攻击
	MagicAttack     int //魔法攻击
	PhysicDefence   int //物理防御
	MagicDefence    int //魔法防御
	MaxHP           int //最大生命
	MaxMP           int //最大魔法值
	Hit             int //命中值
	Dodge           int //闪避值
	Strike          int //暴击值
	StrikeHurt      int //暴击伤害
	Continuous      int //连击值
	ConAttHurt      int //连击伤害
	ConAttTimes     int //连击次数
	BackAttack      int //反击值
	BackAttHurt     int //反击伤害
	AttackSpeed     int //攻击速度
	PhysicHurtDecay int //物理减免
	MagicHurtDecay  int //魔法减免

	Exp       int                     //经验值
	GrowRate  int                     //成长
	BearPoint int                     //负重
	Equip     [MaxEquipNumPerHero]int //装备id
	Color     int                     //颜色
}

func (f *FightDBData) CleanUp() {
	f.Guid = InvalidId
	for i := 0; i < MaxEquipNumPerHero; i++ {
		f.Equip[i] = InvalidId
	}
}

func (f *FightDBData) IsValid() bool {
	if f.Guid == 0 {
		return false
	}
	return f.Guid.IsValid()
}

type FightDB struct {
	HumanGuid   X_GUID //玩家Guid
	FightCount  int
	FightDBData [MaxMatrixCellCount]FightDBData
}

func (f *FightDB) CleanUp() {
	f.HumanGuid = InvalidId
	f.FightCount = 0
}

func (f *FightDB) IsValid() bool {
	return f.HumanGuid.IsValid()
}

func (f *FightDB) AddFightDBData(data FightDBData) {
	if f.FightCount >= 0 && f.FightCount < MaxMatrixCellCount {
		f.FightDBData[f.FightCount] = data
		f.FightCount += 1
	}
}

//impact信息
type ImpactInfo struct {
	SkillID     int
	ImpactID    int
	TargetList  [MaxMatrixCellCount]X_GUID
	TargetCount int
	ConAttTimes int
	Hurts       [MaxMatrixCellCount][MaxConAttackTimes]int
	Mp          [MaxMatrixCellCount]int //本次技能带给的魔法量变化 >0 减蓝 < 0 加蓝
}

func (impact ImpactInfo) GetTargetIndex(targetGuid X_GUID) int {
	if !targetGuid.IsValid() {
		return InvalidId
	}
	for i := 0; i < impact.TargetCount; i++ {
		if impact.TargetList[i] == targetGuid {
			return i
		}
	}
	return InvalidId
}

type AttributeEffect struct {
	AttrValue int
}

//技能信息
type SkillAttack struct {
	SkillID     int                              //技能ID
	SkillTarget X_GUID                           //技能目标
	CostMp      int                              //消耗魔法量
	Impact      [MaxSkillImpactCount]*ImpactInfo //impact列数
	ImpactCount int                              //impact个数
}

func (s *SkillAttack) String() string {
	b, err := json.Marshal(*s)
	if err != nil {
		return fmt.Sprintf("%+v", *s)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *s)
	}
	return out.String()
}
func (s *SkillAttack) AddImpactInfo(info *ImpactInfo) {
	s.Impact[s.ImpactCount] = info
	s.ImpactCount += 1
}

func (s SkillAttack) GetImpactInfo(impactId int) *ImpactInfo {
	for i := 0; i < MaxSkillImpactCount; i++ {
		if s.Impact[i].ImpactID == impactId {
			return s.Impact[i]
		}
	}
	return nil
}
func (s SkillAttack) IsValid() bool {
	return s.SkillID > 0
}

//出手信息
type AttackInfo struct {
	CastGuid X_GUID //对象guid
	Skilled  bool   //魔法攻击还是普通攻击
	//魔法攻击
	SkillAttack      [MaxEquipNumPerHero]*SkillAttack //如果是普通攻击先进行装备技能攻击
	SkillAttackCount int
	//普通攻击
	SkillTarget    X_GUID //普通攻击目标
	BHit           bool   //是否命中
	BStrike        bool   //是否暴击
	Hurt           int    //伤害值
	BBackAttack    bool   //是否有反击
	BackAttackHurt int    //反击伤害
}

func (c AttackInfo) String() string {
	return fmt.Sprintf("攻击方英雄id: %+v  防守方英雄id: %+v, \n"+
		"是否使用了技能攻击 %+v,"+
		"是否命中:%+v,"+
		"是否暴击:%+v,\n"+
		"攻击伤害:%+v, "+
		"是否有反击:%+v, "+
		"反击伤害:%+v\n",
		c.CastGuid, c.SkillTarget, c.Skilled, c.BHit, c.BStrike, c.Hurt, c.BBackAttack, c.BackAttackHurt)
}

func (c *AttackInfo) GetSkillAttack(SkillID int) *SkillAttack {
	for i := 0; i < c.SkillAttackCount; i++ {
		if c.SkillAttack[i].SkillID == SkillID {
			return c.SkillAttack[i]
		}
	}
	return nil
}
func (c AttackInfo) IsValid() bool {
	return c.CastGuid > 0
}

func (c *AttackInfo) CleanUp() {
	if c != nil {
		c.SkillAttackCount = 0
		c.CastGuid = InvalidId
		c.SkillTarget = InvalidId
		c.BHit = true
	}
}

func (c *AttackInfo) AllocSkillAttack() *SkillAttack {
	if c.SkillAttackCount >= MaxEquipNumPerHero {
		return nil
	}
	index := c.SkillAttackCount
	c.SkillAttackCount += 1
	return c.SkillAttack[index]
}

//下发战斗信息
type FObjectData struct {
	Guid       X_GUID //对象guid
	TableID    int    //英雄表格id
	Quality    int    //品质
	Color      int    //颜色
	Profession int    //职业
	Level      int    //等级
	MatrixID   int    //位置
}

func (f *FObjectData) CleanUp() {
	f.Guid = InvalidId
}

type FObjectInfo struct {
	Guid          X_GUID //对象guid
	HP            int    //血量
	MaxHP         int    //最大血量
	MP            int    //魔法值
	MaxMP         int    //最大魔法
	FightDistance int    //战斗条长度
	AttackSpeed   int    //速度
	EndDistance   int    //最后位置
	ImpactCount   int
	ImpactList    [MaxImpactNumber]int //身上impact
	ImpactHurt    [MaxImpactNumber]int //持续impact伤害 + 掉血 - 加血
	ImpactMP      [MaxImpactNumber]int //持续impact 蓝 + 掉蓝 - 加蓝
}

func (f *FObjectInfo) String() string {
	return fmt.Sprintf("英雄id: %+v 信息\n"+
		"HP:%+v MP:%+v MaxHP:%+v MaxMP:%+v \n"+
		"战斗条长度:%+v 速度:%+v 最后战斗条位置:%+v impact列表: %+v \n",
		f.Guid, f.HP, f.MP, f.MaxHP, f.MaxMP, f.FightDistance, f.AttackSpeed, f.EndDistance, f.ImpactList)
}

func (f *FObjectInfo) CleanUp() {
	f.Guid = InvalidId
}

func (f FObjectInfo) AddImpact(impactId int, hurt int, mp int) {
	f.ImpactList[f.ImpactCount] = impactId
	f.ImpactHurt[f.ImpactCount] = hurt
	f.ImpactMP[f.ImpactCount] = mp
	f.ImpactCount += 1
}

type FightRoundInfo struct {
	AttackObjectInfo  [MaxMatrixCellCount]*FObjectInfo //本回合开始前状态数据
	AttackObjectCount int

	DefendObjectInfo  [MaxMatrixCellCount]*FObjectInfo
	DefendObjectCount int

	AttackInfo   [MaxMatrixCellCount * 2]AttackInfo //出手数据
	AttInfoCount int
}

func (f *FightRoundInfo) String() string {

	var s string = "攻击方战前信息:\n"
	for i := 0; i < f.AttackObjectCount; i++ {
		s += f.AttackObjectInfo[i].String() + ""
	}
	s += "防守方战前信息:\n"
	for i := 0; i < f.DefendObjectCount; i++ {
		s += f.DefendObjectInfo[i].String() + ""
	}
	s += "攻击信息\n"
	for i := 0; i < f.AttInfoCount; i++ {
		s += f.AttackInfo[i].String() + ""
	}
	return s
}

func (f *FightRoundInfo) CleanUp() {
	for i := 0; i < MaxMatrixCellCount; i++ {
		if f.AttackObjectInfo[i] != nil && f.DefendObjectInfo[i] != nil {
			f.AttackObjectInfo[i].CleanUp()
			f.DefendObjectInfo[i].CleanUp()
			f.AttackInfo[i].CleanUp()
			f.AttackInfo[i+MaxMatrixCellCount].CleanUp()
		}
	}
	f.AttackObjectCount = 0
	f.DefendObjectCount = 0
	f.AttInfoCount = 0
}

func (f *FightRoundInfo) AddAttackObjectInfo(objInfo *FObjectInfo) {
	f.AttackObjectInfo[f.AttackObjectCount] = objInfo
	f.AttackObjectCount += 1
}

func (f *FightRoundInfo) AddDefendObjectInfo(objInfo *FObjectInfo) {
	f.DefendObjectInfo[f.DefendObjectCount] = objInfo
	f.DefendObjectCount += 1
}

func (f *FightRoundInfo) AddAttackInfo(attackInfo AttackInfo) {
	f.AttackInfo[f.AttInfoCount] = attackInfo
	f.AttInfoCount++
}

func (f *FightRoundInfo) GetFObjectInfoByGuid(guid X_GUID) *FObjectInfo {
	if !guid.IsValid() {
		return nil
	}
	for i := 0; i < f.AttackObjectCount; i++ {
		if f.AttackObjectInfo[i].Guid == guid {
			return f.AttackObjectInfo[i]
		}
	}
	for i := 0; i < f.DefendObjectCount; i++ {
		if f.DefendObjectInfo[i].Guid == guid {
			return f.DefendObjectInfo[i]
		}
	}
	return nil
}

type FightInfo struct {
	AttackObjectData  [MaxMatrixCellCount]FObjectData //攻击方
	AttackObjCount    int                             //攻击方数量
	DefendObjectData  [MaxMatrixCellCount]FObjectData //防守方
	DefendObjectCount int                             //防守方数量
	DefendType        int                             //防守方类型 0 怪物 1人

	RoundInfo [MAX_FIGHT_ROUND]FightRoundInfo //每回合的战斗信息
	Rounds    int                             //总回合

	MaxFightDistance int  //战斗条长度
	BWin             bool //挑战者是否胜利
}

func (f *FightInfo) CleanUp() {
	for i := 0; i < MaxMatrixCellCount; i++ {
		f.AttackObjectData[i].CleanUp()
		f.DefendObjectData[i].CleanUp()
	}
	for i := 0; i < MAX_FIGHT_ROUND; i++ {
		f.RoundInfo[i].CleanUp()
	}
	f.Rounds = 0
	f.AttackObjCount = 0
	f.DefendObjectCount = 0
	f.MaxFightDistance = 0
	f.BWin = false
	f.DefendType = 0
}

func (f *FightInfo) AddAttackObjectData(data FObjectData) {
	f.AttackObjectData[f.AttackObjCount] = data
	f.AttackObjCount += 1
}

func (f *FightInfo) AddDefendObjectData(data FObjectData) {
	f.DefendObjectData[f.DefendObjectCount] = data
	f.DefendObjectCount += 1
}

func (f *FightInfo) AddRoundInfo(roundInfo FightRoundInfo) {
	f.RoundInfo[f.Rounds] = roundInfo
	f.Rounds++
}

func (f *FightInfo) SetWin(Win bool) {
	f.BWin = Win
}

func (f *FightInfo) GetTableIDByGuid(guid X_GUID) int {
	if guid.IsValid() {
		return InvalidId
	}
	for i := 0; i < f.AttackObjCount; i++ {
		if f.AttackObjectData[i].Guid == guid {
			return f.AttackObjectData[i].TableID
		}
	}
	for i := 0; i < f.DefendObjectCount; i++ {
		if f.DefendObjectData[i].Guid == guid {
			return f.DefendObjectData[i].TableID
		}
	}
	return InvalidId
}

type TargetList struct {
	pObjectList [MaxMatrixCellCount]*FightObject
	nCount      int
}

func (t *TargetList) CleanUp() {
	t.nCount = 0
}

func (t *TargetList) GetCount() int {
	return t.nCount
}

func (t *TargetList) GetFightObject(index int) *FightObject {
	if index >= 0 && index < t.nCount {
		return t.pObjectList[index]
	}
	return nil
}

func (t *TargetList) Add(pObject *FightObject) bool {
	if pObject == nil || t.nCount >= MaxMatrixCellCount {
		return false
	}
	t.pObjectList[t.nCount] = pObject
	t.nCount++
	return true
}
