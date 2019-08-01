// Copyright © 2019 Banzai Cloud
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

package clusterfeatureadapter

import (
	"context"
	"fmt"
	"time"

	"emperror.dev/errors"
	"go.uber.org/cadence/client"

	"github.com/banzaicloud/pipeline/internal/clusterfeature"
	"github.com/banzaicloud/pipeline/internal/clusterfeature/clusterfeatureadapter/workflow"
)

type asyncFeatureManager struct {
	clusterfeature.FeatureManager
	cadenceClient client.Client
}

// NewAsyncFeatureManager returns a new, asynchronous feature manager wrapping a synchronous implementation
func NewAsyncFeatureManager(syncFeatureManager clusterfeature.FeatureManager, cadenceClient client.Client) clusterfeature.FeatureManager {
	return asyncFeatureManager{
		FeatureManager: syncFeatureManager,
		cadenceClient:  cadenceClient,
	}
}

// Deploys and activates a feature on the given cluster
func (m asyncFeatureManager) Activate(ctx context.Context, clusterID uint, spec clusterfeature.FeatureSpec) error {
	return m.dispatchAction(ctx, clusterID, workflow.ActionActivate, spec)
}

// Removes feature from the given cluster
func (m asyncFeatureManager) Deactivate(ctx context.Context, clusterID uint) error {
	return m.dispatchAction(ctx, clusterID, workflow.ActionDeactivate, nil)
}

// Updates a feature on the given cluster
func (m asyncFeatureManager) Update(ctx context.Context, clusterID uint, spec clusterfeature.FeatureSpec) error {
	return m.dispatchAction(ctx, clusterID, workflow.ActionUpdate, spec)
}

func (m asyncFeatureManager) dispatchAction(ctx context.Context, clusterID uint, action string, spec clusterfeature.FeatureSpec) error {
	const workflowName = workflow.ClusterFeatureJobWorkflowName
	featureName := m.Name()
	workflowID := getWorkflowID(workflowName, clusterID, featureName)
	const signalName = workflow.ClusterFeatureJobSignalName
	signalArg := workflow.ClusterFeatureJobSignalInput{
		Action:        action,
		FeatureSpec:   spec,
		RetryInterval: 1 * time.Minute,
	}
	options := client.StartWorkflowOptions{}
	workflowInput := workflow.ClusterFeatureJobWorkflowInput{
		ClusterID:   clusterID,
		FeatureName: featureName,
	}
	_, err := m.cadenceClient.SignalWithStartWorkflow(ctx, workflowID, signalName, signalArg, options, workflowName, workflowInput)
	if err != nil {
		return errors.WrapIfWithDetails(err, "signal with start workflow failed", "workflowId", workflowID)
	}
	return nil
}

func getWorkflowID(workflowName string, clusterID uint, featureName string) string {
	return fmt.Sprintf("%s-%d-%s", workflowName, clusterID, featureName)
}
