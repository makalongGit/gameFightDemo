package GameFight

//技能
type TableRowSkill struct {
	SkillId          int                      //技能ID
	Description      string                   //策划注释
	Name             string                   //技能名称
	SkillLevel       int                      //技能等级
	Tips             string                   //技能描述
	SelectTargetOpt  int                      //目标选择方式
	SkillRate        int                      //技能释放概率
	NeedMP           int                      //法力消耗
	CoolDownTime     int                      //技能内置冷却
	CoolDownID       int                      //技能冷却ID
	StartRound       int                      //第一次施放
	SkillType        int                      //技能类型
	ImpactID         [MaxSkillImpactCount]int //附加Impact
	ImpactRate       [MaxSkillImpactCount]int //附加Impact概率
	ImpactTargetType [MaxSkillImpactCount]int //impact选择类型
	ScriptID         int                      //脚本ID
	HeroLevel        int                      //英雄等级
	ConsumeMoney     int                      //银币消耗
	LearnType        int                      //学习类型  0=走commonitem 1=走道具消耗
	//战斗需要才加
}

type TableRowImpact struct {
	ImpactID      int                           //impactID
	Description   string                        //策划描述
	LogicID       int                           //impact逻辑Id
	Param         [MaxImpactLogicParamCount]int //逻辑参数
	Icon          string                        //Buff图标
	ImpactMutexID int                           //Impact互斥ID
	ReplaceLevel  int                           //顶替优先级
	DeadDisappear int                           //死亡后是否消失
	OfflineTimeGO int                           //下线是否计时
	ScriptID      int                           //脚本ID
	Effect        string                        //特效
	Name          string                        //名称
	Desc          string                        //描述
	SkillEffect   string                        //技能效果
}

type TableRowHeroAttr struct {
	ID                        int    //英雄ID
	SpiritID                  int    //英雄魂魄ID
	Name                      string //英雄名称
	Profession                int    //英雄职业
	InitQuality               int    //初始品质
	InitLevel                 int    //初始等级
	LevelLimit                int    //等级上限
	LevelCrossRole            int    //英雄比人物高出等级上限
	InitExp                   int    //初始经验
	TakeLevel                 int    //可携带等级
	EffectAttackByGrowRate    int    //成长值对英雄攻击点数的影响系数
	EffectDefendByGrowRate    int    //成长值对英雄防御点数的影响系数
	EffectHPByGrowRate        int    //成长值对英雄生命上限点数的影响系数
	EffectMPByGrowRate        int    //成长值对英雄魔法上限点数的影响系数
	EffectPhysicAttackByLevel int    //英雄物理攻击随等级的增长系数
	EffectPhysicDefendByLevel int    //英雄物理防御随等级的增长系数
	EffectHpByLevel           int    //英雄生命上限随等级的增长系数
	EffectMpByLevel           int    //英雄魔法上限随等级的增长系数
	InitAttackSpeed           int    //初始速度
	InitPhysicAttack          int    //初始物理攻击
	InitMagicAttack           int    //初始魔法攻击
	InitPhysicDefence         int    //初始物理防御
	InitMagicDefence          int    //初始魔法防御
	InitHP                    int    //初始生命值
	InitMP                    int    //初始魔法值
	InitHit                   int    //初始命中值
	InitDodge                 int    //初始闪避值
	InitStrike                int    //初始暴击值
	InitContinuous            int    //初始连击值	计算发生概率
	InitBackAttack            int    //初始反击值
	InitStrikeHurt            int    //初始暴击伤害
	InitConAttTimes           int    //初始连击次数
	InitConAttHurt            int    //初始连击伤害	计算伤害倍率
	InitBackAttHurt           int    //初始反击伤害
	InitPhysicHurtDecay       int    //初始物理减免
	InitMagicHurtDecay        int    //初始魔法减免
	FloatingHurt              int    //英雄伤害浮动
	PhysicSkillID             int    //普通攻击技能
	MagicSkillID1             int    //魔法技能1
	LearnedLevel1             int    //学会等级1
	MagicSkillID2             int    //魔法技能2
	LearnedLevel2             int    //学会等级2
	MaxGrowPoint              int    //成长值上限
	InitBearPoint             int    //初始承载力
	BearParam                 int    //承载力随等级的增长系数
	RequireHumanLevel         int    //需要召唤师等级
}
