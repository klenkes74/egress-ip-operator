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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
	"strconv"
)

var _ = Describe("CheckIP", func() {
	BeforeEach(func() {
		initMock()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should be fine when the IP is assigned to the specified host", func() {
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

		err := sut.CheckIP(secondaryIPs[1], hostName)

		Expect(err).To(BeNil())
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
