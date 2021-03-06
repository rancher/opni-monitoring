package v1beta1

import (
	"github.com/rancher/opni-monitoring/pkg/config/meta"
)

type GatewayConfig struct {
	meta.TypeMeta `json:",inline"`

	Spec GatewayConfigSpec `json:"spec,omitempty"`
}

type GatewayConfigSpec struct {
	ListenAddress  string         `json:"listenAddress,omitempty"`
	Hostname       string         `json:"hostname,omitempty"`
	MetricsPort    int            `json:"metricsPort,omitempty"`
	Management     ManagementSpec `json:"management,omitempty"`
	EnableMonitor  bool           `json:"enableMonitor,omitempty"`
	TrustedProxies []string       `json:"trustedProxies,omitempty"`
	Cortex         CortexSpec     `json:"cortex,omitempty"`
	AuthProvider   string         `json:"authProvider,omitempty"`
	Storage        StorageSpec    `json:"storage,omitempty"`
	Certs          CertsSpec      `json:"certs,omitempty"`
	Plugins        PluginsSpec    `json:"plugins,omitempty"`
}

type ManagementSpec struct {
	GRPCListenAddress string `json:"grpcListenAddress,omitempty"`
	HTTPListenAddress string `json:"httpListenAddress,omitempty"`
	WebListenAddress  string `json:"webListenAddress,omitempty"`
}

type CortexSpec struct {
	Distributor   DistributorSpec   `json:"distributor,omitempty"`
	Ingester      IngesterSpec      `json:"ingester,omitempty"`
	Alertmanager  AlertmanagerSpec  `json:"alertmanager,omitempty"`
	Ruler         RulerSpec         `json:"ruler,omitempty"`
	QueryFrontend QueryFrontendSpec `json:"queryFrontend,omitempty"`
	Certs         MTLSSpec          `json:"certs,omitempty"`
}

type DistributorSpec struct {
	HTTPAddress string `json:"httpAddress,omitempty"`
	GRPCAddress string `json:"grpcAddress,omitempty"`
}

type IngesterSpec struct {
	HTTPAddress string `json:"httpAddress,omitempty"`
	GRPCAddress string `json:"grpcAddress,omitempty"`
}

type AlertmanagerSpec struct {
	HTTPAddress string `json:"httpAddress,omitempty"`
}

type RulerSpec struct {
	// HTTP address of the cortex ruler
	HTTPAddress string `json:"httpAddress,omitempty"`
}

type QueryFrontendSpec struct {
	// HTTP address of the cortex query frontend
	HTTPAddress string `json:"httpAddress,omitempty"`
	// GRPC address of the cortex query frontend
	GRPCAddress string `json:"grpcAddress,omitempty"`
}

type MTLSSpec struct {
	// Path to the server CA certificate.
	ServerCA string `json:"serverCA,omitempty"`
	// Path to the client CA certificate (not needed in all cases).
	ClientCA string `json:"clientCA,omitempty"`
	// Path to the certificate used for client-cert auth.
	ClientCert string `json:"clientCert,omitempty"`
	// Path to the private key used for client-cert auth.
	ClientKey string `json:"clientKey,omitempty"`
}

type CertsSpec struct {
	// Path to a PEM encoded CA certificate file. Mutually exclusive with CACertData
	CACert *string `json:"caCert,omitempty"`
	// String containing PEM encoded CA certificate data. Mutually exclusive with CACert
	CACertData *string `json:"caCertData,omitempty"`
	// Path to a PEM encoded server certificate file. Mutually exclusive with ServingCertData
	ServingCert *string `json:"servingCert,omitempty"`
	// String containing PEM encoded server certificate data. Mutually exclusive with ServingCert
	ServingCertData *string `json:"servingCertData,omitempty"`
	// Path to a PEM encoded server key file. Mutually exclusive with ServingKeyData
	ServingKey *string `json:"servingKey,omitempty"`
	// String containing PEM encoded server key data. Mutually exclusive with ServingKey
	ServingKeyData *string `json:"servingKeyData,omitempty"`
}

type PluginsSpec struct {
	// Directories to look for plugins in
	Dirs []string `json:"dirs,omitempty"`
}

func (s *GatewayConfigSpec) SetDefaults() {
	if s == nil {
		return
	}
	if s.ListenAddress == "" {
		s.ListenAddress = ":8080"
	}
	if s.Hostname == "" {
		s.Hostname = "localhost"
	}
	if s.MetricsPort == 0 {
		s.MetricsPort = 8086
	}
	if s.Cortex.Distributor.HTTPAddress == "" {
		s.Cortex.Distributor.HTTPAddress = "cortex-distributor:8080"
		s.Cortex.Distributor.GRPCAddress = "cortex-distributor-headless:9095"
	}
	if s.Cortex.Ingester.HTTPAddress == "" {
		s.Cortex.Ingester.HTTPAddress = "cortex-ingester:8080"
		s.Cortex.Ingester.GRPCAddress = "cortex-ingester-headless:9095"
	}
	if s.Cortex.Alertmanager.HTTPAddress == "" {
		s.Cortex.Alertmanager.HTTPAddress = "cortex-alertmanager:8080"
	}
	if s.Cortex.Ruler.HTTPAddress == "" {
		s.Cortex.Ruler.HTTPAddress = "cortex-ruler:8080"
	}
	if s.Cortex.QueryFrontend.HTTPAddress == "" {
		s.Cortex.QueryFrontend.HTTPAddress = "cortex-query-frontend:8080"
		s.Cortex.QueryFrontend.GRPCAddress = "cortex-query-frontend-headless:9095"
	}
}

type StorageType string

const (
	// Use etcd for key-value storage. This is the recommended default.
	StorageTypeEtcd StorageType = "etcd"
	// Use Kubernetes custom resources to store objects. This is experimental,
	// and it is recommended to use the etcd storage type instead for performance
	// reasons.
	StorageTypeCRDs StorageType = "customResources"
)

type StorageSpec struct {
	Type            StorageType                 `json:"type,omitempty"`
	Etcd            *EtcdStorageSpec            `json:"etcd,omitempty"`
	CustomResources *CustomResourcesStorageSpec `json:"customResources,omitempty"`
}

type EtcdStorageSpec struct {
	// List of etcd endpoints to connect to.
	Endpoints []string `json:"endpoints,omitempty"`
	// Configuration for etcd client-cert auth.
	Certs *MTLSSpec `json:"certs,omitempty"`
}

type CustomResourcesStorageSpec struct {
	// Kubernetes namespace where custom resource objects will be stored.
	Namespace string `json:"namespace,omitempty"`
}
