/*
@author: sk
@date: 2022/12/25
*/
package main

import "github.com/hajimehoshi/ebiten/v2"

type PlayerData struct {
	Name, Career             string
	Image                    string          // 地图显示图片
	MaxHp, MaxMp, Atk, Def   float64         // MaxHp,MaxMp 即是最大值也是初始值   Atk ,Def 初始值
	MaxEnergy, RecoverEnergy float64         // MaxEnergy 最大体力值，且是初始体力值  用于移动消耗的
	EnergyCost               map[int]float64 // 对各种地形的消耗 不同
	Skills, Equips           []string        // 技能 与 装备  对应实现  有主动与被动
	AttackRange              []*OffsetPos    //攻击的 偏移 范围
}

type OffsetPos struct {
	OffsetX, OffsetY int
}

func NewOffsetPos(offsetX int, offsetY int) *OffsetPos {
	return &OffsetPos{OffsetX: offsetX, OffsetY: offsetY}
}

type TileData struct {
	Image        *ebiten.Image
	HoverHandler PlayerHandler
	CommandMap   map[string]CommandFactory
}

type EquipmentData struct {
	Name                             string
	Image                            string
	Skills                           []string
	HpBuff, MpBuff, AtkBuff, DefBuff float64
}
