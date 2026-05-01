package upgrade

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/ku9nov/backup/configs"
	"github.com/sirupsen/logrus"
)

func TestRunUpdateAvailable(t *testing.T) {
	originalInstall := installDownloadedArtifactFn
	installDownloadedArtifactFn = func(_ string) error { return nil }
	t.Cleanup(func() {
		installDownloadedArtifactFn = originalInstall
	})

	var gotQuery map[string]string
	gotDeviceID := ""
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/checkVersion":
			gotDeviceID = strings.TrimSpace(r.Header.Get("X-Device-ID"))
			gotQuery = map[string]string{}
			for key, values := range r.URL.Query() {
				if len(values) > 0 {
					gotQuery[key] = values[0]
				}
			}
			_, _ = fmt.Fprintf(w, `{"critical":false,"update_available":true,"update_url":"%s/files/backup-1.0.0"}`, srv.URL)
		case "/files/backup-1.0.0":
			_, _ = fmt.Fprint(w, "binary-content")
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer srv.Close()

	configPath := testConfigPath(t)
	logger, _ := newTestLogger()
	cfg := &configs.Config{}
	cfg.Upgrade.Server = srv.URL
	cfg.Upgrade.Owner = "admin"
	cfg.Upgrade.App = "backup"

	err := Run(Input{
		Logger:  logger,
		Config:  cfg,
		Version: "0.9.0",
		Channel: "stable",
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	assertQuery(t, gotQuery, "app_name", "backup")
	assertQuery(t, gotQuery, "version", "0.9.0")
	assertQuery(t, gotQuery, "channel", "stable")
	assertQuery(t, gotQuery, "platform", runtime.GOOS)
	assertQuery(t, gotQuery, "arch", runtime.GOARCH)
	assertQuery(t, gotQuery, "owner", "admin")
	if gotDeviceID == "" {
		t.Fatal("expected X-Device-ID header to be set")
	}

	tmpFile := filepath.Join(filepath.Dir(configPath), ".backup", "tmp", "backup-1.0.0")
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("expected downloaded file at %s: %v", tmpFile, err)
	}
	if string(content) != "binary-content" {
		t.Fatalf("unexpected downloaded content: %q", string(content))
	}
}

func TestRunNoUpdate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, `{"critical":false,"possible_rollback":false,"update_available":false,"update_url":"https://updates.example/backup-1.0.0"}`)
	}))
	defer srv.Close()

	testConfigPath(t)
	logger, logOutput := newTestLogger()
	cfg := &configs.Config{}
	cfg.Upgrade.Server = srv.URL
	cfg.Upgrade.Owner = "admin"
	cfg.Upgrade.App = "backup"

	err := Run(Input{
		Logger:  logger,
		Config:  cfg,
		Version: "1.0.0",
		Channel: "stable",
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if !strings.Contains(logOutput.String(), "current version is up-to-date") {
		t.Fatalf("expected up-to-date log, got: %s", logOutput.String())
	}
}

func TestRunSendsStableDeviceIDAcrossRuns(t *testing.T) {
	originalInstall := installDownloadedArtifactFn
	installDownloadedArtifactFn = func(_ string) error { return nil }
	t.Cleanup(func() {
		installDownloadedArtifactFn = originalInstall
	})

	receivedIDs := make([]string, 0, 2)
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/checkVersion":
			receivedIDs = append(receivedIDs, strings.TrimSpace(r.Header.Get("X-Device-ID")))
			_, _ = fmt.Fprintf(w, `{"critical":false,"update_available":true,"update_url":"%s/files/backup-1.0.0"}`, srv.URL)
		case "/files/backup-1.0.0":
			_, _ = fmt.Fprint(w, "binary-content")
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer srv.Close()

	testConfigPath(t)
	logger, _ := newTestLogger()
	cfg := &configs.Config{}
	cfg.Upgrade.Server = srv.URL
	cfg.Upgrade.Owner = "admin"
	cfg.Upgrade.App = "backup"

	for i := 0; i < 2; i++ {
		err := Run(Input{
			Logger:  logger,
			Config:  cfg,
			Version: "0.9.0",
			Channel: "stable",
		})
		if err != nil {
			t.Fatalf("Run returned error: %v", err)
		}
	}

	if len(receivedIDs) != 2 {
		t.Fatalf("expected 2 checkVersion requests, got %d", len(receivedIDs))
	}
	if receivedIDs[0] == "" || receivedIDs[1] == "" {
		t.Fatalf("expected non-empty device ids, got %q and %q", receivedIDs[0], receivedIDs[1])
	}
	if receivedIDs[0] != receivedIDs[1] {
		t.Fatalf("expected stable device id across runs, got %q and %q", receivedIDs[0], receivedIDs[1])
	}
}

func newTestLogger() (*logrus.Logger, *bytes.Buffer) {
	buf := bytes.NewBuffer(nil)
	logger := logrus.New()
	logger.SetOutput(buf)
	logger.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})
	logger.SetLevel(logrus.TraceLevel)
	return logger, buf
}

func testConfigPath(t *testing.T) string {
	t.Helper()

	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	return filepath.Join(tempDir, "config.yml")
}

func assertQuery(t *testing.T, query map[string]string, key, want string) {
	t.Helper()

	if got := query[key]; got != want {
		t.Fatalf("unexpected query %s: want %q, got %q", key, want, got)
	}
}
