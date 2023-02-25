package build

var CurrentCommit string

const (
	Version = "0.6.5"
	AppName = "saas"
)

func UserVersion() string {
	return Version + CurrentCommit
}
