package main

import (
	"os"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	"github.com/logikone/autoimagepullsecrets-operator/controllers"
	"github.com/logikone/autoimagepullsecrets-operator/webhooks"
)

func setupControllers(mgr ctrl.Manager) {
	if err := (&controllers.NamespaceReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controller").WithName("Namespace"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "error starting namespace reconciler")
		os.Exit(1)
	}

	if err := (&controllers.SecretReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controller").WithName("Secret"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "error starting secret reconciler")
		os.Exit(1)
	}
}

func setupWebhooks(mgr ctrl.Manager) {
	if err := (&webhooks.ImagePullSecretInjector{
		Client:        mgr.GetClient(),
		EventRecorder: mgr.GetEventRecorderFor("image-pull-secrets-injector"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "error starting image pull secret injector")
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
