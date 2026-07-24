package utils

import (
	"os"
	"path/filepath"
	"testing"
)

// withTempCwd switches to a fresh temp directory for the duration of the
// test so InitLogger's "logs/" directory and log file don't leak into the repo.
func withTempCwd(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(orig); err != nil {
			t.Fatal(err)
		}
	})
}

func TestInitLogger_infoLevel(t *testing.T) {
	withTempCwd(t)
	t.Setenv("LOG_LEVEL", "INFO")

	InitLogger()

	if GetLogLevel() != LogLevelInfo {
		t.Errorf("GetLogLevel() = %v, want LogLevelInfo", GetLogLevel())
	}
	if IsDebugEnabled() {
		t.Error("IsDebugEnabled() = true at INFO level, want false")
	}
	if _, err := os.Stat("logs"); err != nil {
		t.Errorf("expected logs/ directory to be created: %v", err)
	}
}

func TestInitLogger_debugLevel(t *testing.T) {
	withTempCwd(t)
	t.Setenv("LOG_LEVEL", "DEBUG")

	InitLogger()

	if GetLogLevel() != LogLevelDebug {
		t.Errorf("GetLogLevel() = %v, want LogLevelDebug", GetLogLevel())
	}
	if !IsDebugEnabled() {
		t.Error("IsDebugEnabled() = false at DEBUG level, want true")
	}
}

func TestInitLogger_warnAndErrorLevels(t *testing.T) {
	cases := map[string]LogLevel{
		"WARN":    LogLevelWarn,
		"WARNING": LogLevelWarn,
		"ERROR":   LogLevelError,
	}
	for env, want := range cases {
		t.Run(env, func(t *testing.T) {
			withTempCwd(t)
			t.Setenv("LOG_LEVEL", env)
			InitLogger()
			if GetLogLevel() != want {
				t.Errorf("GetLogLevel() = %v, want %v", GetLogLevel(), want)
			}
		})
	}
}

func TestInitLogger_unknownLevelDefaultsToInfo(t *testing.T) {
	withTempCwd(t)
	t.Setenv("LOG_LEVEL", "NOT_A_LEVEL")

	InitLogger()

	if GetLogLevel() != LogLevelInfo {
		t.Errorf("GetLogLevel() = %v, want LogLevelInfo for unknown value", GetLogLevel())
	}
}

func TestInitLogger_logDirBlockedByFile_fallsBackToConsole(t *testing.T) {
	withTempCwd(t)
	// Create a plain file named "logs" so os.MkdirAll("logs", ...) fails,
	// exercising InitLogger's "console only" fallback branch.
	if err := os.WriteFile("logs", []byte("not a directory"), 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LOG_LEVEL", "INFO")

	InitLogger() // must not panic

	if Info == nil || Error == nil {
		t.Error("expected loggers to be initialized even when log dir creation fails")
	}
}

func TestInitLogger_logFileUnwritable_fallsBackToConsole(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("running as root: permission checks don't apply")
	}
	withTempCwd(t)
	if err := os.Mkdir("logs", 0555); err != nil { // read-only dir: MkdirAll succeeds, OpenFile inside it fails
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod("logs", 0755) }) // t.TempDir() cleanup needs write access back

	t.Setenv("LOG_LEVEL", "INFO")

	InitLogger() // must not panic

	if Info == nil || Error == nil {
		t.Error("expected loggers to be initialized even when log file creation fails")
	}
}

func TestSetLogLevel(t *testing.T) {
	withTempCwd(t)
	InitLogger() // ensure loggers are non-nil before SetLogLevel logs through them

	SetLogLevel(LogLevelError)
	if GetLogLevel() != LogLevelError {
		t.Errorf("GetLogLevel() = %v, want LogLevelError", GetLogLevel())
	}

	SetLogLevel(LogLevelInfo)
	if GetLogLevel() != LogLevelInfo {
		t.Errorf("GetLogLevel() = %v, want LogLevelInfo", GetLogLevel())
	}
}

func TestGetLogLevelName_allLevels(t *testing.T) {
	withTempCwd(t)
	InitLogger()

	cases := []struct {
		level LogLevel
		want  string
	}{
		{LogLevelDebug, "DEBUG"},
		{LogLevelInfo, "INFO"},
		{LogLevelWarn, "WARN"},
		{LogLevelError, "ERROR"},
		{LogLevel(99), "UNKNOWN"},
	}
	for _, tc := range cases {
		SetLogLevel(tc.level) // indirectly exercises getLogLevelName via its log line
		if GetLogLevel() != tc.level {
			t.Errorf("GetLogLevel() = %v, want %v", GetLogLevel(), tc.level)
		}
	}
}

func TestInitLogger_writesToDatedLogFile(t *testing.T) {
	withTempCwd(t)
	t.Setenv("LOG_LEVEL", "INFO")

	InitLogger()
	Info.Println("hello from test")

	entries, err := os.ReadDir("logs")
	if err != nil {
		t.Fatalf("reading logs dir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected exactly one log file, got %d", len(entries))
	}
	content, err := os.ReadFile(filepath.Join("logs", entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	if len(content) == 0 {
		t.Error("expected log file to contain the logged message")
	}
}
