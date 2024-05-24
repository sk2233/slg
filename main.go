package main

import (
	"GameBase2/app"
	"GameBase2/config"
	"GameBase2/room"
	"embed"
	R "slg_game/res"
)

var (
	//go:embed res
	files embed.FS

	StackRoom *room.StackUIRoom

	TileManager   *tileManager
	DataManager   *dataManager
	PlayerManager *playerManager

	Cursor *cursor
	Tip    *tip
	Menu   *menu
)

func main() {
	config.ViewSize = complex(640, 384)
	config.Debug = true
	config.ShowFps = true
	config.Files = &files // 先使用内部资源 ，不存在  再寻找外部资源文件
	Font36 = config.FontFactory.CreateFont(R.RAW.IPSX_FONT, 36, 36)
	Font32 = config.FontFactory.CreateFont(R.RAW.IPSX_FONT, 36, 32)
	Font24 = config.FontFactory.CreateFont(R.RAW.IPSX_FONT, 36, 24)
	app.Run(NewMainApp(), 1280, 768)
}

type MainApp struct {
	*app.App
}

// Init 必须先传入实例  初始化使用该方法
func NewMainApp() *MainApp {
	res := &MainApp{}
	res.App = app.NewApp()
	TileManager = NewTileManager()
	DataManager = NewDataManager()
	PlayerManager = NewPlayerManager()
	temp := config.RoomFactory.LoadAndCreate(R.MAP.MAIN)
	StackRoom = room.NewStackUIRoom(temp.(*room.Room))
	StackRoom.AddManager(TileManager)
	StackRoom.AddManager(DataManager)
	StackRoom.AddManager(PlayerManager)
	res.PushRoom(StackRoom)
	return res
}
