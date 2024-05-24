/*
@author: sk
@date: 2022/12/25
*/
package main

import (
	"image/color"

	"golang.org/x/image/font"
)

const ( // 新增 地形  这里一定 要映射 !!!!
	TileTypeUnknown = iota
	TileTypeCannon
	TileTypeDesert
	TileTypeMagma
	TileTypeRiver
	TileTypeHill
	TileTypeForest
	TileTypeGrass
	TileTypeWall
)

var (
	Font36 font.Face
	Font32 font.Face
	Font24 font.Face
)

var (
	Color1_0_98        = color.RGBA{R: 1, B: 98, A: 255}
	Color201_163_66    = color.RGBA{R: 201, G: 163, B: 66, A: 255}
	Color0_255_0_127   = color.RGBA{G: 255, A: 127}
	Color255_0_0_127   = color.RGBA{R: 255, A: 127}
	Color255_255_0_127 = color.RGBA{R: 255, G: 255, A: 127}
)

var (
	Dirs = [][]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
)
