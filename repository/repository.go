package repository

import (
	"ahasuerus/models"
	"encoding/json"
	"fmt"
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"
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

func SaveRectangle(collectionPrefix string, rect *models.Rectangle) {
	r := mapRectangle(rect.Id, rect)
	err := db.Write(formatKey(collectionPrefix, "rectangle"), r.Id, r)
	if err != nil {
		panic(err)
	}
}

func SaveBezier(collectionPrefix string, bez *models.Bezier) {
	r := mapBezier(bez.Id, bez)
	err := db.Write(formatKey(collectionPrefix, "bezier"), r.Id, r)
	if err != nil {
		panic(err)
	}
}

func SaveLine(collectionPrefix string, line *models.Line) {
	r := mapLine(line.Id, line)
	err := db.Write(formatKey(collectionPrefix, "line"), r.Id, r)
	if err != nil {
		panic(err)
	}
}

func SaveImage(collectionPrefix string, img *models.Image) {
	i := mapImage(img.Id, img)
	err := db.Write(formatKey(collectionPrefix, "image"), i.Id, i)
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

func GetAllBeziers(collectionPrefix string) []models.Bezier {
	records, err := db.ReadAll(formatKey(collectionPrefix, "bezier"))
	if err != nil {
		panic(err)
	}

	beziers := []models.Bezier{}
	for _, f := range records {
		bezFound := Bezier{}
		if err := json.Unmarshal([]byte(f), &bezFound); err != nil {
			panic(err)
		}
		c := bezFound.Color
		beziers = append(beziers, *models.NewBezier(
			bezFound.Id,
			rl.NewVector2(float32(bezFound.StartX), float32(bezFound.StartY)),
			rl.NewVector2(float32(bezFound.EndX), float32(bezFound.EndY)),
			float32(bezFound.Thick),
			rl.NewColor(uint8(c.R), uint8(c.G), uint8(c.B), uint8(c.A)),
		))
	}

	return beziers
}

func GetAllLines(collectionPrefix string) []models.Line {
	records, err := db.ReadAll(formatKey(collectionPrefix, "line"))
	if err != nil {
		panic(err)
	}

	lines := []models.Line{}
	for _, f := range records {
		lineFound := Line{}
		if err := json.Unmarshal([]byte(f), &lineFound); err != nil {
			panic(err)
		}
		c := lineFound.Color
		lines = append(lines, *models.NewLine(
			lineFound.Id,
			rl.NewVector2(float32(lineFound.StartX), float32(lineFound.StartY)),
			rl.NewVector2(float32(lineFound.EndX), float32(lineFound.EndY)),
			float32(lineFound.Thick),
			rl.NewColor(uint8(c.R), uint8(c.G), uint8(c.B), uint8(c.A)),
		))
	}

	return lines
}

func GetAllImages(collectionPrefix string) []models.Image {
	records, err := db.ReadAll(formatKey(collectionPrefix, "image"))
	if err != nil {
		panic(err)
	}

	images := []models.Image{}
	for _, f := range records {
		imageFound := Image{}
		if err := json.Unmarshal([]byte(f), &imageFound); err != nil {
			panic(err)
		}
		images = append(images, *models.NewImage(imageFound.DrawIndex, imageFound.Id, imageFound.Path, float32(imageFound.X), float32(imageFound.Y), imageFound.Scale))
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].DrawIndex < images[j].DrawIndex
	})

	return images
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

func mapBezier(id string, bez *models.Bezier) Bezier {
	r := Bezier{
		Id:     id,
		StartX: int(bez.Start.X),
		StartY: int(bez.Start.Y),
		EndX:   int(bez.End.X),
		EndY:   int(bez.End.Y),
		Thick:  int(bez.Thick),
		Color:  Color{R: int(bez.GetColor().R), G: int(bez.GetColor().G), B: int(bez.GetColor().B), A: int(bez.GetColor().A)},
	}
	return r
}

func mapLine(id string, bez *models.Line) Line {
	r := Line{
		Id:     id,
		StartX: int(bez.Start.X),
		StartY: int(bez.Start.Y),
		EndX:   int(bez.End.X),
		EndY:   int(bez.End.Y),
		Thick:  int(bez.Thick),
		Color:  Color{R: int(bez.GetColor().R), G: int(bez.GetColor().G), B: int(bez.GetColor().B), A: int(bez.GetColor().A)},
	}
	return r
}

func mapImage(id string, img *models.Image) Image {
	return Image{
		DrawIndex: img.DrawIndex,
		Id:        id,
		Path:      img.ResourcePath,
		X:         int(img.Pos.X),
		Y:         int(img.Pos.Y),
		Scale:     img.ScaleTex,
	}
}

func formatKey(collectionPrefix, entity string) string {
	return fmt.Sprintf("%s-%s", collectionPrefix, entity)
}
