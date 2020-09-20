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

package openshift

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/klenkes74/egress-ip-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ManageEgressIPFailureDomain(req ctrl.Request, client client.Client, baseLogger logr.Logger) (ctrl.Result, error) {
	ctx := context.Background()
	log := baseLogger.WithValues("egressipfailuredomain", req.NamespacedName)

	instance := &v1alpha1.EgressIPFailureDomain{}
	err := client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("egressIPFailureDomain not found - the request will not be re-queued")

			return ctrl.Result{
				Requeue: false,
			}, err
		}

		log.Info("egressIPFailureDomain could not be loaded - the request will be re-queued in 30 seconds")
		return ctrl.Result{
			RequeueAfter: 30,
		}, err
	}

	// TODO 2020-09-20 rlichti replace this logging by real working code
	log.Info("working on", "egressIPFailureDomain", instance)

	return ctrl.Result{}, nil
}
