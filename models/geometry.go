package models

import (
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
