package repository

import (
	"ahasuerus/collision"
	"ahasuerus/models"
	"encoding/json"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const hitboxDir = "hitbox"

func SaveHitbox(collectionPrefix string, hb *models.Hitbox) {
	i := mapHitbox(hb.Id, hb)
	err := db.Write(formatKey(collectionPrefix, hitboxDir), i.Id, i)
	if err != nil {
		panic(err)
	}
}

func DeleteHitbox(collectionPrefix string, id string) {
	err := db.Delete(formatKey(collectionPrefix, hitboxDir), id)
	if err != nil {
		panic(err)
	}
}

func GetAllHitboxes(collectionPrefix string) []models.Hitbox {
	records, err := db.ReadAll(formatKey(collectionPrefix, hitboxDir))
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
			BaseEditorItem: models.BaseEditorItem{
				Id:         hitboxFound.Id,
				Properties: hitboxFound.Properties,
			},
			Type: models.HitboxType(hitboxFound.Type),
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
		Id:         id,
		Polygons:   []Polygon{},
		Type:       int(hb.Type),
		Properties: hb.Properties,
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
