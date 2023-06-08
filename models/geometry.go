package models

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type CollisionBoxCheck struct {
	Intersected              bool
	Top, Bottom, Left, Right bool
}

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


func CalculatePreviousAndNextPointOfLine(point, start, end rl.Vector2) (rl.Vector2, rl.Vector2) {
	// Calculate the direction vector of the line
	direction := rl.Vector2{end.X - start.X, end.Y - start.Y}

	// Calculate the length of the direction vector
	length := direction.X*direction.X + direction.Y*direction.Y
	length = length * length

	// Calculate the normalized direction vector
	normalizedDirection := rl.Vector2{direction.X / length, direction.Y / length}

	// Calculate the previous point by subtracting the normalized direction vector from the current point
	prevPoint := rl.Vector2{point.X - normalizedDirection.X, point.Y - normalizedDirection.Y}

	// Calculate the next point by adding the normalized direction vector to the current point
	nextPoint := rl.Vector2{point.X + normalizedDirection.X, point.Y + normalizedDirection.Y}

	return prevPoint, nextPoint
}

func CalculateDistance(p1, p2 rl.Vector2) float64 {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y

	return math.Sqrt(float64(dx*dx + dy*dy))
}
