package utils

import (
	"fmt"
	"os"
	"strings"
)

func getK8sEnvVar(serviceName string) string {
	serviceName = strings.Replace(serviceName, " ", "_", -1)
	serviceName = strings.Replace(serviceName, "-", "_", -1)
	return strings.ToUpper(serviceName)
}

func GetK8sAddress(serviceName string) string {
	base := getK8sEnvVar(serviceName)
	host := os.Getenv(base + "_SERVICE_HOST")
	port := os.Getenv(base + "_SERVICE_PORT")
	return fmt.Sprintf("%s:%s", host, port)
}

func InKubernetes() bool {
	return os.Getenv("IN_KUBERNETES") == "true"
}
