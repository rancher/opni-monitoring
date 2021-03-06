package test

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/prometheus/prometheus/model/rulefmt"
	"github.com/rancher/opni-monitoring/pkg/core"
	"github.com/rancher/opni-monitoring/pkg/ident"
	"github.com/rancher/opni-monitoring/pkg/keyring"
	"github.com/rancher/opni-monitoring/pkg/plugins/apis/capability"
	"github.com/rancher/opni-monitoring/pkg/rules"
	"github.com/rancher/opni-monitoring/pkg/storage"
	mock_capability "github.com/rancher/opni-monitoring/pkg/test/mock/capability"
	mock_ident "github.com/rancher/opni-monitoring/pkg/test/mock/ident"
	mock_rules "github.com/rancher/opni-monitoring/pkg/test/mock/rules"
	mock_storage "github.com/rancher/opni-monitoring/pkg/test/mock/storage"
	"github.com/rancher/opni-monitoring/pkg/tokens"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

/******************************************************************************
 * Capabilities                                                               *
 ******************************************************************************/

type CapabilityInfo struct {
	Name              string
	CanInstall        bool
	InstallerTemplate string
}

func (ci *CapabilityInfo) canInstall() error {
	if !ci.CanInstall {
		return errors.New("test error")
	}
	return nil
}

func NewTestCapabilityBackend(
	ctrl *gomock.Controller,
	capBackend *CapabilityInfo,
) capability.Backend {
	backend := mock_capability.NewMockBackend(ctrl)
	backend.EXPECT().
		CanInstall().
		DoAndReturn(capBackend.canInstall).
		AnyTimes()
	backend.EXPECT().
		Install(gomock.Any()).
		Return(nil).
		AnyTimes()
	backend.EXPECT().
		InstallerTemplate().
		Return(capBackend.InstallerTemplate).
		AnyTimes()
	return backend
}

func NewTestCapabilityBackendClient(
	ctrl *gomock.Controller,
	capBackend *CapabilityInfo,
) capability.BackendClient {
	client := mock_capability.NewMockBackendClient(ctrl)
	client.EXPECT().
		Info(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&capability.InfoResponse{
			CapabilityName: capBackend.Name,
		}, nil).
		AnyTimes()
	client.EXPECT().
		CanInstall(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(context.Context, *emptypb.Empty, ...grpc.CallOption) (*emptypb.Empty, error) {
			return nil, capBackend.canInstall()
		}).
		AnyTimes()
	client.EXPECT().
		Install(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(context.Context, *capability.InstallRequest, ...grpc.CallOption) (*emptypb.Empty, error) {
			return nil, nil
		}).
		AnyTimes()
	client.EXPECT().
		Uninstall(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil).
		AnyTimes()
	client.EXPECT().
		InstallerTemplate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&capability.InstallerTemplateResponse{
			Template: capBackend.InstallerTemplate,
		}, nil).
		AnyTimes()
	return client
}

/******************************************************************************
 * Storage                                                                    *
 ******************************************************************************/

