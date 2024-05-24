/*
@author: sk
@date: 2023/1/1
*/
package main

//=====================BaseItem=====================

type BaseItem struct {
	name  string
	Count int // 不太规范 把数量也存储进来
}

func NewBaseItem(name string, num int) *BaseItem {
	return &BaseItem{name: name, Count: num}
}

func (b *BaseItem) GetCount() int {
	return b.Count
}

func (b *BaseItem) GetName() string {
	return b.name
}

//=====================HpPill=======================

type HpPill struct {
	*BaseItem
}

func NewHpPill(num int) *HpPill {
	res := &HpPill{}
	res.BaseItem = NewBaseItem("红瓶", num)
	return res
}

func (h *HpPill) GetCommand(player *Player) ICommand { // 返回nil 就是即时效果
	player.RecoverHp(40, player)
	h.Count--
	return nil
}

//=====================MpPill=======================

type MpPill struct {
	*BaseItem
}

func NewMpPill(num int) *MpPill {
	res := &MpPill{}
	res.BaseItem = NewBaseItem("蓝瓶", num)
	return res
}

func (h *MpPill) GetCommand(player *Player) ICommand { // 返回nil 就是即时效果
	player.RecoverMp(40, player)
	h.Count--
	return nil
}

//=====================AnalepticPill=======================

type AnalepticPill struct {
	*BaseItem
}

func NewAnalepticPill(num int) *AnalepticPill {
	res := &AnalepticPill{}
	res.BaseItem = NewBaseItem("兴奋剂", num)
	return res
}

func (h *AnalepticPill) GetCommand(player *Player) ICommand { // 返回nil 就是即时效果
	player.AddBuff(CreateSkill("兴奋剂", player))
	h.Count--
	return nil
}

//=====================CommandItem=======================

type CommandFactory func(*Player) ICommand

type CommandItem struct {
	*BaseItem
	commandFactory CommandFactory
}

func NewCommandItem(name string, num int, commandFactory CommandFactory) *CommandItem {
	res := &CommandItem{commandFactory: commandFactory}
	res.BaseItem = NewBaseItem(name, num)
	return res
}

func (h *CommandItem) GetCommand(player *Player) ICommand { // 返回nil 就是即时效果
	h.Count--
	return h.commandFactory(player)
}
