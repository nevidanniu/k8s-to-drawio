package k8s

// SupportedResourceKinds lists all Kubernetes resource types supported by the parser
var SupportedResourceKinds = []string{
	"Deployment",
	"StatefulSet",
	"DaemonSet",
	"Service",
	"Ingress",
	"ConfigMap",
	"Secret",
	"PersistentVolume",
	"PersistentVolumeClaim",
	"Namespace",
	"ServiceAccount",
	"Role",
	"RoleBinding",
	"ClusterRole",
	"ClusterRoleBinding",
	"Pod",
	"ReplicaSet",
	"Job",
	"CronJob",
	"Route",
	"ServiceMonitor",
}

// ResourceCategories groups resources by their functional category
var ResourceCategories = map[string]string{
	"Deployment":            "workload",
	"StatefulSet":           "workload",
	"DaemonSet":             "workload",
	"Pod":                   "workload",
	"ReplicaSet":            "workload",
	"Job":                   "workload",
	"CronJob":               "workload",
	"Service":               "networking",
	"Ingress":               "networking",
	"Route":                 "networking",
	"ConfigMap":             "config",
	"Secret":                "config",
	"PersistentVolume":      "storage",
	"PersistentVolumeClaim": "storage",
	"Namespace":             "cluster",
	"ServiceAccount":        "rbac",
	"Role":                  "rbac",
	"RoleBinding":           "rbac",
	"ClusterRole":           "rbac",
	"ClusterRoleBinding":    "rbac",
	"ServiceMonitor":        "monitoring",
}

// IsResourceSupported checks if a given resource kind is supported
func IsResourceSupported(kind string) bool {
	for _, supportedKind := range SupportedResourceKinds {
		if kind == supportedKind {
			return true
		}
	}
	return false
}

// GetResourceCategory returns the category for a given resource kind
func GetResourceCategory(kind string) string {
	if category, exists := ResourceCategories[kind]; exists {
		return category
	}
	return "unknown"
}
