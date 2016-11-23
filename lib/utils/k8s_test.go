package utils

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("ASD_QWE_SERVICE_HOST", "10.0.0.1")
	os.Setenv("ASD_QWE_SERVICE_PORT", "12345")
	m.Run()
}

func TestGetK8sAddress(t *testing.T) {
	serviceName := "asd-qwe"
	baseVar := getK8sEnvVar(serviceName)
	if baseVar != "ASD_QWE" {
		t.Fail()
	}
	addr := GetK8sAddress(serviceName)
	if addr != "10.0.0.1:12345" {
		t.Fail()
	}
}
