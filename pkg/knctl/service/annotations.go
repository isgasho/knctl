/*
Copyright 2018 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cppforlife/knctl/pkg/knctl/util"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	servingclientset "github.com/knative/serving/pkg/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/types"
)

type Annotations struct {
	servingClient servingclientset.Interface
}

func NewAnnotations(servingClient servingclientset.Interface) Annotations {
	return Annotations{servingClient}
}

func (a Annotations) Annotate(revision *v1alpha1.Revision, annotations map[string]interface{}) error {
	mergePatch := map[string]interface{}{
		"metadata": map[string]interface{}{
			"annotations": annotations,
		},
	}

	patchJSON, err := json.Marshal(mergePatch)
	if err != nil {
		return err
	}

	return util.Retry(time.Second, 10*time.Second, func() (bool, error) {
		_, err := a.servingClient.ServingV1alpha1().Revisions(revision.Namespace).Patch(revision.Name, types.MergePatchType, patchJSON)
		if err != nil {
			return false, fmt.Errorf("Annotating revision: %s", err)
		}

		return true, nil
	})
}
