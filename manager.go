/*
@author: sk
@date: 2022/12/25
*/
package main

import (
	"GameBase2/config"
	"GameBase2/utils"
	"reflect"
	R "slg_game/res"
	"strings"
)

//====================tileManager======================

type tileManager struct {
	tileTypes [][]int
}

func (t *tileManager) initTileTypes() {
	mapData := config.MapLoader.LoadMap(R.MAP.MAIN)
	for i := 0; i < len(mapData.Layers); i++ {
		if mapData.Layers[i].Name == R.LAYER.BG {
			t.tileTypes = mapData.Layers[i].Tiles
			return
		}
	}
	utils.LogErr("initTileTypes err")
}

func (t *tileManager) GetTileTypes() [][]int {
	return t.tileTypes
}

func (t *tileManager) GetTileInfo(x int, y int) *TileData {
	return DataManager.GetTileData(t.tileTypes[y][x])
}

func NewTileManager() *tileManager {
	res := &tileManager{}
	res.initTileTypes()
	return res
}

//=========================dataManager============================

type dataManager struct { // 先配置在 这里   最终 通过读取文件的方式 存储
	playerData  map[string]*PlayerData
	tileData    map[int]*TileData
	equipData   map[string]*EquipmentData
	playerItems []IItem // 团队携带的物品   物品都是有使用效果的
	enemyItems  []IItem
}

func (m *dataManager) initPlayerData() {
	m.playerData[R.OBJECT.ALICE] = &PlayerData{
		Name:   "爱丽丝·玛格特洛依德",
		Career: "魔法使",
		Image:  R.OTHER.PLAYER.ALICE,
		Equips: []string{"书", "弓箭", "蝴蝶结"},
	}
	m.playerData[R.OBJECT.REIMU] = &PlayerData{
		Name:   "博丽灵梦",
		Career: "巫女",
		Image:  R.OTHER.PLAYER.REIMU,
		Equips: []string{"斗篷", "帽子", "刀", "灵珠"},
	}
	m.playerData[R.OBJECT.MARISA] = &PlayerData{
		Name:   "雾雨魔理沙",
		Career: "魔法使",
		Image:  R.OTHER.PLAYER.MARISA,
		Equips: []string{"戒指", "法杖", "书"},
	}
	m.playerData[R.OBJECT.REISEN] = &PlayerData{
		Name:   "铃仙·优昙华院·因幡",
		Career: "月兔/妖兽",
		Image:  R.OTHER.PLAYER.REISEN,
		Equips: []string{"弓箭", "蝴蝶结", "斗篷", "帽子"},
	}
	m.playerData[R.OBJECT.SAKUYA] = &PlayerData{
		Name:   "十六夜咲夜",
		Career: "女仆",
		Image:  R.OTHER.PLAYER.SAKUYA,
		Equips: []string{"刀", "灵珠", "戒指"},
	}
	m.playerData[R.OBJECT.YOUMU] = &PlayerData{
		Name:   "魂魄妖梦",
		Career: "半人半灵",
		Image:  R.OTHER.PLAYER.YOUMU,
		Equips: []string{"法杖", "灵珠", "帽子", "蝴蝶结"},
	} // TODO TEST 先 全部统一
	for _, playerData := range m.playerData {
		playerData.MaxHp, playerData.MaxMp = 100, 100
		playerData.Atk, playerData.Def = 15, 5
		playerData.MaxEnergy, playerData.RecoverEnergy = 10, 5
		playerData.EnergyCost = make(map[int]float64)
		for i := 1; i < 9; i++ {
			playerData.EnergyCost[i] = float64(i%3 + 1)
		}
		playerData.EnergyCost[TileTypeMagma] = 4   //岩浆 减少走的情况
		playerData.EnergyCost[TileTypeWall] = 2233 // 不可逾越
		playerData.Skills = []string{"治疗", "火焰", "驱散"}
		playerData.AttackRange = []*OffsetPos{
			{OffsetX: 1}, {OffsetX: -1}, {OffsetY: 1}, {OffsetY: -1},
		}
	}
}

