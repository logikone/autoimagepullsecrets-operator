package webhooks

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
)

func TestWebhooks(t *testing.T) {
	RegisterFailHandler(Fail)

	junitReporter := reporters.NewJUnitReporter("junit.xml")

	RunSpecsWithDefaultAndCustomReporters(t,
		"Webhook Suite",
		[]Reporter{printer.NewlineReporter{}, junitReporter})
}
