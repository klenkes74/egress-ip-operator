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

package aws_provider

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/go-logr/logr"
	"net"
	"reflect"
	"strings"
)

type AwsCloudProvider struct {
	FailureRegion     string
	MaxIPsPerInstance int

	Client AwsDirectCalls

	Log logr.Logger
}

func (a AwsCloudProvider) AddRandomIP(hostName string) (*net.IP, error) {
	instance, err := a.instanceByHostname(hostName)
	if err != nil {
		return nil, err
	}

	err = a.checkValidInstance(instance)
	if err != nil {
		return nil, err
	}

	err = a.checkMaxIPs(instance)
	if err != nil {
		return nil, err
	}

	interfaceID := instance.NetworkInterfaces[0].NetworkInterfaceId

	addressRequest := ec2.AssignPrivateIpAddressesInput{
		NetworkInterfaceId:             interfaceID,
		SecondaryPrivateIpAddressCount: aws.Int64(int64(1)),
	}
	addressResponse, err := a.Client.AssignPrivateIpAddresses(&addressRequest)
	if err != nil {
		return nil, err
	}

	if len(addressResponse.AssignedPrivateIpAddresses) != 1 {
		ips := make([]string, len(addressResponse.AssignedPrivateIpAddresses))
		for i, ipAddress := range addressResponse.AssignedPrivateIpAddresses {
			ips[i] = *ipAddress.PrivateIpAddress
		}
		return nil, fmt.Errorf("there are no or too much IP address assigned to the eni '%v': [%s]",
			*interfaceID, strings.Join(ips, ","))
	}

	ip := net.ParseIP(*addressResponse.AssignedPrivateIpAddresses[0].PrivateIpAddress)
	a.Log.Info("Assigned IP to eni",
		"eni", interfaceID,
		"ip-address", ip.String())

	return &ip, nil
}

func (a AwsCloudProvider) AddSpecifiedIP(ip *net.IP, hostName string) error {
	return a.moveSpecifiedIP(ip, hostName, false)
}

func (a AwsCloudProvider) moveSpecifiedIP(ip *net.IP, hostName string, allowReassignement bool) error {
	instance, err := a.instanceByHostname(hostName)
	if err != nil {
		return err
	}

	err = a.checkValidInstance(instance)
	if err != nil {
		return err
	}

	err = a.checkMaxIPs(instance)
	if err != nil {
		return err
	}

	err = a.checkIP(instance, ip)
	if err != nil {
		return err
	}

	interfaceID := instance.NetworkInterfaces[0].NetworkInterfaceId

	addressRequest := ec2.AssignPrivateIpAddressesInput{
		AllowReassignment:  &allowReassignement,
		NetworkInterfaceId: aws.String(*interfaceID),
		PrivateIpAddresses: aws.StringSlice([]string{ip.String()}),
	}
	addressResponse, err := a.Client.AssignPrivateIpAddresses(&addressRequest)
	if err != nil {
		return err
	}

	if len(addressResponse.AssignedPrivateIpAddresses) != 1 {
		ips := make([]string, len(addressResponse.AssignedPrivateIpAddresses))
		for i, ipAddress := range addressResponse.AssignedPrivateIpAddresses {
			ips[i] = ipAddress.String()
		}
		return fmt.Errorf("there has been no or too much IP address assigned to the eni '%v': [%s]",
			interfaceID, strings.Join(ips, ","))
	}

	a.Log.Info("Assigned IP to eni",
		"eni", interfaceID,
		"ip-address", ip.String())

	return nil
}

func (a AwsCloudProvider) checkValidInstance(instance *ec2.Instance) error {
	if len(instance.NetworkInterfaces) <= 0 {
		return fmt.Errorf(
			"instance '%v' has no network interface",
			*instance.InstanceId,
		)
	}

	return nil
}

