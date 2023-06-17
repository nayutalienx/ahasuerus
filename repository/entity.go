package repository

type Level struct {
	Hitboxes   []Hitbox        `json:"hitboxes"`
	Images     []Image         `json:"images"`
	Properties SceneProperties `json:"properties"`
}

func (l *Level) SaveImage(img Image) {
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

func (l *Level) DeleteImage(id string) {
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
}

func (l *Level) SaveHitbox(hitbox Hitbox) {
	found := false
	foundIndex := 0
	for index, _ := range l.Hitboxes {
		if l.Hitboxes[index].Id == hitbox.Id {
			found = true
			foundIndex = index
			break
		}
	}

	if found {
		l.Hitboxes[foundIndex] = hitbox
	} else {
		l.Hitboxes = append(l.Hitboxes, hitbox)
	}
}

func (l *Level) DeleteHitbox(id string) {
	found := false
	foundIndex := 0
	for index, _ := range l.Hitboxes {
		if l.Hitboxes[index].Id == id {
			found = true
			foundIndex = index
			break
		}
	}
	if found {
		l.Hitboxes = append(l.Hitboxes[:foundIndex], l.Hitboxes[foundIndex+1:]...)
	}
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
	Particles string  `json:"particles"`
}

type Vec2 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Polygon struct {
	Points [3]Vec2 `json:"points"`
}

type Hitbox struct {
	Id         string            `json:"id"`
	Polygons   []Polygon         `json:"polygons"`
	Type       int               `json:"type"`
	Properties map[string]string `json:"properties"`
}

type SceneProperties map[string]interface{}
