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
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// AwsDirectCalls is the interface for accessing AWS services. It is the final interface to be able to mock the AWS
// calls during testing.
type AwsDirectCalls interface {
	AssignPrivateIpAddresses(filter *ec2.AssignPrivateIpAddressesInput) (*ec2.AssignPrivateIpAddressesOutput, error)
	DescribeInstances(filter *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)
	UnassignPrivateIpAddresses(filter *ec2.UnassignPrivateIpAddressesInput) (*ec2.UnassignPrivateIpAddressesOutput, error)
}

var _ AwsDirectCalls = &AwsDirectCallsProd{}

// AwsDirectCallsProd is the working implementation of the AwsDirectCalls interface.
type AwsDirectCallsProd struct {
	Session *session.Session
	Client  *ec2.EC2
}

// AssignPrivateIpAddresses calls assign-private-ip-addresses and returns either the output or an error.
func (a *AwsDirectCallsProd) AssignPrivateIpAddresses(filter *ec2.AssignPrivateIpAddressesInput) (*ec2.AssignPrivateIpAddressesOutput, error) {
	return a.Client.AssignPrivateIpAddresses(filter)
}

// DescribeInstances calls describe-instances at AWS and returns either the output or an error.
func (a *AwsDirectCallsProd) DescribeInstances(filter *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return a.Client.DescribeInstances(filter)
}

// UnassignPrivateIpAddresses calls unassign-private-ip-addresses and returns either the output or an error.
func (a *AwsDirectCallsProd) UnassignPrivateIpAddresses(filter *ec2.UnassignPrivateIpAddressesInput) (*ec2.UnassignPrivateIpAddressesOutput, error) {
	return a.Client.UnassignPrivateIpAddresses(filter)
}
