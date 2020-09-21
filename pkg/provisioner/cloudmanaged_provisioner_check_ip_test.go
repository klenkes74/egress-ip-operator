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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//goland:noinspection GoNilness
var _ = Describe("CheckIP", func() {
	BeforeEach(func() {
		initCloudManagedProvisionerMock()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should return no error when IP is assigned to host in cloudprovider and OpenShiftProvisioner", func() {
		Cloud.EXPECT().
			CheckIP(defaultAddedIP, hostName).
			Return(nil)

		OpenShiftProvisioner.EXPECT().
			CheckIP(ctx, defaultAddedIP, hostName).
			Return(nil)

		err := cloudManagedProvisioner.CheckIP(ctx, defaultAddedIP, hostName)

		Expect(err).To(BeNil())
	})

	It("should return an error when IP is assigned to OpenShiftProvisioner but not in the cloud", func() {
		expectedErr := errors.New("the specified ip is not assigned to the specified host in cloud")

		Cloud.EXPECT().
			CheckIP(defaultAddedIP, hostName).
			Return(expectedErr)

		err := cloudManagedProvisioner.CheckIP(ctx, defaultAddedIP, hostName)

		Expect(err).To(MatchError(expectedErr))
	})

	It("should return an error when IP is assigned to the cloud but not in OpenShiftProvisioner", func() {
		expectedErr := errors.New("the specified ip is not assigned to the specified host in openshift")

		Cloud.EXPECT().
			CheckIP(defaultAddedIP, hostName).
			Return(nil)

		OpenShiftProvisioner.EXPECT().
			CheckIP(ctx, defaultAddedIP, hostName).
			Return(expectedErr)

		err := cloudManagedProvisioner.CheckIP(ctx, defaultAddedIP, hostName)

		Expect(err).To(MatchError(expectedErr))
	})
})
