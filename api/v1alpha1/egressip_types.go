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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FailureDomainEgressIPSpec defines a single IP within a failureDomain
type FailureDomainEgressIPSpec struct {
	// FailureDomain is the defined failuredomain for this EgressIP. Needs to be defined prior to using it.
	FailureDomain string `json:"failure-domain"`
	// +kubebuilder:validation:Pattern=\d+.\d+.\d+.\d+
	// IP is the IP that should be used for this EgressIP.
	IP string `json:"ip,omitempty"`
}

// EgressIPSpec defines the desired state of EgressIP
type EgressIPSpec struct {
	// IPs is an array of defined EgressIPs. You may list all defined failure domains. At least one needs to be listed.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:UniqueItems=true
	IPs []FailureDomainEgressIPSpec `json:"ips"`
}

// EgressIPStatus defines the observed state of EgressIP
type EgressIPStatus struct {
	// +kubebuilder:validation:Enum={"pending","initializing","failed","provisioned","deprovisioned"}
	// Phase is the state of this message. May be pending, initializing, failed, provisioned or deprovisioned
	Phase string `json:"phase"`
	// IP is the ip or cidr for this status.
	IP FailureDomainEgressIPSpec `json:"ip,omitempty"`
	// HostName is the hostname this IP is assigned to
	HostName string `json:"hostname,omitempty"`
	// Message is a human readable message for this state.
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// EgressIP is the Schema for the egressips API
type EgressIP struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EgressIPSpec   `json:"spec,omitempty"`
	Status EgressIPStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EgressIPList contains a list of EgressIP
type EgressIPList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EgressIP `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EgressIP{}, &EgressIPList{})
}
