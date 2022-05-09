package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/uretgec/mydb-server/mystorage/services"

	"github.com/uretgec/mydb-server/mystorage/loggers"

	"github.com/peterbourgon/ff/v3"
	"github.com/tidwall/redcon"
)

var ServiceBuild string
var ServiceCommitId string

var ServiceName string = "mystoragedefault"
var ServiceVersion string = "1.0.0"

var (
	fs = flag.NewFlagSet(ServiceName, flag.ExitOnError)

	bucketList, indexList ArrayFlagString

	syncInterval = fs.Int("sync-interval", 30, "sync interval number for db synced period")
	redisAddr    = fs.String("redis-addr", "localhost:6379", "redis addr for sync data to storage service")
	dbName       = fs.String("db-name", ServiceName, "database name")
	dbPath       = fs.String("db-path", "./", "database path")
	dbReadOnly   = fs.Bool("db-read-only", false, "database read only mode")
	_            = fs.String("env-file", ".env", "env file")
)

var ps redcon.PubSub

func main() {

	// Zap Logger Init
	loggers.SetupSugarLogger(ServiceName, ServiceVersion)

	// Start Logger: Haydaaa
	loggers.Sugar.Info("service started")

	// Flag Parse with Env
	fs.Var(&bucketList, "bucket-list", "db bucket list") // Multiple
	fs.Var(&indexList, "index-list", "db index list")    // Multiple

	// Flag Parse with Env
	err := ff.Parse(
		fs, os.Args[1:],
		ff.WithConfigFileFlag("env-file"),
		ff.WithConfigFileParser(ff.PlainParser),
		ff.WithEnvVarPrefix(strings.ToUpper(ServiceName)),
	)

	if err != nil {
		loggers.Sugar.With("error", err).Fatal("configration error")
	}

	// Database Conn
	err = services.SetupStorage(*dbName, *dbPath, bucketList, indexList, *dbReadOnly)
	if err != nil {
		loggers.Sugar.With("error", err).Fatal("store db error")
	}
	defer services.Store.CloseStore()

	go func() {
		for {
			time.Sleep(time.Duration(*syncInterval) * time.Second)
			services.Store.SyncStore()
		}
	}()

	// Redcon Init
	server := setupRoutes()

	go func() {
		err = redcon.ListenAndServe(*redisAddr, server.ServeRESP, connAccept, connClose)
		if err != nil {
			loggers.Sugar.With("error", err).Fatal("redcon conn error")
			os.Exit(1)
		}
	}()

	// Listen server quit or something happened and notify channel
	close := make(chan os.Signal, 1)
	signal.Notify(close, syscall.SIGINT, syscall.SIGTERM)

	<-close

	// Sync ServiceData
	services.Store.SyncStore()

	// Close ServiceData
	services.Store.CloseStore()

	// Bye bye
	loggers.Sugar.Info("im shutting down. See you later")
}
