package models

type Collision struct {
	Intersected bool
	Top    bool
	Bottom bool
	Left   bool
	Right  bool
	
	X BoxPosition
	Y BoxPosition
}

func DetectCollision(boxPosition1, boxPosition2 BoxPosition) Collision {
	collision := Collision{
		X: boxPosition1,
		Y: boxPosition2,
	}

	pos1 := boxPosition1.GetPos()
	pos2 := boxPosition2.GetPos()

	box1 := boxPosition1.GetBox()
	box2 := boxPosition2.GetBox()

	if pos1.X < pos2.X+box2.X &&
		pos1.X+box1.X > pos2.X &&
		pos1.Y < pos2.Y+box2.Y &&
		pos1.Y+box1.Y > pos2.Y {

		// Обнаружена коллизия
		collision.Intersected = true

		if pos1.X+box1.X > pos2.X && pos1.X < pos2.X {
			collision.Left = true
		}
		if pos1.X < pos2.X+box2.X && pos1.X+box1.X > pos2.X+box2.X {
			collision.Right = true
		}

		// Проверяем, с какой гранью произошла коллизия
		if pos1.Y+box1.Y > pos2.Y && pos1.Y < pos2.Y {
			collision.Top = true
		}
		if pos1.Y < pos2.Y+box2.Y && pos1.Y+box1.Y > pos2.Y+box2.Y {
			collision.Bottom = true
		}
	}

	return collision
}
