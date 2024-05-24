/*
@author: sk
@date: 2022/12/25
*/
package main

import (
	"GameBase2/utils"
	R "slg_game/res"
)

type SkillFactory func(*Player) ISkill

var (
	skillFactories = make(map[string]SkillFactory)
)

func init() {
	skillFactories["治疗"] = createTreatmentSkill
	skillFactories["火焰"] = createFireSkill
	skillFactories["驱散"] = createDisperseSkill
	// 物品效果
	skillFactories["兴奋剂"] = createAnalepticPillSkill
}

func CreateSkill(name string, player *Player) ISkill {
	return skillFactories[name](player)
}

//========================主动效果=========================

//********************DisperseSkill*********************

type DisperseSkill struct {
	player *Player
}

func (d *DisperseSkill) GetName() string {
	return "驱散"
}

func (d *DisperseSkill) GetCommand() ICommand {
	return NewSelectCommand(d.handlePlayer)
}

func (d *DisperseSkill) GetCost() float64 {
	return 10
}

func (d *DisperseSkill) handlePlayer(player *Player) string {
	if player == nil {
		return "必须选择目标!"
	}
	if d.player.IsEnemy == player.IsEnemy {
		player.RemoveDeBuff(nil)
		utils.AddToLayer(R.LAYER.FG, NewSimpleEffect(R.OTHER.EFFECT.RECOVERY, "", player.Pos, nil))
	} else {
		player.RemoveBuff(nil)
		utils.AddToLayer(R.LAYER.FG, NewSimpleEffect(R.OTHER.EFFECT.FIRE, "", player.Pos, nil))
	}
	d.player.MarkAction()
	return ""
}

func createDisperseSkill(player *Player) ISkill {
	return &DisperseSkill{player: player}
}

//*******************FireSkill************************

type FireSkill struct {
	player *Player
}

func (f *FireSkill) GetCost() float64 {
	return 30
}

func (f *FireSkill) GetName() string {
	return "火焰"
}

func (f *FireSkill) GetCommand() ICommand {
	pos1 := make([]*OffsetPos, 0)
	pos2 := make([]*OffsetPos, 0)
	dirs := [][]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for i := 0; i < 4; i++ {
		pos1 = append(pos1, NewOffsetPos(dirs[i][0], dirs[i][1]))
		pos2 = append(pos2, NewOffsetPos(dirs[i][0]*2, dirs[i][1]*2))
		pos2 = append(pos2, NewOffsetPos(dirs[i][0]*3, dirs[i][1]*3))
		pos2 = append(pos2, NewOffsetPos(dirs[i][0]*4, dirs[i][1]*4))
	}
	helpers := make([]*RangeHelper, 0)
	helpers = append(helpers, NewRangeHelper(f.player.X, f.player.Y, pos1, Color255_0_0_127))
	helpers = append(helpers, NewRangeHelper(f.player.X, f.player.Y, pos2, Color255_255_0_127))
	return NewSelectPosCommand(f.handlePos, helpers)
}

func (f *FireSkill) handlePos(x int, y int) string {
	offsetX := x - f.player.X
	offsetY := y - f.player.Y
	for i := 0; i < 4; i++ {
		player := PlayerManager.GetTypePlayer(x, y, !f.player.IsEnemy)
		if player != nil {
			player.Hurt(40, f.player)
		}
		utils.AddToLayer(R.LAYER.FG, NewSimpleEffect(R.OTHER.EFFECT.FIRE, "", Index2Pos(x, y), nil))
		x += offsetX
		y += offsetY
	}
	f.player.MarkAction()
	return ""
}

func createFireSkill(player *Player) ISkill {
	return &FireSkill{player: player}
}

//**********************TreatmentSkill**********************

type TreatmentSkill struct {
	player *Player
}

func (t *TreatmentSkill) GetCost() float64 {
	return 20
}

func (t *TreatmentSkill) GetName() string {
	return "治疗"
}

func (t *TreatmentSkill) GetCommand() ICommand {
	return NewSelectCommand(t.handlePlayer)
}

func (t *TreatmentSkill) handlePlayer(player *Player) string {
	if player == nil {
		return "必须选择目标!"
	}
	if player.IsEnemy != t.player.IsEnemy {
		return "必须选择友方角色!"
	}
	player.RecoverHp(40, t.player)
	t.player.MarkAction()
	return ""
}

func createTreatmentSkill(player *Player) ISkill {
	return &TreatmentSkill{player: player}
}

//=====================被动效果=========================
//******************AtkBuffSkill******************

type AtkBuffSkill struct {
	player  *Player
	round   int
	atkBuff float64
}

func (a *AtkBuffSkill) RoundBegan() {
	a.round--
	if a.round < 0 {
		a.player.RemoveBuff(a) // 超出回合 要移除的
	}
}

func (a *AtkBuffSkill) GetAtkBuff() float64 {
	return a.atkBuff
}

func createAnalepticPillSkill(player *Player) ISkill { // 一个技能可以通过参数对接 对个效果
	return &AtkBuffSkill{player: player, round: 3, atkBuff: 10}
}
