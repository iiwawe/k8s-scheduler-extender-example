package algorithm

import (
	"encoding/json"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/pkg/scheduler/algorithm"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/predicates"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"

	"bonc.com/k8s-scheduler-extender-example/cmd/lvm-scheduler/app/config"
	"bonc.com/k8s-scheduler-extender-example/pkg/lvmutils"
)

const (
	CheckLvmVolumePred = "CheckLvmVolume"
)

type LvmChecker struct {
	optional config.AlgorithmOptionalOptional
	pvcInfo  predicates.PersistentVolumeClaimInfo
}

// NewNodeLabelPredicate creates a predicate which evaluates whether a pod can fit based on the
// node labels which match a filter that it requests.
func NewLvmPredicate(optional config.AlgorithmOptionalOptional, pvcInfo predicates.PersistentVolumeClaimInfo) algorithm.FitPredicate {
	lvmChecker := &LvmChecker{
		optional: optional,
		pvcInfo:  pvcInfo,
	}

	return lvmChecker.CheckNodeLabelPresence
}

//
func (lvm *LvmChecker) CheckNodeLabelPresence(pod *v1.Pod, meta algorithm.PredicateMetadata, nodeInfo *schedulercache.NodeInfo) (bool, []algorithm.PredicateFailureReason, error) {
	nodeName := nodeInfo.Node().Name
	config, err := buildConfig(lvm.optional.Kubecfg)
	if err != nil {
		glog.V(2).Info("can't get kubeconfig")
		return false, []algorithm.PredicateFailureReason{predicates.ErrNodeUnschedulable}, nil
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.V(2).Info("can't connect to k8s")
		return false, []algorithm.PredicateFailureReason{predicates.ErrNodeUnschedulable}, nil
	}

	cm, err := clientset.CoreV1().ConfigMaps(lvm.optional.LvmConfigMapNamespace).Get(lvm.optional.LvmConfigMapName, metav1.GetOptions{})
	if err != nil {
		glog.V(2).Infof("can't get configmap %s:%s", lvm.optional.LvmConfigMapNamespace, lvm.optional.LvmConfigMapName)
		return false, []algorithm.PredicateFailureReason{predicates.ErrNodeUnschedulable}, nil
	}
	data := cm.Data[nodeName]
	LvmResource := &lvmutils.LvmResource{}
	err = json.Unmarshal([]byte(data), LvmResource)
	if err != nil {
		glog.V(2).Info("can't ummarshal json")
		return false, []algorithm.PredicateFailureReason{predicates.ErrNodeUnschedulable}, nil
	}
	vg := LvmResource.Vg
	for _, vol := range pod.Spec.Volumes {
		if vol.PersistentVolumeClaim != nil {
			pvc, err := lvm.pvcInfo.GetPersistentVolumeClaimInfo(pod.Namespace, vol.PersistentVolumeClaim.ClaimName)
			if err != nil {
				glog.V(2).Infof("can't get pvc, %s, %s", pod.Namespace, vol.PersistentVolumeClaim.ClaimName)
				return false, []algorithm.PredicateFailureReason{predicates.ErrNodeUnschedulable}, nil
			}
			capacity := pvc.Spec.Resources.Requests[v1.ResourceStorage]
			glog.V(2).Infof("capacity: %d, node %s vg free %d", capacity.Value(), nodeName, int64(vg.Free))
			if capacity.Value() > int64(vg.Free) {
				return false, []algorithm.PredicateFailureReason{predicates.ErrNodeUnschedulable}, nil
			}
		}
	}

	return true, nil, nil
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}
