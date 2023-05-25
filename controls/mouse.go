package controls

import (
	//"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func SetMousePosition(x, y int, fromLine int) {
	//fmt.Println("set mouse pos ", x, y, fromLine)
	rl.SetMousePosition(x, y)
}

func DisableCursor(fromLine int) {
	//fmt.Println("disable cursor ", fromLine)
	rl.DisableCursor()
}

func EnableCursor(fromLine int) {
	//fmt.Println("enable cursor ", fromLine)
	rl.EnableCursor()
}