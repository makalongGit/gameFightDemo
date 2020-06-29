package GameFight

import (
	"fmt"
	"testing"
)

func TestReadImpactCsv(t *testing.T) {
	ReadImpactCsv("impact.tab.csv")
	fmt.Printf("%+v", G_Impact[1])
}

func TestReadSkillCsv(t *testing.T) {
	ReadSkillCsv("Skill.csv")
	fmt.Printf("%+v", G_Skill[100005])
}

func TestReadHeroAttrCsv(t *testing.T) {
	ReadHeroAttrCsv("HeroAttr.csv")
	fmt.Printf("%+v", G_HeroAttr[1])
}