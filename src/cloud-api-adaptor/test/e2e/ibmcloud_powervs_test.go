//go:build ibmcloud_powervs

// (C) Copyright Confidential Containers Contributors
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"testing"

	_ "github.com/confidential-containers/cloud-api-adaptor/src/cloud-api-adaptor/test/provisioner/ibmcloud-powervs"
)

func TestCreateSimplePod(t *testing.T) {
	assert := IBMCloudPowerVSAssert{}
	DoTestCreateSimplePod(t, testEnv, assert)
}

func TestCreatePodWithConfigMap(t *testing.T) {
	assert := IBMCloudPowerVSAssert{}
	DoTestCreatePodWithConfigMap(t, testEnv, assert)
}

func TestCreatePodWithSecret(t *testing.T) {
	assert := IBMCloudPowerVSAssert{}
	DoTestCreatePodWithSecret(t, testEnv, assert)
}

func TestDeletePod(t *testing.T) {
	assert := IBMCloudPowerVSAssert{}
	DoTestDeleteSimplePod(t, testEnv, assert)
}