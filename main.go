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

package main

import (
	"flag"
	"github.com/klenkes74/egress-ip-operator/pkg/metrics"
	"github.com/klenkes74/egress-ip-operator/pkg/provisioner"
	"os"

	netv1 "github.com/openshift/api/network/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	egressipv1alpha1 "github.com/klenkes74/egress-ip-operator/api/v1alpha1"
	"github.com/klenkes74/egress-ip-operator/controllers"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(netv1.AddToScheme(scheme))

	utilruntime.Must(egressipv1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "416133d7.kaiserpfalz-edv.de",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	alarm := metrics.NewAlarmStore(ctrl.Log.WithName("metrics-based-alarmstore"))

	egressIPProvisioner, err := provisioner.NewEgressIPProvisioner(ctrl.Log)
	if err != nil {
		setupLog.Error(err, "unable to create egress ip provisioner")
		os.Exit(1)
	}

	if err = (&controllers.EgressIPReconciler{
		Client:      mgr.GetClient(),
		Log:         ctrl.Log.WithName("controllers").WithName("egressip-controller"),
		Scheme:      mgr.GetScheme(),
		Provisioner: egressIPProvisioner,
		Alarm:       alarm,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "egressip-controller")
		os.Exit(1)
	}
	if err = (&controllers.EgressIPFailureDomainReconciler{
		Client:      mgr.GetClient(),
		Log:         ctrl.Log.WithName("controllers").WithName("egressip-failuredomain-controller"),
		Scheme:      mgr.GetScheme(),
		Provisioner: egressIPProvisioner,
		Alarm:       alarm,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "egressip-failuredomain-controller")
		os.Exit(1)
	}
	if err = (&controllers.HostSubnetReconciler{
		Client:      mgr.GetClient(),
		Log:         ctrl.Log.WithName("controllers").WithName("HostSubnet"),
		Scheme:      mgr.GetScheme(),
		Provisioner: egressIPProvisioner,
		Alarm:       alarm,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HostSubnet")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