func (a AwsCloudProvider) checkMaxIPs(instance *ec2.Instance) error {
	if len(instance.NetworkInterfaces[0].PrivateIpAddresses) >= a.MaxIPsPerInstance {
		return fmt.Errorf(
			"instance '%v' has already %v IP addresses - maximum of %v reached",
			*instance.InstanceId,
			len(instance.NetworkInterfaces[0].PrivateIpAddresses),
			a.MaxIPsPerInstance,
		)
	}

	return nil
}

func (a AwsCloudProvider) checkIP(instance *ec2.Instance, ip *net.IP) error {
	for _, address := range instance.NetworkInterfaces[0].PrivateIpAddresses {
		if reflect.DeepEqual(net.ParseIP(*address.PrivateIpAddress).String(), ip.String()) {
			return fmt.Errorf(
				"instance '%v' has already a secondary ip '%v' - can not add it again",
				*instance.InstanceId,
				ip.String(),
			)
		}
	}

	return nil
}

func (a AwsCloudProvider) CheckIP(ip *net.IP, hostName string) error {
	instance, err := a.instanceByHostname(hostName)
	if err != nil {
		return err
	}

	return a.checkIPOfInstance(ip, instance)
}

func (a AwsCloudProvider) checkIPOfInstance(ip *net.IP, instance *ec2.Instance) error {
	for _, interfaceIP := range instance.NetworkInterfaces[0].PrivateIpAddresses {
		if ip.String() == *interfaceIP.PrivateIpAddress {
			return nil
		}
	}

	return fmt.Errorf(
		"ip '%v' is not assigned to instance '%v'",
		ip.String(), *instance.InstanceId,
	)
}

func (a AwsCloudProvider) MoveIP(ip *net.IP, oldHostName string, newHostName string) error {
	err := a.CheckIP(ip, oldHostName)
	if err != nil {
		return err
	}

	return a.moveSpecifiedIP(ip, newHostName, true)
}

func (a AwsCloudProvider) RemoveIP(ip *net.IP, hostName string) error {
	instance, err := a.instanceByHostname(hostName)
	if err != nil {
		return err
	}

	err = a.checkValidInstance(instance)
	if err != nil {
		a.Log.Info(
			"host has no network interface or no IPs attached",
			"instance-id", *instance.InstanceId,
			"ip", ip.String(),
		)

		return nil // no network interface means that the ip is not on this host.
	}

	err = a.checkIPOfInstance(ip, instance)
	if err != nil {
		a.Log.Info(
			"ip is not assigned on instance",
			"instance-id", *instance.InstanceId,
			"network-interface-id", *instance.NetworkInterfaces[0].NetworkInterfaceId,
			"ip", ip.String(),
		)

		return nil // since it has not this ip we are fine :-)
	}

	a.Log.Info("removing ip from instance",
		"instance-id", *instance.InstanceId,
		"network-interface-id", *instance.NetworkInterfaces[0].NetworkInterfaceId,
		"ip", ip.String(),
	)
	unAssign := ec2.UnassignPrivateIpAddressesInput{
		NetworkInterfaceId: instance.NetworkInterfaces[0].NetworkInterfaceId,
		PrivateIpAddresses: aws.StringSlice([]string{ip.String()}),
	}

	_, err = a.Client.UnassignPrivateIpAddresses(&unAssign)
	if err != nil {
		return err
	}

	return nil
}

func (a AwsCloudProvider) instanceByHostname(hostName string) (*ec2.Instance, error) {
	filter := ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("private-dns-name"),
				Values: aws.StringSlice([]string{hostName}),
			},
		},
	}

	reservations, err := a.Client.DescribeInstances(&filter)
	if err != nil {
		return nil, err
	}

	if len(reservations.Reservations) > 0 {
		a.Log.Info("found instance",
			"instance-id", reservations.Reservations[0].Instances[0].InstanceId,
		)
	} else {
		return nil, errors.New("no instance found")
	}

	return reservations.Reservations[0].Instances[0], nil
}
