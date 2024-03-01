package main

import (
	initModule "github.com/NeptuneYeh/simpletask/init"
)

func main() {
	initProcess := initModule.NewMainInitProcess()
	initProcess.Run()
}
