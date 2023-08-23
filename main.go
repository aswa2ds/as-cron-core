/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"path"
	"runtime"

	"github.com/aswa2ds/as-cron-core/cmd"
	"github.com/aswa2ds/as-cron-core/config"
	db "github.com/aswa2ds/as-cron-db"
	log "github.com/sirupsen/logrus"
)

func main() {
	config.Init()
	db.Init(config.Config.DatabaseConfig)
	log.SetFormatter(&log.JSONFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
			fileName := path.Base(frame.File)
			return frame.Function, fmt.Sprintf("%s:%d", fileName, frame.Line)
		},
		PrettyPrint: false,
	})
	log.SetReportCaller(true)
	cmd.Execute()
}
