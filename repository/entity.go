package repository

type Color struct {
	R int `json:"R"`
	G int `json:"G"`
	B int `json:"B"`
	A int `json:"A"`
}

type Image struct {
	DrawIndex int    `json:"drawIndex"`
	Id        string `json:"id"`
	Path      string `json:"path"`
	Shader    string `json:"shader"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Rotation  int    `json:"rotation"`
}

type Vec2 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Polygon struct {
	Points [3]Vec2 `json:"points"`
}

type Hitbox struct {
	Id       string    `json:"id"`
	Polygons []Polygon `json:"polygons"`
	Type     int       `json:"type"`
}

type SceneProperties map[string]interface{}
