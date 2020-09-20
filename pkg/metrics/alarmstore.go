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

package metrics

import (
	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	"net"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	"time"
)

// AlarmStore -- the store for keeping alarms of the aws-egress-ip-operator
type AlarmStore interface {
	// Adds a failed namespace to the alarm store
	AddAlarm(namespace string, ips []*net.IP)

	// Removes a recovered namespace from the alarm store
	RemoveAlarm(namespace string)

	RemoveAlarmForIP(namespace string, ip *net.IP)

	// Retrieves all failed namespaces from the alarm store
	GetFailed() map[string]*FailedEgressIP
}

// ensures that the PrometheusLinkedAlarmStore is a valid AlarmStore
var _ AlarmStore = &PrometheusLinkedAlarmStore{}

// PrometheusLinkedAlarmStore -- a simple in memory implementation of the AlarmStore
type PrometheusLinkedAlarmStore struct {
	failures map[string]*FailedEgressIP
	counter  prometheus.GaugeVec

	Log logr.Logger
}

var singletonAlarmStore *PrometheusLinkedAlarmStore

// NewAlarmStore -- creates the default implementation of the alarm store
func NewAlarmStore(logger logr.Logger) *AlarmStore {
	if singletonAlarmStore == nil {
		createAlarmStore(logger)
	}

	result := AlarmStore(singletonAlarmStore)

	return &result
}

func createAlarmStore(logger logr.Logger) {
	counter := *prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "egress_ip",
			Name:      "handling_failures",
			Help:      "Failures while handling egress-ips",
		},
		[]string{"namespace"},
	)
	err := metrics.Registry.Register(counter)
	if err != nil {
		logger.Error(err, "Can't register the new gauge")
	}

	singletonAlarmStore = &PrometheusLinkedAlarmStore{
		failures: make(map[string]*FailedEgressIP),
		counter:  counter,
		Log:      logger,
	}
}

// AddAlarm -- Adds a failed namespace to the alarm store
func (s PrometheusLinkedAlarmStore) AddAlarm(namespace string, ips []*net.IP) {
	if s.failures[namespace] == nil {
		timeStamp := time.Now()

		alarm := FailedEgressIP{
			Namespace:       namespace,
			FailedIPs:       ips,
			FirstOccurrence: timeStamp,
			LastOccurrence:  timeStamp,
			Counter:         float64(1),
		}

		s.failures[namespace] = &alarm

	} else {
		s.failures[namespace].Counter = s.failures[namespace].Counter + 1
		s.failures[namespace].LastOccurrence = time.Now()
		s.failures[namespace].FailedIPs = ips
	}

	s.counter.WithLabelValues(namespace).Set(s.failures[namespace].Counter)
}

// RemoveAlarm -- Removes a recovered namespace from the alarm store
func (s PrometheusLinkedAlarmStore) RemoveAlarm(namespace string) {
	if s.failures[namespace] != nil {
		s.counter.WithLabelValues(namespace).Set(0)

		delete(s.failures, namespace)
	}
}

// RemoveAlarmForIP -- Removes the alarm for a single IP. If there are still IPs in alarm, keep the alarm, if that has
// been the last IP, remove the alarm.
func (s PrometheusLinkedAlarmStore) RemoveAlarmForIP(namespace string, ip *net.IP) {
	if s.failures[namespace] != nil {
		newFailures := make([]*net.IP, 0)

		for _, oldIP := range s.failures[namespace].FailedIPs {
			if !reflect.DeepEqual(oldIP, ip) {
				newFailures = append(newFailures, oldIP)
			}
		}

		if len(newFailures) > 0 {
			s.failures[namespace].FailedIPs = newFailures
		} else {
			s.RemoveAlarm(namespace)
		}
	}
}

// GetFailed -- Retrieves all failed namespaces from the alarm store
func (s PrometheusLinkedAlarmStore) GetFailed() map[string]*FailedEgressIP {
	return s.failures
}

// FailedEgressIP - This is the data for the failure.
type FailedEgressIP struct {
	Namespace       string    // The failed namespace
	FailedIPs       []*net.IP // The failed IPs
	FirstOccurrence time.Time // First occurrence of this failure
	LastOccurrence  time.Time // Last occurrence of this failure
	Counter         float64   // Failure counter
}
