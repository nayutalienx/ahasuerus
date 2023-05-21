package container

import (
	"ahasuerus/models"
)

type ObjectContainer struct {
	objects []models.Object
}

func NewObjectContainer() ObjectContainer {
	return ObjectContainer{
		objects: make([]models.Object, 0),
	}
}

func (w *ObjectContainer) AddObject(obj... models.Object) {
	for i, _ := range obj {
		o := obj[i]
		w.objects = append(w.objects, o)
	}
}

func (w *ObjectContainer) RemoveObject(obj models.Object) {
	for i, _ := range w.objects {
		o := w.objects[i]
		if o == obj {
			w.objects = w.removeObject(i)
		}
	}
}

func (w ObjectContainer) Draw() {
	for _, o := range w.objects {
		o.Draw()
	}
}

func (w ObjectContainer) Update(delta float32) {
	for _, o := range w.objects {
		o.Update(delta)
	}
}

func (w ObjectContainer) Size() int {
	return len(w.objects)
}

func (w ObjectContainer) removeObject(s int) []models.Object {
    return append(w.objects[:s], w.objects[s+1:]...)
}
