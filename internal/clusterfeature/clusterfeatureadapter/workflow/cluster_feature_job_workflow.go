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

package workflow

import (
	"time"

	"emperror.dev/errors"
	"go.uber.org/cadence/workflow"

	"github.com/banzaicloud/pipeline/internal/clusterfeature"
)

// ClusterFeatureJobWorkflowName is the name the ClusterFeatureJobWorkflow is registered under
const ClusterFeatureJobWorkflowName = "cluster-feature-job"

// ClusterFeatureJobSignalName is the name of signal with which jobs can be sent to the workflow
const ClusterFeatureJobSignalName = "job"

const (
	// ActionActivate identifies the cluster feature activation action
	ActionActivate = "activate"
	// ActionDeactivate identifies the cluster feature deactivation action
	ActionDeactivate = "deactivate"
	// ActionUpdate identifies the cluster feature update action
	ActionUpdate = "update"
)

// ClusterFeatureJobWorkflowInput defines the fixed inputs of the ClusterFeatureJobWorkflow
type ClusterFeatureJobWorkflowInput struct {
	ClusterID   uint
	FeatureName string
}

// ClusterFeatureJobSignalInput defines the dynamic inputs of the ClusterFeatureJobWorkflow
type ClusterFeatureJobSignalInput struct {
	Action        string
	FeatureSpec   clusterfeature.FeatureSpec
	RetryInterval time.Duration
}

type branch bool

const (
	newTry branch = true
	newJob branch = false
)

// ClusterFeatureJobWorkflow executes cluster feature jobs
func ClusterFeatureJobWorkflow(ctx workflow.Context, input ClusterFeatureJobWorkflowInput) error {
	jobsChannel := workflow.GetSignalChannel(ctx, ClusterFeatureJobSignalName)

	var signalInput ClusterFeatureJobSignalInput
	_ = jobsChannel.Receive(ctx, &signalInput) // wait until the first jobs arrives

NewJob:
	{
		activityName, activityInput, err := getActivity(input, signalInput)
		if err != nil {
			return err
		}

	NewTry:
		{
			err := workflow.ExecuteActivity(ctx, activityName, activityInput).Get(ctx, nil)
			if shouldRetry(err) {
				// wait and listen for new jobs
				var br branch
				sel := workflow.NewSelector(ctx)
				sel.AddFuture(workflow.NewTimer(ctx, signalInput.RetryInterval), func(f workflow.Future) {
					br = newTry
				})
				sel.AddReceive(jobsChannel, func(c workflow.Channel, more bool) {
					br = newJob
				})
				sel.Select(ctx)

				switch br {
				case newJob:
					goto NewJob
				case newTry:
					goto NewTry
				}

			} else if err != nil {
				return errors.WrapIfWithDetails(err, "activity execution failed", "activityName", activityName, "activityInput", activityInput)
			}
		}

		// activity completed successfully

		if jobsChannel.ReceiveAsync(&signalInput) {
			goto NewJob // got new job, recur
		}
	}

	return nil
}

func getActivity(workflowInput ClusterFeatureJobWorkflowInput, signalInput ClusterFeatureJobSignalInput) (string, interface{}, error) {
	switch signalInput.Action {
	case ActionActivate:
		return ClusterFeatureActivateActivityName, ClusterFeatureActivateActivityInput{
			ClusterID:   workflowInput.ClusterID,
			FeatureName: workflowInput.FeatureName,
			FeatureSpec: signalInput.FeatureSpec,
		}, nil
	case ActionDeactivate:
		return ClusterFeatureDeactivateActivityName, ClusterFeatureDeactivateActivityInput{
			ClusterID:   workflowInput.ClusterID,
			FeatureName: workflowInput.FeatureName,
		}, nil
	case ActionUpdate:
		return ClusterFeatureUpdateActivityName, ClusterFeatureUpdateActivityInput{
			ClusterID:   workflowInput.ClusterID,
			FeatureName: workflowInput.FeatureName,
			FeatureSpec: signalInput.FeatureSpec,
		}, nil
	default:
		return "", nil, errors.NewWithDetails("unsupported action", "action", signalInput.Action)
	}
}

func shouldRetry(err error) bool {
	var sh interface {
		ShouldRetry() bool
	}
	if errors.As(err, &sh) {
		return sh.ShouldRetry()
	}
	return false
}
