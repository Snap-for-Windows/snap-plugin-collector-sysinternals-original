package main

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(t *testing.T) {
	Convey("To ensure plugin loads and responds correctly", t, func() {
		os.Args = []string{"", "{\"NoDaemon\": true}"}
		So(func() { main() }, ShouldNotPanic)
	})
}