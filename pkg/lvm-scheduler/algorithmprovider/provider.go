package algorithmprovider

import (
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/kubernetes/pkg/scheduler/algorithm"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/predicates"
	"k8s.io/kubernetes/pkg/scheduler/factory"

	schedulerserverconfig "bonc.com/k8s-scheduler-extender-example/cmd/lvm-scheduler/app/config"
	lvmalgorithm "bonc.com/k8s-scheduler-extender-example/pkg/lvm-scheduler/algorithm"
)

const (
	// ClusterAutoscalerProvider defines the default autoscaler provider
	ClusterAutoscalerProvider = "ClusterAutoscalerProvider"
)

func ApplyFeatureGates(lo schedulerserverconfig.AlgorithmOptionalOptional) {
	registerAlgorithmProvider(defaultPredicates(lo), defaultPriorities())

	// reset order, add your predicates method in the first place
	order := predicates.Ordering()
	names := append([]string{lvmalgorithm.CheckLvmVolumePred}, order...)
	predicates.SetPredicatesOrdering(names)
}

func registerAlgorithmProvider(predSet, priSet sets.String) {
	// Registers algorithm providers. By default we use 'DefaultProvider', but user can specify one to be used
	// by specifying flag.
	factory.RegisterAlgorithmProvider(factory.DefaultProvider, predSet, priSet)
	// Cluster autoscaler friendly scheduling algorithm.
	factory.RegisterAlgorithmProvider(ClusterAutoscalerProvider, predSet,
		copyAndReplace(priSet, "LeastRequestedPriority", "MostRequestedPriority"))
}

//
func defaultPredicates(lo schedulerserverconfig.AlgorithmOptionalOptional) sets.String {
	return sets.NewString(
		factory.RegisterFitPredicateFactory(
			lvmalgorithm.CheckLvmVolumePred,
			func(args factory.PluginFactoryArgs) algorithm.FitPredicate {
				return lvmalgorithm.NewLvmPredicate(lo, args.PVCInfo)
			},
		),
	)
}

func defaultPriorities() sets.String {
	return nil
}

func copyAndReplace(set sets.String, replaceWhat, replaceWith string) sets.String {
	result := sets.NewString(set.List()...)
	if result.Has(replaceWhat) {
		result.Delete(replaceWhat)
		result.Insert(replaceWith)
	}
	return result
}
