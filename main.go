package main

import (
	"fmt"

	"github.com/xhyonline/xutil/kv"
	_ "github.com/ziutek/mymysql/native"
)

func main() {

	config := &kv.Sentinel{
		MasterName:   "myMaster",
		SentinelIP:   "121.36.253.109",
		SentinelPort: "26379",
	}

	client, err := kv.SentinelGetClient(config)

	if err != nil {
		panic(err)
	}
	fmt.Println(client.SetString("name", "888", 0))
}
