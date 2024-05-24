/*
@author: sk
@date: 2022/12/25
*/
package main

type ICommand interface {
	Move(xOff, yOff int) // 移动的偏移
	Enter(x, y int)      // 确定  最终位置
	Block() bool         // 是否阻塞  玩家行动完 动画时需要阻塞
}

//===================Skill=====================

type ISkill interface { // 包含 主动与被动技能  被动技能靠各种接口处理
}

type IInitiativeSkill interface {
	ISkill

	GetName() string
	GetCommand() ICommand
	GetCost() float64
}

type IAtkBuff interface {
	GetAtkBuff() float64
}

type IRoundBegan interface {
	RoundBegan()
}

//===================Item========================

type IItem interface {
	GetName() string
	GetCount() int
	GetCommand(*Player) ICommand
}
