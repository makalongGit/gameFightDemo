package GameFight

type EmTypeFight int

const (
	EmTypeFightNormal = iota
	EmTypeFightStair  //单排
	EmTypeFightCount
)

const (
	InvalidId    = -1
	InvalidValue = -1

	//技能
	MaxSkillNum int = 2 //技能数量
	//角色相关
	MaxEquipNumPerHero int = 6 //英雄携带最大装备数量
	//技能
	MaxMatrixCellCount       int = 6   //技能最大选择目标
	MaxImpactNumber          int = 20  //impact数量A
	MaxSkillImpactCount      int = 4   //最大技能所带impact
	MaxConAttackTimes        int = 4   //最大连击次数
	MaxImpactLogicParamCount int = 4   //impact参数数量
	MaxBuffNumber            int = 10  //强化类buff数量
	MaxDeBuffNumber          int = 10  //削弱类buff数量
	MAX_FIGHT_ROUND          int = 128 //最大战斗回合
	Distance int = 10

)

type EmSkillType int

const (
	EmSkillTypeInvalid      EmSkillType = iota - 1
	EmSkillTypeHeroActive               //0=英雄主动技能
	EmSkillTypeHeroPassive              //1=英雄被动技能
	EmSkillTypeEquipActive              //2=装备主动技能
	EmSkillTypeEquipPassive             //3=装备被动技能
);

type EmImpactLogic int

const (
	EmImpactLogic0 = iota
	EmImpactLogic1
	EmImpactLogic2 //持续物伤
	EmImpactLogic3 //持续法伤
	EmImpactLogic4
	EmImpactLogic5
	EmImpactLogic6

	EmImpactLogicCount
)

//技能选择目标方式
type EmSkillTargetOpt int

const (
	EmSkillTargetOptAuto  = iota - 1 //自动选择
	EmSkillTargetOptOrder = 0        //0=顺序选择
	EmSkillTargetOptRand  = 1        //1=随机选择
	EmSkillTargetOptSlow  = 2        //2=后排优先
	EmSkillTargetOptNumber
)

//impact选择目标方式
type EmImpactTarget int

const (
	EmImpactTargetOptAuto        = iota - 1
	EmImpactTargetOptSelf        //0=自身；
	EmImpactTargetOwnerSignle    //1=己方个体；
	EmImpactTargetOwnerAll       //2=己方全体；
	EmImpactTargetEnemySignle    //3=敌方个体；
	EmImpactTargetEnemyFront     //4=敌方横排；
	EmImpactTargetEnemyBehind    //5=敌方后排优先；
	EmImpactTargetEnemyAll       //6=敌方全体；
	EmImpactTargetEnemyLine      //7=敌方目标竖排；
	EmImpactTargetEnemyAround    //8=敌方目标及周围；
	EmImpactTargetEnemyBehinDone //9=敌方后排个体
	EmImpactTargetOwnerMinHp     //10=己方血最少
	EmImpactTargetOwnerMinMp     //11=己方蓝最少
	EmImpactTargetOptNumber
)

type EmTypeImpactLogic int

const (
	EmTypeImpactLogicInvalid = iota - 1
	EmTypeImpactLogicSingle  //单次生效
	EmTypeImpactLogicBuff    //强化
	EmTypeImpactLogicDeBuff  //削弱

	EmTypeImpactLogicCount
)

type EmAttribute int

const (
	EmAttributeInvalid              EmAttribute = iota
	EmAttributeMaxHp                            //最大生命
	EmAttributeMoveSpeed                        //移动速度
	EmAttributeAttackSpeed                      //攻击速度
	EmAttributePhysicAttack                     //物理攻击
	EmAttributePhysicDefence                    //物理防御
	EmAttributeHit                              //命中点数
	EmAttributeDodge                            //闪避点数
	EmAttributeStrike                           //暴击
	EmAttributeContinuous                       //连击
	EmAttributeBackAttack                       //反击
	EmAttributeContinuousTimes                  //连击次数
	EmAttributeHurtContinuous                   //连击伤害
	EmAttributeHurtBackAttack                   //反击伤害
	EmAttributeHurtStrike                       //暴击伤害
	EmAttributePhysicHurtDecay                  //物理伤害减免
	EmAttributeStrikeHurtDecay                  //暴击伤害减免
	EmAttributeHurtExtra                        //附加伤害
	EmAttributeHurtPhysic                       //普通攻击伤害
	EmAttributeMagicAttack                      //魔法攻击
	EmAttributeMagicDefence                     //魔法防御
	EmAttributeMagicHurtDecay                   //魔法伤害减免
	EmAttributeMaxMp                            //最大魔法值
	EmAttributePercentAttackSpeed               //攻击速度百分比
	EmAttributePercentPhysicAttack              //物理攻击百分比
	EmAttributePercentMagicAttack               //魔法攻击百分比
	EmAttributePercentPhysicDefence             //物理防御百分比
	EmAttributePercentMagicDefence              //魔法防御百分比
	EmAttributePercentMaxHp                     //最大生命值百分比
	EmAttributePercentMaxMp                     //最大魔法值百分比
	EmAttributeLevel                            //等级
	EmAttributeHp                               //血量
	EmAttributeMp                               //魔法值

	EmAttributeCurrentExp //当前经验
	EmAttributeAction     //行动力
	EmAttributeNumber
)

//impact结果
type EmImpactResult int

const (
	EmImpactResultNormal    = iota //正常加目标身上
	EmImpactResultFail             //不能加在目标身上
	EmImpactResultDisAppear        //抵消
)

//属性值计算公式
func ChkMin(A, B int) (min int) {
	if A < B {
		return B
	} else {
		return A
	}
}

func CalcDamage(A, B, C int) int {
	return int((A - B) * (1 - (C)/100))
}

func CalcAttr1(A, B, C, D, E, F, G int) int {
	return int((A+B*C/100)*(1+D/100+E/100) + F + G)
}

func CalcAttr2(A, B, C int) int {
	return A + B + C
}
func CalcAttr3(A, B, C int) int {
	return 1 - (1-A/100)*(1-B/100)*(1-C/100)*100
}
func GetSkillGroup(A int) int {
	return A / 100
}

func GetSkillLevel(A int) int {
	return A % 100
}
