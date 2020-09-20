package metrics_test

import (
	"github.com/klenkes74/egress-ip-operator/pkg/metrics"
	"net"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"
	"time"
)

var log = zap.New(zap.UseDevMode(true)).WithName("metrics_test")

var expectedNamespace = "test"
var expectedIPs []*net.IP

func prepareStore() (metrics.AlarmStore, time.Time) {
	store := *metrics.NewAlarmStore(log.WithName("alarm-store"))

	expectedIPs = make([]*net.IP, 2)
	ip1 := net.ParseIP("1.1.1.1")
	expectedIPs[0] = &ip1
	ip2 := net.ParseIP("2.2.2.2")
	expectedIPs[1] = &ip2

	store.AddAlarm(expectedNamespace, expectedIPs)

	failed := store.GetFailed()
	return store, failed[expectedNamespace].FirstOccurrence
}

func TestAddingNamespaceToAlarmStore(t *testing.T) {
	store, _ := prepareStore()

	failed := store.GetFailed()

	for _, failedEgress := range failed {
		if failedEgress.Namespace != "test" {
			t.Errorf("Namespace name does not match! expected='%v', current='%v'",
				expectedNamespace,
				failedEgress.Namespace,
			)
		}

		if !reflect.DeepEqual(failedEgress.FailedIPs, expectedIPs) {
			t.Errorf("IPs don't match! expected='%v', current='%v'",
				expectedIPs,
				failedEgress.FailedIPs,
			)
		}
	}

	store.RemoveAlarm(expectedNamespace)
}

func TestRemovingNamespaceFromAlarmStore(t *testing.T) {
	store, _ := prepareStore()

	store.RemoveAlarm(expectedNamespace)

	if len(store.GetFailed()) > 0 {
		t.Error("There should be no failures in the AlarmStore!")
	}
}

func TestRemovingAlarmForSingleIP(t *testing.T) {
	store, _ := prepareStore()

	ip := net.ParseIP("2.2.2.2")
	store.RemoveAlarmForIP(expectedNamespace, &ip)

	if len(store.GetFailed()[expectedNamespace].FailedIPs) > 1 {
		t.Errorf("There should be only one IP listed as failed ip! expected=1, current=%v", len(store.GetFailed()[expectedNamespace].FailedIPs))
	}

	store.RemoveAlarm(expectedNamespace)
}

func TestRemovingAlarmForLastSingleIP(t *testing.T) {
	store, _ := prepareStore()

	ip := net.ParseIP("2.2.2.2")
	store.RemoveAlarmForIP(expectedNamespace, &ip)

	ip = net.ParseIP("1.1.1.1")
	store.RemoveAlarmForIP(expectedNamespace, &ip)

	if len(store.GetFailed()) > 0 {
		t.Errorf("There should be no failure any more! expected=0, current=%v", len(store.GetFailed()))
	}

	store.RemoveAlarm(expectedNamespace)
}

func TestGetFirstAndLastOccurrenceFromAlarm(t *testing.T) {
	store, timeStamp := prepareStore()

	ips := make([]*net.IP, 1)
	ip := net.ParseIP("3.3.3.3")
	ips[0] = &ip

	store.AddAlarm(expectedNamespace, ips)

	failed := store.GetFailed()

	if failed[expectedNamespace].FirstOccurrence != timeStamp {
		t.Errorf("First occurance of failure is not valid. expected='%v', current='%v'",
			timeStamp,
			failed[expectedNamespace].FirstOccurrence,
		)
	}

	if failed[expectedNamespace].LastOccurrence == timeStamp {
		t.Errorf("Last occurance of failure is not valid. expected= junger than '%v', current='%v'",
			timeStamp,
			failed[expectedNamespace].LastOccurrence,
		)
	}

	if !reflect.DeepEqual(failed[expectedNamespace].FailedIPs, ips) {
		t.Errorf("IPs do not match! expected='%v', current='%v'",
			failed[expectedNamespace].FailedIPs,
			ips,
		)
	}

	if failed[expectedNamespace].Counter != 2 {
		t.Errorf("The counter should be 2! expected='2', current='%v'", failed[expectedNamespace].Counter)
	}

	store.RemoveAlarm(expectedNamespace)
}

func TestAddTwoDifferentNamespaces(t *testing.T) {
	store, _ := prepareStore()

	ips := make([]*net.IP, 1)
	ip := net.ParseIP("3.3.3.3")
	ips[0] = &ip

	store.AddAlarm("other", ips)
	store.AddAlarm("other", ips)

	failed := store.GetFailed()

	if len(failed) != 2 {
		t.Errorf("Number of alarms don't match! expected=2, current='%v'", len(failed))
	}

	if failed[expectedNamespace].Counter != 1 {
		t.Errorf("The counter for failed alarm on namespace '%v' should be 1! expected=1, current=%v", expectedNamespace, failed[expectedNamespace].Counter)
	}

	if failed["other"].Counter != 2 {
		t.Errorf("The counter for failed alarm on namespace '%v' should be 2! expected=1, current=%v", "other", failed["other"].Counter)
	}

	store.RemoveAlarm("other")
	store.RemoveAlarm(expectedNamespace)
}
