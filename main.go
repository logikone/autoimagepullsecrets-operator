/*
Copyright 2020 Chris Larsen.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	selectorsv1alpha1 "github.com/logikone/autoimagepullsecrets-operator/api/v1alpha1"
	"github.com/logikone/autoimagepullsecrets-operator/controllers"
	"github.com/logikone/autoimagepullsecrets-operator/webhooks"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = selectorsv1alpha1.AddToScheme(scheme)
}

func main() {
	var certDir string
	var enableLeaderElection bool
	var metricsAddr string
	var zapOptions zap.Options

	zapOptions.BindFlags(flag.CommandLine)

	flag.StringVar(&certDir, "cert-dir", "/tmp/k8s-webhook-server/serving-certs", "The directory webhook certificates will be loaded from")
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&zapOptions)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		CertDir:            certDir,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "8833e298.autoimagepullsecrets.io",
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		Scheme:             scheme,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	checkError(setupControllers(mgr))
	checkError(setupWebhooks(mgr))

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func checkError(err error) {
	if err != nil {
		setupLog.Error(err, "error during setup")
	}
}

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
