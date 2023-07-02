package models

var (
	wannaChangeScene string
)

func ChangeSceneAsync(id string) {
	wannaChangeScene = id
}

func IsWannaChangeScene() (bool, string) {
	if wannaChangeScene != "" {
		savedScene := wannaChangeScene
		wannaChangeScene = ""
		return true, savedScene
	}
	return false, ""
}
