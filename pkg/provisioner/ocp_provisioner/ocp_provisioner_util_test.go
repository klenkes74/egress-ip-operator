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

package ocp_provisioner_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"net"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var log = zap.New(zap.UseDevMode(true)).WithName("ocp_static_provisioner_test")

var (
	mockCtrl *gomock.Controller

	mainIP             *net.IP
	ip                 *net.IP
	hostName           string
	hostId             string
	networkInterfaceId string
	maxIPsPerInstance  int
)

func init() {
}

func initMock() {
	mockCtrl = gomock.NewController(GinkgoT())
}
