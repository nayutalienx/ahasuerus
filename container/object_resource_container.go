package container

import (
	"ahasuerus/models"
)

type ObjectResourceContainer struct {
	ObjectContainer
	objectResources []models.ObjectResource
}

func NewObjectResourceContainer() *ObjectResourceContainer {
	return &ObjectResourceContainer{
		objectResources: make([]models.ObjectResource, 0),
	}
}

func (w *ObjectResourceContainer) AddObjectResource(obj... models.ObjectResource) {
	for i, _ := range obj {
		o := obj[i]
		w.objectResources = append(w.objectResources, o)
		w.AddObject(o)
	}
}

func (w ObjectResourceContainer) Load() {
	for _, o := range w.objectResources {
		o.Load()
	}
}

func (w ObjectResourceContainer) Pause() {
	for _, o := range w.objectResources {
		o.Pause()
	}
}

func (w ObjectResourceContainer) Resume() {
	for _, o := range w.objectResources {
		o.Resume()
	}
}

func (w ObjectResourceContainer) Unload() {
	for _, o := range w.objectResources {
		o.Unload()
	}
}
