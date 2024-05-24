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

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
)

func init() {
	config.ObjectFactory.RegisterPointFactory(R.CLASS.MENU, createMenu)
}

func createMenu(data *model.ObjectData) model.ObjectAble {
	res := &menu{}
	res.PointObject = object.NewPointObject()
	factory.FillPointObject(data, res.PointObject)
	res.Enable = false
	Menu = res
	return res
}

//==================================

type menu struct {
	*object.PointObject
	player        *Player
	mainMenu      []string // 一级菜单
	mainMenuIndex int
	subMenu       []string // 二级菜单  侧边的
	subMenuIndex  int
	arrowImage    *ebiten.Image
	skip          bool // 防止按键事件传递 Enter进来后 继续响应 Enter
}

func (m *menu) Order() int {
	return 33
}

func (m *menu) Init() {
	m.arrowImage = LoadImage(R.SPRITE.OTHER, R.OTHER.OTHER.ARROW)
}

func (m *menu) Move(xOff, yOff int) {}

func (m *menu) Enter(x, y int) {}

func (m *menu) Block() bool {
	return true
}

func (m *menu) Update() {
	if m.skip { // 主要防止按键事件穿透
		m.skip = false
		return
	}
	if m.subMenu != nil {
		m.subMenuUpdate()
	} else {
		m.mainMenuUpdate()
	}
}

func (m *menu) mainMenuUpdate() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) { // 退出菜单
		m.Hide()
		return
	}
	yOffset := GetAxis(ebiten.KeyW, ebiten.KeyS)
	if yOffset != 0 {
		m.mainMenuIndex = (m.mainMenuIndex + yOffset + len(m.mainMenu)) % len(m.mainMenu)
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		temp := m.player.GetSubMenu(m.mainMenu[m.mainMenuIndex])
		if temp != nil {
			m.subMenu = temp
		} else {
			// 触发对应事件
			m.player.StartCommand(m.mainMenu[m.mainMenuIndex], m.mainMenuIndex, m.subMenuIndex)
			m.Hide()
		}
	}
}

func (m *menu) subMenuUpdate() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) { // 退出子菜单
		m.subMenu = nil
		m.subMenuIndex = 0
		return
	}
	l := len(m.subMenu)
	if l <= 0 {
		return
	}
	xOffset := GetAxis(ebiten.KeyA, ebiten.KeyD)
	if xOffset != 0 {
		m.subMenuIndex = (m.subMenuIndex + xOffset + l) % l
		return
	}
	yOffset := GetAxis(ebiten.KeyW, ebiten.KeyS)
	if yOffset != 0 {
		temp := m.subMenuIndex + yOffset*2
		if temp >= 0 && temp < l {
			m.subMenuIndex = temp
		}
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		// 触发对应事件
		m.player.StartCommand(m.subMenu[m.subMenuIndex], m.mainMenuIndex, m.subMenuIndex)
		m.Hide()
	}
}

func (m *menu) UIDraw(screen *ebiten.Image) {
	m.drawMainMenu(screen)
	m.drawSubMenu(screen)
}

func (m *menu) drawMainMenu(screen *ebiten.Image) {
	utils.DrawRect(screen, 64+64i, complex(96, float64(len(m.mainMenu)*32)), Color1_0_98, Color201_163_66)
	for i := 0; i < len(m.mainMenu); i++ {
		utils.DrawAnchorText(screen, m.mainMenu[i], complex(112, float64(80+i*32)), 0.5+0.5i, Font36,
			colornames.White)
	}
	utils.DrawImage(screen, m.arrowImage, complex(72, float64(72+m.mainMenuIndex*32)))
}

func (m *menu) drawSubMenu(screen *ebiten.Image) {
	utils.DrawRect(screen, 224+64i, 352+224i, Color1_0_98, Color201_163_66)
	if m.subMenu == nil { // 绘制基本信息
		player := m.player
		utils.FillRect(screen, 224+64i, 96+96i, colornames.White) // TODO 绘制大头像
		// 名称 职业   TODO 后期各种加成信息也显示出来
		utils.DrawAnchorText(screen, player.Data.Name, 328+80i, 0.5i, Font36, colornames.White)
		utils.DrawAnchorText(screen, player.Data.Career, 328+112i, 0.5i, Font36, colornames.White)
		// HP,MP
		utils.DrawAnchorText(screen, "HP", 256+176i, 0.5+0.5i, Font36, colornames.White)
		utils.DrawAnchorText(screen, fmt.Sprintf("%.1f/%.1f", player.Hp, player.GetMaxHp()), 336+176i, 0.5+0.5i, Font36, colornames.White)
		utils.DrawAnchorText(screen, "MP", 416+176i, 0.5+0.5i, Font36, colornames.White)
		utils.DrawAnchorText(screen, fmt.Sprintf("%.1f/%.1f", player.Mp, player.GetMaxMp()), 496+176i, 0.5+0.5i, Font36, colornames.White)
		// Atk,Def
		utils.DrawAnchorText(screen, "ATK", 256+208i, 0.5+0.5i, Font36, colornames.White)
		utils.DrawAnchorText(screen, FormatFloat(player.GetAtk(), 1), 336+208i, 0.5+0.5i, Font36, colornames.White)
		utils.DrawAnchorText(screen, "DEF", 416+208i, 0.5+0.5i, Font36, colornames.White)
		utils.DrawAnchorText(screen, FormatFloat(player.GetDef(), 1), 496+208i, 0.5+0.5i, Font36, colornames.White)
		// Energy,RecoverEnergy
		utils.DrawAnchorText(screen, "ENE", 256+240i, 0.5+0.5i, Font36, colornames.White)
		utils.DrawAnchorText(screen, fmt.Sprintf("%.1f/%.1f", player.Energy, player.Data.MaxEnergy), 336+240i, 0.5+0.5i, Font36, colornames.White)
		utils.DrawAnchorText(screen, "REC", 416+240i, 0.5+0.5i, Font36, colornames.White)
		utils.DrawAnchorText(screen, FormatFloat(player.Data.RecoverEnergy, 1), 496+240i, 0.5+0.5i, Font36, colornames.White)
	} else { // 有子菜单显示子菜单
		if len(m.subMenu) <= 0 { // 没有请 返回空数组 不要返回nil
			return // 空
		}
		for i := 0; i < len(m.subMenu); i++ {
			utils.DrawAnchorText(screen, m.subMenu[i], complex(float64(256+(i%2)*176), float64(80+(i/2)*32)),
				0.5i, Font36, colornames.White)
		}
		utils.DrawImage(screen, m.arrowImage, complex(float64(232+(m.subMenuIndex%2)*176),
			float64(72+(m.subMenuIndex/2)*32)))
	}
}

func (m *menu) Hide() { // 退出菜单
	Cursor.ClearCommand(m) // 防止 错误移除 (例如已被覆盖的情况)
	m.player = nil
	m.mainMenu, m.subMenu = nil, nil
	m.mainMenuIndex, m.subMenuIndex = 0, 0
	m.Enable = false
}

func (m *menu) Show(player *Player) {
	Cursor.SetCommand(m)
	m.player = player
	m.Enable = true
	m.skip = true
	m.mainMenu = player.GetMainMenu()
}
