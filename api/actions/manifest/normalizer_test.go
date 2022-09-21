package manifest_test

import (
	"code.cloudfoundry.org/korifi/api/actions/manifest"
	"code.cloudfoundry.org/korifi/api/payloads"
	"code.cloudfoundry.org/korifi/api/repositories"
	"code.cloudfoundry.org/korifi/tools"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
)

type prcParams struct {
	Command                      *string
	Memory                       *string
	DiskQuota                    *string
	Instances                    *int
	HealthCheckHTTPEndpoint      *string
	HealthCheckInvocationTimeout *int64
	HealthCheckType              *string
	Timeout                      *int64
}

type (
	appParams prcParams
	expParams prcParams
)

var _ = Describe("Normalizer", func() {
	var (
		normalizer        manifest.Normalizer
		defaultDomainName string
		appInfo           payloads.ManifestApplication
		appState          manifest.AppState

		normalizedAppInfo payloads.ManifestApplication
	)

	BeforeEach(func() {
		defaultDomainName = "my.domain"
		appInfo = payloads.ManifestApplication{
			Name:       "my-app",
			Env:        map[string]string{"FOO": "bar"},
			Buildpacks: []string{"buildpack-one", "buildpack-two"},
		}
		appState = manifest.AppState{
			App:       repositories.AppRecord{},
			Processes: nil,
			Routes:    nil,
		}
		normalizer = manifest.NewNormalizer(defaultDomainName)
	})

	JustBeforeEach(func() {
		normalizedAppInfo = normalizer.Normalize(appInfo, appState)
	})

	Describe("app normalization", func() {
		It("preserves the necessary app fields", func() {
			Expect(normalizedAppInfo.Name).To(Equal(appInfo.Name))
			Expect(normalizedAppInfo.NoRoute).To(Equal(appInfo.NoRoute))
			Expect(normalizedAppInfo.Env).To(Equal(appInfo.Env))
			Expect(normalizedAppInfo.Buildpacks).To(Equal(appInfo.Buildpacks))
		})

		When("no-route is set", func() {
			BeforeEach(func() {
				appInfo.NoRoute = true
			})

			It("propagates it", func() {
				Expect(normalizedAppInfo.NoRoute).To(BeTrue())
			})
		})
	})

	Describe("process normalization", func() {
		BeforeEach(func() {
			appInfo.Processes = []payloads.ManifestApplicationProcess{
				{Type: "bob"},
			}
		})

		It("preserves existing processes", func() {
			Expect(normalizedAppInfo.Processes).To(ConsistOf(payloads.ManifestApplicationProcess{Type: "bob"}))
		})

		DescribeTable("when app-level values are provided",
			func(app appParams, createProc bool, process prcParams, effective expParams) {
				appInfo.Memory = app.Memory
				appInfo.DiskQuota = app.DiskQuota
				appInfo.Instances = app.Instances
				appInfo.Command = app.Command
				appInfo.HealthCheckHTTPEndpoint = app.HealthCheckHTTPEndpoint
				appInfo.HealthCheckType = app.HealthCheckType
				appInfo.HealthCheckInvocationTimeout = app.HealthCheckInvocationTimeout
				appInfo.Timeout = app.Timeout

				if createProc {
					appInfo.Processes = append(appInfo.Processes, payloads.ManifestApplicationProcess{
						Type:                         "web",
						Memory:                       process.Memory,
						DiskQuota:                    process.DiskQuota,
						Instances:                    process.Instances,
						Command:                      process.Command,
						HealthCheckHTTPEndpoint:      process.HealthCheckHTTPEndpoint,
						HealthCheckType:              process.HealthCheckType,
						HealthCheckInvocationTimeout: process.HealthCheckInvocationTimeout,
						Timeout:                      process.Timeout,
					})
				}

				updatedAppInfo := normalizer.Normalize(appInfo, appState)
				webProc := getWebProcess(updatedAppInfo)

				Expect(webProc.Memory).To(Equal(effective.Memory))
				Expect(webProc.DiskQuota).To(Equal(effective.DiskQuota))
				Expect(webProc.Instances).To(Equal(effective.Instances))
				Expect(webProc.Command).To(Equal(effective.Command))
				Expect(webProc.HealthCheckHTTPEndpoint).To(Equal(effective.HealthCheckHTTPEndpoint))
				Expect(webProc.HealthCheckType).To(Equal(effective.HealthCheckType))
				Expect(webProc.HealthCheckInvocationTimeout).To(Equal(effective.HealthCheckInvocationTimeout))
				Expect(webProc.Timeout).To(Equal(effective.Timeout))
			},

			// without an existing web process
			Entry("app-level command only",
				appParams{Command: tools.PtrTo("echo boo")}, false, prcParams{},
				expParams{Command: tools.PtrTo("echo boo")}),
			Entry("app-level memory only",
				appParams{Memory: tools.PtrTo("512M")}, false, prcParams{},
				expParams{Memory: tools.PtrTo("512M")}),
			Entry("app-level disk_quota only",
				appParams{DiskQuota: tools.PtrTo("2G")}, false, prcParams{},
				expParams{DiskQuota: tools.PtrTo("2G")}),
			Entry("app-level instances only",
				appParams{Instances: tools.PtrTo(3)}, false, prcParams{},
				expParams{Instances: tools.PtrTo(3)}),
			Entry("app-level healthcheck endpoint only",
				appParams{HealthCheckHTTPEndpoint: tools.PtrTo("/health")}, false, prcParams{},
				expParams{HealthCheckHTTPEndpoint: tools.PtrTo("/health")}),
			Entry("app-level healthcheck type only",
				appParams{HealthCheckType: tools.PtrTo("typo")}, false, prcParams{},
				expParams{HealthCheckType: tools.PtrTo("typo")}),
			Entry("app-level healthcheck invocation timeout only",
				appParams{HealthCheckInvocationTimeout: tools.PtrTo(int64(64))}, false, prcParams{},
				expParams{HealthCheckInvocationTimeout: tools.PtrTo(int64(64))}),
			Entry("app-level timeout only",
				appParams{Timeout: tools.PtrTo(int64(12))}, false, prcParams{},
				expParams{Timeout: tools.PtrTo(int64(12))}),
			Entry("a combination of fields",
				appParams{Memory: tools.PtrTo("512M"), DiskQuota: tools.PtrTo("2G")}, false, prcParams{},
				expParams{Memory: tools.PtrTo("512M"), DiskQuota: tools.PtrTo("2G")}),

			// with an existing web process without the given value set
			Entry("empty proc with app memory",
				appParams{Memory: tools.PtrTo("512M")}, true, prcParams{},
				expParams{Memory: tools.PtrTo("512M")}),
			Entry("empty proc with disk quota",
				appParams{DiskQuota: tools.PtrTo("2G")}, true, prcParams{},
				expParams{DiskQuota: tools.PtrTo("2G")}),
			Entry("empty proc with instances",
				appParams{Instances: tools.PtrTo(3)}, true, prcParams{},
				expParams{Instances: tools.PtrTo(3)}),
			Entry("empty proc with command",
				appParams{Command: tools.PtrTo("echo foo")}, true, prcParams{},
				expParams{Command: tools.PtrTo("echo foo")}),
			Entry("empty proc with healthcheck endpoint",
				appParams{HealthCheckHTTPEndpoint: tools.PtrTo("/health")}, true, prcParams{},
				expParams{HealthCheckHTTPEndpoint: tools.PtrTo("/health")}),
			Entry("empty proc with healthcheck type",
				appParams{HealthCheckType: tools.PtrTo("type1")}, true, prcParams{},
				expParams{HealthCheckType: tools.PtrTo("type1")}),
			Entry("empty proc with healthcheck invocation timeout",
				appParams{HealthCheckInvocationTimeout: tools.PtrTo(int64(45))}, true, prcParams{},
				expParams{HealthCheckInvocationTimeout: tools.PtrTo(int64(45))}),
			Entry("empty proc with timeout",
				appParams{Timeout: tools.PtrTo(int64(32))}, true, prcParams{},
				expParams{Timeout: tools.PtrTo(int64(32))}),

			// with an existing web process with the given value set
			Entry("value from proc memory used",
				appParams{Memory: tools.PtrTo("256M")}, true,
				prcParams{Memory: tools.PtrTo("512M")},
				expParams{Memory: tools.PtrTo("512M")}),
			Entry("value from proc disk_quota used",
				appParams{DiskQuota: tools.PtrTo("2G")}, true,
				prcParams{DiskQuota: tools.PtrTo("3G")},
				expParams{DiskQuota: tools.PtrTo("3G")}),
			Entry("value from proc instances used",
				appParams{Instances: tools.PtrTo(2)}, true,
				prcParams{Instances: tools.PtrTo(3)},
				expParams{Instances: tools.PtrTo(3)}),
			Entry("value from proc command used",
				appParams{Command: tools.PtrTo("echo bar")}, true,
				prcParams{Command: tools.PtrTo("echo foo")},
				expParams{Command: tools.PtrTo("echo foo")}),
			Entry("value from proc healthcheck endpoint used",
				appParams{HealthCheckHTTPEndpoint: tools.PtrTo("/apphealth")}, true,
				prcParams{HealthCheckHTTPEndpoint: tools.PtrTo("/prchealth")},
				expParams{HealthCheckHTTPEndpoint: tools.PtrTo("/prchealth")}),
			Entry("value from proc healthcheck type used",
				appParams{HealthCheckType: tools.PtrTo("apptype")}, true,
				prcParams{HealthCheckType: tools.PtrTo("proctype")},
				expParams{HealthCheckType: tools.PtrTo("proctype")}),
			Entry("value from proc healthcheck invocation timeout used",
				appParams{HealthCheckInvocationTimeout: tools.PtrTo(int64(345))}, true,
				prcParams{HealthCheckInvocationTimeout: tools.PtrTo(int64(34))},
				expParams{HealthCheckInvocationTimeout: tools.PtrTo(int64(34))}),
			Entry("value from proc timeout used",
				appParams{Timeout: tools.PtrTo(int64(25))}, true,
				prcParams{Timeout: tools.PtrTo(int64(2))},
				expParams{Timeout: tools.PtrTo(int64(2))}),

			Entry("fields are individually defaulted from the app if not set on process",
				appParams{
					Memory:    tools.PtrTo("256M"),
					DiskQuota: tools.PtrTo("2G"),
				}, true,
				prcParams{
					Instances: tools.PtrTo(3),
					Command:   tools.PtrTo("echo boo"),
				},
				expParams{
					Instances: tools.PtrTo(3),
					Memory:    tools.PtrTo("256M"),
					DiskQuota: tools.PtrTo("2G"),
					Command:   tools.PtrTo("echo boo"),
				}),
		)
	})

	Describe("route normalization", func() {
		When("default route is set", func() {
			BeforeEach(func() {
				appInfo.DefaultRoute = true
			})

			It("creates a default route", func() {
				Expect(normalizedAppInfo.Routes).To(ConsistOf(
					payloads.ManifestRoute{
						Route: tools.PtrTo("my-app.my.domain"),
					}),
				)
			})

			When("there is already a route in the manifest", func() {
				BeforeEach(func() {
					appInfo.Routes = []payloads.ManifestRoute{{
						Route: tools.PtrTo("bob"),
					}}
				})

				It("does not add a default route", func() {
					Expect(normalizedAppInfo.Routes).To(ConsistOf(
						payloads.ManifestRoute{
							Route: tools.PtrTo("bob"),
						}),
					)
				})
			})

			When("there is already a route resource in the state", func() {
				BeforeEach(func() {
					appState.Routes = map[string]repositories.RouteRecord{
						"bob": {Host: "bob"},
					}
				})

				It("does not add a default route", func() {
					Expect(normalizedAppInfo.Routes).To(BeEmpty())
				})
			})
		})

		When("random route is set", func() {
			BeforeEach(func() {
				appInfo.RandomRoute = true
			})

			It("creates a random route", func() {
				Expect(normalizedAppInfo.Routes).To(HaveLen(1))
			})

			When("there is already a route in the manifest", func() {
				BeforeEach(func() {
					appInfo.Routes = []payloads.ManifestRoute{{
						Route: tools.PtrTo("bob"),
					}}
				})

				It("does not add a random route", func() {
					Expect(normalizedAppInfo.Routes).To(ConsistOf(
						payloads.ManifestRoute{
							Route: tools.PtrTo("bob"),
						}),
					)
				})
			})

			When("there is already a route resource in the state", func() {
				BeforeEach(func() {
					appState.Routes = map[string]repositories.RouteRecord{
						"bob": {Host: "bob"},
					}
				})

				It("does not add a random route", func() {
					Expect(normalizedAppInfo.Routes).To(BeEmpty())
				})
			})
		})
	})

	Describe("deprecated disk-quota handling", func() {
		When("disk-quota is set on process", func() {
			BeforeEach(func() {
				appInfo.Processes = []payloads.ManifestApplicationProcess{
					{
						Type:         "bob",
						AltDiskQuota: tools.PtrTo("123M"),
					},
				}
			})

			It("sets the value to disk_quota", func() {
				Expect(normalizedAppInfo.Processes[0].DiskQuota).To(gstruct.PointTo(Equal("123M")))
			})
		})

		When("disk-quota is set on app", func() {
			BeforeEach(func() {
				//nolint:staticcheck
				appInfo.AltDiskQuota = tools.PtrTo("123M")
			})

			It("sets the value to disk_quota", func() {
				webProc := getWebProcess(normalizedAppInfo)
				Expect(webProc.DiskQuota).To(gstruct.PointTo(Equal("123M")))
			})
		})
	})
})

func getWebProcess(appInfo payloads.ManifestApplication) payloads.ManifestApplicationProcess {
	for _, proc := range appInfo.Processes {
		if proc.Type == "web" {
			return proc
		}
	}

	Fail("no web process")
	return payloads.ManifestApplicationProcess{}
}
