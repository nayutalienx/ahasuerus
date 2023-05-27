package repository

type Color struct {
	R int `json:"R"`
	G int `json:"G"`
	B int `json:"B"`
	A int `json:"A"`
}

type Rectangle struct {
	Id     string `json:"id"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Color  Color  `json:"color"`
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

type Bezier struct {
	Id     string `json:"id"`
	StartX int    `json:"startX"`
	StartY int    `json:"startY"`
	EndX   int    `json:"endX"`
	EndY   int    `json:"endY"`
	Thick  int    `json:"thick"`
	Color  Color  `json:"color"`
}

type Line struct {
	Id     string `json:"id"`
	StartX int    `json:"startX"`
	StartY int    `json:"startY"`
	EndX   int    `json:"endX"`
	EndY   int    `json:"endY"`
	Thick  int    `json:"thick"`
	Color  Color  `json:"color"`
}
