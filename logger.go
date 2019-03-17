package main

import "log"

type logger struct {
}

var miniLog *logger

func initMiniLog() {
	miniLog = getLogger()
}

func (l *logger) info(v ...interface{}) {
	log.Println("Info(√): ", v)
}

func (l *logger) error(v ...interface{}) {
	log.Println("Error(❌): ", v)
}

func (l *logger) fatal(v ...interface{}) {
	log.Fatalln("Fatal(❌): ", v)
}

func getLogger() *logger {
	return &logger{}
}
