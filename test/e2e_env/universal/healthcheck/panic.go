package healthcheck

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/kumahq/kuma/test/framework"
	"github.com/kumahq/kuma/test/framework/envs/universal"
)

func HealthCheckPanicThreshold() {
	const meshName = "hc-panic"

	healthCheck := fmt.Sprintf(`
type: HealthCheck
name: hc-1
mesh: %s
sources:
- match:
    kuma.io/service: '*'
destinations:
- match:
    kuma.io/service: test-server
conf:
  interval: 10s
  timeout: 2s
  unhealthyThreshold: 3
  healthyThreshold: 1
  healthyPanicThreshold: 61
  failTrafficOnPanic: true
  tcp: {}`, meshName)

	dp := func(idx int) string {
		return fmt.Sprintf(`
type: Dataplane
mesh: %s
name: dp-echo-%d
networking:
  address: 192.168.0.%d
  inbound:
  - port: 8080
    servicePort: 80
    tags:
      kuma.io/service: test-server
      kuma.io/protocol: http`, meshName, idx, idx)
	}

	BeforeAll(func() {
		err := NewClusterSetup().
			Install(MeshUniversal(meshName)).
			Install(YamlUniversal(healthCheck)).
			Setup(universal.Cluster)
		Expect(err).ToNot(HaveOccurred())

		for i := 1; i <= 6; i++ {
			dpName := fmt.Sprintf("dp-echo-%d", i)
			response := fmt.Sprintf("universal-%d", i)
			err = TestServerUniversal(dpName, meshName, WithArgs([]string{"echo", "--instance", response}))(universal.Cluster)
			Expect(err).ToNot(HaveOccurred())
		}
		for i := 7; i <= 10; i++ {
			err := NewClusterSetup().Install(YamlUniversal(dp(i))).Setup(universal.Cluster)
			Expect(err).ToNot(HaveOccurred())
		}

		err = DemoClientUniversal(AppModeDemoClient, meshName, WithTransparentProxy(true))(universal.Cluster)
		Expect(err).ToNot(HaveOccurred())
	})

	E2EAfterAll(func() {
		Expect(universal.Cluster.DeleteMeshApps(meshName)).To(Succeed())
		Expect(universal.Cluster.DeleteMesh(meshName)).To(Succeed())
	})

	It("should switch to panic mode and dismiss all requests", func() {
		Eventually(func(g Gomega) {
			stdout, _, _ := universal.Cluster.Exec("", "", "demo-client",
				"curl", "-v", "test-server.mesh")
			g.Expect(stdout).To(ContainSubstring("no healthy upstream"))
		}, "30s", "500ms").Should(Succeed())
	})
}
