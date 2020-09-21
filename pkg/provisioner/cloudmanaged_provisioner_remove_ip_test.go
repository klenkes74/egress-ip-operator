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
var _ = Describe("RemoveIP", func() {
	BeforeEach(func() {
		initCloudManagedProvisionerMock()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should remove IP from host", func() {
		Cloud.EXPECT().
			RemoveIP(defaultAddedIP, hostName).
			Return(nil)

		OpenShiftProvisioner.EXPECT().
			RemoveIP(ctx, defaultAddedIP, hostName).
			Return(nil)

		err := cloudManagedProvisioner.RemoveIP(ctx, defaultAddedIP, hostName)

		Expect(err).To(BeNil())
	})

	It("should pass the error from the cloudprovider to the caller", func() {
		expectedErr := fmt.Errorf(
			"there are no or too much IP address assigned to the eni '%v': [%v]",
			"eni-1",
			"",
		)

		Cloud.EXPECT().
			RemoveIP(defaultAddedIP, hostName).
			Return(expectedErr)

		err := cloudManagedProvisioner.RemoveIP(ctx, defaultAddedIP, hostName)

		Expect(err).To(MatchError(expectedErr))
	})

	It("should rollback the cloudprovider when OCP fails", func() {
		openShiftErr := errors.New("openshift provisioning failed")
		expectedErr := fmt.Errorf(
			"error while removing IP '%v' from host. IP is still valid: openshift provisioning failed",
			defaultAddedIP,
		)

		Cloud.EXPECT().
			RemoveIP(defaultAddedIP, hostName).
			Return(nil)

		OpenShiftProvisioner.EXPECT().
			RemoveIP(ctx, defaultAddedIP, hostName).
			Return(openShiftErr)

		Cloud.EXPECT().
			AddSpecifiedIP(defaultAddedIP, hostName).
			Return(nil)

		err := cloudManagedProvisioner.RemoveIP(ctx, defaultAddedIP, hostName)

		Expect(err).To(MatchError(expectedErr))
	})

	It("should return error when rollback failed", func() {
		openShiftErr := errors.New("openshift provisioning failed")
		cloudProviderErr := errors.New("could not rollback in cloud")
		expectedErr := fmt.Errorf(
			"error while removing IP '%v' from OpenShiftProvisioner. Re-adding it to the cloudprovider failed: could not rollback in cloud",
			defaultAddedIP,
		)

		Cloud.EXPECT().
			RemoveIP(defaultAddedIP, hostName).
			Return(nil)

		OpenShiftProvisioner.EXPECT().
			RemoveIP(ctx, defaultAddedIP, hostName).
			Return(openShiftErr)

		Cloud.EXPECT().
			AddSpecifiedIP(defaultAddedIP, hostName).
			Return(cloudProviderErr)

		err := cloudManagedProvisioner.RemoveIP(ctx, defaultAddedIP, hostName)

		Expect(err).To(MatchError(expectedErr))
	})
})
