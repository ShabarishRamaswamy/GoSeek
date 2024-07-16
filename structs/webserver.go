package structs

import "context"

type HTTPWebserver struct {
	Ctx            context.Context
	BaseWorkingDir string
}

func GetHTTPWebserver(ctx context.Context, workingDir string) *HTTPWebserver {
	return &HTTPWebserver{
		Ctx:            ctx,
		BaseWorkingDir: workingDir,
	}
}
