package collision

import rl "github.com/gen2brain/raylib-go/raylib"

type Polygon struct {
	Points [3]rl.Vector2
}

type Hitbox struct {
	Polygons []Polygon
	Rotation float32
}

type CollisionDetector struct {
	Hitboxes []*Hitbox
}

func (c *CollisionDetector) AddHitbox(h *Hitbox) {
	c.Hitboxes = append(c.Hitboxes, h)
}

func (c CollisionDetector) Detect(collider Hitbox) (bool, []map[int]float32) {

	collisions := make([]map[int]float32, 0)

	for i, _ := range c.Hitboxes {
		hitbox := c.Hitboxes[i]
		for j, _ := range hitbox.Polygons {

			polygon := hitbox.Polygons[j]
			rotation := hitbox.Rotation
			collisionedPolys := make(map[int]float32, 0)
			for k, _ := range collider.Polygons {
				mainPolygon := collider.Polygons[k]

				if rl.CheckCollisionPointTriangle(mainPolygon.Points[0], polygon.Points[0], polygon.Points[1], polygon.Points[2]) {
					collisionedPolys[k] = rotation
				}
				if rl.CheckCollisionPointTriangle(mainPolygon.Points[1], polygon.Points[0], polygon.Points[1], polygon.Points[2]) {
					collisionedPolys[k] = rotation
				}
				if rl.CheckCollisionPointTriangle(mainPolygon.Points[2], polygon.Points[0], polygon.Points[1], polygon.Points[2]) {
					collisionedPolys[k] = rotation
				}
			}

			if len(collisionedPolys) != 0 {
				collisions = append(collisions, collisionedPolys)
			}

		}
	}

	if len(collisions) != 0 {
		return true, collisions
	}

	return false, nil
}