func (m *dataManager) initItems() {
	m.playerItems = make([]IItem, 0)
	m.enemyItems = make([]IItem, 0)
	m.playerItems = append(m.playerItems, NewHpPill(7))                                     // 补充  40 HP
	m.playerItems = append(m.playerItems, NewMpPill(5))                                     // 补充  40 MP
	m.playerItems = append(m.playerItems, NewAnalepticPill(3))                              // 3回合  Atk+20
	m.playerItems = append(m.playerItems, NewCommandItem("传送卷轴", 4, CreateTeleportCommand)) // 传送到指定位置
	m.enemyItems = append(m.enemyItems, NewHpPill(7))                                       // 补充  40 HP
	m.enemyItems = append(m.enemyItems, NewMpPill(5))                                       // 补充  40 MP
	m.enemyItems = append(m.enemyItems, NewAnalepticPill(3))                                // 3回合  Atk+20
	m.enemyItems = append(m.enemyItems, NewCommandItem("传送卷轴", 4, CreateTeleportCommand))   // 传送到指定位置
}

func (m *dataManager) GetPlayerData(name string) *PlayerData {
	return m.playerData[name]
}

func (m *dataManager) GetPlayerItems(enemy bool) []IItem {
	if enemy {
		return m.enemyItems
	}
	return m.playerItems
}

func (m *dataManager) UsePlayerItem(p *Player, name string, index int) (ICommand, bool) {
	items := m.playerItems
	if p.IsEnemy {
		items = m.enemyItems
	}
	if index >= len(items) {
		return nil, false
	}
	if !strings.HasPrefix(name, items[index].GetName()) {
		return nil, false
	}
	res := items[index].GetCommand(p)
	if items[index].GetCount() <= 0 {
		if p.IsEnemy {
			m.enemyItems = append(m.enemyItems[:index], m.enemyItems[index+1:]...)
		} else {
			m.playerItems = append(m.playerItems[:index], m.playerItems[index+1:]...)
		}
	}
	return res, true
}

func (m *dataManager) GetTileData(tileType int) *TileData {
	return m.tileData[tileType]
}

func (m *dataManager) initTileData() {
	m.tileData = make(map[int]*TileData)
	m.tileData[TileTypeWall] = &TileData{Image: LoadImage(R.SPRITE.TILE, R.TILE.TILE.QIANG)}
	m.tileData[TileTypeDesert] = &TileData{Image: LoadImage(R.SPRITE.TILE, R.TILE.TILE.SHA)}
	m.tileData[TileTypeForest] = &TileData{Image: LoadImage(R.SPRITE.TILE, R.TILE.TILE.SHEN)}
	m.tileData[TileTypeGrass] = &TileData{Image: LoadImage(R.SPRITE.TILE, R.TILE.TILE.CHAO)}
	m.tileData[TileTypeHill] = &TileData{Image: LoadImage(R.SPRITE.TILE, R.TILE.TILE.SHAN)}
	m.tileData[TileTypeRiver] = &TileData{Image: LoadImage(R.SPRITE.TILE, R.TILE.TILE.HE)}
	m.tileData[TileTypeCannon] = &TileData{Image: LoadImage(R.SPRITE.TILE, R.TILE.TILE.PAO)}
	m.tileData[TileTypeMagma] = &TileData{Image: LoadImage(R.SPRITE.TILE, R.TILE.TILE.YAN)}
	// 岩浆 大炮 都是特殊地形
	m.tileData[TileTypeMagma].HoverHandler = MagmaHoverHandler
	m.tileData[TileTypeCannon].CommandMap = make(map[string]CommandFactory)
	m.tileData[TileTypeCannon].CommandMap["开炮"] = CreateOpenFireCommand
}

func MagmaHoverHandler(player *Player) string {
	// 每走一下 烧 一部分血
	player.Hurt(7, nil)
	return ""
}

