package logger

import (
	"fmt"
	"io"
	"os"
	"log"
	"suglider-auth/configs"
	"log/slog"
)

// No used, just backup
func LogMultiWriter(ServiceName string) {
	if configs.ApplicationConfig.Log.Filelog.Filename == "" {
		configs.ApplicationConfig.Log.Filelog.Filename = fmt.Sprintf("/var/log/%s/%s.log", ServiceName, ServiceName)
	}
	
	srvName := fmt.Sprintf("[%s]", ServiceName)
	flog := configs.ApplicationConfig.Log.Filelog
	mw := io.MultiWriter(os.Stdout, flog)
	
	if configs.Args.Syslog {
		logger, err := configs.ApplicationConfig.Log.Syslog.CreateSyslogLogger(srvName)
		if err != nil {
			log.Printf("Remote or local syslog server error: %v\n", err)
		} else {
			mlw := io.MultiWriter(os.Stdout, flog, logger)
			log.SetOutput(mlw)
		}
	}
	log.SetPrefix(srvName)
	log.SetFlags(log.LstdFlags|log.Lshortfile)
	log.SetOutput(mw)
}

func SlogMultiWriter(ServiceName string) {
	if configs.ApplicationConfig.Log.Filelog.Filename == "" {
		configs.ApplicationConfig.Log.Filelog.Filename = fmt.Sprintf("/var/log/%s/%s.log", ServiceName, ServiceName)
	}
	
	srvName := fmt.Sprintf("[%s]", ServiceName)
	flog := configs.ApplicationConfig.Log.Filelog
	mw := io.MultiWriter(os.Stdout, flog)
	
	if configs.Args.Syslog {
		logger, err := configs.ApplicationConfig.Log.Syslog.CreateSyslogLogger(srvName)
		if err != nil {
			log.Printf("Remote or local syslog server error: %v\n", err)
		} else {
			mlw := io.MultiWriter(os.Stdout, flog, logger)

			// Add attributes to all logs
			jsonHandler := slog.NewJSONHandler(mlw, nil).
				WithAttrs([]slog.Attr{slog.String("app", ServiceName)})
			logger := slog.New(jsonHandler)
			slog.SetDefault(logger)
		}
	}
	
	// Add attributes to all logs
	jsonHandler := slog.NewJSONHandler(mw, nil).
		WithAttrs([]slog.Attr{slog.String("app", ServiceName)})

	logger := slog.New(jsonHandler)
	slog.SetDefault(logger)
}