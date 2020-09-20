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

package aws_provider_test

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
	"strconv"
)

var _ = Describe("AssignSpecifiedIP", func() {
	BeforeEach(func() {
		initMock()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should add a specified ip when ip is unused", func() {
		awsDirect.
			EXPECT().DescribeInstances(createDescribeInstancesInput(hostName)).
			Return(
				createDescribeInstancesOutput(hostName, hostId, networkInterfaceId, mainIP, []*net.IP{}),
				nil,
			)

		allowReassignement := false
		awsDirect.
			EXPECT().AssignPrivateIpAddresses(&ec2.AssignPrivateIpAddressesInput{
			AllowReassignment:  &allowReassignement,
			NetworkInterfaceId: aws.String(networkInterfaceId),
			PrivateIpAddresses: aws.StringSlice([]string{ip.String()}),
		}).
			Return(
				&ec2.AssignPrivateIpAddressesOutput{
					AssignedPrivateIpAddresses: []*ec2.AssignedPrivateIpAddress{
						{
							PrivateIpAddress: aws.String(ip.String()),
						},
					},
					NetworkInterfaceId: aws.String(networkInterfaceId),
				},
				nil,
			)

		err := sut.AddSpecifiedIP(ip, hostName)

		Expect(err).To(BeNil())
	})

	It("should throw an error when specified ip is already used", func() {
		expectedErr := fmt.Errorf("instance '%v' has already a secondary ip '%v' - can not add it again", hostId, ip.String())

		secondaryIPs := make([]*net.IP, 1)
		secondaryIPs[0] = ip

		awsDirect.
			EXPECT().DescribeInstances(createDescribeInstancesInput(hostName)).
			Return(
				createDescribeInstancesOutput(hostName, hostId, networkInterfaceId, mainIP, secondaryIPs),
				nil,
			)

		err := sut.AddSpecifiedIP(ip, hostName)

		Expect(err).To(MatchError(expectedErr))
	})

	It("should return an error when cloud instance can't be found", func() {
		expectedErr := fmt.Errorf("host '%v' not found", hostName)

		awsDirect.
			EXPECT().DescribeInstances(createDescribeInstancesInput(hostName)).
			Return(
				nil,
				expectedErr,
			)

		err := sut.AddSpecifiedIP(ip, hostName)

		Expect(err).To(MatchError(expectedErr))
	})

	It("should return failure when specified ip is already used", func() {
		awsDirect.
			EXPECT().DescribeInstances(createDescribeInstancesInput(hostName)).
			Return(
				createDescribeInstancesOutput(hostName, hostId, networkInterfaceId, mainIP, []*net.IP{}),
				nil,
			)

		expectedErr := fmt.Errorf("ip '%v' already in use", ip)

		allowReassignement := false
		awsDirect.
			EXPECT().AssignPrivateIpAddresses(&ec2.AssignPrivateIpAddressesInput{
			AllowReassignment:  &allowReassignement,
			NetworkInterfaceId: aws.String(networkInterfaceId),
			PrivateIpAddresses: aws.StringSlice([]string{ip.String()}),
		}).
			Return(
				nil,
				expectedErr,
			)

		err := sut.AddSpecifiedIP(ip, hostName)

		Expect(err).To(MatchError(expectedErr))
	})

	It("should return an error when cloud instance has already too many IPs", func() {
		expectedErr := fmt.Errorf("instance '%v' has already %v IP addresses - maximum of %v reached", hostId, maxIPsPerInstance, maxIPsPerInstance)

		secondaryIPs := make([]*net.IP, sut.MaxIPsPerInstance-1)
		for i := 0; i < sut.MaxIPsPerInstance-1; i++ {
			s := strconv.Itoa(20 + i)
			ip := net.ParseIP("10.0.1." + s)
			secondaryIPs[i] = &ip
		}

		awsDirect.
			EXPECT().DescribeInstances(createDescribeInstancesInput(hostName)).
			Return(
				createDescribeInstancesOutput(hostName, hostId, networkInterfaceId, mainIP, secondaryIPs),
				nil,
			)

		err := sut.AddSpecifiedIP(ip, hostName)

		Expect(err).To(MatchError(expectedErr))
	})

	It("should return an error when cloud instance has no interface", func() {
		expectedErr := fmt.Errorf("instance '%v' has no network interface", hostId)
		awsDirect.
			EXPECT().DescribeInstances(createDescribeInstancesInput(hostName)).
			Return(
				createDescribeInstancesOutput(hostName, hostId, "", mainIP, []*net.IP{}),
				nil,
			)

		err := sut.AddSpecifiedIP(ip, hostName)

		Expect(err).To(MatchError(expectedErr))
	})

})
