package models

import rl "github.com/gen2brain/raylib-go/raylib"

type Collision struct {
	Intersected              bool
	Top, Bottom, Left, Right bool
	X                        BoxPosition
	Y                        BoxPosition
}

func DetectCollision(boxPosition1, boxPosition2 BoxPosition) Collision {
	pos1 := boxPosition1.GetPos()
	pos2 := boxPosition2.GetPos()

	box1 := boxPosition1.GetBox()
	box2 := boxPosition2.GetBox()

	rect1 := rl.NewRectangle(pos1.X, pos1.Y, box1.X, box1.Y)
	rect2 := rl.NewRectangle(pos2.X, pos2.Y, box2.X, box2.Y)

	collision := Collision{
		X: boxPosition1,
		Y: boxPosition2,
	}

	collision.Intersected = rl.CheckCollisionRecs(rect1, rect2)

	if collision.Intersected {

		if pos1.X+box1.X > pos2.X && pos1.X < pos2.X {
			collision.Left = true
		}
		if pos1.X < pos2.X+box2.X && pos1.X+box1.X > pos2.X+box2.X {
			collision.Right = true
		}

		if pos1.Y+box1.Y > pos2.Y && pos1.Y < pos2.Y {
			collision.Top = true
		}
		if pos1.Y < pos2.Y+box2.Y && pos1.Y+box1.Y > pos2.Y+box2.Y {
			collision.Bottom = true
		}

	}

	return collision
}
