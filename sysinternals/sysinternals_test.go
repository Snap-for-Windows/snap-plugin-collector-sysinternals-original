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
            metrics := plugin.Metric{
                //Create fake values to test to see if these numbers come back
                plugin.Metric {
                    Namespace: plugin.NewNamespace("intel", "sysinternals", "processCount"),
                    Config: map[string]interface{}{"testint": int},
                    Data: 42,
                    Unit: "int",
                    Timestamp: time.Now(),
                },
            }
            mts, er := pm.CollectMetrics(metrics)
            So(mts, ShouldBeEmpty)
            So(mts, ShouldBeNil)
            So(mts[0].Data, ShouldEqual, 0)
        }
    })
}