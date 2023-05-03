package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// func syncLoggerCycle() {
// 	for {
// 		_ = Logger.Sync()
// 		time.Sleep(5 * time.Millisecond)
// 	}
// }

func initloggertofile(logfile string) {
	var err error
	ConfigLog := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			TimeKey:       "time",
			NameKey:       "logger",
			CallerKey:     "caller",
			FunctionKey:   "function",
			StacktraceKey: "stacktrace",
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{logfile},
		ErrorOutputPaths: []string{logfile},
	}

	// Create file output
	file, err := os.Create(logfile)
	if err != nil {
		panic(err)
	}

	// Override console encoder with a more human-readable format
	ConfigLog.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	ConfigLog.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	// Create logger
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(ConfigLog.EncoderConfig),
		zapcore.AddSync(file),
		ConfigLog.Level,
	)
	Logger = zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
}

func Init(logfile string) {
	initloggertofile(logfile)
	// go syncLoggerCycle()
	createappforlogs(logfile)
}

func main() {
	Init("test.log")
	Logger.Info("test", zap.String("test", "test"))
}

func createappforlogs(logfile string) {
	// Create a new Fyne application
	myApp := app.New()

	// Create a new window
	myWindow := myApp.NewWindow("Test Log Viewer")

	// Create a text box widget to show the log content
	logText := widget.NewMultiLineEntry()
	logText.SetReadOnly(true)

	// Create a scroll container to hold the text box widget
	scrollContainer := widget.NewScrollContainer(logText)

	// Create a container to hold the scroll container widget
	container := fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		scrollContainer,
	)

	// Set the content of the window to the container widget
	myWindow.SetContent(container)

	// Start a goroutine to read the log file and update the text box widget every 1 second
	go func() {
		for {
			// Open the test.log file
			file, err := os.Open(logfile)
			if err != nil {
				logText.SetText(fmt.Sprintf("Error: %v", err))
			} else {
				// Create a scanner to read the file line by line
				scanner := bufio.NewScanner(file)

				// Create a string builder to build the log content
				var sb strings.Builder

				// Read each line from the file and append it to the string builder
				for scanner.Scan() {
					sb.WriteString(scanner.Text())
					sb.WriteString("\n")
				}

				// Set the text of the text box widget to the log content
				logText.SetText(sb.String())

				// Close the file
				file.Close()
			}

			// Wait for 1 second before reading the file again
			time.Sleep(time.Second)
		}
	}()

	// Set the size of the window to look like a terminal
	myWindow.Resize(fyne.NewSize(800, 600))

	// Show the window and run the app
	myWindow.ShowAndRun()
}
