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
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type SecretReconciler struct {
	client.Client

	Log logr.Logger
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

	return reconcile.Result{}, nil
}

func (r SecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	ctx := context.Background()

	predicateFuncs := predicate.Funcs{
		CreateFunc: func(event event.CreateEvent) bool {
			return IsManagedSecret(event.Meta) || IsSourceSecret(event.Meta)
		},
		DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
			return IsManagedSecret(deleteEvent.Meta) || IsSourceSecret(deleteEvent.Meta)
		},
		UpdateFunc: func(updateEvent event.UpdateEvent) bool {
			return (IsManagedSecret(updateEvent.MetaOld) || IsManagedSecret(updateEvent.MetaNew)) ||
				(IsSourceSecret(updateEvent.MetaOld) || IsSourceSecret(updateEvent.MetaNew))
		},
		GenericFunc: func(genericEvent event.GenericEvent) bool {
			return IsManagedSecret(genericEvent.Meta) || IsSourceSecret(genericEvent.Meta)
		},
	}

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

		if IsSourceSecret(metaObject) {
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
		For(&corev1.Secret{}, builder.WithPredicates(predicateFuncs)).
		Complete(r)
}
