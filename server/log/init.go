package log

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/utils/fs"
)

var theAccessLog, theErrorLog *os.File
var accessLogLock, errorLogLock sync.RWMutex

func getAccessLog() *os.File {
	accessLogLock.RLock()
	defer accessLogLock.RUnlock()
	return theAccessLog
}

func getErrorLog() *os.File {
	errorLogLock.RLock()
	defer errorLogLock.RUnlock()
	return theErrorLog
}

func init() {
	if isDevMode {
		theAccessLog, theErrorLog = os.Stdout, os.Stderr
		return
	}

	if err := setupLogger(); err != nil {
		log.Fatal(err)
	}
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGUSR1)

	go func() {
		for {
			<-ch
			if err := setupLogger(); err != nil {
				log.Println(err)
			}
		}
	}()
}

func setupLogger() error {
	logDir := filepath.Join(config.Root(), `log`)
	if err := os.MkdirAll(logDir, 0775); err != nil {
		return err
	}
	accessLogPath := filepath.Join(logDir, `app.log`)
	if accessLog, err := fs.OpenAppend(accessLogPath); err != nil {
		return err
	} else {
		accessLogLock.Lock()
		if theAccessLog != nil {
			theAccessLog.Close()
		}
		theAccessLog = accessLog
		accessLogLock.Unlock()
	}
	errorLogPath := filepath.Join(logDir, `app.err`)
	if errorLog, err := fs.OpenAppend(errorLogPath); err != nil {
		return err
	} else {
		errorLogLock.Lock()
		if theErrorLog != nil {
			theErrorLog.Close()
		}
		theErrorLog = errorLog
		errorLogLock.Unlock()
	}
	return nil
}
