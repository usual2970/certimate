package common

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
)

func IsEnvironmentPublic(env string) bool {
	switch strings.ToLower(env) {
	case "", "default", "public", "azurecloud":
		return true
	default:
		return false
	}
}

func IsEnvironmentGovernment(env string) bool {
	switch strings.ToLower(env) {
	case "usgovernment", "government", "azureusgovernment", "azuregovernment":
		return true
	default:
		return false
	}
}

func IsEnvironmentChina(env string) bool {
	switch strings.ToLower(env) {
	case "china", "chinacloud", "azurechina", "azurechinacloud":
		return true
	default:
		return false
	}
}

func GetCloudEnvironmentConfiguration(env string) (cloud.Configuration, error) {
	if IsEnvironmentPublic(env) {
		return cloud.AzurePublic, nil
	} else if IsEnvironmentGovernment(env) {
		return cloud.AzureGovernment, nil
	} else if IsEnvironmentChina(env) {
		return cloud.AzureChina, nil
	}

	return cloud.Configuration{}, fmt.Errorf("unknown azure cloud environment %s", env)
}
