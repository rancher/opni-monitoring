package crds

import (
	"context"

	"github.com/rancher/opni-monitoring/pkg/core"
	"github.com/rancher/opni-monitoring/pkg/keyring"
	"github.com/rancher/opni-monitoring/pkg/sdk/api/v1beta1"
	"github.com/rancher/opni-monitoring/pkg/storage"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type crdKeyringStore struct {
	CRDStoreOptions
	client client.Client
	ref    *core.Reference
	prefix string
}

func (ks *crdKeyringStore) Put(ctx context.Context, keyring keyring.Keyring) error {
	data, err := keyring.Marshal()
	if err != nil {
		return err
	}
	kr := &v1beta1.Keyring{}
	if err := ks.client.Get(ctx, types.NamespacedName{
		Name:      ks.ref.Id,
		Namespace: ks.namespace,
	}, kr); err != nil {
		if !k8serrors.IsNotFound(err) {
			return err
		}
		kr = &v1beta1.Keyring{
			ObjectMeta: metav1.ObjectMeta{
				Name:      ks.ref.Id,
				Namespace: ks.namespace,
			},
			Data: data,
		}
		return ks.client.Create(ctx, kr)
	}
	return retry.OnError(defaultBackoff, k8serrors.IsConflict, func() error {
		kr := &v1beta1.Keyring{}
		if err := ks.client.Get(ctx, types.NamespacedName{
			Name:      ks.ref.Id,
			Namespace: ks.namespace,
		}, kr); err != nil {
			return err
		}
		kr.Data = data
		return ks.client.Update(ctx, kr)
	})
}

func (ks *crdKeyringStore) Get(ctx context.Context) (keyring.Keyring, error) {
	kr := &v1beta1.Keyring{}
	if err := ks.client.Get(ctx, types.NamespacedName{
		Name:      ks.ref.Id,
		Namespace: ks.namespace,
	}, kr); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}
	return keyring.Unmarshal(kr.Data)
}
