package kubernetes

import (
	"context"
	"fmt"

	"github.com/rancher/opni-monitoring/pkg/core"
	"github.com/rancher/opni-monitoring/pkg/logger"
	"github.com/rancher/opni-monitoring/pkg/storage"
	opniv2beta1 "github.com/rancher/opni/apis/v2beta1"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	opensearchv1 "opensearch.opster.io/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KubernetesStore implements LoggingClusterStore
type KubernetesStore struct {
	KubernetesStoreOptions
	client client.Client
	logger *zap.SugaredLogger
}

type KubernetesStoreOptions struct {
	systemNamespace string
}

type KubernetesStoreOption func(*KubernetesStoreOptions)

func (o *KubernetesStoreOptions) Apply(opts ...KubernetesStoreOption) {
	for _, op := range opts {
		op(o)
	}
}

func WithNamespace(namespace string) KubernetesStoreOption {
	return func(o *KubernetesStoreOptions) {
		o.systemNamespace = namespace
	}
}

func NewKubernetesStore(opts ...KubernetesStoreOption) *KubernetesStore {
	options := KubernetesStoreOptions{}
	options.Apply(opts...)

	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(opensearchv1.AddToScheme(scheme))
	utilruntime.Must(opniv2beta1.AddToScheme(scheme))

	lg := logger.New().Named("kubernetes")

	cli, err := client.New(ctrl.GetConfigOrDie(), client.Options{
		Scheme: scheme,
	})
	if err != nil {
		lg.With(
			zap.Error(err),
		).Fatal("failed to create kubernetes client")
	}

	return &KubernetesStore{
		KubernetesStoreOptions: options,
		client:                 cli,
		logger:                 lg,
	}
}

type opensearchUserStore struct {
	client    client.Client
	namespace string
}

func (k *KubernetesStore) OpensearchUserStore(ctx context.Context) (storage.OpensearchUserStore, error) {
	return &opensearchUserStore{
		client:    k.client,
		namespace: k.KubernetesStoreOptions.systemNamespace,
	}, nil
}

func (o *opensearchUserStore) Put(ctx context.Context, user *core.OpensearchUser) error {
	indexUserSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      user.Id,
			Namespace: o.namespace,
		},
		Data: map[string][]byte{
			"password": []byte(user.Secret),
		},
	}
	err := o.client.Create(ctx, indexUserSecret)

	if k8serrors.IsAlreadyExists(err) {
		err = o.client.Get(ctx, client.ObjectKeyFromObject(indexUserSecret), indexUserSecret)
		if err != nil {
			return err
		}
		indexUserSecret.Data["password"] = []byte(user.Secret)
		err = o.client.Update(ctx, indexUserSecret)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func (o *opensearchUserStore) Get(ctx context.Context, ref *core.Reference) (*core.OpensearchUser, error) {
	indexUserSecret := &corev1.Secret{}
	err := o.client.Get(ctx, types.NamespacedName{
		Name:      ref.Id,
		Namespace: o.namespace,
	}, indexUserSecret)
	if err != nil {
		return nil, err
	}

	password, ok := indexUserSecret.Data["password"]
	if !ok {
		return nil, fmt.Errorf("missing secret data")
	}

	return &core.OpensearchUser{
		Id:     indexUserSecret.Name,
		Secret: string(password),
	}, nil
}
