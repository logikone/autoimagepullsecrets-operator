package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/logikone/autoimagepullsecrets-operator/webhooks"
)

type PodReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;update;patch;
// +kubebuilder:rbac:groups="",resources=pods/status,verbs=get;update;patch;

func (r PodReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	var pod corev1.Pod

	if err := r.Get(ctx, req.NamespacedName, &pod); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	var secret corev1.Secret
	secret.Name = fmt.Sprintf(webhooks.IPSInjectionSecretNameFormat, pod.Name)
	secret.Namespace = pod.Namespace

	if _, err := ctrl.CreateOrUpdate(ctx, r, &secret, func() error {
		if err := r.MutateSecret(pod, &secret); err != nil {
			return err
		}

		return ctrl.SetControllerReference(&pod, &secret, r.Scheme)
	}); err != nil {
		return reconcile.Result{}, fmt.Errorf(
			"error creating or updatding secret [%s]: %w", secret.Name, err)
	}

	return reconcile.Result{}, nil
}

func (r PodReconciler) MutateSecret(_ corev1.Pod, secret *corev1.Secret) error {
	if secret.Annotations == nil {
		secret.Annotations = map[string]string{}
	}

	secret.Annotations[ManagedSecretAnnotation] = "true"

	// TODO: setup pod annotations in webhooks to be clusterregistries and registries and comma separated list
	// TODO: use those annotations here to lookup the resources and render the secret, preferably serialized them

	(&corev1.ObjectReference{}).SwaggerDoc()

	return nil
}

func (r PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}, builder.WithPredicates(podPredicates)).
		Complete(r)
}
