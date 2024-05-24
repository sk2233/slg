/*
@author: sk
@date: 2022/12/25
*/
package main

import (
	"GameBase2/model"
	"GameBase2/object"
	"GameBase2/utils"
	"image/color"
	R "slg_game/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//======================MoveCommand==========================

type Grid struct {
	LastEnergy   float64 // 到达该步骤的 最大剩余能量
	IsSource     bool    // 是否为 源头 源头 没有上一个对象
	LastX, LastY int     // 上一步的位置
	X, Y         int     // 当前格子的位置
}

type MoveCommand struct {
	*object.PointObject
	move   bool
	player *Player
	die    bool
	grids  [][]*Grid
	timer  int
	step   *model.Stack[*Grid]
}

func (m *MoveCommand) Update() {
	if !m.move {
		m.timer++
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			m.die = true
			Cursor.ClearCommand(m)
		}
		return
	}
	if m.timer > 0 {
		m.timer--
	} else {
		if m.step.IsEmpty() {
			m.die = true
			Cursor.ClearCommand(m)
			m.player.Rank = 0 // 恢复层级
		} else {
			m.timer = 15
			grid := m.step.Pop()
			m.player.SetPos(grid.X, grid.Y)
			tileData := TileManager.GetTileInfo(grid.X, grid.Y)
			if tileData.HoverHandler != nil { // 处理经过土块事件
				tileData.HoverHandler(m.player)
			}
		}
	}
}

func (m *MoveCommand) IsDie() bool {
	return m.die
}

func NewMoveCommand(player *Player) *MoveCommand {
	res := &MoveCommand{player: player, move: false, die: false}
	res.PointObject = object.NewPointObject()
	utils.AddToLayer(R.LAYER.COMMAND, res)
	return res
}

func (m *MoveCommand) Draw(screen *ebiten.Image) {
	if m.move {
		return
	}
	if m.timer%30 < 15 {
		return // 闪烁的效果
	}
	for i := 0; i < len(m.grids); i++ {
		for j := 0; j < len(m.grids[i]); j++ {
			if m.grids[i][j] != nil {
				utils.FillRect(screen, complex(float64(i*32+2), float64(j*32+2)), 28+28i, Color0_255_0_127)
			}
		}
	}
}

func (m *MoveCommand) Move(xOff, yOff int) {}

func (m *MoveCommand) Enter(x, y int) {
	if m.grids[x][y] == nil {
		Tip.SetMsg("精力不足!")
		return
	}
	if PlayerManager.GetPlayer(x, y) != nil {
		Tip.SetMsg("目标处已存在对象!")
		return
	}
	m.move = true
	m.step = model.NewStack[*Grid]()
	grid := m.grids[x][y]
	m.player.Energy = grid.LastEnergy // 扣减精力
	for !grid.IsSource {
		m.step.Push(grid)
		grid = m.grids[grid.LastX][grid.LastY]
	}
	m.player.Rank = 2233 // 防止 遮挡
	m.timer = 15
}

func (m *MoveCommand) Block() bool {
	return m.move
}

func (m *MoveCommand) Init() {
	m.initGrids()
}

func (m *MoveCommand) initGrids() {
	m.grids = make([][]*Grid, 40)
	for i := 0; i < 40; i++ {
		m.grids[i] = make([]*Grid, 22)
	}
	p := m.player
	m.grids[p.X][p.Y] = &Grid{LastEnergy: p.Energy, IsSource: true, X: p.X, Y: p.Y}
	queue := model.NewQueue[*Grid]()
	queue.Add(m.grids[p.X][p.Y])
	dirs := [][]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	tileType := TileManager.GetTileTypes()
	has := false
	for !queue.IsEmpty() {
		grid := queue.Poll()
		for i := 0; i < len(dirs); i++ {
			x := grid.X + dirs[i][0]
			y := grid.Y + dirs[i][1]
			if x >= 0 && x < 40 && y >= 0 && y < 22 {
				energy := grid.LastEnergy - p.GetEnergyCost(tileType[y][x])
				currentEnergy := 0.0
				if m.grids[x][y] != nil {
					currentEnergy = m.grids[x][y].LastEnergy
				}
				if energy >= currentEnergy {
					m.grids[x][y] = &Grid{LastEnergy: energy, LastX: grid.X, LastY: grid.Y, X: x, Y: y}
					queue.Add(m.grids[x][y])
					has = true
				}
			}
		}
	}
	if !has { // 没有可以移动的 地方 提前结束
		Tip.SetMsg("没有精力移动到任何位置!")
		m.die = true
		Cursor.ClearCommand(m)
	}
}

