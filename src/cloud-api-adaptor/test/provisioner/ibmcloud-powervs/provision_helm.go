//go:build ibmcloud_powervs

// (C) Copyright Confidential Containers Contributors
// SPDX-License-Identifier: Apache-2.0

package ibmcloud_powervs

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	pv "github.com/confidential-containers/cloud-api-adaptor/src/cloud-api-adaptor/test/provisioner"
	byomprov "github.com/confidential-containers/cloud-api-adaptor/src/cloud-api-adaptor/test/provisioner/byom"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func NewIBMCloudPowerVSInstallChart(installDir, provider string) (pv.InstallChart, error) {
	chartPath := filepath.Join(installDir, "charts", "peerpods")
	namespace := pv.GetCAANamespace()
	releaseName := "peerpods"
	debug := false

	// The providers file on disk is ibmcloudpowervs.yaml and the chart provider value is
	// also "ibmcloudpowervs", matching what CLOUD_PROVIDER is set to when running tests.
	helm, err := pv.NewHelm(chartPath, namespace, releaseName, "ibmcloudpowervs", debug)
	if err != nil {
		return nil, err
	}

	return &IBMCloudPowerVSInstallChart{
		Helm: helm,
	}, nil
}

type IBMCloudPowerVSInstallChart struct {
	Helm *pv.Helm
}

func (ic *IBMCloudPowerVSInstallChart) GetHelm() *pv.Helm {
	return ic.Helm
}

// Install creates the SSH key secret required by CAA, then installs the Helm
// chart and waits for the worker node to receive the kata-runtime label.
func (ic *IBMCloudPowerVSInstallChart) Install(ctx context.Context, cfg *envconf.Config) error {
	if err := ic.createSSHKeySecret(ctx, cfg); err != nil {
		return fmt.Errorf("failed to create SSH key secret: %w", err)
	}

	if err := ic.Helm.Install(ctx, cfg); err != nil {
		return err
	}
	return nil
}

// createSSHKeySecret creates the ssh-key-secret in the CAA namespace from the
// SSH key paths configured in the BYOM properties.
func (ic *IBMCloudPowerVSInstallChart) createSSHKeySecret(ctx context.Context, cfg *envconf.Config) error {
	privKeyPath := byomprov.ByomProps.SSHSecretPrivKeyPath
	pubKeyPath := byomprov.ByomProps.SSHSecretPubKeyPath

	if privKeyPath == "" || pubKeyPath == "" {
		return fmt.Errorf("SSH_SECRET_PRIV_KEY_PATH and SSH_SECRET_PUB_KEY_PATH must be set")
	}

	// Ensure the namespace exists before creating the secret.
	if err := pv.CreateAndWaitForNamespace(ctx, cfg.Client(), ic.Helm.Namespace); err != nil {
		return fmt.Errorf("failed to create namespace: %w", err)
	}

	secretName := "ssh-key-secret"
	log.Infof("Creating SSH key secret from %s and %s", privKeyPath, pubKeyPath)

	args := []string{
		"create", "secret", "generic", secretName,
		"--from-file=id_rsa=" + privKeyPath,
		"--from-file=id_rsa.pub=" + pubKeyPath,
		"-n", ic.Helm.Namespace,
	}
	cmd := exec.Command("kubectl", args...)
	cmd.Env = append(os.Environ(), "KUBECONFIG="+cfg.KubeconfigFile())
	stdoutStderr, err := cmd.CombinedOutput()
	log.Tracef("%v, output: %s", cmd, stdoutStderr)
	if err != nil {
		return fmt.Errorf("failed to create ssh-key-secret: %w, output: %s", err, string(stdoutStderr))
	}

	log.Infof("Created ssh-key-secret successfully")
	return nil
}

func (ic *IBMCloudPowerVSInstallChart) Uninstall(ctx context.Context, cfg *envconf.Config) error {
	return ic.Helm.Uninstall(ctx, cfg)
}

func (ic *IBMCloudPowerVSInstallChart) Configure(ctx context.Context, cfg *envconf.Config, properties map[string]string) error {
	// Must match the provider/providerConfigs/providerSecrets key in ibmcloudpowervs.yaml.
	const chartProvider = "ibmcloud-powervs"

	// Set IBM Cloud API key as a provider secret.
	if key := properties["IBMCLOUD_API_KEY"]; key != "" {
		ic.Helm.OverrideValues[fmt.Sprintf("providerSecrets.%s.IBMCLOUD_API_KEY", chartProvider)] = key
	}

	// Set PowerVS-specific provider config values from properties.
	providerKeys := []string{
		"POWERVS_ZONE",
		"POWERVS_SERVICE_INSTANCE_ID",
		"POWERVS_IMAGE_ID",
		"POWERVS_NETWORK_ID",
		"POWERVS_SSH_KEY_NAME",
		"POWERVS_MEMORY",
		"POWERVS_PROCESSORS",
		"POWERVS_PROCESSOR_TYPE",
		"POWERVS_SYSTEM_TYPE",
	}
	for _, k := range providerKeys {
		if v := properties[k]; v != "" {
			ic.Helm.OverrideValues[fmt.Sprintf("providerConfigs.%s.%s", chartProvider, k)] = v
		}
	}

	return nil
}
