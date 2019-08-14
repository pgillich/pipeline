// Copyright © 2018 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cluster

import (
	"errors"
	"time"

	"emperror.dev/emperror"

	"github.com/banzaicloud/pipeline/model"
	pkgCluster "github.com/banzaicloud/pipeline/pkg/cluster"
	pkgCommon "github.com/banzaicloud/pipeline/pkg/common"
	"github.com/banzaicloud/pipeline/secret"
)

// DigitalOceanCluster struct for DC
type DigitalOceanCluster struct {
	modelCluster *model.ClusterModel
	APIEndpoint  string
}

// CreateDigitalOceanClusterFromRequest creates ClusterModel struct from the request
func CreateDigitalOceanClusterFromRequest(request *pkgCluster.CreateClusterRequest, orgID uint, userID uint) (*DigitalOceanCluster, error) {
	var nodePools = make([]*model.DOKNodePoolModel, 0, len(request.Properties.CreateClusterDOK.NodePools))
	for name, np := range request.Properties.CreateClusterDOK.NodePools {
		nodePools = append(nodePools, &model.DOKNodePoolModel{
			CreatedBy:    userID,
			Name:         name,
			Count:        np.Count,
			InstanceType: np.InstanceType,
		})
	}

	var cluster DigitalOceanCluster

	cluster.modelCluster = &model.ClusterModel{
		Name:           request.Name,
		Location:       request.Location,
		Cloud:          request.Cloud,
		OrganizationId: orgID,
		CreatedBy:      userID,
		SecretId:       request.SecretId,
		Distribution:   pkgCluster.DigitalOcean,
		DOK: model.DOKClusterModel{
			Region:            request.Properties.CreateClusterDOK.Region,
			KubernetesVersion: request.Properties.CreateClusterDOK.Version,
			NodePools:         nodePools,
		},
		TtlMinutes: request.TtlMinutes,
	}
	return &cluster, nil
}

// CreateCluster creates a new cluster
func (c *DigitalOceanCluster) CreateCluster() error {
	return nil
}

// Persist save the cluster model
// Deprecated: Do not use.
func (c *DigitalOceanCluster) Persist() error {
	return emperror.Wrap(c.modelCluster.Save(), "failed to persist cluster")
}

// DownloadK8sConfig downloads the kubeconfig file from cloud
func (c *DigitalOceanCluster) DownloadK8sConfig() ([]byte, error) {
	return nil, errors.New("DigitalOceanCluster.DownloadK8sConfig is not implemented")
}

// GetName returns the name of the cluster
func (c *DigitalOceanCluster) GetName() string {
	return c.modelCluster.Name
}

// GetCloud returns the cloud type of the cluster
func (c *DigitalOceanCluster) GetCloud() string {
	return pkgCluster.DigitalOcean
}

// GetDistribution returns the distribution type of the cluster
func (c *DigitalOceanCluster) GetDistribution() string {
	return c.modelCluster.Distribution
}

// GetStatus gets cluster status
func (c *DigitalOceanCluster) GetStatus() (*pkgCluster.GetClusterStatusResponse, error) {

	return &pkgCluster.GetClusterStatusResponse{
		Status:            c.modelCluster.Status,
		StatusMessage:     c.modelCluster.StatusMessage,
		Name:              c.modelCluster.Name,
		Location:          c.modelCluster.Location,
		Cloud:             pkgCluster.DigitalOcean,
		Distribution:      pkgCluster.DigitalOcean,
		ResourceID:        c.GetID(),
		Logging:           c.GetLogging(),
		Monitoring:        c.GetMonitoring(),
		ServiceMesh:       c.GetServiceMesh(),
		SecurityScan:      c.GetSecurityScan(),
		CreatorBaseFields: *NewCreatorBaseFields(c.modelCluster.CreatedAt, c.modelCluster.CreatedBy),
		NodePools:         nil,
		Region:            c.modelCluster.Location,
		TtlMinutes:        c.modelCluster.TtlMinutes,
		StartedAt:         c.modelCluster.StartedAt,
	}, nil
}

