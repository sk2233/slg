/*
@author: sk
@date: 2022/12/25
*/
package main

import (
	"GameBase2/config"
	"GameBase2/factory"
	"GameBase2/model"
	"GameBase2/object"
	"GameBase2/utils"
	"fmt"
	"math"
	R "slg_game/res"

	"golang.org/x/image/colornames"
)

func init() {
	config.ObjectFactory.RegisterPointFactory(R.CLASS.PLAYER, createPlayer)
}

func createPlayer(data *model.ObjectData) model.ObjectAble {
	playerData := DataManager.GetPlayerData(data.Name)
	res := &Player{Data: playerData, Rank: 0, skills: make([]ISkill, 0), buffs: make([]ISkill, 0),
		deBuffs: make([]ISkill, 0), Equips: make([]*Equipment, 0)}
	res.DrawObject = object.NewDrawObject()
	factory.FillPointObject(data, res.PointObject)
	res.IsEnemy = res.GetBool(R.PROP.ENEMY, false)
	res.X, res.Y = Pos2Index(res.Pos)
	res.BindSprite(utils.LoadSprite(R.SPRITE.OTHER, playerData.Image), res)
	res.applyData()
	return res
}

//=================Player(一个棋子)==================

type Player struct {
	*object.DrawObject
	Data           *PlayerData
	skills         []ISkill
	Equips         []*Equipment
	buffs          []ISkill
	deBuffs        []ISkill
	IsEnemy        bool
	Action         bool // 是否行动过
	Hp, Mp, Energy float64
	X, Y           int
	Rank           int
}

func (p *Player) Order() int {
	return p.Rank
}

func (p *Player) CanSelect(x int, y int) bool {
	return !p.Action && x == p.X && y == p.Y
}

func (p *Player) GetMainMenu() []string {
	arr := []string{"移动", "攻击", "技能", "物品", "结束"}
	tileData := TileManager.GetTileInfo(p.X, p.Y)
	for name := range tileData.CommandMap {
		arr = append(arr, name) // 添加地图位置支持的特殊操作
	}
	return arr
}

func (p *Player) GetSubMenu(name string) []string { // 子菜单
	switch name {
	case "技能":
		skillNames := make([]string, 0)
		skills := p.GetSkills()
		for i := 0; i < len(skills); i++ {
			if initiativeSkill, ok := skills[i].(IInitiativeSkill); ok {
				skillNames = append(skillNames, initiativeSkill.GetName())
			}
		}
		return skillNames
	case "物品":
		items := DataManager.GetPlayerItems(p.IsEnemy)
		arr := make([]string, 0)
		for i := 0; i < len(items); i++ {
			arr = append(arr, fmt.Sprintf("%s(*%d)", items[i].GetName(), items[i].GetCount()))
		}
		return arr
	}
	return nil // 返回nil 会触发指令执行  若是确实没有东西 返回 []string{}
}

// 直接行动   技能     物品
func (p *Player) StartCommand(name string, mainIndex, subIndex int) {
	// 直接行动
	if p.mainCommand(name) {
		return
	}
	// 使用技能
	if p.skillCommand(name) {
		return
	}
	if p.itemCommand(name, subIndex) {
		return
	} // 特殊地理 可执行的操作
	if p.tileCommand(name) {
		return
	}
	utils.LogErr("StartCommand err name %s", name)
}

func (p *Player) tileCommand(name string) bool {
	tileData := TileManager.GetTileInfo(p.X, p.Y)
	if commandFactory, ok := tileData.CommandMap[name]; ok {
		Cursor.SetCommand(commandFactory(p))
		return true
	}
	return false
}

func (p *Player) mainCommand(name string) bool {
	switch name {
	case "移动":
		Cursor.SetCommand(NewMoveCommand(p))
		return true
	case "攻击":
		Cursor.SetCommand(NewAttackCommand(p))
		return true
	case "结束":
		p.MarkAction()
		return true
	default:
		return false
	}
}

func (p *Player) skillCommand(name string) bool {
	skills := p.GetSkills()
	for i := 0; i < len(skills); i++ {
		if initiativeSkill, ok := skills[i].(IInitiativeSkill); ok {
			if initiativeSkill.GetName() == name {
				if p.Mp < initiativeSkill.GetCost() {
					Tip.SetMsg("MP不足!")
					return true
				}
				p.Mp -= initiativeSkill.GetCost()
				command := initiativeSkill.GetCommand()
				Cursor.SetCommand(command)
				return true
			}
		}
	}
	return false
}

func (p *Player) itemCommand(name string, index int) bool {
	res, ok := DataManager.UsePlayerItem(p, name, index)
	if !ok {
		return false
	}
	Cursor.SetCommand(res)
	return true
}

func (p *Player) MarkAction() {
	p.Action = true
	PlayerManager.UpdateAction()
}

