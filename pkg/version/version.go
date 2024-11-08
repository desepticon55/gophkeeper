package version

import (
	_ "embed"
	"go.uber.org/zap"
	"html/template"
	"strings"
)

//go:embed version.txt
var greeting string

var (
	Version   string
	BuildTime string
)

var (
	// Build version
	BuildVersion string
	// Build date
	BuildDate string
	// Build commit hash
	BuildCommit string
)

type buildInfo struct {
	Version string
	Date    string
	Commit  string
}

// Return pretty printed build info
func MakeBuildInfo(log *zap.Logger) string {
	builder := &strings.Builder{}
	info := buildInfo{
		Version: "N/A",
		Date:    "N/A",
		Commit:  "N/A",
	}

	if BuildVersion != "" {
		info.Version = BuildVersion
	}
	if BuildDate != "" {
		info.Date = BuildDate
	}
	if BuildCommit != "" {
		info.Commit = BuildCommit
	}

	tmpl := template.Must(template.New("version").Parse(greeting))
	if err := tmpl.Execute(builder, info); err != nil {
		log.Error("Error during make buildInfo", zap.Error(err))
	}
	return builder.String()
}