func NewTestClusterStore(ctrl *gomock.Controller) storage.ClusterStore {
	mockClusterStore := mock_storage.NewMockClusterStore(ctrl)

	clusters := map[string]*core.Cluster{}
	mu := sync.Mutex{}

	mockClusterStore.EXPECT().
		CreateCluster(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, cluster *core.Cluster) error {
			mu.Lock()
			defer mu.Unlock()
			clusters[cluster.Id] = cluster
			return nil
		}).
		AnyTimes()
	mockClusterStore.EXPECT().
		DeleteCluster(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, ref *core.Reference) error {
			mu.Lock()
			defer mu.Unlock()
			if _, ok := clusters[ref.Id]; !ok {
				return storage.ErrNotFound
			}
			delete(clusters, ref.Id)
			return nil
		}).
		AnyTimes()
	mockClusterStore.EXPECT().
		ListClusters(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, matchLabels *core.LabelSelector, matchOptions core.MatchOptions) (*core.ClusterList, error) {
			mu.Lock()
			defer mu.Unlock()
			clusterList := &core.ClusterList{}
			selectorPredicate := storage.ClusterSelector{
				LabelSelector: matchLabels,
				MatchOptions:  matchOptions,
			}.Predicate()
			for _, cluster := range clusters {
				if selectorPredicate(cluster) {
					clusterList.Items = append(clusterList.Items, cluster)
				}
			}
			return clusterList, nil
		}).
		AnyTimes()
	mockClusterStore.EXPECT().
		GetCluster(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, ref *core.Reference) (*core.Cluster, error) {
			mu.Lock()
			defer mu.Unlock()
			if _, ok := clusters[ref.Id]; !ok {
				return nil, storage.ErrNotFound
			}
			return clusters[ref.Id], nil
		}).
		AnyTimes()
	mockClusterStore.EXPECT().
		UpdateCluster(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, ref *core.Reference, mutator storage.MutatorFunc[*core.Cluster]) (*core.Cluster, error) {
			mu.Lock()
			defer mu.Unlock()
			if _, ok := clusters[ref.Id]; !ok {
				return nil, storage.ErrNotFound
			}
			cluster := clusters[ref.Id]
			cloned := proto.Clone(cluster).(*core.Cluster)
			mutator(cloned)
			if _, ok := clusters[ref.Id]; !ok {
				return nil, storage.ErrNotFound
			}
			clusters[ref.Id] = cloned
			return cloned, nil
		}).
		AnyTimes()
	return mockClusterStore
}

type KeyringStoreHandler = func(_ context.Context, prefix string, ref *core.Reference) (storage.KeyringStore, error)

func NewTestKeyringStoreBroker(ctrl *gomock.Controller, handler ...KeyringStoreHandler) storage.KeyringStoreBroker {
	mockKeyringStoreBroker := mock_storage.NewMockKeyringStoreBroker(ctrl)
	keyringStores := map[string]storage.KeyringStore{}
	defaultHandler := func(_ context.Context, prefix string, ref *core.Reference) (storage.KeyringStore, error) {
		if keyringStore, ok := keyringStores[prefix+ref.Id]; !ok {
			s := NewTestKeyringStore(ctrl, prefix, ref)
			keyringStores[prefix+ref.Id] = s
			return s, nil
		} else {
			return keyringStore, nil
		}
	}

	var h KeyringStoreHandler
	if len(handler) > 0 {
		h = handler[0]
	} else {
		h = defaultHandler
	}

	mockKeyringStoreBroker.EXPECT().
		KeyringStore(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, prefix string, ref *core.Reference) (storage.KeyringStore, error) {
			if prefix == "gateway-internal" {
				return defaultHandler(ctx, prefix, ref)
			}
			return h(ctx, prefix, ref)
		}).
		AnyTimes()
	return mockKeyringStoreBroker
}

func NewTestKeyringStore(ctrl *gomock.Controller, prefix string, ref *core.Reference) storage.KeyringStore {
	mockKeyringStore := mock_storage.NewMockKeyringStore(ctrl)
	keyrings := map[string]keyring.Keyring{}
	mockKeyringStore.EXPECT().
		Put(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, keyring keyring.Keyring) error {
			keyrings[prefix+ref.Id] = keyring
			return nil
		}).
		AnyTimes()
	mockKeyringStore.EXPECT().
		Get(gomock.Any()).
		DoAndReturn(func(_ context.Context) (keyring.Keyring, error) {
			keyring, ok := keyrings[prefix+ref.Id]
			if !ok {
				return nil, storage.ErrNotFound
			}
			return keyring, nil
		}).
		AnyTimes()
	return mockKeyringStore
}

func NewTestKeyValueStoreBroker(ctrl *gomock.Controller) storage.KeyValueStoreBroker {
	mockKvStoreBroker := mock_storage.NewMockKeyValueStoreBroker(ctrl)
	kvStores := map[string]storage.KeyValueStore{}
	mockKvStoreBroker.EXPECT().
		KeyValueStore(gomock.Any()).
		DoAndReturn(func(namespace string) (storage.KeyValueStore, error) {
			if kvStore, ok := kvStores[namespace]; !ok {
				s := NewTestKeyValueStore(ctrl)
				kvStores[namespace] = s
				return s, nil
			} else {
				return kvStore, nil
			}
		}).
		AnyTimes()
	return mockKvStoreBroker
}

