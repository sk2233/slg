/*
@author: sk
@date: 2023/1/1
*/
package main

type Equipment struct {
	player *Player
	Data   *EquipmentData
	Skills []ISkill
}

func (e *Equipment) initSkill() {
	for i := 0; i < len(e.Data.Skills); i++ {
		e.Skills = append(e.Skills, CreateSkill(e.Data.Skills[i], e.player))
	}
}

func NewEquipment(name string, player *Player) *Equipment {
	data := DataManager.GetEquipData(name)
	res := &Equipment{Data: data, player: player, Skills: make([]ISkill, 0)}
	res.initSkill()
	return res
}
