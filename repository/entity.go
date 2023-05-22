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

type Bezier struct {
	Id     string `json:"id"`
	StartX int    `json:"startX"`
	StartY int    `json:"startY"`
	EndX   int    `json:"endX"`
	EndY   int    `json:"endY"`
	Thick  int    `json:"thick"`
	Color  Color  `json:"color"`
}