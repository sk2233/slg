/*
@author: sk
@date: 2022/12/25
*/
package main

import (
	"GameBase2/config"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func Pos2Index(pos complex128) (int, int) {
	return int(real(pos) / 32), int(imag(pos) / 32)
}

func Index2Pos(x, y int) complex128 {
	return complex(float64(x*32), float64(y*32))
}

func LoadImage(path, name string) *ebiten.Image {
	return config.SpritesLoader.LoadStaticSprite(path, name).Image
}

func LoadImages(path, name string) []*ebiten.Image {
	return config.SpritesLoader.LoadDynamicSprite(path, name).Images
}

func GetAxis(negative, positive ebiten.Key) int {
	if inpututil.IsKeyJustPressed(negative) {
		return -1
	}
	if inpututil.IsKeyJustPressed(positive) {
		return 1
	}
	return 0
}

func FormatFloat(value float64, prec int) string {
	return strconv.FormatFloat(value, 'f', prec, 64)
}

func InvokeRoundBegan(src any) {
	if roundBegan, ok := src.(IRoundBegan); ok {
		roundBegan.RoundBegan()
	}
}

func InvokeGetAtkBuff(src any) float64 {
	if atkBuff, ok := src.(IAtkBuff); ok {
		return atkBuff.GetAtkBuff()
	}
	return 0
}
