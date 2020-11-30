package webhooks

import (
	"crypto/tls"
	"net"
	"path"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	aipsv1alpha1 "github.com/logikone/autoimagepullsecrets-operator/api/v1alpha1"
)

var cfg *rest.Config
var k8sClient client.Client
var k8sManager ctrl.Manager
var webhookServer *webhook.Server
var testEnv *envtest.Environment

func TestWebhooks(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Webhook Suite")
}

var _ = BeforeSuite(func(done Done) {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("Bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			path.Join("..", "deploy", "crd", "bases"),
		},
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = aipsv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).ToNot(HaveOccurred())

	k8sManager, err = ctrl.NewManager(cfg, ctrl.Options{
		MetricsBindAddress: "127.0.0.1:60555",
		Scheme:             scheme.Scheme,
		CertDir:            path.Join("..", "test", "certs"),
		Host:               "127.0.0.1",
		Port:               10443,
	})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sManager).ToNot(BeNil())

	aipsv1alpha1.Client = k8sManager.GetClient()
	webhookServer = k8sManager.GetWebhookServer()

	err = (&ImagePullSecretPodInjector{
		Client:        k8sManager.GetClient(),
		EventRecorder: k8sManager.GetEventRecorderFor("ips-pod-injector"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = k8sManager.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred())
	}()

	k8sClient = k8sManager.GetClient()
	Expect(k8sClient).ToNot(BeNil())

	By("Running webhook server")

	d := &net.Dialer{Timeout: time.Second * 3}
	Eventually(func() error {
		conn, err := tls.DialWithDialer(d, "tcp", "127.0.0.1:10443", &tls.Config{
			InsecureSkipVerify: true,
		})
		if err != nil {
			return err
		}

		Expect(conn.Close()).To(Succeed())
		return nil
	}).Should(Succeed())

	close(done)
}, 60)

var _ = AfterSuite(func() {
	By("Tearing down test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})