func (p *Player) RoundBegan() {
	p.Action = false
	p.Energy += p.Data.RecoverEnergy // 恢复体力
	if p.Energy > p.Data.MaxEnergy {
		p.Energy = p.Data.MaxEnergy
	}
	skills := p.GetSkills()
	for i := 0; i < len(skills); i++ {
		InvokeRoundBegan(skills[i])
	}
}

func (p *Player) GetEnergyCost(tileType int) float64 {
	return p.Data.EnergyCost[tileType]
}

func (p *Player) SetPos(x int, y int) {
	p.X, p.Y = x, y
	p.Pos = Index2Pos(x, y)
}

func (p *Player) Attack(target *Player) {
	atk := p.GetAtk()
	target.Hurt(atk, p)
}

func (p *Player) GetAtk() float64 {
	atk := p.Data.Atk
	for i := 0; i < len(p.Equips); i++ {
		atk += p.Equips[i].Data.AtkBuff
	}
	skills := p.GetSkills()
	for i := 0; i < len(skills); i++ {
		atk += InvokeGetAtkBuff(skills[i])
	}
	return atk
}

func (p *Player) GetDef() float64 {
	def := p.Data.Def
	for i := 0; i < len(p.Equips); i++ {
		def += p.Equips[i].Data.DefBuff
	}
	return def
}

func (p *Player) GetMaxHp() float64 {
	maxHp := p.Data.MaxHp
	for i := 0; i < len(p.Equips); i++ {
		maxHp += p.Equips[i].Data.HpBuff
	}
	return maxHp
}

func (p *Player) GetMaxMp() float64 {
	maxMp := p.Data.MaxMp
	for i := 0; i < len(p.Equips); i++ {
		maxMp += p.Equips[i].Data.MpBuff
	}
	return maxMp
}

func (p *Player) Hurt(atk float64, source *Player) {
	def := p.GetDef()
	hurt := math.Max(atk-def, 0)
	p.Hp -= hurt
	utils.AddToLayer(R.LAYER.FG, NewSimpleEffect(R.OTHER.EFFECT.ATTACK, fmt.Sprintf("-%.1f", hurt), p.Pos, colornames.Red))
}

func (p *Player) RecoverHp(value float64, source *Player) {
	utils.AddToLayer(R.LAYER.FG, NewSimpleEffect(R.OTHER.EFFECT.RECOVERY, fmt.Sprintf("+%.1f", value), p.Pos, colornames.Green))
	p.Hp = math.Min(p.Hp+value, p.GetMaxHp())
}

func (p *Player) RecoverMp(value float64, source *Player) {
	utils.AddToLayer(R.LAYER.FG, NewSimpleEffect(R.OTHER.EFFECT.RECOVERY, fmt.Sprintf("+%.1f", value), p.Pos, colornames.Blue))
	p.Mp = math.Min(p.Mp+value, p.GetMaxMp())
}

func (p *Player) initSkill() {
	for i := 0; i < len(p.Data.Skills); i++ {
		p.skills = append(p.skills, CreateSkill(p.Data.Skills[i], p))
	}
}

func (p *Player) applyData() {
	p.initSkill()
	p.initEquip()
	p.Hp = p.GetMaxHp()
	p.Mp = p.GetMaxMp()
	p.Energy = p.Data.MaxEnergy
}

func (p *Player) initEquip() {
	for i := 0; i < len(p.Data.Equips); i++ {
		p.Equips = append(p.Equips, NewEquipment(p.Data.Equips[i], p))
	}
}

func (p *Player) GetSkills() []ISkill {
	res := make([]ISkill, 0)
	res = append(res, p.skills...) // 必须解包
	res = append(res, p.buffs...)  // 正面效果 与 负面效果 也是 技能一部分
	res = append(res, p.deBuffs...)
	for i := 0; i < len(p.Equips); i++ {
		res = append(res, p.Equips[i].Skills...)
	}
	return res
}

func (p *Player) AddBuff(buff ISkill) {
	p.buffs = append(p.buffs, buff)
}

func (p *Player) RemoveBuff(buff ISkill) {
	if buff == nil { // 不指定 全部移除
		p.buffs = make([]ISkill, 0)
		return
	}
	for i := 0; i < len(p.buffs); i++ {
		if p.buffs[i] == buff {
			p.buffs = append(p.buffs[:i], p.buffs[i+1:]...)
			return
		}
	}
	utils.LogErr("RemoveBuff没有移除成功")
}

func (p *Player) RemoveDeBuff(deBuff ISkill) {
	if deBuff == nil { // 不指定 全部移除
		p.deBuffs = make([]ISkill, 0)
		return
	}
	for i := 0; i < len(p.deBuffs); i++ {
		if p.deBuffs[i] == deBuff {
			p.deBuffs = append(p.deBuffs[:i], p.deBuffs[i+1:]...)
			return
		}
	}
	utils.LogErr("RemoveDeBuff没有移除成功")
}