//==================RangeHelper==================

type RangeHelper struct {
	x, y   int
	range0 []*OffsetPos
	color0 color.Color
}

func (h *RangeHelper) Draw(screen *ebiten.Image) {
	for i := 0; i < len(h.range0); i++ {
		x := h.x + h.range0[i].OffsetX
		y := h.y + h.range0[i].OffsetY
		if x >= 0 && x < 40 && y >= 0 && y < 22 {
			utils.FillRect(screen, complex(float64(x*32+2), float64(y*32+2)), 28+28i, h.color0)
		}
	}
}

func (h *RangeHelper) Has(x int, y int) bool {
	for i := 0; i < len(h.range0); i++ {
		if x == h.x+h.range0[i].OffsetX && y == h.y+h.range0[i].OffsetY {
			return true
		}
	}
	return false
}

func (h *RangeHelper) Add(offsetX, offsetY int) {
	h.x += offsetX
	h.y += offsetY
}

func NewRangeHelper(x int, y int, range0 []*OffsetPos, color0 color.Color) *RangeHelper {
	return &RangeHelper{x: x, y: y, range0: range0, color0: color0}
}

//===================AttackCommand===================

type AttackCommand struct { // 动画 通过 对象生成
	*object.PointObject
	player      *Player
	rangeHelper *RangeHelper
	timer       int
	die         bool
}

func (a *AttackCommand) IsDie() bool {
	return a.die
}

func (a *AttackCommand) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		a.die = true
		Cursor.ClearCommand(a)
	}
	a.timer++
}

func (a *AttackCommand) Draw(screen *ebiten.Image) {
	if a.timer%30 < 15 {
		return
	}
	a.rangeHelper.Draw(screen)
}

func (a *AttackCommand) Move(xOff, yOff int) {}

func (a *AttackCommand) Enter(x, y int) {
	if !a.rangeHelper.Has(x, y) {
		Tip.SetMsg("选择目标超出攻击范围")
		return
	}
	player := PlayerManager.GetTypePlayer(x, y, !a.player.IsEnemy)
	if player == nil {
		Tip.SetMsg("必须选择攻击目标")
		return
	}
	a.player.Attack(player)
	a.player.MarkAction() // 标记行动
	a.die = true
	Cursor.ClearCommand(a)
}

func (a *AttackCommand) Block() bool {
	return false
}

func NewAttackCommand(player *Player) *AttackCommand {
	res := &AttackCommand{player: player, die: false}
	res.PointObject = object.NewPointObject()
	res.rangeHelper = NewRangeHelper(player.X, player.Y, player.Data.AttackRange, Color255_0_0_127)
	utils.AddToLayer(R.LAYER.COMMAND, res)
	return res
}

//=======================TeleportCommand========================

type TeleportCommand struct {
	player *Player
}

func (t *TeleportCommand) Move(xOff, yOff int) {}

func (t *TeleportCommand) Enter(x, y int) {
	if PlayerManager.GetPlayer(x, y) != nil {
		Tip.SetMsg("传送目的地非空!")
		return
	}
	t.player.SetPos(x, y)
	Cursor.ClearCommand(t)
	utils.AddToLayer(R.LAYER.FG, NewSimpleEffect(R.OTHER.EFFECT.DELIVERY, "", t.player.Pos-4-4i, nil))
}

func (t *TeleportCommand) Block() bool {
	return false
}

func CreateTeleportCommand(player *Player) ICommand {
	return &TeleportCommand{player: player}
}

//=======================SelectPlayerCommand========================

type PlayerHandler func(player *Player) string // 若是返回 消息为"" 代表成功了

type SelectPlayerCommand struct {
	playerHandler PlayerHandler
}

func NewSelectCommand(playerHandler PlayerHandler) *SelectPlayerCommand {
	return &SelectPlayerCommand{playerHandler: playerHandler}
}

func (t *SelectPlayerCommand) Move(xOff, yOff int) {}

