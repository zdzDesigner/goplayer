package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"player/util"
)

// var DIR_ASSETS = fmt.Sprintf("%s/source/assets/enya", strings.TrimRight(util.Dir(), "/palyer"))

var DIR_ASSETS = "/home/zdz/temp/music/xuwei-resource/默认"

type Assets struct {
	DIR_ASSETS string `json:"DIR_ASSETS"`
}

func init() {
	f, err := os.Open(fmt.Sprintf("%s/%s", util.Dir(), "conf/config.json"))
	util.MustNoErr(err)

	bts, err := ioutil.ReadAll(f)
	util.MustNoErr(err)
	// fmt.Println(bts)

	assets := &Assets{}
	err = json.Unmarshal(bts, assets)
	util.MustNoErr(err)

	DIR_ASSETS = assets.DIR_ASSETS
}

// var DIR_ASSETS = "/home/zdz/temp/music/like-resource-one"

// var DIR_ASSETS = "/home/zdz/temp/music/wsgs-resource"

// var DIR_ASSETS = "/home/zdz/temp/music/lj-resource"

// var DIR_ASSETS = "/home/zdz/temp/music/lj-resource/李健-2003 似水流年"

// var DIR_ASSETS = "/home/zdz/temp/music/ape-resource2"

// var DIR_ASSETS = "/home/zdz/temp/music/ape-resource"

// var DIR_ASSETS = "/home/zdz/temp/music/qyy-resource"
// var DIR_ASSETS = "/home/zdz/temp/music/xuwei-resource"


// var DIR_ASSETS = "/home/zdz/temp/music/zhoujielun-resource"

// var DIR_ASSETS = "/home/zdz/temp/music/like"

// var DIR_ASSETS = "/home/zdz/temp/music/like-resource"
// var DIR_ASSETS = "/home/zdz/temp/music/lry"
// var DIR_ASSETS = "/home/zdz/temp/music/ywjd"
// var DIR_ASSETS = "/home/zdz/temp/music/piano"
