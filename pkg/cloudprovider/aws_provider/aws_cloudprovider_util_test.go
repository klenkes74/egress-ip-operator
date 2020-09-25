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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/klenkes74/egress-ip-operator/pkg/cloudprovider/aws_provider"
	. "github.com/onsi/ginkgo"
	"net"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var log = zap.New(zap.UseDevMode(true)).WithName("cloudprovider_test")

var (
	mockCtrl  *gomock.Controller
	awsDirect *MockAwsDirectCalls
	sut       *aws_provider.AwsCloudProvider

	mainIP             *net.IP
	ip                 *net.IP
	hostName           string
	hostId             string
	networkInterfaceId string
	maxIPsPerInstance  int
)

func init() {
	tempIP := net.ParseIP("10.0.1.8")
	mainIP = &tempIP
	tempIP2 := net.ParseIP("10.0.1.42")
	ip = &tempIP2
	hostName = "host"
	hostId = "vm-1"
	networkInterfaceId = "eni-1"
	maxIPsPerInstance = 4
}

func initMock() {
	mockCtrl = gomock.NewController(GinkgoT())
	awsDirect = NewMockAwsDirectCalls(mockCtrl)
	sut = &aws_provider.AwsCloudProvider{
		FailureRegion:     "region",
		MaxIPsPerInstance: maxIPsPerInstance,
		Client:            awsDirect,
		Log:               log,
	}
}

func createDescribeInstancesInput(hostName string) *ec2.DescribeInstancesInput {
	return &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("private-dns-name"),
				Values: aws.StringSlice([]string{hostName}),
			},
		},
	}
}

func createDescribeInstancesOutput(hostName, hostId, networkInterfaceId string, mainIP *net.IP, secondaryIPs []*net.IP) *ec2.DescribeInstancesOutput {
	log.Info(
		"create instance",
		"host-name", hostName,
		"instance-id", hostId,
		"network-interface", networkInterfaceId,
		"main-ip", mainIP.String(),
		"secondary-ips", secondaryIPs,
	)

	primary := true
	privateIPAddresses := make([]*ec2.InstancePrivateIpAddress, 1+len(secondaryIPs))
	privateIPAddresses[0] = &ec2.InstancePrivateIpAddress{
		PrivateDnsName:   aws.String(hostName),
		Primary:          &primary,
		PrivateIpAddress: aws.String(mainIP.String()),
	}

	if len(secondaryIPs) >= 1 {
		secondary := false
		for i, ip := range secondaryIPs {
			privateIPAddresses[i+1] = &ec2.InstancePrivateIpAddress{
				PrivateDnsName:   aws.String(hostName),
				Primary:          &secondary,
				PrivateIpAddress: aws.String(ip.String()),
			}
		}
	}

	var networkInterfaces []*ec2.InstanceNetworkInterface

	if networkInterfaceId != "" {
		networkInterfaces = []*ec2.InstanceNetworkInterface{
			{
				Attachment: &ec2.InstanceNetworkInterfaceAttachment{
					AttachmentId: aws.String(networkInterfaceId),
				},
				NetworkInterfaceId: aws.String(networkInterfaceId),
				PrivateDnsName:     aws.String(hostName),
				PrivateIpAddress:   aws.String(mainIP.String()),
				PrivateIpAddresses: privateIPAddresses,
			},
		}
	}

	return &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:        aws.String(hostId),
						PrivateIpAddress:  aws.String(mainIP.String()),
						PrivateDnsName:    aws.String(hostName),
						NetworkInterfaces: networkInterfaces,
					},
				},
				ReservationId: nil,
			},
		},
	}
}
