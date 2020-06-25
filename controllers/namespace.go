package controllers

import (
	"context"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type NamespaceReconciler struct {
	client.Client

	Log logr.Logger
}

// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;update;patch;
// +kubebuilder:rbac:groups="",resources=namespaces/status,verbs=get;update;patch;

func (r NamespaceReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	var namespace corev1.Namespace

	if err := r.Get(ctx, req.NamespacedName, &namespace); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	return reconcile.Result{}, nil
}

func (r NamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Namespace{}).
		Complete(r)
}
