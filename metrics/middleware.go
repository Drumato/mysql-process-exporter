package metrics

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

// OndemandUpdateMetricsMiddleware updates metrics on demand when a request is received.
// Although updating metrics periodically is an idea, it would increase requests to a MySQL server.
// Therefore, we update the metrics only when a request is received.
func OndemandUpdateMetricsMiddleware(
	logger *slog.Logger,
	dbHost string,
	dbConn *sql.DB,
) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.DebugContext(c.Request().Context(), "trying to update metrics")
			rows, err := dbConn.Query("SHOW FULL PROCESSLIST")
			if err != nil {
				log.Fatalf("Failed to execute query: %v", err)
			}
			defer func() {
				if err := rows.Close(); err != nil {
					slog.ErrorContext(c.Request().Context(), "failed to close rows", slog.String("error", err.Error()))
				}
				c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "failed to close rows",
				})
			}()

			for rows.Next() {
				var (
					idC      sql.NullInt64
					userC    sql.NullString
					hostC    sql.NullString
					dbNameC  sql.NullString
					commandC sql.NullString
					timeC    sql.NullInt64
					stateC   sql.NullString
					infoC    sql.NullString
				)

				if err := rows.Scan(&idC, &userC, &hostC, &dbNameC, &commandC, &timeC, &stateC, &infoC); err != nil {
					logger.ErrorContext(c.Request().Context(), "Failed to scan row", slog.Any("error", err))
					return c.JSON(http.StatusInternalServerError, map[string]string{
						"error": "failed to close rows",
					})
				}
				if !timeC.Valid {
					logger.ErrorContext(c.Request().Context(), "Invalid time value", slog.Any("time", timeC))
					return c.JSON(http.StatusInternalServerError, map[string]string{
						"error": "failed to close rows",
					})
				}

				id := fmt.Sprintf("%d", lo.If(idC.Valid, idC.Int64).Else(0))
				user := lo.If(userC.Valid, userC.String).Else("")
				host := lo.If(hostC.Valid, hostC.String).Else("")
				dbName := lo.If(dbNameC.Valid, dbNameC.String).Else("")
				command := lo.If(commandC.Valid, commandC.String).Else("")
				state := lo.If(stateC.Valid, stateC.String).Else("")
				info := lo.If(infoC.Valid, infoC.String).Else("")
				UpdateMySQLProcessSecondsGaugeVec(MySQLProcessSecondsGaugeVecLabels{
					DBHost:  dbHost,
					ID:      id,
					User:    user,
					Host:    host,
					DB:      dbName,
					Command: command,
					State:   state,
					Info:    info,
				}, float64(timeC.Int64))
			}

			if err := rows.Err(); err != nil {
				logger.ErrorContext(c.Request().Context(), "Row iteration error", slog.Any("error", err))
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "failed to close rows",
				})
			}
			return next(c)
		}
	}
}
