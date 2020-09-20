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
	"strings"
)

var _ = Describe("AssignRandomIP", func() {
	BeforeEach(func() {
		initMock()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should add a random ip", func() {
		awsDirect.
			EXPECT().DescribeInstances(createDescribeInstancesInput(hostName)).
			Return(
				createDescribeInstancesOutput(hostName, hostId, networkInterfaceId, mainIP, []*net.IP{}),
				nil,
			)

		awsDirect.
			EXPECT().AssignPrivateIpAddresses(&ec2.AssignPrivateIpAddressesInput{
			NetworkInterfaceId:             aws.String(networkInterfaceId),
			SecondaryPrivateIpAddressCount: aws.Int64(int64(1)),
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

		ip, err := sut.AddRandomIP(hostName)

		Expect(ip).ToNot(BeNil())
		Expect(ip.String()).To(Equal(ip.String()))
		Expect(err).To(BeNil())
	})

	It("should throw an error when there is not exactely one IP assigned", func() {
		ips := []*ec2.AssignedPrivateIpAddress{
			{
				PrivateIpAddress: aws.String(ip.String()),
			},
			{
				PrivateIpAddress: aws.String("10.1.21.2"),
			},
		}

		ipStrings := make([]string, len(ips))
		for i, ipAddress := range ips {
			ipStrings[i] = *ipAddress.PrivateIpAddress
		}

		expectedErr := fmt.Errorf(
			"there are no or too much IP address assigned to the eni '%v': [%v]",
			networkInterfaceId,
			strings.Join(ipStrings, ","),
		)

		awsDirect.
			EXPECT().DescribeInstances(createDescribeInstancesInput(hostName)).
			Return(
				createDescribeInstancesOutput(hostName, hostId, networkInterfaceId, mainIP, []*net.IP{}),
				nil,
			)

		awsDirect.
			EXPECT().AssignPrivateIpAddresses(&ec2.AssignPrivateIpAddressesInput{
			NetworkInterfaceId:             aws.String(networkInterfaceId),
			SecondaryPrivateIpAddressCount: aws.Int64(int64(1)),
		}).
			Return(
				&ec2.AssignPrivateIpAddressesOutput{
					AssignedPrivateIpAddresses: ips,
					NetworkInterfaceId:         aws.String(networkInterfaceId),
				},
				nil,
			)

		ip, err := sut.AddRandomIP(hostName)

		Expect(ip).To(BeNil())
		Expect(err).To(MatchError(expectedErr))

	})

	It("should throw an error when there are no IPs left", func() {
		expectedErr := fmt.Errorf("no ips available to instance '%v'", hostId)

		awsDirect.
			EXPECT().DescribeInstances(createDescribeInstancesInput(hostName)).
			Return(
				createDescribeInstancesOutput(hostName, hostId, networkInterfaceId, mainIP, []*net.IP{}),
				nil,
			)

		awsDirect.
			EXPECT().AssignPrivateIpAddresses(&ec2.AssignPrivateIpAddressesInput{
			NetworkInterfaceId:             aws.String(networkInterfaceId),
			SecondaryPrivateIpAddressCount: aws.Int64(int64(1)),
		}).
			Return(
				nil,
				expectedErr,
			)

		ip, err := sut.AddRandomIP(hostName)

		Expect(ip).To(BeNil())
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

		ip, err := sut.AddRandomIP(hostName)

		Expect(ip).To(BeNil())
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

		ip, err := sut.AddRandomIP(hostName)

		Expect(ip).To(BeNil())
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

		ip, err := sut.AddRandomIP(hostName)

		Expect(ip).To(BeNil())
		Expect(err).To(MatchError(expectedErr))
	})

})
