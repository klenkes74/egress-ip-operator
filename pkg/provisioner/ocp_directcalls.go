//go:generate go run github.com/golang/mock/mockgen -package provisioner_test -destination ./mock_ocp_directcalls_test.go github.com/klenkes74/egress-ip-operator/pkg/provisioner OCPDirectCalls

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

package provisioner

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AwsDirectCalls is the interface for accessing AWS services. It is the final interface to be able to mock the AWS
// calls during testing. Basically it is a delegate for the client.Reader and client.Writer interface of the k8s client.
type OCPDirectCalls interface {
	Create(ctx context.Context, obj runtime.Object, opts ...client.CreateOption) error
	Delete(ctx context.Context, obj runtime.Object, opts ...client.DeleteOption) error
	Update(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error
	Patch(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error
	DeleteAllOf(ctx context.Context, obj runtime.Object, opts ...client.DeleteAllOfOption) error
	Get(ctx context.Context, key client.ObjectKey, obj runtime.Object) error
	List(ctx context.Context, list runtime.Object, opts ...client.ListOption) error
}

var (
	_ client.Reader  = &OCPDirectCallsProd{}
	_ client.Writer  = &OCPDirectCallsProd{}
	_ OCPDirectCalls = &OCPDirectCallsProd{}
)

// AwsDirectCallsProd is the working implementation of the AwsDirectCalls interface.
type OCPDirectCallsProd struct {
	Client client.Client
}

func (a *OCPDirectCallsProd) Create(ctx context.Context, obj runtime.Object, opts ...client.CreateOption) error {
	return a.Client.Create(ctx, obj, opts...)
}

func (a *OCPDirectCallsProd) Delete(ctx context.Context, obj runtime.Object, opts ...client.DeleteOption) error {
	return a.Client.Delete(ctx, obj, opts...)
}

func (a *OCPDirectCallsProd) Update(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
	return a.Client.Update(ctx, obj, opts...)
}

func (a *OCPDirectCallsProd) Patch(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error {
	return a.Client.Patch(ctx, obj, patch, opts...)
}

func (a *OCPDirectCallsProd) DeleteAllOf(ctx context.Context, obj runtime.Object, opts ...client.DeleteAllOfOption) error {
	return a.Client.DeleteAllOf(ctx, obj, opts...)
}

func (a *OCPDirectCallsProd) Get(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
	return a.Client.Get(ctx, key, obj)
}

func (a *OCPDirectCallsProd) List(ctx context.Context, list runtime.Object, opts ...client.ListOption) error {
	return a.Client.List(ctx, list, opts...)
}
