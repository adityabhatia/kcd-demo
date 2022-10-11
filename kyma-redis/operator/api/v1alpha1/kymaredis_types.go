/*
Copyright 2022.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kyma-project/module-manager/operator/pkg/types"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KymaRedisSpec defines the desired state of KymaRedis
type KymaRedisSpec struct {
	ReleaseName string `json:"releaseName,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="State",type=string,JSONPath=".status.state"

// KymaRedis is the Schema for the kymaredis API
type KymaRedis struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KymaRedisSpec `json:"spec,omitempty"`
	Status types.Status  `json:"status,omitempty"`
}

var _ types.CustomObject = &KymaRedis{}

func (in *KymaRedis) ComponentName() string {
	return "kyma-redis-module"
}

func (in *KymaRedis) GetStatus() types.Status {
	return in.Status
}

func (in *KymaRedis) SetStatus(status types.Status) {
	in.Status = status
}

//+kubebuilder:object:root=true

// KymaRedisList contains a list of KymaRedis
type KymaRedisList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KymaRedis `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KymaRedis{}, &KymaRedisList{})
}