// DeleteCluster deletes cluster
func (c *DigitalOceanCluster) DeleteCluster() error {
	return nil
}

// UpdateNodePools updates nodes pools of a cluster
func (c *DigitalOceanCluster) UpdateNodePools(request *pkgCluster.UpdateNodePoolsRequest, userId uint) error {
	return nil
}

// UpdateCluster updates the dummy cluster
func (c *DigitalOceanCluster) UpdateCluster(r *pkgCluster.UpdateClusterRequest, _ uint) error {
	// TODO @pgillich implement

	return nil
}

// GetID returns the specified cluster id
func (c *DigitalOceanCluster) GetID() uint {
	return c.modelCluster.ID
}

func (c *DigitalOceanCluster) GetUID() string {
	return c.modelCluster.UID
}

// GetModel returns the whole clusterModel
func (c *DigitalOceanCluster) GetModel() *model.ClusterModel {
	return c.modelCluster
}

// CheckEqualityToUpdate validates the update request
func (c *DigitalOceanCluster) CheckEqualityToUpdate(r *pkgCluster.UpdateClusterRequest) error {
	return nil
}

// AddDefaultsToUpdate adds defaults to update request
func (c *DigitalOceanCluster) AddDefaultsToUpdate(r *pkgCluster.UpdateClusterRequest) {

}

// GetAPIEndpoint returns the Kubernetes Api endpoint
func (c *DigitalOceanCluster) GetAPIEndpoint() (string, error) {
	c.APIEndpoint = "http://cow.org:8080"
	return c.APIEndpoint, nil
}

// DeleteFromDatabase deletes model from the database
func (c *DigitalOceanCluster) DeleteFromDatabase() error {
	return c.modelCluster.Delete()
}

// GetOrganizationId gets org where the cluster belongs
func (c *DigitalOceanCluster) GetOrganizationId() uint {
	return c.modelCluster.OrganizationId
}

// GetLocation gets where the cluster is.
func (c *DigitalOceanCluster) GetLocation() string {
	return c.modelCluster.Location
}

// GetSecretId retrieves the secret id
func (c *DigitalOceanCluster) GetSecretId() string {
	return c.modelCluster.SecretId
}

// GetSshSecretId retrieves the ssh secret id
func (c *DigitalOceanCluster) GetSshSecretId() string {
	return c.modelCluster.SshSecretId
}

// SaveSshSecretId saves the ssh secret id to database
func (c *DigitalOceanCluster) SaveSshSecretId(sshSecretId string) error {
	c.modelCluster.SshSecretId = sshSecretId
	return nil
}

// RequiresSshPublicKey returns false
func (c *DigitalOceanCluster) RequiresSshPublicKey() bool {
	return true
}

// CreateDigitalOceanClusterFromModel creates the cluster from the model
func CreateDigitalOceanClusterFromModel(clusterModel *model.ClusterModel) (*DigitalOceanCluster, error) {
	dummyCluster := DigitalOceanCluster{
		modelCluster: clusterModel,
	}
	return &dummyCluster, nil
}

// SetStatus sets the cluster's status
func (c *DigitalOceanCluster) SetStatus(status, statusMessage string) error {
	return c.modelCluster.UpdateStatus(status, statusMessage)
}

// NodePoolExists returns true if node pool with nodePoolName exists
func (c *DigitalOceanCluster) NodePoolExists(nodePoolName string) bool {
	return false
}

// IsReady checks if the cluster is running according to the cloud provider.
func (c *DigitalOceanCluster) IsReady() (bool, error) {
	return true, nil
}

// ValidateCreationFields validates all field
func (c *DigitalOceanCluster) ValidateCreationFields(r *pkgCluster.CreateClusterRequest) error {
	return nil
}

// GetSecretWithValidation returns secret from vault
func (c *DigitalOceanCluster) GetSecretWithValidation() (*secret.SecretItemResponse, error) {
	return &secret.SecretItemResponse{
		Type: pkgCluster.DigitalOcean,
	}, nil
}