func (m *dataManager) initEquipData() {
	m.equipData = make(map[string]*EquipmentData)
	// 各种装备增益效果 暂时 都为0
	m.equipData["书"] = &EquipmentData{Name: "书", Image: R.OTHER.EQUIP.BOOK}
	m.equipData["弓箭"] = &EquipmentData{Name: "弓箭", Image: R.OTHER.EQUIP.BOW}
	m.equipData["蝴蝶结"] = &EquipmentData{Name: "蝴蝶结", Image: R.OTHER.EQUIP.BOWKNOT}
	m.equipData["斗篷"] = &EquipmentData{Name: "斗篷", Image: R.OTHER.EQUIP.CLOAK}
	m.equipData["帽子"] = &EquipmentData{Name: "帽子", Image: R.OTHER.EQUIP.HAT}
	m.equipData["刀"] = &EquipmentData{Name: "刀", Image: R.OTHER.EQUIP.KNIFE}
	m.equipData["灵珠"] = &EquipmentData{Name: "灵珠", Image: R.OTHER.EQUIP.PILL}
	m.equipData["戒指"] = &EquipmentData{Name: "戒指", Image: R.OTHER.EQUIP.RING}
	m.equipData["法杖"] = &EquipmentData{Name: "法杖", Image: R.OTHER.EQUIP.STAFF}
}

func (m *dataManager) GetEquipData(name string) *EquipmentData {
	return m.equipData[name]
}

func NewDataManager() *dataManager {
	res := &dataManager{playerData: make(map[string]*PlayerData)}
	res.initPlayerData()
	res.initTileData()
	res.initEquipData()
	res.initItems()
	return res
}

//===========================playerManager=============================

type playerManager struct {
	players      []*Player
	enemies      []*Player
	playerAction bool
	round        int
}

func NewPlayerManager() *playerManager {
	return &playerManager{players: make([]*Player, 0), enemies: make([]*Player, 0), playerAction: true, round: 1}
}

func (p *playerManager) Init() {
	objects := utils.GetObjectLayer(R.LAYER.OBJECT).GetObjectsByType(reflect.TypeOf(&Player{}))
	for i := 0; i < len(objects); i++ {
		player := objects[i].(*Player)
		if player.IsEnemy {
			p.enemies = append(p.enemies, player)
		} else {
			p.players = append(p.players, player)
		}
	}
	p.ResetAction(false)
}

func (p *playerManager) ResetAction(isEnemy bool) {
	p.playerAction = !isEnemy
	players := p.enemies
	if p.playerAction {
		players = p.players
		p.round++ // 累计回合数
	}
	for i := 0; i < len(players); i++ {
		players[i].RoundBegan()
	}
}

func (p *playerManager) GetActivePlayer(x, y int) *Player {
	players := p.GetPlayers()
	for i := 0; i < len(players); i++ {
		if players[i].CanSelect(x, y) {
			return players[i]
		}
	}
	return nil
}

func (p *playerManager) UpdateAction() {
	players := p.GetPlayers()
	for i := 0; i < len(players); i++ {
		if !players[i].Action {
			return
		}
	}
	p.ResetAction(p.playerAction)
}

func (p *playerManager) GetPlayers() []*Player { // 获取当前激活的 玩家数组
	if p.playerAction {
		return p.players
	}
	return p.enemies
}

func (p *playerManager) GetTypePlayer(x int, y int, isEnemy bool) *Player {
	players := p.players
	if isEnemy {
		players = p.enemies
	}
	for i := 0; i < len(players); i++ {
		if players[i].X == x && players[i].Y == y {
			return players[i]
		}
	}
	return nil
}

func (p *playerManager) GetPlayer(x int, y int) *Player {
	for i := 0; i < len(p.players); i++ {
		if p.players[i].X == x && p.players[i].Y == y {
			return p.players[i]
		}
	}
	for i := 0; i < len(p.enemies); i++ {
		if p.enemies[i].X == x && p.enemies[i].Y == y {
			return p.enemies[i]
		}
	}
	return nil
}

func (p *playerManager) GetRound() int {
	return p.round
}
