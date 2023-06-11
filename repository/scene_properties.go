package repository

import "encoding/json"

const propertiesDir = "properties"

func GetSceneProperties(collectionPrefix string) SceneProperties {
	records, err := db.ReadAll(formatKey(collectionPrefix, propertiesDir))
	if err != nil {
		panic(err)
	}
	properties := map[string]interface{}{}
	if err := json.Unmarshal([]byte(records[0]), &properties); err != nil {
		panic(err)
	}
	return properties
}
