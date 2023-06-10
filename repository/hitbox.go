package repository

import (
	"ahasuerus/collision"
	"ahasuerus/models"
	"encoding/json"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func SaveHitbox(collectionPrefix string, container string, hb *models.Hitbox) {
	i := mapHitbox(hb.Id, hb)
	err := db.Write(formatKey(collectionPrefix, fmt.Sprintf("hitbox-%s", container)), i.Id, i)
	if err != nil {
		panic(err)
	}
}

func DeleteHitbox(collectionPrefix string, container string, hb *models.Hitbox) {
	err := db.Delete(formatKey(collectionPrefix, fmt.Sprintf("hitbox-%s", container)), hb.Id)
	if err != nil {
		panic(err)
	}
}

func GetAllHitboxes(collectionPrefix string, container string) []models.Hitbox {
	records, err := db.ReadAll(formatKey(collectionPrefix, fmt.Sprintf("hitbox-%s", container)))
	if err != nil {
		panic(err)
	}

	hitboxes := []models.Hitbox{}
	for _, f := range records {
		hitboxFound := Hitbox{}
		if err := json.Unmarshal([]byte(f), &hitboxFound); err != nil {
			panic(err)
		}
		hb := models.Hitbox{
			Id: hitboxFound.Id,
		}

		polys := [2]collision.Polygon{}

		for i, _ := range hitboxFound.Polygons {
			polys[i] = collision.Polygon{
				Points: [3]rl.Vector2{
					{hitboxFound.Polygons[i].Points[0].X, hitboxFound.Polygons[i].Points[0].Y},
					{hitboxFound.Polygons[i].Points[1].X, hitboxFound.Polygons[i].Points[1].Y},
					{hitboxFound.Polygons[i].Points[2].X, hitboxFound.Polygons[i].Points[2].Y},
				},
			}
		}

		hb.SetPolygons(polys)

		hitboxes = append(hitboxes, hb)
	}
	return hitboxes
}

func mapHitbox(id string, hb *models.Hitbox) Hitbox {
	hitbox := Hitbox{
		Id:       id,
		Polygons: []Polygon{},
	}

	polygons := hb.Polygons()

	for i, _ := range polygons {
		p := polygons[i]
		hitbox.Polygons = append(hitbox.Polygons, Polygon{
			Points: [3]Vec2{
				{p.Points[0].X, p.Points[0].Y},
				{p.Points[1].X, p.Points[1].Y},
				{p.Points[2].X, p.Points[2].Y},
			},
		})
	}

	return hitbox
}
