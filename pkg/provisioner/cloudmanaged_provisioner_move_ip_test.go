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
var _ = Describe("MoveIP", func() {
	BeforeEach(func() {
		initCloudManagedProvisionerMock()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should move ip to another host", func() {
		Cloud.EXPECT().
			MoveIP(defaultAddedIP, hostName, targetName).
			Return(nil)

		OpenShiftProvisioner.EXPECT().
			MoveIP(ctx, defaultAddedIP, hostName, targetName).
			Return(nil)

		err := cloudManagedProvisioner.MoveIP(ctx, defaultAddedIP, hostName, targetName)

		Expect(err).To(BeNil())
	})

	It("should pass the error from the cloudprovider to the caller", func() {
		expectedErr := fmt.Errorf(
			"there are no or too much IP address assigned to the eni '%v': [%v]",
			"eni-1",
			"",
		)

		Cloud.EXPECT().
			MoveIP(defaultAddedIP, hostName, targetName).
			Return(expectedErr)

		err := cloudManagedProvisioner.MoveIP(ctx, defaultAddedIP, hostName, targetName)

		Expect(err).To(MatchError(expectedErr))
	})

	It("should rollback the cloudprovider when OCP fails", func() {
		openShiftErr := errors.New("openshift provisioning failed")
		expectedErr := fmt.Errorf(
			"error while moving IP '%v' from '%v' to '%v'. Change reverted: openshift provisioning failed",
			defaultAddedIP,
			hostName,
			targetName,
		)

		Cloud.EXPECT().
			MoveIP(defaultAddedIP, hostName, targetName).
			Return(nil)

		OpenShiftProvisioner.EXPECT().
			MoveIP(ctx, defaultAddedIP, hostName, targetName).
			Return(openShiftErr)

		Cloud.EXPECT().
			MoveIP(defaultAddedIP, targetName, hostName).
			Return(nil)

		err := cloudManagedProvisioner.MoveIP(ctx, defaultAddedIP, hostName, targetName)

		Expect(err).To(MatchError(expectedErr))
	})

	It("should return error when rollback failed", func() {
		openShiftErr := errors.New("openshift provisioning failed")
		cloudProviderErr := errors.New("could not rollback in cloud")
		expectedErr := fmt.Errorf(
			"error while moving IP '%v' from '%v' to '%v': could not rollback in cloud",
			defaultAddedIP,
			hostName,
			targetName,
		)

		Cloud.EXPECT().
			MoveIP(defaultAddedIP, hostName, targetName).
			Return(nil)

		OpenShiftProvisioner.EXPECT().
			MoveIP(ctx, defaultAddedIP, hostName, targetName).
			Return(openShiftErr)

		Cloud.EXPECT().
			MoveIP(defaultAddedIP, targetName, hostName).
			Return(cloudProviderErr)

		err := cloudManagedProvisioner.MoveIP(ctx, defaultAddedIP, hostName, targetName)

		Expect(err).To(MatchError(expectedErr))
	})
})
