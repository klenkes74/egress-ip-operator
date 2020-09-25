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

// ocp_static_provisioner provisions the EgressIP according to the documentation "Configuring manually assigned egress
// IP addresses for a namespace"
// (https://docs.openshift.com/container-platform/4.5/networking/openshift_sdn/assigning-egress-ips.html#nw-egress-ips-static_egress-ips
// for OCP 4.x) or "Enabling Static IPs for External Project Traffic"
// (https://docs.openshift.com/container-platform/3.11/admin_guide/managing_networking.html#enabling-static-ips-for-external-project-traffic
// for OCP 3.11).
//
// It manages the CIDR defined in the EgressIPFailureDomain objects and will select IPs from that range when calling
// AddRandomIP.
// The range needs to be in ownership of this operator - the cluster admin is responsible to select a cidr that is not
// managed by other means.
//
// The IPs itself are managed by OCP. OCP will manage virtual network interfaces on the main interface of the host as
// configured on the HostSubnet within OCP.
//
// In addition the cloudmanaged_provisioner uses the ocp_static_provisioner as backend to handle the OCP configuration.
package provisioner

import (
	"context"
	"github.com/go-logr/logr"
	v1 "github.com/openshift/api/network/v1"
	"k8s.io/apimachinery/pkg/types"
	"net"
)

// The OcpStaticEgressIPProvisioner will manage the IP on the hosts by assigning free IPs from the failure-domain.
type OcpStaticEgressIPProvisioner struct {
	Client OCPDirectCalls

	Log logr.Logger
}

// AddSpecifiedIP assignes a random IP to the host specified by its hostname. The IP will be configured on the
// HostSubnet of the matching host. The IP will be taken from any matching failure domain of the HostSubnet.
func (o OcpStaticEgressIPProvisioner) AddRandomIP(ctx context.Context, hostName string, failureDomain string) (*net.IP, string, error) {
	hostName, ip, err := o.FindHostForNewIP(ctx, failureDomain)
	if err != nil {
		return nil, hostName, err
	}

	hostSubnet, err := o.loadHostSubNet(ctx, hostName)
	if err != nil {
		return nil, hostName, err
	}

	hostSubnet.EgressIPs = append(hostSubnet.EgressIPs, ip.String())

	return &ip, hostName, o.Client.Update(ctx, hostSubnet)
}

func (o OcpStaticEgressIPProvisioner) loadHostSubNet(ctx context.Context, hostName string) (*v1.HostSubnet, error) {
	name := o.createNamespacedForHostSubnet(hostName)

	hostSubNet := &v1.HostSubnet{}
	err := o.Client.Get(ctx, *name, hostSubNet)
	if err != nil {
		return nil, err
	}

	return hostSubNet, nil
}

func (o OcpStaticEgressIPProvisioner) createNamespacedForHostSubnet(hostName string) *types.NamespacedName {
	name := types.NamespacedName{
		Namespace: "",
		Name:      hostName,
	}
	return &name
}

// AddSpecifiedIP assignes the specified IP to the host specified by its hostname. The IP will be configured on the
// HostSubnet of the matching host.
func (o OcpStaticEgressIPProvisioner) AddSpecifiedIP(ctx context.Context, ip *net.IP, hostName string) error {
	// TODO 2020-09-19 rlichti implement AddSpecifiedIP in OcpStaticEgressIPProvisioner
	panic("implement me")
}

// AssignCIDR is used by the dynamic provisioner to manage the EgressCIDR on the HostSubntets. It's a no-op here.
func (o OcpStaticEgressIPProvisioner) AssignCIDR(_ context.Context, _ string) error {
	return nil
}

// CheckIP checks if the specified IP is assigned on the specified host.
func (o OcpStaticEgressIPProvisioner) CheckIP(ctx context.Context, ip *net.IP, hostName string) error {
	// TODO 2020-09-19 rlichti implement CheckIP in OcpStaticEgressIPProvisioner
	panic("implement me")
}

// FindHostForNewIP selects the host within the failureDomain to get a new IP assigned to. Will give either a hostName
// or an error if there are no eligible hosts within the specified failureDomain.
func (o OcpStaticEgressIPProvisioner) FindHostForNewIP(ctx context.Context, failureDomain string) (string, net.IP, error) {
	// TODO 2020-09-19 rlichti implement FindHostForNewIP in OcpStaticEgressIPProvisioner
	panic("implement me")
}

// MoveIP moves the specified IP from the old host to a new host.
func (o OcpStaticEgressIPProvisioner) MoveIP(ctx context.Context, ip *net.IP, oldHostName string, newHostName string) error {
	// TODO 2020-09-19 rlichti implement MoveIP in OcpStaticEgressIPProvisioner
	panic("implement me")
}

// RemoveIP removes the specified IP from the specified host. If the IP is already not at the host, no error is raised.
func (o OcpStaticEgressIPProvisioner) RemoveIP(ctx context.Context, ip *net.IP, hostName string) error {
	// TODO 2020-09-19 rlichti implement RemoveIP in OcpStaticEgressIPProvisioner
	panic("implement me")
}
