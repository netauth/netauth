package health

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
)

var (
	checks map[string]SubsystemCheck

	logger = hclog.L().Named("health")
)

// SubsystemStatus contains the information needed to be returned by
// callback checks to determine worthyness to serve.
type SubsystemStatus struct {
	OK     bool
	Name   string
	Status string
}

// String provides the string representation of the SubsystemStatus
func (s SubsystemStatus) String() string {
	statword := "FAIL"
	if s.OK {
		statword = "PASS"
	}
	return fmt.Sprintf("[%s]  %s\t%s",
		statword,
		s.Name,
		s.Status)
}

// A SubsystemCheck is a function that is supplied by a particular
// subsystem of the server which then provides informatino about that
// subsystem during a status poll.
type SubsystemCheck func() SubsystemStatus

// A SystemStatus is returned containing the status of the entire
// server, and in the event of a failure will have the FirstFailure
// encountered called out.
type SystemStatus struct {
	OK           bool
	FirstFailure SubsystemStatus
	Subsystems   []SubsystemStatus
}

// String provides the string representation of the SystemStatus
func (s SystemStatus) String() string {
	out := ""

	statword := "FAIL"
	if s.OK {
		statword = "PASS"
	}
	out += fmt.Sprintf("System Check Status: %s\n", statword)

	if s.FirstFailure != (SubsystemStatus{}) {
		out += fmt.Sprintf("\nFirst Failure:\n%s\n", s.FirstFailure)
	}

	if len(s.Subsystems) > 0 {
		out += "\nSubsystems:\n"
	}

	for _, st := range s.Subsystems {
		out += fmt.Sprintf("%s\n", st)
	}

	return out
}

func init() {
	checks = make(map[string]SubsystemCheck)
}

// RegisterCheck allows an interested subsystem to register a check
// that will be called when health status is requested.
func RegisterCheck(name string, check SubsystemCheck) {
	if _, ok := checks[name]; ok {
		logger.Warn("Refusing to overwrite existing check", "check", name)
		return
	}
	checks[name] = check
	return
}

// Check runs all the health checks and returns the aggregate status
// to the caller.
func Check() SystemStatus {
	status := SystemStatus{
		OK: true,
	}
	logger.Debug("Running health check")
	for name, check := range checks {
		logger.Trace("Polling subsystem", "check", name)
		result := check()
		status.Subsystems = append(status.Subsystems, result)
		status.OK = status.OK && result.OK
		if !result.OK && status.FirstFailure == (SubsystemStatus{}) {
			status.FirstFailure = result
		}
	}
	return status
}
