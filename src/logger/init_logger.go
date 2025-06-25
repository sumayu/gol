package logger

import (
	configs "main/src/backend/backendCFG"
	"os"
	"strings"
)

func Init() {
InitLogger("prod")
}
func getConfigPath() (string) {
	if IsDocker()  == true {
return	"/app/configYML/config.yaml"
	}
		return  "../../../configYML/config.yaml"
	
}
func LoadConfig() (*configs.Сonfig, error) {
    path := getConfigPath()
    return configs.LoadConfigs(path)
}

func IsDocker() bool {
    if _, err := os.Stat("/.dockerenv"); err == nil {
        return true
    }

    data, err := os.ReadFile("/proc/1/cgroup")
    if err == nil && (strings.Contains(string(data), "docker") || strings.Contains(string(data), "kubepods")) {
        return true
    }

    return os.Getenv("IN_DOCKER") == "true"
}