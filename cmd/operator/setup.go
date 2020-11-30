package main

import (
	"os"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	aipsv1alpha1 "github.com/logikone/autoimagepullsecrets-operator/api/v1alpha1"
	"github.com/logikone/autoimagepullsecrets-operator/controllers"
	"github.com/logikone/autoimagepullsecrets-operator/webhooks"
)

func setupControllers(mgr ctrl.Manager) {
	if err := (&controllers.PodReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controller").WithName("Pod"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "error starting pod reconciler")
		os.Exit(1)
	}

	if err := (&controllers.SecretReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controller").WithName("Secret"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "error starting secret reconciler")
		os.Exit(1)
	}
}

func setupWebhooks(mgr ctrl.Manager) {
	if err := (&webhooks.ImagePullSecretPodInjector{
		Client:        mgr.GetClient(),
		EventRecorder: mgr.GetEventRecorderFor("image-pull-secrets-injector"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "error starting image pull secret injector")
		os.Exit(1)
	}

	if err := (&aipsv1alpha1.ClusterDockerRegistry{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", aipsv1alpha1.ResourceClusterDockerRegistry)
		os.Exit(1)
	}

	if err := (&aipsv1alpha1.DockerRegistry{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", aipsv1alpha1.ResourceDockerRegistry)
		os.Exit(1)
	}
}

func setupHealthChecks(mgr ctrl.Manager) {
	if err := mgr.AddHealthzCheck("ping", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to add health check")
		os.Exit(1)
	}

	if err := mgr.AddReadyzCheck("ready", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to add ready check")
		os.Exit(1)
	}
}

func setupClients(mgr ctrl.Manager) {
	aipsv1alpha1.Client = mgr.GetClient()
}
