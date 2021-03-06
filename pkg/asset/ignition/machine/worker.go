package machine

import (
	"encoding/json"

	igntypes "github.com/coreos/ignition/config/v2_2/types"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
)

// Worker is an asset that generates the ignition config for worker nodes.
type Worker struct {
	config *igntypes.Config
	file   *asset.File
}

var _ asset.WritableAsset = (*Worker)(nil)

// Dependencies returns the assets on which the Worker asset depends.
func (a *Worker) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&tls.RootCA{},
	}
}

// Generate generates the ignition config for the Worker asset.
func (a *Worker) Generate(dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	rootCA := &tls.RootCA{}
	dependencies.Get(installConfig, rootCA)

	a.config = pointerIgnitionConfig(installConfig.Config, rootCA.Cert(), "worker", "")

	data, err := json.Marshal(a.config)
	if err != nil {
		return errors.Wrap(err, "failed to get InstallConfig from parents")
	}
	a.file = &asset.File{
		Filename: "worker.ign",
		Data:     data,
	}

	return nil
}

// Name returns the human-friendly name of the asset.
func (a *Worker) Name() string {
	return "Worker Ignition Config"
}

// Files returns the files generated by the asset.
func (a *Worker) Files() []*asset.File {
	return []*asset.File{a.file}
}
