// Package example contains log statements for loglint testing.
// Run: go-loglint ./testdata/example
// Covers: log, slog, logrus, zap.SugaredLogger
package example

import (
	"log"
	"log/slog"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func example(password, apiKey, token string) {
	// --- log ---
	log.Print("Starting server")   // wrong
	log.Print("starting server")   // correct
	log.Printf("Error: %s", "x")  // wrong (Printf)
	log.Println("done")           // correct

	// --- slog ---
	slog.Info("Starting server")  // wrong
	slog.Info("starting server")   // correct
	slog.Error("Failed to connect")
	slog.Debug("debug message")
	slog.Warn("something wrong")

	// --- logrus ---
	logrus.Info("Starting server")  // wrong
	logrus.Info("starting server")  // correct
	logrus.Error("Failed to connect")
	logrus.Warn("warning message")
	logrus.Debug("debug info")

	// --- zap.SugaredLogger ---
	logger := zap.NewExample().Sugar()
	logger.Info("Starting server")  // wrong
	logger.Info("starting server")  // correct
	logger.Error("Failed to connect")
	logger.Debug("debug message")
	logger.Warn("warning")

	// Rule 2: english_only
	slog.Info("запуск сервера")     // wrong
	logrus.Info("ошибка")           // wrong

	// Rule 3: no_special_chars
	slog.Info("server started!")   // wrong
	log.Printf("failed!!!")        // wrong

	// Rule 4: no_sensitive_data
	slog.Info("user password " + password)  // wrong
	logrus.Debug("api_key" + apiKey)       // wrong
	logger.Info("token " + token)           // wrong
}
