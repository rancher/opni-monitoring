package kubernetes

import (
	"context"
	"fmt"

	"github.com/rancher/opni-monitoring/pkg/core"
	"github.com/rancher/opni-monitoring/pkg/storage"
	opniv2beta1 "github.com/rancher/opni/apis/v2beta1"
	"github.com/rancher/opni/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	opensearchClusterName = "opni-logging"
)

func (k *KubernetesStore) CreateCluster(ctx context.Context, cluster *core.LoggingCluster) error {
	labels := cluster.Labels
	labels[resources.OpniClusterID] = cluster.Id
	loggingCluster := &opniv2beta1.LoggingCluster{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "downstream-",
			Namespace:    k.KubernetesStoreOptions.systemNamespace,
			Labels:       labels,
		},
		Spec: opniv2beta1.LoggingClusterSpec{
			OpensearchClusterRef: &opniv2beta1.OpensearchClusterRef{
				Name:      opensearchClusterName,
				Namespace: k.KubernetesStoreOptions.systemNamespace,
			},
			IndexUserSecret: &corev1.LocalObjectReference{
				Name: cluster.OpensearchUserID,
			},
			FriendlyName: cluster.Name,
		},
	}

	err := k.client.Create(ctx, loggingCluster)
	if err != nil {
		return fmt.Errorf("failed to create cluster: %w", err)
	}

	return nil
}

func (k *KubernetesStore) DeleteCluster(ctx context.Context, ref *core.Reference) error {
	loggingCluster := &opniv2beta1.LoggingCluster{}
	err := k.client.DeleteAllOf(ctx, loggingCluster, client.InNamespace(k.KubernetesStoreOptions.systemNamespace), client.MatchingLabels{resources.OpniClusterID: ref.Id})
	if err != nil {
		return fmt.Errorf("failed to delete cluster: %w", err)
	}
	return nil
}

func (k *KubernetesStore) GetCluster(ctx context.Context, ref *core.Reference) (*core.LoggingCluster, error) {
	loggingClusterList := &opniv2beta1.LoggingClusterList{}

	err := k.client.List(ctx, loggingClusterList, client.InNamespace(k.KubernetesStoreOptions.systemNamespace), client.MatchingLabels{resources.OpniClusterID: ref.Id})
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster: %w", err)
	}

	if len(loggingClusterList.Items) != 1 {
		return nil, fmt.Errorf("fetched %d clusters, expected 1", len(loggingClusterList.Items))
	}

	cluster := loggingClusterList.Items[0]
	return &core.LoggingCluster{
		Id:               ref.Id,
		Name:             cluster.Name,
		OpensearchUserID: derefLocalObjectName(cluster.Spec.IndexUserSecret, ""),
		Labels:           stripIDLabel(cluster.Labels),
	}, nil
}

func (k *KubernetesStore) UpdateCluster(ctx context.Context, cluster *core.LoggingCluster) (*core.LoggingCluster, error) {
	loggingClusterList := &opniv2beta1.LoggingClusterList{}

	err := k.client.List(ctx, loggingClusterList, client.InNamespace(k.KubernetesStoreOptions.systemNamespace), client.MatchingLabels{resources.OpniClusterID: cluster.Id})
	if err != nil {
		return nil, fmt.Errorf("failed to update cluster: %w", storage.ErrNotFound)
	}
	if len(loggingClusterList.Items) != 1 {
		return nil, fmt.Errorf("fetched %d clusters, expected 1", len(loggingClusterList.Items))
	}

	labels := cluster.Labels
	labels[resources.OpniClusterID] = cluster.Id

	loggingCluster := &loggingClusterList.Items[0]
	loggingCluster.Spec.FriendlyName = cluster.Name
	loggingCluster.Spec.IndexUserSecret = &corev1.LocalObjectReference{
		Name: cluster.OpensearchUserID,
	}
	loggingCluster.Labels = labels

	err = k.client.Update(ctx, loggingCluster)
	if err != nil {
		return nil, fmt.Errorf("failed to update cluster: %w", err)
	}

	return cluster, nil
}

func (k *KubernetesStore) ListClusters(ctx context.Context, matchLabels *core.LabelSelector, matchOptions core.MatchOptions) (*core.LoggingClusterList, error) {
	k8sSelector := labels.NewSelector()
	if matchLabels != nil {
		for _, selector := range matchLabels.MatchExpressions {
			req, err := labels.NewRequirement(selector.Key, selection.Operator(selector.Operator), selector.Values)
			if err != nil {
				return nil, fmt.Errorf("error converting selector: %w", err)
			}
			k8sSelector.Add(*req)
		}
	}

	loggingClusterList := &opniv2beta1.LoggingClusterList{}
	err := k.client.List(ctx, loggingClusterList, client.InNamespace(k.KubernetesStoreOptions.systemNamespace), client.MatchingLabelsSelector{Selector: k8sSelector})
	if err != nil {
		return nil, fmt.Errorf("failed to list clusters: %w", err)
	}

	clusterList := []*core.LoggingCluster{}
	for _, loggingCluster := range loggingClusterList.Items {
		clusterList = append(clusterList, &core.LoggingCluster{
			Id:               loggingCluster.Labels[resources.OpniClusterID],
			Name:             loggingCluster.Spec.FriendlyName,
			OpensearchUserID: derefLocalObjectName(loggingCluster.Spec.IndexUserSecret, ""),
			Labels:           stripIDLabel(loggingCluster.Labels),
		})
	}

	return &core.LoggingClusterList{
		Items: clusterList,
	}, nil
}

func stripIDLabel(labels map[string]string) map[string]string {
	stripped := map[string]string{}
	for k, v := range labels {
		if k != resources.OpniClusterID {
			stripped[k] = v
		}
	}

	return stripped
}

func derefLocalObjectName(objectRef *corev1.LocalObjectReference, defaultName string) string {
	if objectRef == nil {
		return defaultName
	}
	return objectRef.Name
}
