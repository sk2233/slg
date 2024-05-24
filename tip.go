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
	R "slg_game/res"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func init() {
	config.ObjectFactory.RegisterPointFactory(R.CLASS.TIP, createTip)
}

func createTip(data *model.ObjectData) model.ObjectAble {
	res := &tip{showTime: 60}
	res.PointObject = object.NewPointObject()
	factory.FillPointObject(data, res.PointObject)
	Tip = res
	return res
}

//=====================tip========================

type tip struct {
	*object.PointObject
	showMsg             string
	showTime, showTimer int
}

func (t *tip) Order() int {
	return 2233
}

func (t *tip) SetMsg(msg string) {
	t.showMsg = msg
	t.showTimer = t.showTime
}

func (t *tip) UIDraw(screen *ebiten.Image) {
	if t.showTimer > 0 {
		utils.DrawAnchorText(screen, t.showMsg, complex(360, float64(t.showTimer)*64/60), 0.5, Font36, colornames.White)
		t.showTimer--
	}
}
