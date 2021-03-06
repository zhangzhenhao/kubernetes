/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kubeadm

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeletconfigv1alpha1 "k8s.io/kubernetes/pkg/kubelet/apis/kubeletconfig/v1alpha1"
	kubeproxyconfigv1alpha1 "k8s.io/kubernetes/pkg/proxy/apis/kubeproxyconfig/v1alpha1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MasterConfiguration contains a list of elements which make up master's
// configuration object.
type MasterConfiguration struct {
	metav1.TypeMeta

	API                  API
	KubeProxy            KubeProxy
	Etcd                 Etcd
	KubeletConfiguration KubeletConfiguration
	Networking           Networking
	KubernetesVersion    string
	CloudProvider        string
	NodeName             string
	AuthorizationModes   []string

	Token    string
	TokenTTL *metav1.Duration

	APIServerExtraArgs         map[string]string
	ControllerManagerExtraArgs map[string]string
	SchedulerExtraArgs         map[string]string

	APIServerExtraVolumes         []HostPathMount
	ControllerManagerExtraVolumes []HostPathMount
	SchedulerExtraVolumes         []HostPathMount

	// APIServerCertSANs sets extra Subject Alternative Names for the API Server signing cert
	APIServerCertSANs []string
	// CertificatesDir specifies where to store or look for all required certificates
	CertificatesDir string

	// ImageRepository what container registry to pull control plane images from
	ImageRepository string

	// Container registry for core images generated by CI
	// +k8s:conversion-gen=false
	CIImageRepository string

	// UnifiedControlPlaneImage specifies if a specific container image should be used for all control plane components
	UnifiedControlPlaneImage string

	// FeatureGates enabled by the user
	FeatureGates map[string]bool
}

// API struct contains elements of API server address.
type API struct {
	// AdvertiseAddress sets the address for the API server to advertise.
	AdvertiseAddress string
	// BindPort sets the secure port for the API Server to bind to
	BindPort int32
}

// TokenDiscovery contains elements needed for token discovery
type TokenDiscovery struct {
	ID        string
	Secret    string
	Addresses []string
}

// Networking contains elements describing cluster's networking configuration
type Networking struct {
	ServiceSubnet string
	PodSubnet     string
	DNSDomain     string
}

// Etcd contains elements describing Etcd configuration
type Etcd struct {
	Endpoints []string
	CAFile    string
	CertFile  string
	KeyFile   string
	DataDir   string
	ExtraArgs map[string]string
	// Image specifies which container image to use for running etcd. If empty, automatically populated by kubeadm using the image repository and default etcd version
	Image      string
	SelfHosted *SelfHostedEtcd
}

// SelfHostedEtcd describes options required to configure self-hosted etcd
type SelfHostedEtcd struct {
	// CertificatesDir represents the directory where all etcd TLS assets are stored. By default this is
	// a dir names "etcd" in the main CertificatesDir value.
	CertificatesDir string
	// ClusterServiceName is the name of the service that load balances the etcd cluster
	ClusterServiceName string
	// EtcdVersion is the version of etcd running in the cluster.
	EtcdVersion string
	// OperatorVersion is the version of the etcd-operator to use.
	OperatorVersion string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeConfiguration contains elements describing a particular node
type NodeConfiguration struct {
	metav1.TypeMeta

	CACertPath     string
	DiscoveryFile  string
	DiscoveryToken string
	// Currently we only pay attention to one api server but hope to support >1 in the future
	DiscoveryTokenAPIServers []string
	NodeName                 string
	TLSBootstrapToken        string
	Token                    string

	// DiscoveryTokenCACertHashes specifies a set of public key pins to verify
	// when token-based discovery is used. The root CA found during discovery
	// must match one of these values. Specifying an empty set disables root CA
	// pinning, which can be unsafe. Each hash is specified as "<type>:<value>",
	// where the only currently supported type is "sha256". This is a hex-encoded
	// SHA-256 hash of the Subject Public Key Info (SPKI) object in DER-encoded
	// ASN.1. These hashes can be calculated using, for example, OpenSSL:
	// openssl x509 -pubkey -in ca.crt openssl rsa -pubin -outform der 2>&/dev/null | openssl dgst -sha256 -hex
	DiscoveryTokenCACertHashes []string

	// DiscoveryTokenUnsafeSkipCAVerification allows token-based discovery
	// without CA verification via DiscoveryTokenCACertHashes. This can weaken
	// the security of kubeadm since other nodes can impersonate the master.
	DiscoveryTokenUnsafeSkipCAVerification bool

	// FeatureGates enabled by the user
	FeatureGates map[string]bool
}

// KubeletConfiguration contains elements describing initial remote configuration of kubelet
type KubeletConfiguration struct {
	BaseConfig *kubeletconfigv1alpha1.KubeletConfiguration
}

// GetControlPlaneImageRepository returns name of image repository
// for control plane images (API,Controller Manager,Scheduler and Proxy)
// It will override location with CI registry name in case user requests special
// Kubernetes version from CI build area.
// (See: kubeadmconstants.DefaultCIImageRepository)
func (cfg *MasterConfiguration) GetControlPlaneImageRepository() string {
	if cfg.CIImageRepository != "" {
		return cfg.CIImageRepository
	}
	return cfg.ImageRepository
}

// HostPathMount contains elements describing volumes that are mounted from the
// host
type HostPathMount struct {
	Name      string
	HostPath  string
	MountPath string
}

// KubeProxy contains elements describing the proxy configuration
type KubeProxy struct {
	Config *kubeproxyconfigv1alpha1.KubeProxyConfiguration
}
