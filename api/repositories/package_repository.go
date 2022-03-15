package repositories

import (
	"context"
	"fmt"
	"sort"

	"code.cloudfoundry.org/cf-k8s-controllers/api/apierrors"
	"code.cloudfoundry.org/cf-k8s-controllers/api/authorization"
	workloadsv1alpha1 "code.cloudfoundry.org/cf-k8s-controllers/controllers/apis/workloads/v1alpha1"

	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	kind = "CFPackage"

	PackageStateAwaitingUpload = "AWAITING_UPLOAD"
	PackageStateReady          = "READY"

	PackageResourceType = "Package"
)

//+kubebuilder:rbac:groups=workloads.cloudfoundry.org,resources=cfpackages,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=workloads.cloudfoundry.org,resources=cfpackages/status,verbs=get

//+kubebuilder:rbac:groups="",resources=serviceaccounts;secrets,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=serviceaccounts/status;secrets/status,verbs=get

type PackageRepo struct {
	privilegedClient   client.Client
	namespaceRetriever NamespaceRetriever
	userClientFactory  UserK8sClientFactory
}

func NewPackageRepo(privilegedClient client.Client, namespaceRetriever NamespaceRetriever, userClientFactory UserK8sClientFactory) *PackageRepo {
	return &PackageRepo{
		privilegedClient:   privilegedClient,
		namespaceRetriever: namespaceRetriever,
		userClientFactory:  userClientFactory,
	}
}

type PackageRecord struct {
	GUID      string
	UID       types.UID
	Type      string
	AppGUID   string
	SpaceGUID string
	State     string
	CreatedAt string
	UpdatedAt string
}

type ListPackagesMessage struct {
	AppGUIDs        []string
	SortBy          string
	DescendingOrder bool
	States          []string
}

type CreatePackageMessage struct {
	Type      string
	AppGUID   string
	SpaceGUID string
	OwnerRef  metav1.OwnerReference
}

func (message CreatePackageMessage) toCFPackage() workloadsv1alpha1.CFPackage {
	guid := uuid.NewString()
	return workloadsv1alpha1.CFPackage{
		TypeMeta: metav1.TypeMeta{
			Kind:       kind,
			APIVersion: workloadsv1alpha1.GroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            guid,
			Namespace:       message.SpaceGUID,
			OwnerReferences: []metav1.OwnerReference{message.OwnerRef},
		},
		Spec: workloadsv1alpha1.CFPackageSpec{
			Type: workloadsv1alpha1.PackageType(message.Type),
			AppRef: corev1.LocalObjectReference{
				Name: message.AppGUID,
			},
		},
	}
}

type UpdatePackageSourceMessage struct {
	GUID               string
	SpaceGUID          string
	ImageRef           string
	RegistrySecretName string
}

func (r *PackageRepo) CreatePackage(ctx context.Context, authInfo authorization.Info, message CreatePackageMessage) (PackageRecord, error) {
	userClient, err := r.userClientFactory.BuildClient(authInfo)
	if err != nil {
		return PackageRecord{}, fmt.Errorf("failed to build user client: %w", err)
	}

	cfPackage := message.toCFPackage()
	err = userClient.Create(ctx, &cfPackage)
	if err != nil {
		return PackageRecord{}, apierrors.FromK8sError(err, PackageResourceType)
	}

	return cfPackageToPackageRecord(cfPackage), nil
}

func (r *PackageRepo) GetPackage(ctx context.Context, authInfo authorization.Info, guid string) (PackageRecord, error) {
	ns, err := r.namespaceRetriever.NamespaceFor(ctx, guid, PackageResourceType)
	if err != nil {
		return PackageRecord{}, err
	}

	userClient, err := r.userClientFactory.BuildClient(authInfo)
	if err != nil {
		return PackageRecord{}, fmt.Errorf("failed to build user k8s client: %w", err)
	}

	cfpackage := workloadsv1alpha1.CFPackage{}
	if err := userClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: guid}, &cfpackage); err != nil {
		return PackageRecord{}, fmt.Errorf("failed to get package %q: %w", guid, apierrors.FromK8sError(err, PackageResourceType))
	}

	return cfPackageToPackageRecord(cfpackage), nil
}

