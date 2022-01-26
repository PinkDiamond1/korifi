package repositories_test

import (
	"context"

	. "github.com/onsi/gomega/gstruct"

	"code.cloudfoundry.org/cf-k8s-controllers/api/repositories"
	. "code.cloudfoundry.org/cf-k8s-controllers/api/repositories"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	buildv1alpha2 "github.com/pivotal/kpack/pkg/apis/build/v1alpha2"
	buildv1alpha1 "github.com/pivotal/kpack/pkg/apis/core/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("BuildpackRepository", func() {
	var (
		beforeCtx                 context.Context
		buildpackRepo             *BuildpackRepository
		clientFactory             repositories.UserK8sClientFactory
		spaceDeveloperClusterRole *rbacv1.ClusterRole
	)

	BeforeEach(func() {
		beforeCtx = context.Background()
		clientFactory = repositories.NewUnprivilegedClientFactory(k8sConfig)
	})

	Describe("GetBuildpacksForBuilder", func() {
		var clusterBuilder *buildv1alpha2.ClusterBuilder

		BeforeEach(func() {
			clusterBuilder = &buildv1alpha2.ClusterBuilder{
				ObjectMeta: metav1.ObjectMeta{
					Name: generateGUID(),
				},
				Spec: buildv1alpha2.ClusterBuilderSpec{
					BuilderSpec: buildv1alpha2.BuilderSpec{
						Tag: "registry/builder-image",
						Stack: corev1.ObjectReference{
							Kind: "ClusterStack",
							Name: "some-cluster-stack",
						},
						Store: corev1.ObjectReference{
							Kind: "ClusterStore",
							Name: "some-cluster-store",
						},
						Order: []buildv1alpha1.OrderEntry{
							{
								Group: []buildv1alpha1.BuildpackRef{
									newBuildpackRef("paketo-buildpacks/buildpack-1-1"),
								},
							},
							{
								Group: []buildv1alpha1.BuildpackRef{
									newBuildpackRef("paketo-buildpacks/buildpack-2-1"),
									newBuildpackRef("paketo-buildpacks/buildpack-2-2"),
									newBuildpackRef("paketo-buildpacks/buildpack-2-3"),
								},
							},
							{
								Group: []buildv1alpha1.BuildpackRef{
									newBuildpackRef("paketo-buildpacks/buildpack-3-1"),
								},
							},
						},
					},
					ServiceAccountRef: corev1.ObjectReference{
						Namespace: "some-namespace",
						Name:      "some-service-account",
					},
				},
			}

			Expect(k8sClient.Create(beforeCtx, clusterBuilder)).To(Succeed())

			clusterBuilderOrderStatus := []buildv1alpha1.OrderEntry{
				{
					Group: []buildv1alpha1.BuildpackRef{
						newBuildpackRef("paketo-buildpacks/buildpack-1-1", "1.1"),
					},
				},
				{
					Group: []buildv1alpha1.BuildpackRef{
						newBuildpackRef("paketo-buildpacks/buildpack-2-1", "2.1"),
						newBuildpackRef("paketo-buildpacks/buildpack-2-2", "2.2"),
						newBuildpackRef("paketo-buildpacks/buildpack-2-3", "2.3"),
					},
				},
				{
					Group: []buildv1alpha1.BuildpackRef{
						newBuildpackRef("paketo-buildpacks/buildpack-3-1", "3.1"),
					},
				},
			}
			clusterBuilder.Status.Order = clusterBuilderOrderStatus

			Expect(k8sClient.Status().Update(beforeCtx, clusterBuilder)).To(Succeed())
			DeferCleanup(func() {
				_ = k8sClient.Delete(context.Background(), clusterBuilder)
			})
		})

		Describe("List", func() {
			When("auth is enabled", func() {
				It("returns records matching the buildpacks of the ClusterBuilder and no error", func() {
					buildpackRepo = NewBuildpackRepository(k8sClient, clientFactory, true)
					spaceDeveloperClusterRole = createClusterRole(beforeCtx, repositories.SpaceDeveloperClusterRoleRules)
					createClusterRoleBinding(beforeCtx, userName, spaceDeveloperClusterRole.Name)

					buildpackRecords, err := buildpackRepo.GetBuildpacksForBuilder(context.Background(), authInfo, clusterBuilder.Name)
					Expect(err).NotTo(HaveOccurred())
					Expect(buildpackRecords).To(ConsistOf(
						MatchFields(IgnoreExtras, Fields{
							"Name":     Equal("paketo-buildpacks/buildpack-1-1"),
							"Position": Equal(1),
							"Stack":    Equal(clusterBuilder.Spec.Stack.Name),
							"Version":  Equal("1.1"),
						}),
						MatchFields(IgnoreExtras, Fields{
							"Name":     Equal("paketo-buildpacks/buildpack-2-1"),
							"Position": Equal(2),
							"Stack":    Equal(clusterBuilder.Spec.Stack.Name),
							"Version":  Equal("2.1"),
						}),
						MatchFields(IgnoreExtras, Fields{
							"Name":     Equal("paketo-buildpacks/buildpack-3-1"),
							"Position": Equal(3),
							"Stack":    Equal(clusterBuilder.Spec.Stack.Name),
							"Version":  Equal("3.1"),
						}),
					))
				})
			})

			When("auth is enabled, but insufficient perm", func() {
				It("fails to retrieve buildpack records", func() {
					buildpackRepo = NewBuildpackRepository(k8sClient, clientFactory, true)
					_, err := buildpackRepo.GetBuildpacksForBuilder(context.Background(), authInfo, clusterBuilder.Name)
					Expect(err).To(HaveOccurred())
				})
			})

			When("auth is disabled", func() {
				It("returns records matching the buildpacks of the ClusterBuilder and no error", func() {
					authDisabledBuildpackRepo := NewBuildpackRepository(k8sClient, clientFactory, false)
					buildpackRecords, err := authDisabledBuildpackRepo.GetBuildpacksForBuilder(context.Background(), authInfo, clusterBuilder.Name)

					Expect(err).NotTo(HaveOccurred())
					Expect(buildpackRecords).To(ConsistOf(
						MatchFields(IgnoreExtras, Fields{
							"Name":     Equal("paketo-buildpacks/buildpack-1-1"),
							"Position": Equal(1),
							"Stack":    Equal(clusterBuilder.Spec.Stack.Name),
							"Version":  Equal("1.1"),
						}),
						MatchFields(IgnoreExtras, Fields{
							"Name":     Equal("paketo-buildpacks/buildpack-2-1"),
							"Position": Equal(2),
							"Stack":    Equal(clusterBuilder.Spec.Stack.Name),
							"Version":  Equal("2.1"),
						}),
						MatchFields(IgnoreExtras, Fields{
							"Name":     Equal("paketo-buildpacks/buildpack-3-1"),
							"Position": Equal(3),
							"Stack":    Equal(clusterBuilder.Spec.Stack.Name),
							"Version":  Equal("3.1"),
						}),
					))
				})
			})
		})
	})
})

func newBuildpackRef(id string, version ...string) buildv1alpha1.BuildpackRef {
	toReturn := buildv1alpha1.BuildpackRef{
		BuildpackInfo: buildv1alpha1.BuildpackInfo{
			Id: id,
		},
		Optional: true,
	}

	if len(version) > 0 {
		toReturn.BuildpackInfo.Version = version[0]
	}
	return toReturn
}