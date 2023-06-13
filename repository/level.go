package repository

const dataId = "data"

func GetLevel(levelName string) Level {
	var level Level
	err := db.Read(levelName, dataId, &level)
	if err != nil {
		panic(err)
	}
	return level
}

func SaveLevel(levelName string, level Level) {
	err := db.Write(levelName, dataId, level)
	if err != nil {
		panic(err)
	}
}
