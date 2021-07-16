package util

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRandom(t *testing.T) {
	fmt.Println("+++++++++++++++ TestRandom 随机数 +++++++++++++++")

	Convey("TestRandom 范围随机数 => \n", t, func() {
		Convey(`随机数 Random(12, 13) =>`, func() {
			So(Random(12, 13), ShouldEqual, 12)
			fmt.Println(Random(12, 13))
		})
	})

	Convey("TestRandom 随机数 => \n", t, func() {
		Convey(`随机数 Random(3) => 0,1,2`, func() {
			v := Random(3)
			ok := Contains([]int{0, 1, 2}, v)
			So(ok, ShouldEqual, true)
			fmt.Println(v)
		})
	})

	Convey("TestRandom 1 随机数 => \n", t, func() {
		Convey(`随机数 Random(1) => 0`, func() {
			v := Random(1)
			So(v, ShouldEqual, 0)
			fmt.Println(v)
		})
	})

}

func TestRandomMutil(t *testing.T) {
	fmt.Println("+++++++++++++++ TestRandomMutil 随机数多个 +++++++++++++++")

	Convey("TestRandomMutil 随机数 => \n", t, func() {
		Convey(`随机数 TestRandomMutil(12, 13) 5, 需求值大于容量, 取交集 =>`, func() {
			arr := RandomMutil(5, 12, 13)
			So(len(arr), ShouldEqual, 1)
			fmt.Println(arr)
		})
	})

	Convey("TestRandomMutil 随机数 => \n", t, func() {
		Convey(`随机数 TestRandomMutil(12, 43) 5=>`, func() {
			arr := RandomMutil(5, 12, 43)
			So(len(arr), ShouldEqual, 5)
			fmt.Println(arr)
		})
	})

	Convey("TestRandomMutil 随机数 => \n", t, func() {
		Convey(`随机数 TestRandomMutil(0, 0) 5=>`, func() {
			arr := RandomMutil(5, 0, 0)
			So(len(arr), ShouldEqual, 0)
			fmt.Println(arr)
		})
	})

	Convey("TestRandomMutil 随机数 => \n", t, func() {
		Convey(`随机数 TestRandomMutil(3, 3) => `, func() {
			arr := RandomMutil(3, 3)
			So(len(arr), ShouldEqual, 3)
			fmt.Println(arr)
		})
	})

	Convey("TestRandomMutil 随机数 => \n", t, func() {
		Convey(`随机数 TestRandomMutil(20, 12) => `, func() {
			arr := RandomMutil(20, 12)
			So(len(arr), ShouldEqual, 12)
			fmt.Println(arr)
		})
	})

}
