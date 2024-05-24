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
	R "slg_game/res"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
)

func init() {
	config.ObjectFactory.RegisterPointFactory(R.CLASS.CURSOR, createCursor)
}

func createCursor(data *model.ObjectData) model.ObjectAble {
	res := &cursor{}
	res.DrawObject = object.NewDrawObject()
	factory.FillPointObject(data, res.PointObject)
	res.x, res.y = Pos2Index(res.Pos)
	res.BindSprite(utils.LoadSprite(R.SPRITE.OTHER, R.OTHER.OTHER.POINT), res)
	config.Camera.SetTarget(res)
	Cursor = res
	return res
}

//=================cursor====================

var (
	keys    = []ebiten.Key{ebiten.KeyW, ebiten.KeyS, ebiten.KeyA, ebiten.KeyD}
	offsets = [][]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
)

type cursor struct {
	*object.DrawObject
	command ICommand
	x, y    int
}

func (c *cursor) UIDraw(screen *ebiten.Image) {
	player := PlayerManager.GetPlayer(c.x, c.y) // 为了显示 边框 并未完全对齐
	utils.DrawRect(screen, 1+320i, 639+63i, Color1_0_98, Color201_163_66)
	tileInfo := TileManager.GetTileInfo(c.x, c.y)
	utils.DrawImage(screen, tileInfo.Image, 16+336i)
	if player != nil {
		c.drawPlayer(player, screen)
	} else {
		utils.DrawAnchorText(screen, "回合:"+strconv.Itoa(PlayerManager.GetRound()), 576+352i, 0.5+0.5i, Font36, colornames.White)
	}
}

func (c *cursor) SetCommand(command ICommand) {
	c.command = command // 一次仅能执行 一个指令
}

func (c *cursor) ClearCommand(command ICommand) {
	if command != nil && c.command != command { // 传入对象不为nil时必须对比 一致才移除 防止错误移除 否则直接移除
		return
	}
	c.command = nil
}

func (c *cursor) Update() {
	c.DrawObject.Update()
	if c.command != nil && c.command.Block() {
		return // 命令 阻塞  防止  玩家乱操作
	}
	c.move()
	c.enter()
}

func (c *cursor) move() {
	for i := 0; i < 4; i++ {
		if inpututil.IsKeyJustPressed(keys[i]) {
			c.x = utils.Clamp(c.x+offsets[i][0], 0, 39)
			c.y = utils.Clamp(c.y+offsets[i][1], 0, 21)
			c.Pos = Index2Pos(c.x, c.y)
			if c.command != nil {
				c.command.Move(offsets[i][0], offsets[i][1])
			}
			break
		}
	}
}

func (c *cursor) enter() {
	if !inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return
	}
	if c.command != nil {
		c.command.Enter(c.x, c.y)
	} else {
		player := PlayerManager.GetActivePlayer(c.x, c.y)
		if player != nil {
			Menu.Show(player)
		} else {
			Tip.SetMsg("选择的对象无法操作,请重新选择")
		}
	}
}

func (c *cursor) drawPlayer(p *Player, screen *ebiten.Image) {
	img := LoadImage(R.SPRITE.OTHER, p.Data.Image)
	utils.DrawImage(screen, img, 80+336i)
	// 名称称号
	utils.DrawAnchorText(screen, p.Data.Name, 144+336i, 0.5i, Font36, colornames.White)
	utils.DrawAnchorText(screen, p.Data.Career, 432+336i, 1+0.5i, Font36, colornames.White)
	// 装备  448 320
	for i := 0; i < len(p.Equips); i++ {
		temp := LoadImage(R.SPRITE.OTHER, p.Equips[i].Data.Image)
		utils.DrawImage(screen, temp, complex(float64(452+i*32), 324))
	}
	// HP,MP,ATK,DEF,ENE
	utils.DrawAnchorText(screen, fmt.Sprintf("HP:%.1f", p.Hp), 144+368i, 0.5i, Font32, colornames.White)
	utils.DrawAnchorText(screen, fmt.Sprintf("MP:%.1f", p.Mp), 240+368i, 0.5i, Font32, colornames.White)
	utils.DrawAnchorText(screen, fmt.Sprintf("ATK:%.1f", p.GetAtk()), 336+368i, 0.5i, Font32, colornames.White)
	utils.DrawAnchorText(screen, fmt.Sprintf("DEF:%.1f", p.GetDef()), 432+368i, 0.5i, Font32, colornames.White)
	utils.DrawAnchorText(screen, fmt.Sprintf("ENE:%.1f", p.Energy), 528+368i, 0.5i, Font32, colornames.White)
}
