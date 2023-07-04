package models

import (
	"ahasuerus/collision"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func RotateTriangleByA(A *rl.Vector2, B *rl.Vector2, C *rl.Vector2, degrees float64) {
	radians := degrees * math.Pi / 180.0

	// Вычисляем синус и косинус угла поворота
	sin := math.Sin(radians)
	cos := math.Cos(radians)

	// Перемещаем вершины B и C относительно вершины A
	B.X, B.Y = rotatePoint(A, B, float32(sin), float32(cos))
	C.X, C.Y = rotatePoint(A, C, float32(sin), float32(cos))
}

type DynamicHitboxMap struct {
	topLeftOne rl.Vector2
	topLeftTwo rl.Vector2

	topMiddle rl.Vector2

	topRightOne rl.Vector2
	topRightTwo rl.Vector2

	leftMiddle  rl.Vector2
	center      rl.Vector2
	rightMiddle rl.Vector2

	bottomLeftOne rl.Vector2
	bottomLeftTwo rl.Vector2

	bottomMiddle rl.Vector2

	bottomRightOne rl.Vector2
	bottomRightTwo rl.Vector2
}

func Vec2Rotate(v rl.Vector2, degrees float64) rl.Vector2 {
	// Преобразование угла в радианы
	radians := degrees * math.Pi / 180

	// Поворот вектора с использованием тригонометрических функций
	result :=  rl.Vector2{
		X: (v.X*float32(math.Cos(radians)) - v.Y*float32(math.Sin(radians))),
		Y: (v.X*float32(math.Sin(radians)) + v.Y*float32(math.Cos(radians))),
	}
	return result
}

func GetDynamicHitboxMap(pos rl.Vector2, width, height float32) DynamicHitboxMap {
	topCornerOffset := float32(20)
	middleBottomCornerOffset := float32(30)
	bottomCornerOffset := float32(0)
	return DynamicHitboxMap{
		topLeftOne: rl.Vector2{pos.X + topCornerOffset, pos.Y},
		topLeftTwo: rl.Vector2{pos.X, pos.Y + topCornerOffset},

		topMiddle: rl.Vector2{pos.X + width/2, pos.Y},

		topRightOne: rl.Vector2{pos.X + width - topCornerOffset, pos.Y},
		topRightTwo: rl.Vector2{pos.X + width, pos.Y + topCornerOffset},

		leftMiddle:  rl.Vector2{pos.X, pos.Y + height/2},
		center:      rl.Vector2{pos.X + width/2, pos.Y + height/2},
		rightMiddle: rl.Vector2{pos.X + width, pos.Y + height/2},

		bottomRightOne: rl.Vector2{pos.X + width, pos.Y + height - middleBottomCornerOffset},
		bottomRightTwo: rl.Vector2{pos.X + width - bottomCornerOffset, pos.Y + height},

		bottomMiddle: rl.Vector2{pos.X + width/2, pos.Y + height},

		bottomLeftOne: rl.Vector2{pos.X + bottomCornerOffset, pos.Y + height},
		bottomLeftTwo: rl.Vector2{pos.X, pos.Y + height - middleBottomCornerOffset},
	}
}

func GetDynamicHitboxFromMap(m DynamicHitboxMap) collision.Hitbox {
	return collision.Hitbox{
		Polygons: []collision.Polygon{
			{
				Points: [3]rl.Vector2{
					m.topLeftOne, m.topMiddle, m.center,
				},
			},
			{
				Points: [3]rl.Vector2{
					m.topMiddle, m.topRightOne, m.center,
				},
			},
			{
				Points: [3]rl.Vector2{
					m.topRightTwo, m.rightMiddle, m.center,
				},
			},
			{
				Points: [3]rl.Vector2{
					m.rightMiddle, m.bottomRightOne, m.center,
				},
			},
			{
				Points: [3]rl.Vector2{
					m.bottomRightTwo, m.bottomMiddle, m.center,
				},
			},
			{
				Points: [3]rl.Vector2{
					m.bottomMiddle, m.bottomLeftOne, m.center,
				},
			},
			{
				Points: [3]rl.Vector2{
					m.bottomLeftTwo, m.leftMiddle, m.center,
				},
			},
			{
				Points: [3]rl.Vector2{
					m.leftMiddle, m.topLeftTwo, m.center,
				},
			},
		},
	}
}

func rotatePoint(origin, point *rl.Vector2, sin, cos float32) (float32, float32) {
	// Переносим точку в начало координат путем вычитания координат вершины A
	translatedX := point.X - origin.X
	translatedY := point.Y - origin.Y

	// Применяем поворот
	rotatedX := translatedX*cos - translatedY*sin
	rotatedY := translatedX*sin + translatedY*cos

	// Возвращаем точку с обратным переносом обратно в исходное положение
	return rotatedX + origin.X, rotatedY + origin.Y
}
