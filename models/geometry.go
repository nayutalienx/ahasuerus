package models

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type CollisionBoxCheck struct {
	Intersected              bool
	Top, Bottom, Left, Right bool
	X                        BoxPosition
	Y                        BoxPosition
}

func DetectBoxCollision(boxPosition1, boxPosition2 BoxPosition) CollisionBoxCheck {
	pos1 := boxPosition1.GetPos()
	pos2 := boxPosition2.GetPos()

	box1 := boxPosition1.GetBox()
	box2 := boxPosition2.GetBox()

	rect1 := rl.NewRectangle(pos1.X, pos1.Y, box1.X, box1.Y)
	rect2 := rl.NewRectangle(pos2.X, pos2.Y, box2.X, box2.Y)

	collision := CollisionBoxCheck{
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

type CollisionBezierCheck struct {
	Point     rl.Vector2
	Colliding bool
	Curve     *Bezier
}

func CheckCollisionLineBezier(lineStartPos, lineEndPos, bezierStartPos, bezierEndPos rl.Vector2, thickness float64) CollisionBezierCheck {
	// Calculate the control points of the cubic-bezier curve with ease-in-out motion
	cp1 := rl.Vector2{X: bezierStartPos.X + (bezierEndPos.X-bezierStartPos.X)/3, Y: bezierStartPos.Y}
	cp2 := rl.Vector2{X: bezierStartPos.X + (bezierEndPos.X-bezierStartPos.X)*2/3, Y: bezierEndPos.Y}

	// Check collision between the line segment and the cubic-bezier curve
	return CheckLineBezierCollision(lineStartPos, lineEndPos, bezierStartPos, cp1, cp2, bezierEndPos, thickness)
}

func CheckLineBezierCollision(lineStartPos, lineEndPos, p0, p1, p2, p3 rl.Vector2, thickness float64) CollisionBezierCheck {
	// Calculate the line direction and length
	lineDir := rl.Vector2{X: lineEndPos.X - lineStartPos.X, Y: lineEndPos.Y - lineStartPos.Y}
	lineLength := math.Sqrt(math.Pow(float64(lineDir.X), 2) + math.Pow(float64(lineDir.Y), 2))

	// Normalize the line direction
	lineDir.X /= float32(lineLength)
	lineDir.Y /= float32(lineLength)

	// Perform collision check using ray casting algorithm
	for t := 0.0; t <= 1.0; t += 0.01 {
		// Calculate the point on the cubic-bezier curve
		bezierPoint := CalculateCubicBezierPoint(t, p0, p1, p2, p3)

		// Calculate the vector from the line start position to the bezier point
		lineToPoint := rl.Vector2{X: bezierPoint.X - lineStartPos.X, Y: bezierPoint.Y - lineStartPos.Y}

		// Calculate the dot product between the line direction and line-to-point vector
		dotProduct := lineDir.X*lineToPoint.X + lineDir.Y*lineToPoint.Y

		// Check if the point is in front of the line segment
		if dotProduct >= 0 && dotProduct <= float32(lineLength) {
			// Calculate the distance between the line and the bezier point
			distance := math.Abs(float64(lineDir.X*lineToPoint.Y - lineDir.Y*lineToPoint.X))

			// Check if the distance is within a threshold (collision threshold)
			if distance <= thickness/2 { // Adjust the threshold as needed
				return CollisionBezierCheck{Point: bezierPoint, Colliding: true}
			}
		}
	}

	return CollisionBezierCheck{Colliding: false}
}

func CalculateCubicBezierPoint(t float64, p0, p1, p2, p3 rl.Vector2) rl.Vector2 {
	u := 1 - t
	tt := t * t
	uu := u * u
	uuu := uu * u
	ttt := tt * t

	p := rl.Vector2{
		X: float32(uuu*float64(p0.X) + 3*uu*t*float64(p1.X) + 3*u*tt*float64(p2.X) + ttt*float64(p3.X)),
		Y: float32(uuu*float64(p0.Y) + 3*uu*t*float64(p1.Y) + 3*u*tt*float64(p2.Y) + ttt*float64(p3.Y)),
	}

	return p
}

func CalculatePreviousNextPoints(concretePoint, bezierStartPos, bezierEndPos rl.Vector2) (rl.Vector2, rl.Vector2) {

	// Calculate the control points of the cubic-bezier curve with ease-in-out motion
	cp1 := rl.Vector2{X: bezierStartPos.X + (bezierEndPos.X-bezierStartPos.X)/3, Y: bezierStartPos.Y}
	cp2 := rl.Vector2{X: bezierStartPos.X + (bezierEndPos.X-bezierStartPos.X)*2/3, Y: bezierEndPos.Y}

	return calculatePreviousNextPoints(concretePoint, bezierStartPos, cp1, cp2, bezierEndPos)
}

func calculatePreviousNextPoints(concretePoint, p0, p1, p2, p3 rl.Vector2) (rl.Vector2, rl.Vector2) {
    // Calculate the t parameter for the concrete point on the curve
    t := FindTParameter(concretePoint, p0, p1, p2, p3)

    // Calculate the t parameter for the previous point
    tPrev := math.Max(t-0.01, 0.0) // Adjust the step size as needed

    // Calculate the t parameter for the next point
    tNext := math.Min(t+0.01, 1.0) // Adjust the step size as needed

    // Calculate the previous and next points on the cubic-bezier curve
    prevPoint := CalculateCubicBezierPoint(tPrev, p0, p1, p2, p3)
    nextPoint := CalculateCubicBezierPoint(tNext, p0, p1, p2, p3)

    return prevPoint, nextPoint
}

func FindTParameter(concretePoint, p0, p1, p2, p3 rl.Vector2) float64 {
    // Iterate and find the t parameter that yields the closest point to the concrete point
    minDistance := math.Inf(1)
    bestT := 0.0

    for t := 0.0; t <= 1.0; t += 0.001 { // Adjust the step size as needed
        bezierPoint := CalculateCubicBezierPoint(t, p0, p1, p2, p3)
        distance := CalculateDistance(bezierPoint, concretePoint)

        if distance < minDistance {
            minDistance = distance
            bestT = t
        }
    }

    return bestT
}

func CalculateDistance(p1, p2 rl.Vector2) float64 {
    dx := p2.X - p1.X
    dy := p2.Y - p1.Y

    return math.Sqrt(float64(dx*dx + dy*dy))
}