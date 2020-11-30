package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	aipsv1alpha1 "github.com/logikone/autoimagepullsecrets-operator/api/v1alpha1"
)

type SecretReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;update;patch;
// +kubebuilder:rbac:groups="",resources=secrets/status,verbs=get;update;patch;

// +kubebuilder:rbac:groups="autoimagepullsecrets.io",resources=dockerregistries,verbs=get;list;watch;
// +kubebuilder:rbac:groups="autoimagepullsecrets.io",resources=clusterdockerregistries,verbs=get;list;watch;

func (r SecretReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	var secret corev1.Secret

	if err := r.Get(ctx, req.NamespacedName, &secret); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	if IsManagedSecret(&secret) {
		namespacedName, err := ParseNamespacedName(secret.Annotations[SourceAnnotation])
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("error parsing namespaced name: %w", err)
		}

		sourceType, ok := secret.Annotations[SourceTypeAnnotation]
		if !ok {
			return reconcile.Result{Requeue: false}, fmt.Errorf("missing source type annotation [%s]", SourceTypeAnnotation)
		}

		switch sourceType {
		case aipsv1alpha1.ResourceClusterDockerRegistry:
			var clusterDockerRegistry aipsv1alpha1.ClusterDockerRegistry

			if err := r.Get(ctx, namespacedName, &clusterDockerRegistry); err != nil {
				return reconcile.Result{}, err
			}

			if _, err := ctrl.CreateOrUpdate(ctx, r, &secret, func() error {
				return r.MutateSecret(&secret, &clusterDockerRegistry)
			}); err != nil {
				return reconcile.Result{}, err
			}
		case aipsv1alpha1.ResourceDockerRegistry:
			var dockerRegistry aipsv1alpha1.DockerRegistry

			if err := r.Get(ctx, namespacedName, &dockerRegistry); err != nil {
				return reconcile.Result{}, err
			}

			if _, err := ctrl.CreateOrUpdate(ctx, r, &secret, func() error {
				return r.MutateSecret(&secret, &dockerRegistry)
			}); err != nil {
				return reconcile.Result{}, err
			}
		default:
			return reconcile.Result{Requeue: false}, fmt.Errorf("unhandled source type [%s]", sourceType)
		}

	}

	return reconcile.Result{}, nil
}

func (r SecretReconciler) MutateSecret(secret *corev1.Secret, registry aipsv1alpha1.Registry) error {
	if secret.Data == nil {
		secret.Data = map[string][]byte{}
	}

	var rendered []byte
	var err error

	currentConfig, ok := secret.Data[corev1.DockerConfigJsonKey]
	if ok {
		rendered, err = registry.Render(currentConfig)
	} else {
		rendered, err = registry.Render(nil)
	}

	if err != nil {
		return fmt.Errorf("error rendering docker config to json: %w", err)
	}

	secret.Data[corev1.DockerConfigJsonKey] = rendered

	return nil
}

func (r SecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	ctx := context.Background()

	if err := mgr.GetFieldIndexer().IndexField(ctx, &corev1.Secret{}, ManagedSecretIndex, func(object runtime.Object) []string {
		metaObject, ok := object.(metav1.Object)
		if !ok {
			return nil
		}

		if IsManagedSecret(metaObject) {
			namespaceName, err := GetNamespacedName(metaObject)
			if err != nil {
				return nil
			}

			return []string{namespaceName.Name}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error indexing %s field: %w", ManagedSecretIndex, err)
	}

	if err := mgr.GetFieldIndexer().IndexField(ctx, &corev1.Secret{}, SourceSecretIndex, func(object runtime.Object) []string {
		metaObject, ok := object.(metav1.Object)
		if !ok {
			return nil
		}

		if IsSource(metaObject) {
			namespaceName, err := GetNamespacedName(metaObject)
			if err != nil {
				return nil
			}

			return []string{namespaceName.Name}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error indexing %s field: %w", ManagedSecretIndex, err)
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Secret{}, builder.WithPredicates(secretPredicates)).
		Complete(r)
}
