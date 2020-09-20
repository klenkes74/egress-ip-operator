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
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
	"strconv"
)

var _ = Describe("RemoveIP", func() {
	BeforeEach(func() {
		initMock()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should remove the IP from the instance", func() {
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

		awsDirect.
			EXPECT().UnassignPrivateIpAddresses(&ec2.UnassignPrivateIpAddressesInput{
			NetworkInterfaceId: aws.String(networkInterfaceId),
			PrivateIpAddresses: aws.StringSlice([]string{secondaryIPs[1].String()}),
		}).
			Return(
				&ec2.UnassignPrivateIpAddressesOutput{},
				nil,
			)

		err := sut.RemoveIP(secondaryIPs[1], hostName)

		Expect(err).To(BeNil())
	})

	It("should throw an error when removing the ip fails", func() {
		expectedErr := errors.New("could not remove ip from interface")

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

		awsDirect.
			EXPECT().UnassignPrivateIpAddresses(&ec2.UnassignPrivateIpAddressesInput{
			NetworkInterfaceId: aws.String(networkInterfaceId),
			PrivateIpAddresses: aws.StringSlice([]string{secondaryIPs[1].String()}),
		}).
			Return(
				nil,
				expectedErr,
			)

		err := sut.RemoveIP(secondaryIPs[1], hostName)

		Expect(err).To(MatchError(expectedErr))
	})

	It("should ignore missing IPs when removing IP from an instance", func() {
		awsDirect.
			EXPECT().DescribeInstances(createDescribeInstancesInput(hostName)).
			Return(
				createDescribeInstancesOutput(hostName, hostId, networkInterfaceId, mainIP, []*net.IP{}),
				nil,
			)

		ip := net.ParseIP("9.9.9.9")
		err := sut.RemoveIP(&ip, hostName)

		Expect(err).To(BeNil())
	})

	It("should ignore missing network interface when removing IP from an instance", func() {
		awsDirect.
			EXPECT().DescribeInstances(createDescribeInstancesInput(hostName)).
			Return(
				createDescribeInstancesOutput(hostName, hostId, "", mainIP, []*net.IP{}),
				nil,
			)

		ip := net.ParseIP("9.9.9.9")
		err := sut.RemoveIP(&ip, hostName)

		Expect(err).To(BeNil())
	})
})
