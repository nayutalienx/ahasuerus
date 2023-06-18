package repository

import (
	"ahasuerus/collision"
	"ahasuerus/models"
	"ahasuerus/resources"
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const dataId = "data"

type Level struct {
	Name               string
	Characters         []Character            `json:"characters"`
	Lights             []Light                `json:"lights"`
	CollissionHitboxes []CollisionHitbox      `json:"collissionHitboxes"`
	Images             []Image                `json:"images"`
	Properties         map[string]interface{} `json:"properties"`
}

func GetLevel(levelName string) Level {
	var level Level
	err := db.Read(levelName, dataId, &level)
	if err != nil {
		panic(err)
	}
	level.Name = levelName
	return level
}

func (level *Level) SaveLevel() {
	err := db.Write(level.Name, dataId, level)
	if err != nil {
		panic(err)
	}
}

func (l *Level) SaveImage(imageModel *models.Image) {
	img := l.mapImage(imageModel)
	found := false
	foundIndex := 0
	for index, _ := range l.Images {
		if l.Images[index].Id == img.Id {
			found = true
			foundIndex = index
			break
		}
	}

	if found {
		l.Images[foundIndex] = img
	} else {
		l.Images = append(l.Images, img)
	}
}

func (l *Level) DeleteImage(id string) *Level{
	found := false
	foundIndex := 0
	for index, _ := range l.Images {
		if l.Images[index].Id == id {
			found = true
			foundIndex = index
			break
		}
	}
	if found {
		l.Images = append(l.Images[:foundIndex], l.Images[foundIndex+1:]...)
	}
	return l
}

func (l *Level) SaveCollisionHitbox(hitboxModel *models.CollisionHitbox) {
	hitbox := l.mapCollisionHitbox(hitboxModel)
	found := false
	foundIndex := 0
	for index, _ := range l.CollissionHitboxes {
		if l.CollissionHitboxes[index].Id == hitbox.Id {
			found = true
			foundIndex = index
			break
		}
	}

	if found {
		l.CollissionHitboxes[foundIndex] = hitbox
	} else {
		l.CollissionHitboxes = append(l.CollissionHitboxes, hitbox)
	}
}

func (l *Level) DeleteCollisionHitbox(id string) *Level {
	found := false
	foundIndex := 0
	for index, _ := range l.CollissionHitboxes {
		if l.CollissionHitboxes[index].Id == id {
			found = true
			foundIndex = index
			break
		}
	}
	if found {
		l.CollissionHitboxes = append(l.CollissionHitboxes[:foundIndex], l.CollissionHitboxes[foundIndex+1:]...)
	}
	return l
}

func (l *Level) DeleteLight(id string) *Level {
	found := false
	foundIndex := 0
	for index, _ := range l.Lights {
		if l.Lights[index].Id == id {
			found = true
			foundIndex = index
			break
		}
	}
	if found {
		l.Lights = append(l.Lights[:foundIndex], l.Lights[foundIndex+1:]...)
	}
	return l
}

func (l *Level) DeleteCharacter(id string) *Level {
	found := false
	foundIndex := 0
	for index, _ := range l.Characters {
		if l.Characters[index].Id == id {
			found = true
			foundIndex = index
			break
		}
	}
	if found {
		l.Characters = append(l.Characters[:foundIndex], l.Characters[foundIndex+1:]...)
	}
	return l
}

func (level *Level) GetAllImages() []models.Image {
	images := []models.Image{}
	for i, _ := range level.Images {
		imageFound := level.Images[i]

		imageModel := models.NewImage(imageFound.DrawIndex, imageFound.Id, resources.GameTexture(imageFound.Path), float32(imageFound.X), float32(imageFound.Y), float32(imageFound.Width), float32(imageFound.Height), float32(imageFound.Rotation))
		imageModel.Parallax = float32(imageFound.Parallax)
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

func (level *Level) GetAllCollisionHitboxes() []models.CollisionHitbox {
	result := []models.CollisionHitbox{}
	for i, _ := range level.CollissionHitboxes {
		element := models.CollisionHitbox{
			BaseEditorItem: level.getBaseEditorItem(level.CollissionHitboxes[i].BaseEditorItem),
		}
		result = append(result, element)
	}
	return result
}

func (level *Level) GetAllLights() []models.Light {
	result := []models.Light{}
	for i, _ := range level.Lights {
		element := models.Light{
			BaseEditorItem: level.getBaseEditorItem(level.Lights[i].BaseEditorItem),
		}
		result = append(result, element)
	}
	return result
}

func (level *Level) GetAllCharacters() []models.Npc {
	result := []models.Npc{}
	for i, _ := range level.Characters {
		element := models.Npc{
			CollisionHitbox: models.CollisionHitbox{
				BaseEditorItem: level.getBaseEditorItem(level.Characters[i].BaseEditorItem),
			},
		}
		result = append(result, element)
	}
	return result
}

func (l *Level) getBaseEditorItem(item BaseEditorItem) models.BaseEditorItem {
	bei := models.BaseEditorItem{
		Id:         item.Id,
		Properties: item.Properties,
	}
	if bei.Properties == nil {
		bei.Properties = map[string]string{}
	}

	polys := [2]collision.Polygon{}

	for i, _ := range item.Polygons {
		polys[i] = collision.Polygon{
			Points: [3]rl.Vector2{
				{item.Polygons[i].Points[0].X, item.Polygons[i].Points[0].Y},
				{item.Polygons[i].Points[1].X, item.Polygons[i].Points[1].Y},
				{item.Polygons[i].Points[2].X, item.Polygons[i].Points[2].Y},
			},
		}
	}

	bei.SetPolygons(polys)
	return bei
}

func (level *Level) mapCollisionHitbox(hb *models.CollisionHitbox) CollisionHitbox {
	return CollisionHitbox{
		BaseEditorItem: level.mapBaseEditorItem(&hb.BaseEditorItem),
	}
}

func (level *Level) mapBaseEditorItem(bei *models.BaseEditorItem) BaseEditorItem {
	result := BaseEditorItem{
		Id:         bei.Id,
		Polygons:   []Polygon{},
		Properties: bei.Properties,
	}

	polygons := bei.Polygons()

	for i, _ := range polygons {
		p := polygons[i]
		result.Polygons = append(result.Polygons, Polygon{
			Points: [3]Vec2{
				{p.Points[0].X, p.Points[0].Y},
				{p.Points[1].X, p.Points[1].Y},
				{p.Points[2].X, p.Points[2].Y},
			},
		})
	}
	return result
}

func (l *Level) mapImage(img *models.Image) Image {
	i := Image{
		DrawIndex: img.DrawIndex,
		Id:        img.Id,
		Path:      string(img.ImageTexture),
		Shader:    string(img.ImageShader),
		X:         int(img.Pos.X),
		Y:         int(img.Pos.Y),
		Rotation:  int(img.Rotation),
		Parallax:  img.Parallax,
	}

	if img.WidthHeight.X > 0 && img.WidthHeight.Y > 0 {
		i.Width = int(img.WidthHeight.X)
		i.Height = int(img.WidthHeight.Y)
	}

	return i
}

type Color struct {
	R int `json:"R"`
	G int `json:"G"`
	B int `json:"B"`
	A int `json:"A"`
}

type Image struct {
	DrawIndex int     `json:"drawIndex"`
	Id        string  `json:"id"`
	Path      string  `json:"path"`
	Shader    string  `json:"shader"`
	X         int     `json:"x"`
	Y         int     `json:"y"`
	Width     int     `json:"width"`
	Height    int     `json:"height"`
	Rotation  int     `json:"rotation"`
	Parallax  float32 `json:"parallax"`
}

type Vec2 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Polygon struct {
	Points [3]Vec2 `json:"points"`
}

type BaseEditorItem struct {
	Id         string            `json:"id"`
	Polygons   []Polygon         `json:"polygons"`
	Properties map[string]string `json:"properties"`
}

type CollisionHitbox struct {
	BaseEditorItem
}

type Light struct {
	BaseEditorItem
}

type Character struct {
	BaseEditorItem
}
