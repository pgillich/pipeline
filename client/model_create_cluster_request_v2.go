/*
 * Pipeline API
 *
 * Pipeline v0.3.0 swagger
 *
 * API version: 0.3.0
 * Contact: info@banzaicloud.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

type CreateClusterRequestV2 struct {
	Name         string                     `json:"name"`
	Features     []Feature                  `json:"features,omitempty"`
	SecretId     string                     `json:"secretId,omitempty"`
	SecretName   string                     `json:"secretName,omitempty"`
	SshSecretId  string                     `json:"sshSecretId,omitempty"`
	ScaleOptions ScaleOptions               `json:"scaleOptions,omitempty"`
	Type         string                     `json:"type"`
	Kubernetes   CreatePkeClusterKubernetes `json:"kubernetes"`
	// Non-existent resources will be created in this location. Existing resources that must have the same location as the cluster will be validated against this.
	Location string `json:"location,omitempty"`
	// Required resources will be created in this resource group.
	ResourceGroup string                   `json:"resourceGroup"`
	Network       PkeOnAzureClusterNetwork `json:"network,omitempty"`
	Nodepools     []PkeOnAzureNodePool     `json:"nodepools,omitempty"`
}
