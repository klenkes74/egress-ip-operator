/*
 * Copyright 2020 Kaiserpfalz EDV-Service, Roland T. Lichti.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FailureDomainSpec defines the desired state of FailureDomain
type EgressIPFailureDomainSpec struct {
	// +kubebuilder:validation:Pattern=\d+.\d+.\d+.\d+/\d+
	// Network is the CIDR of the network. Only needed for provisioner 'operator'
	Cidr         string              `json:"cidr,omitempty"`
	// NodeSelector is the nodeselector of all nodes eligible to get egress ips assigned to.
	NodeSelector corev1.NodeSelector `json:"nodeSelector,omitempty"`
}

// FailureDomainStatus defines the observed state of FailureDomain
type EgressIPFailureDomainStatus struct {
	// +kubebuilder:validation:Enum={"pending","initializing","failed","provisioned","deprovisioned"}
	// Phase is the state of this message. May be pending, initializing, failed or deprovisioned
	Phase string `json:"phase"`
	// Message is a human readable message for this state.
	Message string `json:"message,omitempty"`
	// +kubebuilder:validation:Pattern=\d+.\d+.\d+.\d+(/\d+)?
	// IP is the ip or cidr for this status.
	IP string `json:"ip,omitempty"`
	// Namespace is the namespace this IP belongs to.
	Namespace string `json:"namespace,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// FailureDomain is the Schema for the failuredomains API
type EgressIPFailureDomain struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EgressIPFailureDomainSpec   `json:"spec,omitempty"`
	Status EgressIPFailureDomainStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FailureDomainList contains a list of FailureDomain
type EgressIPFailureDomainList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EgressIPFailureDomain `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EgressIPFailureDomain{}, &EgressIPFailureDomainList{})
}
