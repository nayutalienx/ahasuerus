package repository

import (
	"ahasuerus/models"
	"ahasuerus/resources"
	"encoding/json"
	"sort"
)

const imageDir = "image"

func SaveImage(collectionPrefix string, img *models.Image) {
	i := mapImage(img.Id, img)
	err := db.Write(formatKey(collectionPrefix, imageDir), i.Id, i)
	if err != nil {
		panic(err)
	}
}

func DeleteImage(collectionPrefix string, img *models.Image) {
	err := db.Delete(formatKey(collectionPrefix, imageDir), img.Id)
	if err != nil {
		panic(err)
	}
}

func GetAllImages(collectionPrefix string) []models.Image {
	records, err := db.ReadAll(formatKey(collectionPrefix, imageDir))
	if err != nil {
		panic(err)
	}

	images := []models.Image{}
	for _, f := range records {
		imageFound := Image{}
		if err := json.Unmarshal([]byte(f), &imageFound); err != nil {
			panic(err)
		}

		imageModel := models.NewImage(imageFound.DrawIndex, imageFound.Id, resources.GameTexture(imageFound.Path), float32(imageFound.X), float32(imageFound.Y), float32(imageFound.Width), float32(imageFound.Height), float32(imageFound.Rotation))
		if imageFound.Shader != string(resources.UndefinedShader) {
			imageModel.WithShader(resources.GameShader(imageFound.Shader))
		}

		images = append(images, *imageModel)
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].DrawIndex < images[j].DrawIndex
	})

	return images
}

func mapImage(id string, img *models.Image) Image {
	i := Image{
		DrawIndex: img.DrawIndex,
		Id:        id,
		Path:      string(img.ImageTexture),
		Shader:    string(img.ImageShader),
		X:         int(img.Pos.X),
		Y:         int(img.Pos.Y),
		Rotation:  int(img.Rotation),
	}

	if img.WidthHeight.X > 0 && img.WidthHeight.Y > 0 {
		i.Width = int(img.WidthHeight.X)
		i.Height = int(img.WidthHeight.Y)
	}

	return i
}