// SaveConfigSecretId saves the config secret id in database
func (c *DigitalOceanCluster) SaveConfigSecretId(configSecretId string) error {
	return c.modelCluster.UpdateConfigSecret(configSecretId)
}

// GetConfigSecretId return config secret id
func (c *DigitalOceanCluster) GetConfigSecretId() string {
	return c.modelCluster.ConfigSecretId
}

func (c *DigitalOceanCluster) GetK8sIpv4Cidrs() (*pkgCluster.Ipv4Cidrs, error) {
	return nil, errors.New("not implemented")
}

// GetK8sConfig returns the Kubernetes config
func (c *DigitalOceanCluster) GetK8sConfig() ([]byte, error) {
	return c.DownloadK8sConfig()
}

// ListNodeNames returns node names to label them
func (c *DigitalOceanCluster) ListNodeNames() (nodeNames pkgCommon.NodeNames, err error) {
	return
}

// RbacEnabled returns true if rbac enabled on the cluster
func (c *DigitalOceanCluster) RbacEnabled() bool {
	return c.modelCluster.RbacEnabled
}

// GetSecurityScan returns true if security scan enabled on the cluster
func (c *DigitalOceanCluster) GetSecurityScan() bool {
	return c.modelCluster.SecurityScan
}

// SetSecurityScan returns true if security scan enabled on the cluster
func (c *DigitalOceanCluster) SetSecurityScan(scan bool) {
	c.modelCluster.SecurityScan = scan
}

// GetLogging returns true if logging enabled on the cluster
func (c *DigitalOceanCluster) GetLogging() bool {
	return c.modelCluster.Logging
}

// SetLogging returns true if logging enabled on the cluster
func (c *DigitalOceanCluster) SetLogging(l bool) {
	c.modelCluster.Logging = l
}

// GetMonitoring returns true if momnitoring enabled on the cluster
func (c *DigitalOceanCluster) GetMonitoring() bool {
	return c.modelCluster.Monitoring
}

// SetMonitoring returns true if monitoring enabled on the cluster
func (c *DigitalOceanCluster) SetMonitoring(l bool) {
	c.modelCluster.Monitoring = l
}

// getScaleOptionsFromModelV1 returns scale options for the cluster
func (c *DigitalOceanCluster) GetScaleOptions() *pkgCluster.ScaleOptions {
	return getScaleOptionsFromModel(c.modelCluster.ScaleOptions)
}

// SetScaleOptions sets scale options for the cluster
func (c *DigitalOceanCluster) SetScaleOptions(scaleOptions *pkgCluster.ScaleOptions) {
	updateScaleOptions(&c.modelCluster.ScaleOptions, scaleOptions)
}

// GetServiceMesh returns true if service mesh is enabled on the cluster
func (c *DigitalOceanCluster) GetServiceMesh() bool {
	return c.modelCluster.ServiceMesh
}

// SetServiceMesh sets service mesh flag on the cluster
func (c *DigitalOceanCluster) SetServiceMesh(m bool) {
	c.modelCluster.ServiceMesh = m
}

// NeedAdminRights returns true if rbac is enabled and need to create a cluster role binding to user
func (c *DigitalOceanCluster) NeedAdminRights() bool {
	return false
}

// GetKubernetesUserName returns the user ID which needed to create a cluster role binding which gives admin rights to the user
func (c *DigitalOceanCluster) GetKubernetesUserName() (string, error) {
	return "", nil
}

// GetCreatedBy returns cluster create userID.
func (c *DigitalOceanCluster) GetCreatedBy() uint {
	return c.modelCluster.CreatedBy
}

// GetTTL retrieves the TTL of the cluster
func (c *DigitalOceanCluster) GetTTL() time.Duration {
	return time.Duration(c.modelCluster.TtlMinutes) * time.Minute
}

// SetTTL sets the lifespan of a cluster
func (c *DigitalOceanCluster) SetTTL(ttl time.Duration) {
	c.modelCluster.TtlMinutes = uint(ttl.Minutes())
}
