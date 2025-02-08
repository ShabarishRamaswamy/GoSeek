package structs

import (
	"context"
	"database/sql"
)

type HTTPWebserver struct {
	Ctx            context.Context
	BaseWorkingDir string
	DB             *sql.DB
}

func GetHTTPWebserver(ctx context.Context, workingDir string, db *sql.DB) *HTTPWebserver {
	return &HTTPWebserver{
		Ctx:            ctx,
		BaseWorkingDir: workingDir,
		DB:             db,
	}
}
