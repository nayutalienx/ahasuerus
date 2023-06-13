package repository

import (
	"ahasuerus/models"
	"ahasuerus/resources"
	"sort"
)

const imageDir = "image"

func SaveImage(levelName string, img *models.Image) {
	i := mapImage(img.Id, img)
	level := GetLevel(levelName)
	level.SaveImage(i)
	SaveLevel(levelName, level)
}

func DeleteImage(levelName string, img *models.Image) {
	level := GetLevel(levelName)
	level.DeleteImage(img.Id)
	SaveLevel(levelName, level)
}

func GetAllImages(levelName string) []models.Image {
	level := GetLevel(levelName)

	images := []models.Image{}
	for i, _ := range level.Images {
		imageFound := level.Images[i]

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
