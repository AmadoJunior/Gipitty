package utils

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GenerateError(err error, origin error) error {
	pc, file, line, _ := runtime.Caller(2)
	fn := runtime.FuncForPC(pc)
	functionName := strings.Split(fn.Name(), "/")
	trace := fmt.Sprintf("\n\tat %s (%s:%d)", functionName[len(functionName)-1], file, line)
	return fmt.Errorf("%w:"+trace+"\n%w", err, origin)
}

func LogError(err error, ctx *gin.Context) {
	t := time.Now()
	logFile, er := os.OpenFile("/home/amado/Documents/Gipitty/logs/"+t.Format("2006-01-02")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if er != nil {
		fmt.Println(er)
	}
	log.SetOutput(logFile)
	logger := log.New(logFile, "ERROR: ", log.LstdFlags|log.Lshortfile)
	logger.Println("Error:\n" + err.Error())
	logger.Println("Method:", ctx.Request.Method)
	logger.Println("URI:", ctx.Request.RequestURI)
	logger.Println("User Agent:", ctx.Request.UserAgent())
	logger.Println("IP:", ctx.ClientIP())
	logger.Println("Status:", ctx.Writer.Status())
}
