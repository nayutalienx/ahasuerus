package repository

import (
	"ahasuerus/models"
	"encoding/json"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
	"github.com/sdomino/scribble"
)

var (
	db *scribble.Driver
)

func init() {
	var err error
	db, err = scribble.New("data", nil)
	if err != nil {
		panic(err)
	}
}

func AddNewRectangle(collectionPrefix string, rect *models.Rectangle) {
	r := mapRectangle(uuid.NewString(), rect)
	err := db.Write(formatKey(collectionPrefix, "rectangle"), r.Id, r)
	if err != nil {
		panic(err)
	}
}

func SaveRectangle(collectionPrefix string, rect *models.Rectangle) {
	r := mapRectangle(rect.Id, rect)
	err := db.Write(formatKey(collectionPrefix, "rectangle"), r.Id, r)
	if err != nil {
		panic(err)
	}
}

func GetAllRectangles(collectionPrefix string) []models.Rectangle {
	records, err := db.ReadAll(formatKey(collectionPrefix, "rectangle"))
	if err != nil {
		panic(err)
	}

	rectangles := []models.Rectangle{}
	for _, f := range records {
		rectFound := Rectangle{}
		if err := json.Unmarshal([]byte(f), &rectFound); err != nil {
			panic(err)
		}
		c := rectFound.Color
		rectangles = append(rectangles, *models.NewRectangle(
			rectFound.Id,
			float32(rectFound.X),
			float32(rectFound.Y),
			float32(rectFound.Width),
			float32(rectFound.Height),
			rl.NewColor(uint8(c.R), uint8(c.G), uint8(c.B), uint8(c.A)),
		))
	}

	return rectangles
}

func mapRectangle(id string, rect *models.Rectangle) Rectangle {
	r := Rectangle{
		Id:     id,
		X:      int(rect.GetPos().X),
		Y:      int(rect.GetPos().Y),
		Width:  int(rect.GetBox().X),
		Height: int(rect.GetBox().Y),
		Color: Color{
			R: int(rect.GetColor().R),
			G: int(rect.GetColor().G),
			B: int(rect.GetColor().B),
			A: int(rect.GetColor().A),
		},
	}
	return r
}

func formatKey(collectionPrefix, entity string) string {
	return fmt.Sprintf("%s-%s", collectionPrefix, entity)
}