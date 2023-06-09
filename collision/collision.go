package collision

import rl "github.com/gen2brain/raylib-go/raylib"

type Polygon struct {
	Points [3]rl.Vector2
}

type Hitbox struct {
	Polygons []Polygon
}

type CollisionDetector struct {
	Hitboxes []Hitbox
}

func (c *CollisionDetector) AddHitbox(h Hitbox) {
	c.Hitboxes = append(c.Hitboxes, h)
}

func (c CollisionDetector) Detect(collider Hitbox) (bool, map[int]rl.Vector2) {

	collisionedPolys := make(map[int]rl.Vector2, 0)

	for i, _ := range c.Hitboxes {
		hitbox := c.Hitboxes[i]
		for j, _ := range hitbox.Polygons {
			polygon := hitbox.Polygons[j]
			for k, _ := range collider.Polygons {
				mainPolygon := collider.Polygons[k]

				if rl.CheckCollisionPointTriangle(mainPolygon.Points[0], polygon.Points[0], polygon.Points[1], polygon.Points[2]) {
					collisionedPolys[k] = mainPolygon.Points[0]
				}
				if rl.CheckCollisionPointTriangle(mainPolygon.Points[1], polygon.Points[0], polygon.Points[1], polygon.Points[2]) {
					collisionedPolys[k] = mainPolygon.Points[1]
				}
				if rl.CheckCollisionPointTriangle(mainPolygon.Points[2], polygon.Points[0], polygon.Points[1], polygon.Points[2]) {
					collisionedPolys[k] = mainPolygon.Points[2]
				}
			}
		}
	}

	if len(collisionedPolys) != 0 {
		return true, collisionedPolys
	}

	return false, nil
}
