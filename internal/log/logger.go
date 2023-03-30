package log

import "log"

func Info(v ...interface{}) {
	log.Println(v...)
}

func Infof(format string, v ...interface{}) {
	log.Printf(format+"\n", v...)
}

func Error(v ...interface{}) {
	log.Println(v...)
}

func Errorf(format string, v ...interface{}) {
	log.Printf(format+"\n", v...)
}
