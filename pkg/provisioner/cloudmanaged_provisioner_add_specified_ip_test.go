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
var _ = Describe("AssignSpecifiedIP", func() {
	BeforeEach(func() {
		initCloudManagedProvisionerMock()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should add a specified ip", func() {
		Cloud.EXPECT().
			AddSpecifiedIP(defaultAddedIP, hostName).
			Return(nil)

		OpenShiftProvisioner.EXPECT().
			AddSpecifiedIP(ctx, defaultAddedIP, hostName).
			Return(nil)

		err := cloudManagedProvisioner.AddSpecifiedIP(ctx, defaultAddedIP, hostName)

		Expect(err).To(BeNil())
	})

	It("should pass the error from the cloudprovider to the caller", func() {
		expectedErr := fmt.Errorf(
			"there are no or too much IP address assigned to the eni '%v': [%v]",
			"eni-1",
			"",
		)

		Cloud.EXPECT().
			AddSpecifiedIP(defaultAddedIP, hostName).
			Return(expectedErr)

		err := cloudManagedProvisioner.AddSpecifiedIP(ctx, defaultAddedIP, hostName)

		Expect(err).To(MatchError(expectedErr))
	})

	It("should rollback the cloudprovider when OCP fails", func() {
		expectedErr := errors.New("openshift provisioning failed")

		Cloud.EXPECT().
			AddSpecifiedIP(defaultAddedIP, hostName).
			Return(nil)

		OpenShiftProvisioner.EXPECT().
			AddSpecifiedIP(ctx, defaultAddedIP, hostName).
			Return(expectedErr)

		Cloud.EXPECT().
			RemoveIP(defaultAddedIP, hostName).
			Return(nil)

		err := cloudManagedProvisioner.AddSpecifiedIP(ctx, defaultAddedIP, hostName)

		Expect(err).To(MatchError(expectedErr))
	})

	It("should return error when rollback failed", func() {
		openShiftErr := errors.New("could not add ip to OpenShiftProvisioner")
		cloudProviderErr := errors.New("could not remove ip")

		expectedErr := fmt.Errorf(
			"error while rolling back adding ip '%v' to host '%v': could not add ip to OpenShiftProvisioner",
			defaultAddedIP.String(),
			hostName,
		)

		Cloud.EXPECT().
			AddSpecifiedIP(defaultAddedIP, hostName).
			Return(nil)

		OpenShiftProvisioner.EXPECT().
			AddSpecifiedIP(ctx, defaultAddedIP, hostName).
			Return(openShiftErr)

		Cloud.EXPECT().
			RemoveIP(defaultAddedIP, hostName).
			Return(cloudProviderErr)

		err := cloudManagedProvisioner.AddSpecifiedIP(ctx, defaultAddedIP, hostName)

		Expect(err).To(MatchError(expectedErr))
	})
})
