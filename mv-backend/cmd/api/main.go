package main

import (
	"log/slog"
	"os"

	"github.com/JoseLuis21/mv-backend/internal/server"
	"github.com/JoseLuis21/mv-backend/internal/shared/utils"
)

func main() {
	// Initialize structured logging for MisViaticos
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(handler))

	slog.Info("🚀 Starting MisViaticos Backend API")

	// Initialize control database connection
	dbControl, err := utils.InitDatabaseControl()
	if err != nil {
		slog.Error("Failed to initialize control database", "error", err.Error())
		os.Exit(1)
	}
	defer dbControl.Pool.Close()
	slog.Info("✅ Control database connected")

	// Create server instance with MisViaticos configuration
	host := utils.GetEnvOrDefault("HOST", "0.0.0.0")
	port := utils.GetEnvOrDefault("PORT", "8080")

	serverApi := server.NewServer(host, port, dbControl)

	slog.Info("🌟 MisViaticos API starting", "host", host, "port", port)

	if err := serverApi.Start(); err != nil {
		slog.Error("Failed to start MisViaticos server", "error", err)
		os.Exit(1)
	}
}
