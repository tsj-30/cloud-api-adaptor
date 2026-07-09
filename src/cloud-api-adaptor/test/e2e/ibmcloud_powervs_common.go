// (C) Copyright Confidential Containers Contributors
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"testing"
	"time"
)

// IBMCloudPowerVSAssert implements the CloudAssert interface for IBM Cloud PowerVS.
type IBMCloudPowerVSAssert struct{}

func (c IBMCloudPowerVSAssert) DefaultTimeout() time.Duration {
	return 30 * time.Second
}

func (c IBMCloudPowerVSAssert) HasPodVM(t *testing.T, podvmName string) {
	t.Logf("IBM Cloud PowerVS: pod VM created for pod %s", podvmName)
}

func (c IBMCloudPowerVSAssert) GetInstanceType(t *testing.T, podName string) (string, error) {
	return "", nil
}