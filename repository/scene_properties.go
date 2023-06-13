package repository

func GetSceneProperties(levelName string) SceneProperties {
	return GetLevel(levelName).Properties
}
