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

// cloudprovider is the abstraction layer for communicating with the IP management of the cloud OCP is running in.
package cloudprovider

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/go-logr/logr"
	"github.com/klenkes74/egress-ip-operator/pkg/cloudprovider/aws_provider"
	"net"
	"os"
	"strconv"
)

type CloudProvider interface {
	// AddRandomIP adds a random IP to the specified host.
	// It will return the IP or the error.
	AddRandomIP(hostName string) (*net.IP, error)
	// AddSpecifiedIP adds a predefined IP to the specified host.
	// It will return an error or nil.
	AddSpecifiedIP(ip *net.IP, hostName string) error
	// CheckIP will check if the specified IP is assigned on the specified host.
	// it will return an error or nil.
	CheckIP(ip *net.IP, hostName string) error
	// MoveIP will move the specified IP from oldHost to newHost.
	// It will return an error or nil.
	MoveIP(ip *net.IP, oldHostName string, newHostName string) error
	// RemoveIP will remove the given IP from the specified host.
	// It will return an error or nil.
	RemoveIP(ip *net.IP, hostName string) error
}

var _ CloudProvider = &aws_provider.AwsCloudProvider{}

const (
	DefaultFailureRegion     = "Kunchom"
	DefaultMaxIPsPerInstance = 8
)

var (
	FailureRegion     string
	MaxIPsPerInstance int
)

func init() {
	var found bool
	var err error

	FailureRegion, found = os.LookupEnv("CLOUD_FAILURE_REGION")
	if !found {
		FailureRegion = DefaultFailureRegion
	}
	maxIPs, found := os.LookupEnv("CLOUD_MAX_IPS_PER_INSTANCE")
	if found {
		MaxIPsPerInstance, err = strconv.Atoi(maxIPs)
	}

	if !found || err != nil {
		MaxIPsPerInstance = DefaultMaxIPsPerInstance
	}
}

// NewCloudProvider initializes the cloudprovider configured for this system.
func NewCloudProvider(cloudProviderType string, logger logr.Logger) (*CloudProvider, error) {
	var result CloudProvider

	switch cloudProviderType {
	case "aws":
		awsSession := session.Must(session.NewSession())
		client := ec2.New(awsSession, aws.NewConfig().WithRegion(FailureRegion))

		awsProvider := aws_provider.AwsDirectCallsProd{
			Session: awsSession,
			Client:  client,
		}

		provider := &aws_provider.AwsCloudProvider{
			FailureRegion:     FailureRegion,
			MaxIPsPerInstance: MaxIPsPerInstance,
			Client:            &awsProvider,
			Log:               logger.WithName("aws"),
		}
		result = CloudProvider(provider)
	default:
		return nil, fmt.Errorf("cloudprovider type '%v' is not defined - please use one of: 'aws'", cloudProviderType)
	}

	return &result, nil
}
