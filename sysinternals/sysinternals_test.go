package sysinternals

import (
	"testing"
	"time"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSysinternalsCollector(t *testing.T) {
	pm := SysinternalsCollector{}

	Convey("Test SysinternalsCollector", t, func() {
		Convey("Collect number of processes", func() {
			metrics := []plugin.Metric{
				//Create fake values to test to see if these numbers come back
				plugin.Metric{
					Namespace: plugin.NewNamespace("intel", "sysinternals", "processCount"),
					Config:    map[string]interface{}{"testint": int64(42)},
					Data:      42,
					Unit:      "int",
					Timestamp: time.Now(),
				},
			}
			mts, err := pm.CollectMetrics(metrics)
			So(mts, ShouldNotBeEmpty)
			So(err, ShouldBeNil)
			So(mts[0].Data, ShouldEqual, 42)
		})
	})

	Convey("Test GetMetricTypes", t, func() {
		pm := SysinternalsCollector{}

		Convey("Collect All Metrics String", func() {
			mt, err := pm.GetMetricTypes(nil)
			So(err, ShouldBeNil)
			So(len(mt), ShouldEqual, 3)
		})
	})

	Convey("Test GetConfigPolicy", t, func() {
		pm := SysinternalsCollector{}
		_, err := pm.GetConfigPolicy()

		Convey("No error returned", func() {
			So(err, ShouldBeNil)
		})
	})
}
