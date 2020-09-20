//go:generate go run github.com/golang/mock/mockgen -package provisioner_test -destination ./mock_provisioner_test.go github.com/klenkes74/egress-ip-operator/pkg/provisioner EgressIPProvisioner

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

// provisioner contains the low level provisioners for handling the IP provisioning.
//
// Currently there are three strategies defined:
// - 'aws' - the operator will call AWS for random IP assignement.
// - 'ocp-static' - the operator will manage which IPs are configured on which node.
// - 'ocp-dynamic' - the operator will add the CIDR range to all matching hosts and OpenShift will manage the host
//   IP networking.
package provisioner

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/klenkes74/egress-ip-operator/pkg/cloudprovider"
	"github.com/klenkes74/egress-ip-operator/pkg/provisioner/cloudmanaged_provisioner"
	"github.com/klenkes74/egress-ip-operator/pkg/provisioner/ocp_dynamic_provisioner"
	"github.com/klenkes74/egress-ip-operator/pkg/provisioner/ocp_static_provisioner"
	"net"
	"os"
)

// EgressIPProvisioner is the low level IP manager for
type EgressIPProvisioner interface {
	// FindHostForNewIP searches for a host in the failure domain to add an IP to.
	// Will return the hostname or an error.
	FindHostForNewIP(ctx context.Context, failureDomain string) (string, error)
	// AddSpecifiedIP adds a predefined IP to the specified host.
	// It will return an error or nil.
	AddSpecifiedIP(ctx context.Context, ip *net.IP, hostName string) error
	// AddRandomIP adds a random IP to the specified host.
	// It will return the IP or the error.
	AddRandomIP(ctx context.Context, hostName string) (*net.IP, error)
	// RemoveIP will remove the given IP from the specified host.
	// It will return an error or nil.
	RemoveIP(ctx context.Context, ip *net.IP, hostName string) error
	// MoveIP will move the specified IP from oldHost to newHost.
	// It will return an error or nil.
	MoveIP(ctx context.Context, ip *net.IP, oldHostName string, newHostName string) error
	// CheckIP will check if the specified IP is assigned on the specified host.
	// it will return an error or nil.
	CheckIP(ctx context.Context, ip *net.IP, hostName string) error
	// AssignCIDR will assign the cidr range to a host.
	// Basically it is only needed by the provisioner 'ocp-dynamic'. The other provisioners will be no-ops.
	AssignCIDR(ctx context.Context, hostName string) error
}

var _ EgressIPProvisioner = &cloudmanaged_provisioner.CloudManagedEgressIPProvisioner{}
var _ EgressIPProvisioner = &ocp_dynamic_provisioner.OcpDynamicEgressIPProvisioner{}
var _ EgressIPProvisioner = &ocp_static_provisioner.OcpStaticEgressIPProvisioner{}

func NewEgressIPProvisioner(logger logr.Logger) (*EgressIPProvisioner, error) {
	var result EgressIPProvisioner

	provisionerType, found := os.LookupEnv("EGRESSIP_PROVISIONER")
	if !found {
		return nil, errors.New("no provisioner defined - please set environment 'EGRESSIP_PROVISIONER")
	}

	switch provisionerType {
	case "cloud":
		cloudProviderType, found := os.LookupEnv("CLOUD_PROVIDER")
		if !found {
			return nil, errors.New("no cloud provider type defined - please set environment 'CLOUD_PROVIDER")
		}

		cloud, err := cloudprovider.NewCloudProvider(cloudProviderType, logger.WithName("cloud"))
		if err != nil {
			return nil, err
		}
		provider := &cloudmanaged_provisioner.CloudManagedEgressIPProvisioner{
			Cloud: *cloud,
			Log:   logger,
		}
		result = EgressIPProvisioner(provider)
	case "ocp-dynamic":
		provider := &ocp_dynamic_provisioner.OcpDynamicEgressIPProvisioner{
			Log: logger.WithName("ocp-dynamic"),
		}
		result = EgressIPProvisioner(provider)
	case "ocp-static":
		provider := &ocp_static_provisioner.OcpStaticEgressIPProvisioner{
			Log: logger.WithName("ocp-static"),
		}
		result = EgressIPProvisioner(provider)
	default:
		return nil, fmt.Errorf("cloudprovider type '%v' is not defined - please use one of: 'cloud', 'ocp-dynamic', or 'ocp-static", provisionerType)
	}

	return &result, nil
}
