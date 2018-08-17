package options

import (
	"bonc.com/lvm-scheduler/cmd/lvm-scheduler/app/config"
	"github.com/spf13/pflag"
)

type LvmConfigMap struct {
	Kubecfg               string
	LvmConfigMapName      string
	LvmConfigMapNamespace string
}

// AddFlags adds flags for the deprecated options.
func (lvm *LvmConfigMap) AddFlags(fs *pflag.FlagSet) {
	if lvm == nil {
		return
	}

	fs.StringVar(&lvm.Kubecfg, "kubecfg", lvm.Kubecfg, "define the kubeconfig of the k8s-cluster if run this plugin out of cluster.")
	fs.StringVar(&lvm.LvmConfigMapName, "lvm-configmap-name", lvm.LvmConfigMapName, "define the name of the lvm resource ConfigMap.")
	fs.StringVar(&lvm.LvmConfigMapNamespace, "lvm-configmap-namespace", lvm.LvmConfigMapNamespace, "define the namespace of the lvm resource ConfigMap.")
}

// Validate validates the deprecated scheduler options.
func (lvm *LvmConfigMap) Validate() []error {
	if lvm == nil {
		return nil
	}
	return nil
}

func (lvm *LvmConfigMap) ApplyTo(lo *config.AlgorithmOptionalOptional) error {
	if lvm == nil {
		return nil
	}
	if lvm.LvmConfigMapName != "" {
		lo.LvmConfigMapName = lvm.LvmConfigMapName
	}
	if lvm.LvmConfigMapNamespace != "" {
		lo.LvmConfigMapNamespace = lvm.LvmConfigMapNamespace
	}
	lo.Kubecfg = lvm.Kubecfg

	return nil
}
