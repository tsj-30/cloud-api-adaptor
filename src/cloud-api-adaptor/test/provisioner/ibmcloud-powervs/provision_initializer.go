// (C) Copyright Confidential Containers Contributors
// SPDX-License-Identifier: Apache-2.0

package ibmcloud_powervs

import (
	"fmt"
	"os"
	"strings"
)

type IBMPowerVSProperties struct {
	ApiKey                   string
	TestPodVMImage           string
	PowerVSZone              string
	PowerVSRegion            string
	PowerVSServiceInstanceID string
	PowerVSImageID           string
	PowerVSNetworkID         string
	PowerVSSSHKeyName        string
	PowerVSProcessorType     string
	PowerVSSystemType        string
	PowerVSMemory            string
	PowerVSProcessors        string
}

var IBMPowerVSProps = &IBMPowerVSProperties{}

func InitIBMCloudPowerVSProperties(properties map[string]string) error {

	memory := properties["POWERVS_MEMORY"]
	if memory == "" {
		memory = "2"
	}

	processors := properties["POWERVS_PROCESSORS"]
	if processors == "" {
		processors = "0.25"
	}

	IBMPowerVSProps = &IBMPowerVSProperties{
		ApiKey:                   properties["IBMCLOUD_API_KEY"],
		PowerVSRegion:            properties["POWERVS_REGION"],
		PowerVSZone:              properties["POWERVS_ZONE"],
		PowerVSServiceInstanceID: properties["POWERVS_SERVICE_INSTANCE_ID"],
		PowerVSImageID:           properties["POWERVS_IMAGE_ID"],
		PowerVSNetworkID:         properties["POWERVS_NETWORK_ID"],
		PowerVSSSHKeyName:        properties["POWERVS_SSH_KEY_NAME"],
		PowerVSProcessorType:     properties["POWERVS_PROCESSOR_TYPE"],
		PowerVSSystemType:        properties["POWERVS_SYSTEM_TYPE"],
		PowerVSMemory:            memory,
		PowerVSProcessors:        processors,
	}

	needProvisionStr := os.Getenv("TEST_PROVISION")
	if strings.EqualFold(needProvisionStr, "yes") || strings.EqualFold(needProvisionStr, "true") {
		required := []string{
			"IBMCLOUD_API_KEY",
			"POWERVS_ZONE",
			"POWERVS_SERVICE_INSTANCE_ID",
			"POWERVS_IMAGE_ID",
			"POWERVS_NETWORK_ID",
			"POWERVS_SSH_KEY_NAME",
		}
		for _, key := range required {
			if properties[key] == "" {
				return fmt.Errorf("required property %s is not set", key)
			}
		}
	}
	return nil
}