func NewTestKeyValueStore(ctrl *gomock.Controller) storage.KeyValueStore {
	mockKvStore := mock_storage.NewMockKeyValueStore(ctrl)
	kvs := map[string][]byte{}
	mockKvStore.EXPECT().
		Put(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, key string, value []byte) error {
			kvs[key] = value
			return nil
		}).
		AnyTimes()
	mockKvStore.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, key string) ([]byte, error) {
			v, ok := kvs[key]
			if !ok {
				return nil, storage.ErrNotFound
			}
			return v, nil
		}).
		AnyTimes()
	return mockKvStore
}

func NewTestRBACStore(ctrl *gomock.Controller) storage.RBACStore {
	mockRBACStore := mock_storage.NewMockRBACStore(ctrl)

	roles := map[string]*core.Role{}
	rbs := map[string]*core.RoleBinding{}
	mu := sync.Mutex{}

	mockRBACStore.EXPECT().
		CreateRole(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, role *core.Role) error {
			mu.Lock()
			defer mu.Unlock()
			roles[role.Id] = role
			return nil
		}).
		AnyTimes()
	mockRBACStore.EXPECT().
		DeleteRole(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, ref *core.Reference) error {
			mu.Lock()
			defer mu.Unlock()
			if _, ok := roles[ref.Id]; !ok {
				return storage.ErrNotFound
			}
			delete(roles, ref.Id)
			return nil
		}).
		AnyTimes()
	mockRBACStore.EXPECT().
		GetRole(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, ref *core.Reference) (*core.Role, error) {
			mu.Lock()
			defer mu.Unlock()
			if _, ok := roles[ref.Id]; !ok {
				return nil, storage.ErrNotFound
			}
			return roles[ref.Id], nil
		}).
		AnyTimes()
	mockRBACStore.EXPECT().
		ListRoles(gomock.Any()).
		DoAndReturn(func(_ context.Context) (*core.RoleList, error) {
			mu.Lock()
			defer mu.Unlock()
			roleList := &core.RoleList{}
			for _, role := range roles {
				roleList.Items = append(roleList.Items, role)
			}
			return roleList, nil
		}).
		AnyTimes()
	mockRBACStore.EXPECT().
		CreateRoleBinding(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, rb *core.RoleBinding) error {
			mu.Lock()
			defer mu.Unlock()
			rbs[rb.Id] = rb
			return nil
		}).
		AnyTimes()
	mockRBACStore.EXPECT().
		DeleteRoleBinding(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, ref *core.Reference) error {
			mu.Lock()
			defer mu.Unlock()
			if _, ok := rbs[ref.Id]; !ok {
				return storage.ErrNotFound
			}
			delete(rbs, ref.Id)
			return nil
		}).
		AnyTimes()
	mockRBACStore.EXPECT().
		GetRoleBinding(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, ref *core.Reference) (*core.RoleBinding, error) {
			mu.Lock()
			if _, ok := rbs[ref.Id]; !ok {
				mu.Unlock()
				return nil, storage.ErrNotFound
			}
			cloned := proto.Clone(rbs[ref.Id]).(*core.RoleBinding)
			mu.Unlock()
			storage.ApplyRoleBindingTaints(ctx, mockRBACStore, cloned)
			return cloned, nil
		}).
		AnyTimes()
	mockRBACStore.EXPECT().
		ListRoleBindings(gomock.Any()).
		DoAndReturn(func(ctx context.Context) (*core.RoleBindingList, error) {
			mu.Lock()
			rbList := &core.RoleBindingList{}
			for _, rb := range rbs {
				cloned := proto.Clone(rb).(*core.RoleBinding)
				rbList.Items = append(rbList.Items, cloned)
			}
			mu.Unlock()
			for _, rb := range rbList.Items {
				storage.ApplyRoleBindingTaints(ctx, mockRBACStore, rb)
			}
			return rbList, nil
		}).
		AnyTimes()
	return mockRBACStore
}

