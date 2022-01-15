package test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path"
	"sync"
	"text/template"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kralicky/opni-monitoring/pkg/agent"
	"github.com/kralicky/opni-monitoring/pkg/auth"
	"github.com/kralicky/opni-monitoring/pkg/bootstrap"
	"github.com/kralicky/opni-monitoring/pkg/config/v1beta1"
	"github.com/kralicky/opni-monitoring/pkg/gateway"
	"github.com/kralicky/opni-monitoring/pkg/management"
	"github.com/kralicky/opni-monitoring/pkg/pkp"
	mock_ident "github.com/kralicky/opni-monitoring/pkg/test/mock/ident"
	"github.com/kralicky/opni-monitoring/pkg/tokens"
	"github.com/onsi/ginkgo/v2"
	"github.com/phayes/freeport"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type servicePorts struct {
	Etcd       int
	Gateway    int
	Management int
	Cortex     int
}

type Environment struct {
	TestBin string
	Logger  *zap.SugaredLogger

	waitGroup *sync.WaitGroup
	mockCtrl  *gomock.Controller
	ctx       context.Context
	cancel    context.CancelFunc
	tempDir   string
	ports     servicePorts

	gatewayConfig *v1beta1.GatewayConfig
}

func (e *Environment) Start() error {
	e.ctx, e.cancel = context.WithCancel(context.Background())
	e.mockCtrl = gomock.NewController(ginkgo.GinkgoT())
	e.waitGroup = &sync.WaitGroup{}

	if _, err := auth.GetMiddleware("test"); err != nil {
		if err := auth.InstallMiddleware("test", &TestAuthMiddleware{
			Strategy: AuthStrategyUserIDInAuthHeader,
		}); err != nil {
			return fmt.Errorf("failed to install test auth middleware: %w", err)
		}
	}
	ports, err := freeport.GetFreePorts(4)
	if err != nil {
		panic(err)
	}
	e.ports = servicePorts{
		Etcd:       ports[0],
		Gateway:    ports[1],
		Management: ports[2],
		Cortex:     ports[3],
	}
	e.tempDir, err = os.MkdirTemp("", "opni-monitoring-test-*")
	if err != nil {
		return err
	}
	if err := os.Mkdir(path.Join(e.tempDir, "etcd"), 0700); err != nil {
		return err
	}
	cortexTempDir := path.Join(e.tempDir, "cortex")
	if err := os.MkdirAll(path.Join(cortexTempDir, "rules"), 0700); err != nil {
		return err
	}
	entries, _ := fs.ReadDir(TestDataFS, "testdata/cortex")
	fmt.Printf("Copying %d files from embedded testdata/cortex to %s\n", len(entries), cortexTempDir)
	for _, entry := range entries {
		if err := os.WriteFile(path.Join(cortexTempDir, entry.Name()), TestData("cortex/"+entry.Name()), 0644); err != nil {
			return err
		}
	}

	e.startEtcd()
	e.startGateway()
	go e.startCortex()
	return nil
}

func (e *Environment) Stop() error {
	e.cancel()
	e.mockCtrl.Finish()
	e.waitGroup.Wait()
	os.RemoveAll(e.tempDir)
	return nil
}

func (e *Environment) startEtcd() {
	e.waitGroup.Add(1)
	defaultArgs := []string{
		fmt.Sprintf("--listen-client-urls=http://localhost:%d", e.ports.Etcd),
		fmt.Sprintf("--advertise-client-urls=http://localhost:%d", e.ports.Etcd),
		"--listen-peer-urls=http://localhost:0",
		"--log-level=error",
		fmt.Sprintf("--data-dir=%s", path.Join(e.tempDir, "etcd")),
	}
	etcdBin := path.Join(e.TestBin, "etcd")
	cmd := exec.CommandContext(e.ctx, etcdBin, defaultArgs...)
	cmd.Env = []string{"ALLOW_NONE_AUTHENTICATION=yes"}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		if !errors.Is(e.ctx.Err(), context.Canceled) {
			panic(err)
		} else {
			return
		}
	}
	fmt.Println("Waiting for etcd to start...")
	for e.ctx.Err() == nil {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health", e.ports.Etcd))
		if err == nil {
			if resp.StatusCode == http.StatusOK {
				break
			}
		}
		time.Sleep(time.Second)
	}
	fmt.Println("Etcd started")
	go func() {
		defer e.waitGroup.Done()
		<-e.ctx.Done()
		cmd.Wait()
	}()
}

