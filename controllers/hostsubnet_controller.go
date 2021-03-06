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

package controllers

import (
	"github.com/klenkes74/egress-ip-operator/pkg/metrics"
	"github.com/klenkes74/egress-ip-operator/pkg/openshift"
	"github.com/klenkes74/egress-ip-operator/pkg/provisioner"

	"github.com/go-logr/logr"
	netv1 "github.com/openshift/api/network/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EgressIPReconciler reconciles a EgressIP object
type HostSubnetReconciler struct {
	client.Client
	Log         logr.Logger
	Scheme      *runtime.Scheme
	Provisioner *provisioner.EgressIPProvisioner
	Alarm       *metrics.AlarmStore
}

// +kubebuilder:rbac:groups=network.openshift.io,resources=hostsubnets,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=network.openshift.io,resources=hostsubnets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=egressip.kaiserpfalz-edv.de,resources=egressipfailuredomains/status,verbs=get;update;patch;create;delete
// +kubebuilder:rbac:groups=egressip.kaiserpfalz-edv.de,resources=egressips/status,verbs=get;update;patch;create;delete

func (r *HostSubnetReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	return openshift.ManageHostSubnet(req, r.Client, r.Log)
}

func (r *HostSubnetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&netv1.HostSubnet{}).
		Complete(r)
}
