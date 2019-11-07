package health

import (
	"testing"

	"github.com/golang/protobuf/proto"
)

func subsystemSuccess() SubsystemStatus {
	return SubsystemStatus{
		OK:     true,
		Name:   "Successful Subsystem",
		Status: "Working as expected",
	}
}

func subsystemFailure() SubsystemStatus {
	return SubsystemStatus{
		OK:     false,
		Name:   "Failed Subsystem",
		Status: "Internal Error",
	}
}

func subsystemFailureTwo() SubsystemStatus {
	return SubsystemStatus{
		OK:     false,
		Name:   "Failed Subsystem (1)",
		Status: "Internal Error Again",
	}
}

func TestHealthCheck(t *testing.T) {
	// This is a somewhat annoying hack needed to make the checks
	// work given that there's a single package scoped variable
	// that makes the imports work.
	checks = make(map[string]SubsystemCheck)

	// Register some checks with the first one attempted to
	// register twice.  The system must return OK on the first
	// check.
	RegisterCheck("one", subsystemSuccess)
	RegisterCheck("one", subsystemFailure)
	status := Check()
	if !status.OK {
		t.Error("Test was overwritten")
	}

	// Now we register a failing test and make sure it fails.
	RegisterCheck("two", subsystemFailure)
	status = Check()
	if status.OK || len(status.Subsystems) != 2 {
		t.Error("Wrong number of results")
	}

	// Now add a third test, that fails, and make sure
	// FirstFailure is set correctly.
	RegisterCheck("three", subsystemFailureTwo)
	status = Check()
	if status.OK || len(status.Subsystems) != 3 {
		t.Error("Wrong number of results")
	}
}

func TestSubsystemStatusString(t *testing.T) {
	cases := []struct {
		status SubsystemStatus
		text   string
	}{
		{
			status: SubsystemStatus{
				OK:     true,
				Name:   "Sub1",
				Status: "System is Ready",
			},
			text: "[PASS]  Sub1	System is Ready",
		},
		{
			status: SubsystemStatus{
				OK:     false,
				Name:   "Sub1",
				Status: "System is Broken",
			},
			text: "[FAIL]  Sub1	System is Broken",
		},
	}

	for i, c := range cases {
		if c.status.String() != c.text {
			t.Errorf("%d: Got %s Want %s", i, c.status, c.text)
		}
	}
}

func TestSubsystemStatusProto(t *testing.T) {
	status := SubsystemStatus{
		OK:     true,
		Name:   "example",
		Status: "everything is broken",
	}

	p := status.Proto()
	if p.GetName() != status.Name || p.GetOK() != status.OK || p.GetFaultMessage() != status.Status {
		t.Error("Proto is not identical")
	}

}

func TestSystemStatusString(t *testing.T) {
	cases := []struct {
		status SystemStatus
		text   string
	}{
		{
			status: SystemStatus{
				OK: true,
				Subsystems: []SubsystemStatus{
					{
						OK:     true,
						Name:   "Sub1",
						Status: "System is Ready",
					},
				},
			},
			text: "System Check Status: PASS\n\nSubsystems:\n[PASS]  Sub1	System is Ready\n",
		},
		{
			status: SystemStatus{
				OK: false,
				FirstFailure: SubsystemStatus{
					OK:     false,
					Name:   "Sub2",
					Status: "System is Broken",
				},
				Subsystems: []SubsystemStatus{
					{
						OK:     true,
						Name:   "Sub1",
						Status: "System is Ready",
					},
					{
						OK:     false,
						Name:   "Sub2",
						Status: "System is Broken",
					},
				},
			},
			text: "System Check Status: FAIL\n\nFirst Failure:\n[FAIL]  Sub2	System is Broken\n\nSubsystems:\n[PASS]  Sub1	System is Ready\n[FAIL]  Sub2	System is Broken\n",
		},
	}

	for i, c := range cases {
		if c.status.String() != c.text {
			t.Errorf("%d: Got %s Want %s", i, c.status, c.text)
		}
	}
}

func TestSystemStatusProto(t *testing.T) {
	status := SystemStatus{
		OK: false,
		FirstFailure: SubsystemStatus{
			OK:     false,
			Name:   "Sub2",
			Status: "System is Broken",
		},
		Subsystems: []SubsystemStatus{
			{
				OK:     true,
				Name:   "Sub1",
				Status: "System is Ready",
			},
			{
				OK:     false,
				Name:   "Sub2",
				Status: "System is Broken",
			},
		},
	}

	p := status.Proto()

	if p.GetSystemOK() != status.OK {
		t.Error("Overall status is incorrect")
	}

	if !proto.Equal(p.GetFirstFailure(), status.FirstFailure.Proto()) {
		t.Error("First failure is not correct")
	}
}