type cortexTemplateOptions struct {
	HttpListenPort int
	StorageDir     string
}

func (e *Environment) startCortex() {
	e.waitGroup.Add(1)
	configTemplate := TestData("cortex/config.yaml")
	t := template.Must(template.New("config").Parse(string(configTemplate)))
	configFile, err := os.Create(path.Join(e.tempDir, "cortex", "config.yaml"))
	if err != nil {
		panic(err)
	}
	if err := t.Execute(configFile, cortexTemplateOptions{
		HttpListenPort: e.ports.Cortex,
		StorageDir:     path.Join(e.tempDir, "cortex"),
	}); err != nil {
		panic(err)
	}
	configFile.Close()
	cortexBin := path.Join(e.TestBin, "cortex")
	defaultArgs := []string{
		fmt.Sprintf("-config.file=%s", path.Join(e.tempDir, "cortex/config.yaml")),
	}
	cmd := exec.CommandContext(e.ctx, cortexBin, defaultArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		if !errors.Is(e.ctx.Err(), context.Canceled) {
			panic(err)
		}
	}
	fmt.Println("Waiting for cortex to start...")
	for e.ctx.Err() == nil {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://localhost:%d/ready", e.ports.Gateway), nil)
		client := http.Client{
			Transport: &http.Transport{
				TLSClientConfig: e.GatewayTLSConfig(),
			},
		}
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(time.Second)
	}
	fmt.Println("Cortex started")
	go func() {
		defer e.waitGroup.Done()
		<-e.ctx.Done()
		cmd.Wait()
	}()
}

func (e *Environment) newGatewayConfig() *v1beta1.GatewayConfig {
	caCertData := string(TestData("root_ca.crt"))
	servingCertData := string(TestData("localhost.crt"))
	servingKeyData := string(TestData("localhost.key"))
	return &v1beta1.GatewayConfig{
		Spec: v1beta1.GatewayConfigSpec{
			ListenAddress:           fmt.Sprintf("localhost:%d", e.ports.Gateway),
			ManagementListenAddress: fmt.Sprintf("tcp:///localhost:%d", e.ports.Management),
			AuthProvider:            "test",
			Certs: v1beta1.CertsSpec{
				CACertData:      &caCertData,
				ServingCertData: &servingCertData,
				ServingKeyData:  &servingKeyData,
			},
			Cortex: v1beta1.CortexSpec{
				Distributor: v1beta1.DistributorSpec{
					Address: fmt.Sprintf("localhost:%d", e.ports.Cortex),
				},
				Ingester: v1beta1.IngesterSpec{
					Address: fmt.Sprintf("localhost:%d", e.ports.Cortex),
				},
				Alertmanager: v1beta1.AlertmanagerSpec{
					Address: fmt.Sprintf("localhost:%d", e.ports.Cortex),
				},
				Ruler: v1beta1.RulerSpec{
					Address: fmt.Sprintf("localhost:%d", e.ports.Cortex),
				},
				QueryFrontend: v1beta1.QueryFrontendSpec{
					Address: fmt.Sprintf("localhost:%d", e.ports.Cortex),
				},
				Certs: v1beta1.CortexCertsSpec{
					ServerCA:   path.Join(e.tempDir, "cortex/root.crt"),
					ClientCA:   path.Join(e.tempDir, "cortex/root.crt"),
					ClientCert: path.Join(e.tempDir, "cortex/client.crt"),
					ClientKey:  path.Join(e.tempDir, "cortex/client.key"),
				},
			},
			Storage: v1beta1.StorageSpec{
				Type: v1beta1.StorageTypeEtcd,
				Etcd: &v1beta1.EtcdStorageSpec{
					Config: clientv3.Config{
						Endpoints: []string{fmt.Sprintf("http://localhost:%d", e.ports.Etcd)},
					},
				},
			},
		},
	}
}

