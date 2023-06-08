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
