package version

import "fmt"

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// String renders the build's version, commit and date on one line.
// The three vars are set at link time by -ldflags -X (see release.yml and
// the environment/*.dockerfile builds); an unstamped build reports the
// hard-coded defaults.
func String() string {
	return fmt.Sprintf("%s (commit %s, built %s)", Version, Commit, Date)
}
