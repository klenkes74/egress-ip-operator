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

package ocp_static_provisioner

import (
	"context"
	"github.com/go-logr/logr"
	"net"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// The OcpStaticEgressIPProvisioner will manage the IP on the hosts by assigning free IPs from the failure-domain.
type OcpStaticEgressIPProvisioner struct {
	client.Client

	Log logr.Logger
}

func (o OcpStaticEgressIPProvisioner) AddSpecifiedIP(ctx context.Context, ip *net.IP, hostName string) error {
	// TODO 2020-09-19 rlichti implement AddSpecifiedIP in OcpStaticEgressIPProvisioner
	panic("implement me")
}

func (o OcpStaticEgressIPProvisioner) AddRandomIP(ctx context.Context, hostName string) (*net.IP, error) {
	// TODO 2020-09-19 rlichti implement AddRandomIP in OcpStaticEgressIPProvisioner
	panic("implement me")
}

func (o OcpStaticEgressIPProvisioner) AssignCIDR(_ context.Context, _ string) error {
	return nil
}

func (o OcpStaticEgressIPProvisioner) CheckIP(ctx context.Context, ip *net.IP, hostName string) error {
	// TODO 2020-09-19 rlichti implement CheckIP in OcpStaticEgressIPProvisioner
	panic("implement me")
}

func (o OcpStaticEgressIPProvisioner) FindHostForNewIP(ctx context.Context, failureDomain string) (string, error) {
	// TODO 2020-09-19 rlichti implement FindHostForNewIP in OcpStaticEgressIPProvisioner
	panic("implement me")
}

func (o OcpStaticEgressIPProvisioner) MoveIP(ctx context.Context, ip *net.IP, oldHostName string, newHostName string) error {
	// TODO 2020-09-19 rlichti implement MoveIP in OcpStaticEgressIPProvisioner
	panic("implement me")
}

func (o OcpStaticEgressIPProvisioner) RemoveIP(ctx context.Context, ip *net.IP, hostName string) error {
	// TODO 2020-09-19 rlichti implement RemoveIP in OcpStaticEgressIPProvisioner
	panic("implement me")
}
