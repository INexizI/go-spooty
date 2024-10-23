package logg

import (
	"log"
	"os"
)

var (
	SystemLogger  *log.Logger
	CommandLogger *log.Logger // TODO(D): create log for only bot commands
	MessageLogger *log.Logger
)

func init() {
	sysFile, err := os.OpenFile("internal/log/system.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(sysFile)
	SystemLogger = log.New(sysFile, "SYS: ", log.Ldate|log.Ltime|log.Lshortfile)

	cmdFile, err := os.OpenFile("internal/log/command.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(cmdFile)
	CommandLogger = log.New(cmdFile, "CMD: ", log.Ldate|log.Ltime|log.Lshortfile)

	msgFile, err := os.OpenFile("internal/log/message.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(msgFile)
	MessageLogger = log.New(msgFile, "MSG: ", log.Ldate|log.Ltime|log.Lshortfile)
}
