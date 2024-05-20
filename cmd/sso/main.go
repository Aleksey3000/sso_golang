package main

import (
	"SSO/internal/app"
	"SSO/internal/config"
	"fmt"
	"log/slog"
	"os"
)

const ConfigPath = "config/config.yaml"

func main() {
	cnf, err := config.GetConfig(ConfigPath)
	if err != nil {
		panic(err)
	}
	l := SetupLogger()
	l.Info(fmt.Sprintf("Config: %+v", cnf))
	l.Info(fmt.Sprintf("PATH: %+v", os.Args[0]))
	App := app.New(l, cnf)
	defer App.GRPCApp.Stop()

	go func() {
		if err := App.HTTPApp.Run(); err != nil {
			l.Error(err.Error())
			panic(err)
		}
	}()

	if err := App.GRPCApp.Run(); err != nil {
		l.Error(err.Error())
		panic(err)
	}

}

func SetupLogger() *slog.Logger {
	var log *slog.Logger
	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))
	return log
}
