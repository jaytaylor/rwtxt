package main

import (
	"flag"
	"fmt"

	log "github.com/cihub/seelog"
	"github.com/jaytaylor/rwtxt"
	"github.com/jaytaylor/rwtxt/pkg/db"
)

var (
	dbName  string
	Version string
)

func main() {
	var (
		err         error
		debug       = flag.Bool("debug", false, "debug mode")
		showVersion = flag.Bool("v", false, "show version")
		database    = flag.String("db", "rwtxt.db", "name of the database")
		listen      = flag.String("listen", rwtxt.DefaultBind, "interface:port to listen on")
	)
	flag.Parse()

	if *showVersion {
		fmt.Println(Version)
		return
	}
	if *debug {
		err = setLogLevel("debug")
		db.SetLogLevel("debug")
	} else {
		err = setLogLevel("info")
		db.SetLogLevel("info")
	}
	if err != nil {
		panic(err)
	}
	dbName = *database
	defer log.Flush()

	fs, err := db.New(dbName)
	if err != nil {
		panic(err)
	}

	rwt, err := rwtxt.New(fs)
	if err != nil {
		panic(err)
	}
	if listen != nil && *listen != "" {
		rwt.Bind = *listen
	}

	err = rwt.Serve()
	if err != nil {
		log.Error(err)
	}
}

// setLogLevel determines the log level
func setLogLevel(level string) (err error) {

	// https://en.wikipedia.org/wiki/ANSI_escape_code#3/4_bit
	// https://github.com/cihub/seelog/wiki/Log-levels
	appConfig := `
	<seelog minlevel="` + level + `">
	<outputs formatid="stdout">
	<filter levels="debug,trace">
		<console formatid="debug"/>
	</filter>
	<filter levels="info">
		<console formatid="info"/>
	</filter>
	<filter levels="critical,error">
		<console formatid="error"/>
	</filter>
	<filter levels="warn">
		<console formatid="warn"/>
	</filter>
	</outputs>
	<formats>
		<format id="stdout"   format="%Date %Time [%LEVEL] %File %FuncShort:%Line %Msg %n" />
		<format id="debug"   format="%Date %Time %EscM(37)[%LEVEL]%EscM(0) %File %FuncShort:%Line %Msg %n" />
		<format id="info"    format="%Date %Time %EscM(36)[%LEVEL]%EscM(0) %File %FuncShort:%Line %Msg %n" />
		<format id="warn"    format="%Date %Time %EscM(33)[%LEVEL]%EscM(0) %File %FuncShort:%Line %Msg %n" />
		<format id="error"   format="%Date %Time %EscM(31)[%LEVEL]%EscM(0) %File %FuncShort:%Line %Msg %n" />
	</formats>
	</seelog>
	`
	logger, err := log.LoggerFromConfigAsBytes([]byte(appConfig))
	if err != nil {
		return
	}
	log.ReplaceLogger(logger)
	return
}