func (r *PackageRepo) ListPackages(ctx context.Context, authInfo authorization.Info, message ListPackagesMessage) ([]PackageRecord, error) {
	packageList := &workloadsv1alpha1.CFPackageList{}
	err := r.privilegedClient.List(ctx, packageList)
	if err != nil {
		return []PackageRecord{}, apierrors.FromK8sError(err, PackageResourceType)
	}

	orderedPackages := orderPackages(packageList.Items, message)

	packageRecords := convertToPackageRecords(orderedPackages)

	return applyPackageFilter(packageRecords, message), nil
}

func orderPackages(packages []workloadsv1alpha1.CFPackage, message ListPackagesMessage) []workloadsv1alpha1.CFPackage {
	sort.Slice(packages, func(i, j int) bool {
		if message.SortBy == "created_at" && message.DescendingOrder {
			return !packages[i].CreationTimestamp.Before(&packages[j].CreationTimestamp)
		}
		// For now, we order by created_at by default- if you really want to optimize runtime you can use bucketsort
		return packages[i].CreationTimestamp.Before(&packages[j].CreationTimestamp)
	})

	return packages
}

func applyPackageFilter(packages []PackageRecord, message ListPackagesMessage) []PackageRecord {
	var appFiltered []PackageRecord
	if len(message.AppGUIDs) > 0 {
		for _, currentPackage := range packages {
			for _, appGUID := range message.AppGUIDs {
				if currentPackage.AppGUID == appGUID {
					appFiltered = append(appFiltered, currentPackage)
					break
				}
			}
		}
	} else {
		appFiltered = packages
	}

	var stateFiltered []PackageRecord
	if len(message.States) > 0 {
		for _, currentPackage := range appFiltered {
			for _, state := range message.States {
				if currentPackage.State == state {
					stateFiltered = append(stateFiltered, currentPackage)
					break
				}
			}
		}
	} else {
		stateFiltered = appFiltered
	}

	return stateFiltered
}

func (r *PackageRepo) UpdatePackageSource(ctx context.Context, authInfo authorization.Info, message UpdatePackageSourceMessage) (PackageRecord, error) {
	baseCFPackage := &workloadsv1alpha1.CFPackage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      message.GUID,
			Namespace: message.SpaceGUID,
		},
	}
	cfPackage := baseCFPackage.DeepCopy()
	cfPackage.Spec.Source.Registry.Image = message.ImageRef
	cfPackage.Spec.Source.Registry.ImagePullSecrets = []corev1.LocalObjectReference{{Name: message.RegistrySecretName}}

	userClient, err := r.userClientFactory.BuildClient(authInfo)
	if err != nil {
		return PackageRecord{}, fmt.Errorf("failed to build user k8s client: %w", err)
	}

	err = userClient.Patch(ctx, cfPackage, client.MergeFrom(baseCFPackage))
	if err != nil {
		return PackageRecord{}, fmt.Errorf("failed to update package source: %w", apierrors.FromK8sError(err, PackageResourceType))
	}

	record := cfPackageToPackageRecord(*cfPackage)
	return record, nil
}

func cfPackageToPackageRecord(cfPackage workloadsv1alpha1.CFPackage) PackageRecord {
	updatedAtTime, _ := getTimeLastUpdatedTimestamp(&cfPackage.ObjectMeta)
	state := PackageStateAwaitingUpload
	if cfPackage.Spec.Source.Registry.Image != "" {
		state = PackageStateReady
	}
	return PackageRecord{
		GUID:      cfPackage.Name,
		UID:       cfPackage.UID,
		SpaceGUID: cfPackage.Namespace,
		Type:      string(cfPackage.Spec.Type),
		AppGUID:   cfPackage.Spec.AppRef.Name,
		State:     state,
		CreatedAt: formatTimestamp(cfPackage.CreationTimestamp),
		UpdatedAt: updatedAtTime,
	}
}

func convertToPackageRecords(packages []workloadsv1alpha1.CFPackage) []PackageRecord {
	packageRecords := make([]PackageRecord, 0, len(packages))

	for _, currentPackage := range packages {
		packageRecords = append(packageRecords, cfPackageToPackageRecord(currentPackage))
	}
	return packageRecords
}