func (e *Environment) NewManagementClient() management.ManagementClient {
	c, err := management.NewClient(management.WithListenAddress(
		fmt.Sprintf("localhost:%d", e.ports.Management)))
	if err != nil {
		panic(err)
	}
	return c
}

func (e *Environment) startGateway() {
	e.waitGroup.Add(1)
	e.gatewayConfig = e.newGatewayConfig()
	g := gateway.NewGateway(e.gatewayConfig,
		gateway.WithAuthMiddleware(e.gatewayConfig.Spec.AuthProvider),
	)
	go func() {
		if err := g.Listen(); err != nil {
			fmt.Println("gateway error:", err)
		}
	}()
	fmt.Println("Waiting for gateway to start...")
	for i := 0; i < 10; i++ {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s/healthz",
			e.gatewayConfig.Spec.ListenAddress), nil)
		client := http.Client{
			Transport: &http.Transport{
				TLSClientConfig: e.GatewayTLSConfig(),
			},
		}
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
	}
	fmt.Println("Gateway started")
	go func() {
		defer e.waitGroup.Done()
		<-e.ctx.Done()
		if err := g.Shutdown(); err != nil {
			fmt.Println("gateway error:", err)
		}
	}()
}

func (e *Environment) StartAgent(id string, token string, pins []string) {
	e.waitGroup.Add(1)
	defer e.waitGroup.Done()
	port, err := freeport.GetFreePort()
	if err != nil {
		panic(err)
	}

	ident := mock_ident.NewMockProvider(e.mockCtrl)
	ident.EXPECT().
		UniqueIdentifier(gomock.Any()).
		Return(id, nil).
		AnyTimes()

	agentConfig := &v1beta1.AgentConfig{
		Spec: v1beta1.AgentConfigSpec{
			ListenAddress:  fmt.Sprintf("localhost:%d", port),
			GatewayAddress: fmt.Sprintf("localhost:%d", e.ports.Gateway),
			Storage: v1beta1.StorageSpec{
				Type: v1beta1.StorageTypeEtcd,
				Etcd: &v1beta1.EtcdStorageSpec{
					Config: clientv3.Config{
						Endpoints: []string{fmt.Sprintf("http://localhost:%d", e.ports.Etcd)},
					},
				},
			},
		},
	}

	publicKeyPins := []*pkp.PublicKeyPin{}
	for _, pin := range pins {
		d, err := pkp.DecodePin(pin)
		if err != nil {
			panic(err)
		}
		publicKeyPins = append(publicKeyPins, d)
	}

	t, err := tokens.ParseHex(token)
	if err != nil {
		panic(err)
	}
	agent := agent.New(agentConfig,
		agent.WithBootstrapper(&bootstrap.ClientConfig{
			Token:    t,
			Pins:     publicKeyPins,
			Endpoint: fmt.Sprintf("http://localhost:%d", e.ports.Gateway),
		}))
	go func() {
		if err := agent.ListenAndServe(); err != nil {
			fmt.Println("agent error:", err)
		}
	}()
	<-e.ctx.Done()
	if err := agent.Shutdown(); err != nil {
		fmt.Println("agent error:", err)
	}
}

func (e *Environment) GatewayTLSConfig() *tls.Config {
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(*e.gatewayConfig.Spec.Certs.CACertData))
	return &tls.Config{
		RootCAs: pool,
	}
}

func (e *Environment) GatewayConfig() *v1beta1.GatewayConfig {
	return e.gatewayConfig
}
