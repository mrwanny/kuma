package get_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	gomega_types "github.com/onsi/gomega/types"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	"github.com/kumahq/kuma/app/kumactl/cmd"
	"github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	core_model "github.com/kumahq/kuma/pkg/core/resources/model"
	core_store "github.com/kumahq/kuma/pkg/core/resources/store"
	memory_resources "github.com/kumahq/kuma/pkg/plugins/resources/memory"
	test_kumactl "github.com/kumahq/kuma/pkg/test/kumactl"
	test_model "github.com/kumahq/kuma/pkg/test/resources/model"
)

var _ = Describe("kumactl get circuit-breakers", func() {

	circuitBreakerResources := []*mesh.CircuitBreakerResource{
		{
			Spec: &mesh_proto.CircuitBreaker{
				Sources: []*mesh_proto.Selector{
					{
						Match: map[string]string{
							"service": "frontend",
							"version": "0.1",
						},
					},
				},
				Destinations: []*mesh_proto.Selector{
					{
						Match: map[string]string{
							"service": "backend",
						},
					},
				},
				Conf: &mesh_proto.CircuitBreaker_Conf{
					Interval:                    &durationpb.Duration{Seconds: 5},
					BaseEjectionTime:            &durationpb.Duration{Seconds: 5},
					MaxEjectionPercent:          &wrapperspb.UInt32Value{Value: 50},
					SplitExternalAndLocalErrors: false,
					Detectors: &mesh_proto.CircuitBreaker_Conf_Detectors{
						TotalErrors:       &mesh_proto.CircuitBreaker_Conf_Detectors_Errors{},
						GatewayErrors:     &mesh_proto.CircuitBreaker_Conf_Detectors_Errors{},
						LocalErrors:       &mesh_proto.CircuitBreaker_Conf_Detectors_Errors{},
						StandardDeviation: &mesh_proto.CircuitBreaker_Conf_Detectors_StandardDeviation{},
						Failure:           &mesh_proto.CircuitBreaker_Conf_Detectors_Failure{},
					},
				},
			},
			Meta: &test_model.ResourceMeta{
				Mesh: "default",
				Name: "cb1",
			},
		},
		{
			Spec: &mesh_proto.CircuitBreaker{
				Sources: []*mesh_proto.Selector{
					{
						Match: map[string]string{
							"service": "web",
							"version": "0.1",
						},
					},
				},
				Destinations: []*mesh_proto.Selector{
					{
						Match: map[string]string{
							"service": "redis",
						},
					},
				},
				Conf: &mesh_proto.CircuitBreaker_Conf{
					Interval:                    &durationpb.Duration{Seconds: 5},
					BaseEjectionTime:            &durationpb.Duration{Seconds: 5},
					MaxEjectionPercent:          &wrapperspb.UInt32Value{Value: 50},
					SplitExternalAndLocalErrors: false,
					Detectors: &mesh_proto.CircuitBreaker_Conf_Detectors{
						TotalErrors:   &mesh_proto.CircuitBreaker_Conf_Detectors_Errors{Consecutive: &wrapperspb.UInt32Value{Value: 20}},
						GatewayErrors: &mesh_proto.CircuitBreaker_Conf_Detectors_Errors{Consecutive: &wrapperspb.UInt32Value{Value: 10}},
						LocalErrors:   &mesh_proto.CircuitBreaker_Conf_Detectors_Errors{Consecutive: &wrapperspb.UInt32Value{Value: 2}},
						StandardDeviation: &mesh_proto.CircuitBreaker_Conf_Detectors_StandardDeviation{
							RequestVolume: &wrapperspb.UInt32Value{Value: 20},
							MinimumHosts:  &wrapperspb.UInt32Value{Value: 3},
							Factor:        &wrapperspb.DoubleValue{Value: 1.9},
						},
						Failure: &mesh_proto.CircuitBreaker_Conf_Detectors_Failure{
							RequestVolume: &wrapperspb.UInt32Value{Value: 20},
							MinimumHosts:  &wrapperspb.UInt32Value{Value: 3},
							Threshold:     &wrapperspb.UInt32Value{Value: 85},
						},
					},
				},
			},
			Meta: &test_model.ResourceMeta{
				Mesh: "default",
				Name: "cb2",
			},
		},
	}

	Describe("GetCircuitBreakerCmd", func() {

		var rootCmd *cobra.Command
		var buf *bytes.Buffer
		var store core_store.ResourceStore
		rootTime, _ := time.Parse(time.RFC3339, "2008-04-27T16:05:36.995Z")
		BeforeEach(func() {
			// setup
			store = core_store.NewPaginationStore(memory_resources.NewStore())

			rootCtx, err := test_kumactl.MakeRootContext(rootTime, store, mesh.CircuitBreakerResourceTypeDescriptor)
			Expect(err).ToNot(HaveOccurred())

			for _, cb := range circuitBreakerResources {
				err := store.Create(context.Background(), cb, core_store.CreateBy(core_model.MetaToResourceKey(cb.GetMeta())))
				Expect(err).ToNot(HaveOccurred())
			}

			rootCmd = cmd.NewRootCmd(rootCtx)
			buf = &bytes.Buffer{}
			rootCmd.SetOut(buf)
		})

		type testCase struct {
			outputFormat string
			goldenFile   string
			pagination   string
			matcher      func(interface{}) gomega_types.GomegaMatcher
		}

		DescribeTable("kumactl get circuit-breakers -o table|json|yaml",
			func(given testCase) {
				// given
				rootCmd.SetArgs(append([]string{
					"--config-file", filepath.Join("..", "testdata", "sample-kumactl.config.yaml"),
					"get", "circuit-breakers"}, given.outputFormat, given.pagination))

				// when
				err := rootCmd.Execute()
				// then
				Expect(err).ToNot(HaveOccurred())

				// when
				expected, err := ioutil.ReadFile(filepath.Join("testdata", given.goldenFile))
				// then
				Expect(err).ToNot(HaveOccurred())
				// and
				Expect(buf.String()).To(given.matcher(expected))
			},
			Entry("should support Table output by default", testCase{
				outputFormat: "",
				goldenFile:   "get-circuit-breakers.golden.txt",
				matcher: func(expected interface{}) gomega_types.GomegaMatcher {
					return WithTransform(strings.TrimSpace, Equal(strings.TrimSpace(string(expected.([]byte)))))
				},
			}),
			Entry("should support Table output explicitly", testCase{
				outputFormat: "-otable",
				goldenFile:   "get-circuit-breakers.golden.txt",
				matcher: func(expected interface{}) gomega_types.GomegaMatcher {
					return WithTransform(strings.TrimSpace, Equal(strings.TrimSpace(string(expected.([]byte)))))
				},
			}),
			Entry("should support pagination", testCase{
				outputFormat: "-otable",
				goldenFile:   "get-circuit-breakers.pagination.golden.txt",
				pagination:   "--size=1",
				matcher: func(expected interface{}) gomega_types.GomegaMatcher {
					return WithTransform(strings.TrimSpace, Equal(strings.TrimSpace(string(expected.([]byte)))))
				},
			}),
			Entry("should support JSON output", testCase{
				outputFormat: "-ojson",
				goldenFile:   "get-circuit-breakers.golden.json",
				matcher:      MatchJSON,
			}),
			Entry("should support YAML output", testCase{
				outputFormat: "-oyaml",
				goldenFile:   "get-circuit-breakers.golden.yaml",
				matcher:      MatchYAML,
			}),
		)
	})
})
