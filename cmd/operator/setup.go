package main

import (
	"os"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/logikone/autoimagepullsecrets-operator/controllers"
	"github.com/logikone/autoimagepullsecrets-operator/webhooks"
)

func setupControllers(mgr ctrl.Manager) error {
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

	return nil
}

func setupWebhooks(mgr ctrl.Manager) error {
	if err := (&webhooks.ImagePullSecretInjector{
		Client:        mgr.GetClient(),
		EventRecorder: mgr.GetEventRecorderFor("image-pull-secrets-injector"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "error starting image pull secret injector")
		os.Exit(1)
	}

	return nil
}
