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

package provisioner_test

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//goland:noinspection GoNilness
var _ = Describe("FindHostForNewIP", func() {
	BeforeEach(func() {
		initCloudManagedProvisionerMock()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should return the default host", func() {
		OpenShiftProvisioner.EXPECT().
			FindHostForNewIP(ctx, defaultFailureDomain).
			Return(hostName, *defaultAddedIP, nil)

		host, ip, err := cloudManagedProvisioner.FindHostForNewIP(ctx, defaultFailureDomain)

		Expect(host).To(Equal(hostName))
		Expect(ip).To(Equal(*defaultAddedIP))
		Expect(err).To(BeNil())
	})

	It("should return the failure from the OpenShift static provisioner", func() {
		openShiftErr := errors.New("openshift could not find host")
		expectedErr := fmt.Errorf("openshift could not find host")

		OpenShiftProvisioner.EXPECT().
			FindHostForNewIP(ctx, defaultFailureDomain).
			Return("", nil, openShiftErr)

		host, ip, err := cloudManagedProvisioner.FindHostForNewIP(ctx, defaultFailureDomain)

		Expect(host).To(Equal(""))
		Expect(ip).To(BeNil())
		Expect(err).To(MatchError(expectedErr))
	})
})
