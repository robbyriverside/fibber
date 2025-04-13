package fibber

var (
	// Version is set at build time with -ldflags "-X github.com/robbyriverside/fibber.Version=1.2.3"
	Version = "dev"

	// Commit is an optional git commit SHA set at build time
	Commit = ""

	// BuildTime is set at build time with -ldflags "-X github.com/robbyriverside/fibber.BuildTime=2025-04-11T15:04:05Z"
	BuildTime = ""
)
