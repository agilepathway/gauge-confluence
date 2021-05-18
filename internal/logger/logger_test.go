//nolint:errcheck,gosec
package logger

import (
	"path/filepath"
	"testing"

	"os"

	logging "github.com/op/go-logging"
)

func TestLoggerInitWithInfoLevel(t *testing.T) {
	Initialize("info")

	if !loggersMap.getLogger().IsEnabledFor(logging.INFO) {
		t.Error("Expected gaugeConfluenceLog to be enabled for INFO")
	}
}

func TestLoggerInitWithDefaultLevel(t *testing.T) {
	Initialize("")

	if !loggersMap.getLogger().IsEnabledFor(logging.INFO) {
		t.Error("Expected gaugeConfluenceLog to be enabled for default log level")
	}
}

func TestLoggerInitWithDebugLevel(t *testing.T) {
	Initialize("debug")

	if !loggersMap.getLogger().IsEnabledFor(logging.DEBUG) {
		t.Error("Expected gaugeConfluenceLog to be enabled for DEBUG")
	}
}

func TestLoggerInitWithWarningLevel(t *testing.T) {
	Initialize("warning")

	if !loggersMap.getLogger().IsEnabledFor(logging.WARNING) {
		t.Error("Expected gaugeConfluenceLog to be enabled for WARNING")
	}
}

func TestLoggerInitWithErrorLevel(t *testing.T) {
	Initialize("error")

	if !loggersMap.getLogger().IsEnabledFor(logging.ERROR) {
		t.Error("Expected gaugeConfluenceLog to be enabled for ERROR")
	}
}

func TestGetLogFileWhenLogsDirNotSet(t *testing.T) {
	want, _ := filepath.Abs(filepath.Join(logs, gaugeConfluenceLogFileName))

	got := getLogFile()
	if got != want {
		t.Errorf("Got %s, want %s", got, want)
	}
}

func TestGetLogFileWhenRelativeCustomLogsDirIsSet(t *testing.T) {
	myLogsDir := "my_logs"
	os.Setenv(logsDirectory, myLogsDir)

	defer os.Unsetenv(logsDirectory)

	want, _ := filepath.Abs(filepath.Join(myLogsDir, gaugeConfluenceLogFileName))

	got := getLogFile()

	if got != want {
		t.Errorf("Got %s, want %s", got, want)
	}
}

//nolint:errcheck,gosec
func TestGetLogFileInGaugeProjectWhenAbsoluteCustomLogsDirIsSet(t *testing.T) {
	myLogsDir, err := filepath.Abs("my_logs")
	if err != nil {
		t.Errorf("Unable to convert to absolute path, %s", err.Error())
	}

	os.Setenv(logsDirectory, myLogsDir)
	defer os.Unsetenv(logsDirectory)

	want := filepath.Join(myLogsDir, gaugeConfluenceLogFileName)

	got := getLogFile()

	if got != want {
		t.Errorf("Got %s, want %s", got, want)
	}
}
