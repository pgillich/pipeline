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

package dok

/*
https://developers.digitalocean.com/documentation/v2/#create-a-new-kubernetes-cluster
*/

// CreateClusterDigitalOcean describes Pipeline's DigitalOcean fields of a CreateCluster request
type CreateClusterDOK struct {
	// TODO Name is auto-generated?
	// Name                 string `json:"name" yaml:"name"`
	Region  string `json:"region" yaml:"region"`
	Version string `json:"version" yaml:"version"`
	/* TODO
	AutoUpgrade bool   `json:"auto_upgrade,omitempty" yaml:"auto_upgrade,omitempty"`
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`
	MaintenancePolicy MaintenancePolicy `json:"maintenance_policy,omitempty" yaml:"maintenance_policy,omitempty"`
	*/
	NodePools map[string]*NodePoolCreate `json:"node_pools,omitempty" yaml:"node_pools,omitempty"`
}

// NodePoolCreate describes Azure's node fields of a CreateCluster request
type NodePoolCreate struct {
	InstanceType string `json:"size" yaml:"size"`
	Count        int    `json:"count" yaml:"count"`
	// TODO ORM string list
	// Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// Validate validates cluster create request
func (d *CreateClusterDOK) Validate() error {
	// TODO implement

	return nil
}

/*
// Validate validates the update request
func (r *UpdateClusterDOK) Validate() error {
	// TODO implement

	return nil
}
*/
