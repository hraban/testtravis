package foo

import (
	"bitbucket.org/kardianos/osext"
)

// does one thing and returns another
func DoSomething() string {
	osext.Executable()
	return "batman"
}
