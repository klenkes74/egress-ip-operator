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

package ocp_provisioner

import (
	"context"
	"github.com/go-logr/logr"
	"net"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// OcpDynamicEgressIPProvisioner is basically a no-op provisioner since it only has to handle
type OcpDynamicEgressIPProvisioner struct {
	client.Client

	Log logr.Logger
}

func (o OcpDynamicEgressIPProvisioner) AssignCIDR(ctx context.Context, hostName string) error {
	// TODO 2020-09-19 rlichti Implement the AssignCIDR for dynamic provisioner. Need to find the matching failure-domain and assign the cidr-range to this host.
	panic("implement me")
}

func (o OcpDynamicEgressIPProvisioner) AddSpecifiedIP(_ context.Context, _ *net.IP, _ string) error {
	return nil
}

func (o OcpDynamicEgressIPProvisioner) AddRandomIP(ctx context.Context, hostName string) (*net.IP, error) {
	// TODO 2020-09-19 rlichti Implement the AddRandomIP for dynamic provisioner. Need to find a free IP in the failure-domain and return it to the caller.
	panic("implement me")
}

func (o OcpDynamicEgressIPProvisioner) CheckIP(_ context.Context, _ *net.IP, _ string) error {
	return nil
}

func (o OcpDynamicEgressIPProvisioner) FindHostForNewIP(_ context.Context, _ string) (string, error) {
	return "-no host needed-", nil
}

func (o OcpDynamicEgressIPProvisioner) MoveIP(_ context.Context, ip *net.IP, _ string, _ string) error {
	return nil
}

func (o OcpDynamicEgressIPProvisioner) RemoveIP(_ context.Context, _ *net.IP, _ string) error {
	return nil
}
