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

package cloudmanaged_provisioner

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/klenkes74/egress-ip-operator/pkg/cloudprovider"
	"github.com/klenkes74/egress-ip-operator/pkg/provisioner/ocp_static_provisioner"
	"net"
)

type CloudManagedEgressIPProvisioner struct {
	Log logr.Logger
	// Cloud is the low level interface to the cloudprovider for managing IPs on instances.
	Cloud cloudprovider.CloudProvider
	// OpenShift is the static Egress IP provisioner which is called internally for managing the OCP part of the IP management.
	OpenShift ocp_static_provisioner.OcpStaticEgressIPProvisioner
}

func (a CloudManagedEgressIPProvisioner) AddRandomIP(ctx context.Context, hostName string) (*net.IP, error) {
	ip, err := a.Cloud.AddRandomIP(hostName)
	if err != nil {
		return nil, err
	}

	err = a.OpenShift.AddSpecifiedIP(ctx, ip, hostName)
	if err != nil {
		redoErr := a.Cloud.RemoveIP(ip, hostName)
		if redoErr != nil {
			return nil, fmt.Errorf(
				"error while rolling back adding random ip to host '%v': %v",
				hostName,
				err.Error(),
			)
		}

		return nil, err
	}

	return ip, err
}

func (a CloudManagedEgressIPProvisioner) AddSpecifiedIP(ctx context.Context, ip *net.IP, hostName string) error {
	err := a.Cloud.AddSpecifiedIP(ip, hostName)
	if err != nil {
		return err
	}

	err = a.OpenShift.AddSpecifiedIP(ctx, ip, hostName)
	if err != nil {
		redoErr := a.Cloud.RemoveIP(ip, hostName)
		if redoErr != nil {
			return fmt.Errorf(
				"error while rolling back adding ip '%v' to host '%v': %v",
				ip.String(),
				hostName,
				err.Error(),
			)
		}
	}

	return nil
}

func (a CloudManagedEgressIPProvisioner) AssignCIDR(_ context.Context, _ string) error {
	return nil
}

func (a CloudManagedEgressIPProvisioner) CheckIP(ctx context.Context, ip *net.IP, hostName string) error {
	err := a.Cloud.CheckIP(ip, hostName)
	if err != nil {
		return err
	}

	return a.OpenShift.CheckIP(ctx, ip, hostName)
}

func (a CloudManagedEgressIPProvisioner) FindHostForNewIP(ctx context.Context, failureDomain string) (string, error) {
	return a.OpenShift.FindHostForNewIP(ctx, failureDomain)
}

func (a CloudManagedEgressIPProvisioner) MoveIP(ctx context.Context, ip *net.IP, oldHostName string, newHostName string) error {
	err := a.Cloud.MoveIP(ip, oldHostName, newHostName)
	if err != nil {
		return err
	}

	err = a.OpenShift.MoveIP(ctx, ip, oldHostName, newHostName)
	if err != nil {
		redoErr := a.Cloud.MoveIP(ip, newHostName, oldHostName)
		if redoErr != nil {
			return fmt.Errorf("error while moving IP '%v' from '%v' to '%v': %v", ip.String(), oldHostName, newHostName, redoErr.Error())
		}

		return fmt.Errorf("error while moving IP '%v' from '%v' to '%v'. Change reverted: %v", ip.String(), oldHostName, newHostName, err.Error())
	}

	return nil
}

func (a CloudManagedEgressIPProvisioner) RemoveIP(ctx context.Context, ip *net.IP, hostName string) error {
	err := a.Cloud.RemoveIP(ip, hostName)
	if err != nil {
		return err
	}

	err = a.OpenShift.RemoveIP(ctx, ip, hostName)
	if err != nil {
		redoErr := a.Cloud.AddSpecifiedIP(ip, hostName)
		if redoErr != nil {
			return fmt.Errorf("error while removing IP '%v' from OpenShift. Re-adding it to the cloudprovider failed: %v", ip.String(), redoErr.Error())
		}

		return fmt.Errorf("error while removing IP '%v' from host. IP is still valid: %v", ip.String(), err.Error())
	}

	return nil
}
