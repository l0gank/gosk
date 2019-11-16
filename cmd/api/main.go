package main

import (
	"github.com/vickydk/gosk/interface/api"
	"github.com/vickydk/gosk/utl/config"
)

func main() {
	config.LoadEnv()
	checkErr(api.Start())
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
