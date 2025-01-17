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

package cluster

import (
	"net/http"

	"emperror.dev/emperror"
	"github.com/gin-gonic/gin"
	"github.com/moogar0880/problems"

	"github.com/banzaicloud/pipeline/internal/cluster"
	ginutils "github.com/banzaicloud/pipeline/internal/platform/gin/utils"
	"github.com/banzaicloud/pipeline/pkg/ctxutil"
)

// NewClusterCheckMiddleware returns a new gin middleware that checks cluster is exists in the current org.
func NewClusterCheckMiddleware(manager *Manager, errorHandler emperror.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterID, ok := ginutils.UintParam(c, "id")
		if !ok {
			return
		}

		orgID, ok := ginutils.UintParam(c, "orgid")
		if !ok {
			return
		}

		_, err := manager.GetClusterByID(c, orgID, clusterID)
		if err != nil && cluster.IsClusterNotFoundError(err) {
			problem := problems.NewDetailedProblem(http.StatusNotFound, err.Error())
			c.AbortWithStatusJSON(http.StatusNotFound, problem)

			return
		}
		if err != nil {
			errorHandler.Handle(err)

			problem := problems.NewDetailedProblem(http.StatusInternalServerError, "internal server error")
			c.AbortWithStatusJSON(http.StatusInternalServerError, problem)

			return
		}

		c.Request = c.Request.WithContext(ctxutil.WithClusterID(c.Request.Context(), clusterID))
	}
}
