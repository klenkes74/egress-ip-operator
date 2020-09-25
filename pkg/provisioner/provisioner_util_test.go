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
	"context"
	"github.com/golang/mock/gomock"
	"github.com/klenkes74/egress-ip-operator/pkg/cloudprovider"
	"github.com/klenkes74/egress-ip-operator/pkg/provisioner"
	. "github.com/onsi/ginkgo"
	"net"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var log = zap.New(zap.UseDevMode(true)).WithName("cloudmanaged_provisioner_test")

var (
	mockCtrl             *gomock.Controller
	Cloud                *cloudprovider.MockCloudProvider
	OpenShiftProvisioner *provisioner.MockEgressIPProvisioner

	cloudManagedProvisioner *provisioner.CloudManagedEgressIPProvisioner
	OcpStaticProvisioner    *provisioner.OcpStaticEgressIPProvisioner
	OcpDynamicProvisioner   *provisioner.OcpDynamicEgressIPProvisioner
	OpenShiftClient         *MockOCPDirectCalls

	ctx context.Context

	hostName     string
	mainIP       *net.IP
	targetName   string
	targetMainIP *net.IP

	defaultFailureDomain string
	defaultCIDR          *net.IPNet
	defaultAddedIP       *net.IP
)

func init() {
	ctx = context.TODO()

	hostName = "ip-10-1-1-8.kunchom.compute.internal"
	mainIP = parseIP("10.1.1.18")

	targetName = "ip-10-1-1-72.kunchom.compute.internal"
	targetMainIP = parseIP("10.1.1.72")

	defaultFailureDomain = "kunchom"
	_, defaultCIDR, _ = net.ParseCIDR("10.1.1.0/24")
	defaultAddedIP = parseIP("10.1.1.183")
}

func parseIP(ip string) *net.IP {
	result := net.ParseIP(ip)
	return &result
}

func initCloudManagedProvisionerMock() {
	mockCtrl = gomock.NewController(GinkgoT())

	Cloud = cloudprovider.NewMockCloudProvider(mockCtrl)
	OpenShiftProvisioner = provisioner.NewMockEgressIPProvisioner(mockCtrl)

	cloudManagedProvisioner = &provisioner.CloudManagedEgressIPProvisioner{
		Log:       log.WithName("cloudManagedProvisioner"),
		Cloud:     Cloud,
		OpenShift: OpenShiftProvisioner,
	}
}

func initOcpStaticProvisionerMock() {
	mockCtrl = gomock.NewController(GinkgoT())

	OpenShiftClient = NewMockOCPDirectCalls(mockCtrl)

	OcpStaticProvisioner = &provisioner.OcpStaticEgressIPProvisioner{
		Client: OpenShiftClient,
		Log:    log.WithName("ocpStaticProvisioner"),
	}
}

func initOcpDynamicProvisionerMock() {
	mockCtrl = gomock.NewController(GinkgoT())

	OpenShiftClient = NewMockOCPDirectCalls(mockCtrl)

	OcpDynamicProvisioner = &provisioner.OcpDynamicEgressIPProvisioner{
		Client: OpenShiftClient,
		Log:    log.WithName("cloudManagedProvisioner"),
	}
}
