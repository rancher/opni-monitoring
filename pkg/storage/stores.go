package storage

import (
	"context"
	"time"

	"github.com/rancher/opni-monitoring/pkg/core"
	"github.com/rancher/opni-monitoring/pkg/keyring"
)

type Backend interface {
	TokenStore
	ClusterStore
	RBACStore
	KeyringStoreBroker
	KeyValueStoreBroker
}

type MutatorFunc[T any] func(T)

type TokenMutator = MutatorFunc[*core.BootstrapToken]
type ClusterMutator = MutatorFunc[*core.Cluster]

type TokenStore interface {
	CreateToken(ctx context.Context, ttl time.Duration, opts ...TokenCreateOption) (*core.BootstrapToken, error)
	DeleteToken(ctx context.Context, ref *core.Reference) error
	GetToken(ctx context.Context, ref *core.Reference) (*core.BootstrapToken, error)
	UpdateToken(ctx context.Context, ref *core.Reference, mutator TokenMutator) (*core.BootstrapToken, error)
	ListTokens(ctx context.Context) ([]*core.BootstrapToken, error)
}

type ClusterStore interface {
	CreateCluster(ctx context.Context, cluster *core.Cluster) error
	DeleteCluster(ctx context.Context, ref *core.Reference) error
	GetCluster(ctx context.Context, ref *core.Reference) (*core.Cluster, error)
	UpdateCluster(ctx context.Context, ref *core.Reference, mutator ClusterMutator) (*core.Cluster, error)
	ListClusters(ctx context.Context, matchLabels *core.LabelSelector, matchOptions core.MatchOptions) (*core.ClusterList, error)
}

type RBACStore interface {
	CreateRole(context.Context, *core.Role) error
	DeleteRole(context.Context, *core.Reference) error
	GetRole(context.Context, *core.Reference) (*core.Role, error)
	CreateRoleBinding(context.Context, *core.RoleBinding) error
	DeleteRoleBinding(context.Context, *core.Reference) error
	GetRoleBinding(context.Context, *core.Reference) (*core.RoleBinding, error)
	ListRoles(context.Context) (*core.RoleList, error)
	ListRoleBindings(context.Context) (*core.RoleBindingList, error)
}

type KeyringStore interface {
	Put(ctx context.Context, keyring keyring.Keyring) error
	Get(ctx context.Context) (keyring.Keyring, error)
}

type KeyValueStore interface {
	Put(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	ListKeys(ctx context.Context, prefix string) ([]string, error)
}

type KeyringStoreBroker interface {
	KeyringStore(ctx context.Context, namespace string, ref *core.Reference) (KeyringStore, error)
}

type KeyValueStoreBroker interface {
	KeyValueStore(namespace string) (KeyValueStore, error)
}

// A store that can be used to compute subject access rules
type SubjectAccessCapableStore interface {
	ListClusters(ctx context.Context, matchLabels *core.LabelSelector, matchOptions core.MatchOptions) (*core.ClusterList, error)
	GetRole(ctx context.Context, ref *core.Reference) (*core.Role, error)
	ListRoleBindings(ctx context.Context) (*core.RoleBindingList, error)
}