func NewTestTokenStore(ctx context.Context, ctrl *gomock.Controller) storage.TokenStore {
	mockTokenStore := mock_storage.NewMockTokenStore(ctrl)

	leaseStore := NewLeaseStore(ctx)
	tks := map[string]*core.BootstrapToken{}
	mu := sync.Mutex{}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case tokenID := <-leaseStore.LeaseExpired():
				mockTokenStore.DeleteToken(ctx, &core.Reference{
					Id: tokenID,
				})
			}
		}
	}()

	mockTokenStore.EXPECT().
		CreateToken(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, ttl time.Duration, opts ...storage.TokenCreateOption) (*core.BootstrapToken, error) {
			mu.Lock()
			defer mu.Unlock()
			options := storage.NewTokenCreateOptions()
			options.Apply(opts...)
			t := tokens.NewToken().ToBootstrapToken()
			lease := leaseStore.New(t.TokenID, ttl)
			t.Metadata = &core.BootstrapTokenMetadata{
				LeaseID:      int64(lease.ID),
				Ttl:          int64(ttl),
				UsageCount:   0,
				Labels:       options.Labels,
				Capabilities: options.Capabilities,
			}
			tks[t.TokenID] = t
			return t, nil
		}).
		AnyTimes()
	mockTokenStore.EXPECT().
		DeleteToken(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, ref *core.Reference) error {
			mu.Lock()
			defer mu.Unlock()
			if _, ok := tks[ref.Id]; !ok {
				return storage.ErrNotFound
			}
			delete(tks, ref.Id)
			return nil
		}).
		AnyTimes()
	mockTokenStore.EXPECT().
		GetToken(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, ref *core.Reference) (*core.BootstrapToken, error) {
			mu.Lock()
			defer mu.Unlock()
			if _, ok := tks[ref.Id]; !ok {
				return nil, storage.ErrNotFound
			}
			return tks[ref.Id], nil
		}).
		AnyTimes()
	mockTokenStore.EXPECT().
		ListTokens(gomock.Any()).
		DoAndReturn(func(_ context.Context) ([]*core.BootstrapToken, error) {
			mu.Lock()
			defer mu.Unlock()
			tokens := make([]*core.BootstrapToken, 0, len(tks))
			for _, t := range tks {
				tokens = append(tokens, t)
			}
			return tokens, nil
		}).
		AnyTimes()
	mockTokenStore.EXPECT().
		UpdateToken(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, ref *core.Reference, mutator storage.MutatorFunc[*core.BootstrapToken]) (*core.BootstrapToken, error) {
			mu.Lock()
			defer mu.Unlock()
			if _, ok := tks[ref.Id]; !ok {
				return nil, storage.ErrNotFound
			}
			token := tks[ref.Id]
			cloned := proto.Clone(token).(*core.BootstrapToken)
			mutator(cloned)
			if _, ok := tks[ref.Id]; !ok {
				return nil, storage.ErrNotFound
			}
			tks[ref.Id] = cloned
			return cloned, nil
		}).
		AnyTimes()

	return mockTokenStore
}

/******************************************************************************
 * Ident                                                                      *
 ******************************************************************************/

func NewTestIdentProvider(ctrl *gomock.Controller, id string) ident.Provider {
	mockIdent := mock_ident.NewMockProvider(ctrl)
	mockIdent.EXPECT().
		UniqueIdentifier(gomock.Any()).
		Return(id, nil).
		AnyTimes()
	return mockIdent
}

/******************************************************************************
 * Rules                                                                      *
 ******************************************************************************/

func NewTestRuleFinder(ctrl *gomock.Controller, groups func() []rulefmt.RuleGroup) rules.RuleFinder {
	mockRuleFinder := mock_rules.NewMockRuleFinder(ctrl)
	mockRuleFinder.EXPECT().
		FindGroups(gomock.Any()).
		DoAndReturn(func(ctx context.Context) ([]rulefmt.RuleGroup, error) {
			return groups(), nil
		}).
		AnyTimes()
	return mockRuleFinder
}