func (t *SelectPlayerCommand) Enter(x, y int) {
	if msg := t.playerHandler(PlayerManager.GetPlayer(x, y)); len(msg) > 0 {
		Tip.SetMsg(msg)
	} else {
		Cursor.ClearCommand(t)
	}
}

func (t *SelectPlayerCommand) Block() bool {
	return false
}

//=======================SelectPosCommand========================

type PosHandler func(x, y int) string // 若是返回 消息为"" 代表成功了

type SelectPosCommand struct {
	*object.PointObject
	posHandler   PosHandler
	rangeHelpers []*RangeHelper
	timer        int
	die          bool
}

func (t *SelectPosCommand) Draw(screen *ebiten.Image) {
	if t.timer%30 < 15 {
		return
	}
	for i := 0; i < len(t.rangeHelpers); i++ {
		t.rangeHelpers[i].Draw(screen)
	}
}

func (t *SelectPosCommand) Update() {
	t.timer++
}

func (t *SelectPosCommand) IsDie() bool {
	return t.die
}

func NewSelectPosCommand(posHandler PosHandler, rangeHelpers []*RangeHelper) *SelectPosCommand {
	res := &SelectPosCommand{posHandler: posHandler, rangeHelpers: rangeHelpers}
	res.PointObject = object.NewPointObject()
	utils.AddToLayer(R.LAYER.FG, res)
	return res
}

func (t *SelectPosCommand) Move(xOff, yOff int) {}

func (t *SelectPosCommand) Enter(x, y int) {
	if !t.rangeHelpers[0].Has(x, y) { // 第一个是辅助范围，其他是 辅助范围
		Tip.SetMsg("超出选择范围")
		return
	}
	if msg := t.posHandler(x, y); len(msg) > 0 {
		Tip.SetMsg(msg)
		return
	}
	t.die = true
	Cursor.ClearCommand(t)
}

func (t *SelectPosCommand) Block() bool {
	return false
}

//===================OpenFireCommand=========================

type OpenFireCommand struct {
	*object.PointObject
	player       *Player
	die          bool
	timer        int
	rangeHelpers []*RangeHelper
}

func (o *OpenFireCommand) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		o.die = true
		Cursor.ClearCommand(o)
	}
	o.timer++
}

func (o *OpenFireCommand) Draw(screen *ebiten.Image) {
	if o.timer%30 < 15 {
		return
	}
	for i := 0; i < len(o.rangeHelpers); i++ {
		o.rangeHelpers[i].Draw(screen)
	}
}

func (o *OpenFireCommand) IsDie() bool {
	return o.die
}

func (o *OpenFireCommand) Move(xOff, yOff int) {
	for i := 0; i < len(o.rangeHelpers); i++ {
		o.rangeHelpers[i].Add(xOff, yOff)
	}
}

func (o *OpenFireCommand) Enter(x, y int) {
	player := PlayerManager.GetPlayer(x, y)
	if player != nil {
		player.Hurt(40, o.player)
	}
	utils.AddToLayer(R.LAYER.FG, NewSimpleEffect(R.OTHER.EFFECT.FIRE, "", Index2Pos(x, y), nil))
	for i := 0; i < 4; i++ {
		tx, ty := x+Dirs[i][0], y+Dirs[i][1]
		player = PlayerManager.GetPlayer(tx, ty)
		if player != nil {
			player.Hurt(20, o.player)
		}
		utils.AddToLayer(R.LAYER.FG, NewSimpleEffect(R.OTHER.EFFECT.FIRE, "", Index2Pos(tx, ty), nil))
	}
	o.die = true
	Cursor.ClearCommand(o)
	o.player.MarkAction()
}

func (o *OpenFireCommand) Block() bool {
	return false
}

func CreateOpenFireCommand(player *Player) ICommand {
	res := &OpenFireCommand{player: player, die: false}
	res.PointObject = object.NewPointObject()
	res.rangeHelpers = make([]*RangeHelper, 0)
	pos1 := []*OffsetPos{{0, 0}}
	pos2 := []*OffsetPos{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	res.rangeHelpers = append(res.rangeHelpers, NewRangeHelper(player.X, player.Y, pos1, Color255_0_0_127))
	res.rangeHelpers = append(res.rangeHelpers, NewRangeHelper(player.X, player.Y, pos2, Color255_255_0_127))
	utils.AddToLayer(R.LAYER.FG, res)
	return res
}
