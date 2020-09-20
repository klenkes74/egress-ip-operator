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
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
	"strconv"
)

var _ = Describe("MoveIP", func() {
	BeforeEach(func() {
		initMock()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should move the IP from old instance to new instance", func() {
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

		targetHostName := "target"
		targetMainIP := net.ParseIP("10.2.1.33")
		targetSecondaryIPs := make([]*net.IP, 2)
		for i := 0; i < 2; i++ {
			s := strconv.Itoa(50 + i)
			ip := net.ParseIP("10.243.1." + s)
			targetSecondaryIPs[i] = &ip
		}
		awsDirect.
			EXPECT().DescribeInstances(createDescribeInstancesInput(targetHostName)).
			Return(
				createDescribeInstancesOutput(targetHostName, "vm-2", "eni-2", &targetMainIP, targetSecondaryIPs),
				nil,
			)

		allowReassignement := true
		awsDirect.
			EXPECT().AssignPrivateIpAddresses(&ec2.AssignPrivateIpAddressesInput{
			AllowReassignment:  &allowReassignement,
			NetworkInterfaceId: aws.String("eni-2"),
			PrivateIpAddresses: aws.StringSlice([]string{secondaryIPs[1].String()}),
		}).
			Return(
				&ec2.AssignPrivateIpAddressesOutput{
					AssignedPrivateIpAddresses: []*ec2.AssignedPrivateIpAddress{
						{
							PrivateIpAddress: aws.String(secondaryIPs[1].String()),
						},
					},
					NetworkInterfaceId: aws.String("eni-2"),
				},
				nil,
			)

		err := sut.MoveIP(secondaryIPs[1], hostName, targetHostName)

		Expect(err).To(BeNil())
	})

	It("should throw an error when the moving of the IP failed", func() {
		expectedErr := errors.New("can not move IP to target host")

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

		targetHostName := "target"
		targetMainIP := net.ParseIP("10.2.1.33")
		targetSecondaryIPs := make([]*net.IP, 2)
		for i := 0; i < 2; i++ {
			s := strconv.Itoa(50 + i)
			ip := net.ParseIP("10.243.1." + s)
			targetSecondaryIPs[i] = &ip
		}
		awsDirect.
			EXPECT().DescribeInstances(createDescribeInstancesInput(targetHostName)).
			Return(
				createDescribeInstancesOutput(targetHostName, "vm-2", "eni-2", &targetMainIP, targetSecondaryIPs),
				nil,
			)

		allowReassignement := true
		awsDirect.
			EXPECT().AssignPrivateIpAddresses(&ec2.AssignPrivateIpAddressesInput{
			AllowReassignment:  &allowReassignement,
			NetworkInterfaceId: aws.String("eni-2"),
			PrivateIpAddresses: aws.StringSlice([]string{secondaryIPs[1].String()}),
		}).
			Return(
				nil,
				errors.New("can not move IP to target host"),
			)

		err := sut.MoveIP(secondaryIPs[1], hostName, targetHostName)

		Expect(err).To(MatchError(expectedErr))
	})

	It("should throw an error when the IP is not assigned to the specified host", func() {
		unassignedIP := net.ParseIP("10.1.2.3")
		expectedErr := fmt.Errorf(
			"ip '%v' is not assigned to instance '%v'",
			unassignedIP.String(), hostId,
		)

		awsDirect.
			EXPECT().DescribeInstances(createDescribeInstancesInput(hostName)).
			Return(
				createDescribeInstancesOutput(hostName, hostId, networkInterfaceId, mainIP, []*net.IP{}),
				nil,
			)

		err := sut.CheckIP(&unassignedIP, hostName)

		Expect(err).To(MatchError(expectedErr))
	})
})
