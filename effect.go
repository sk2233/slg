/*
@author: sk
@date: 2023/1/1
*/
package main

import (
	"GameBase2/config"
	"GameBase2/object"
	"GameBase2/utils"
	"image/color"
	R "slg_game/res"

	"github.com/hajimehoshi/ebiten/v2"
)

type SimpleEffect struct {
	*object.TimeObject
	msg      string
	yOffset  float64 // 最多 64
	msgColor color.Color
}

func NewSimpleEffect(name, msg string, pos complex128, msgColor color.Color) *SimpleEffect {
	res := &SimpleEffect{msg: msg, msgColor: msgColor}
	data := config.SpritesLoader.LoadDynamicSprite(R.SPRITE.OTHER, name)
	animTime := float64(len(data.Images)) / data.AnimSpeed // 取 动画 与文本 最长的生存时间
	res.TimeObject = object.NewTimeObject(utils.Max(64, int(animTime)))
	res.Pos = pos
	temp := config.SpriteFactory.CreateDynamicSprite(R.SPRITE.OTHER, name)
	res.BindSprite(temp, res)
	return res
}

func (s *SimpleEffect) Update() {
	s.TimeObject.Update()
	s.yOffset += 0.5
}

func (s *SimpleEffect) Draw(screen *ebiten.Image) {
	s.TimeObject.Draw(screen)
	if s.yOffset > 32 || len(s.msg) == 0 || s.msgColor == nil {
		return // 文本绘制提前结束
	}
	utils.DrawAnchorText(screen, s.msg, s.Pos+complex(16, 16-s.yOffset), 0.5+0.5i, Font24, s.msgColor)
}
