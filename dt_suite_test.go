package dt_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dt Suite")
}
