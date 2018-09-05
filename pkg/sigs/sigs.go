package sigs

type SIG string

const (
	APIMachinery SIG = "sig-api-machinery"
	Apps SIG = "sig-apps"
	Architecture SIG = "sig-architecture"
	Auth SIG = "sig-auth"
	Autoscaling SIG = "sig-autoscaling"
	AWS SIG = "sig-aws"
	Azure SIG = "sig-azure"
	BigData SIG = "sig-big-data"
	CLI SIG = "sig-cli"
	CloudProvider SIG = "sig-cloud-provider"
	ClusterLifecycle SIG = "sig-cluster-lifecycle"
	ClusterOps SIG = "sig-cluster-ops"
	ContributorExperience SIG = "sig-contributor-experience"
	Docs SIG = "sig-docs"
	GCP SIG = "sig-gcp"
	IBMCloud SIG = "sig-ibmcloud"
	Instrumentation SIG = "sig-instrumentation"
	Multicluster SIG = "sig-multicluster"
	Network SIG = "sig-network"
	Node SIG = "sig-node"
	OpenStack SIG = "sig-openstack"
	PM SIG = "sig-pm"
	Release SIG = "sig-release"
	Scalability SIG = "sig-scalability"
	Scheduling SIG = "sig-scheduling"
	ServiceCatalog SIG = "sig-service-catalog"
	Storage SIG = "sig-storage"
	Testing SIG = "sig-testing"
	UI SIG = "sig-ui"
	VMWare SIG = "sig-vmware"
	Windows SIG = "sig-windows"
)

var All []SIG = []SIG{
	APIMachinery,
	Apps,
	Architecture,
	Auth,
	Autoscaling,
	AWS,
	Azure,
	BigData,
	CLI,
	CloudProvider,
	ClusterLifecycle,
	ClusterOps,
	ContributorExperience,
	Docs,
	GCP,
	IBMCloud,
	Instrumentation,
	Multicluster,
	Network,
	Node,
	OpenStack,
	PM,
	Release,
	Scalability,
	Scheduling,
	ServiceCatalog,
	Storage,
	Testing,
	UI,
	VMWare,
	Windows,
}
