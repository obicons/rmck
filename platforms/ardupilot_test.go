package platforms

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/obicons/avis/util"
)

func TestUnitNewArduPilotFromEnvNoEnvVar(t *testing.T) {
	prev := os.Getenv("ARDUPILOT_SRC_PATH")
	defer os.Setenv("ARDUPILOT_SRC_PATH", prev)

	// the same as having no value
	err := os.Setenv("ARDUPILOT_SRC_PATH", "")
	if err != nil {
		t.Fatalf("Setenv returned unexpected error: %s", err)
	}

	_, err = NewArduPilotFromEnv()
	if err == nil {
		t.Fatal("NewArduPilotFromEnv did not return an error when it should have")
	}
}

func TestUnitNewArduPilotFromEnvBadPath(t *testing.T) {
	prev := os.Getenv("ARDUPILOT_SRC_PATH")
	defer os.Setenv("ARDUPILOT_SRC_PATH", prev)

	err := os.Setenv("ARDUPILOT_SRC_PATH", "/there/is/no/way/this/path/exists")
	if err != nil {
		t.Fatalf("Setenv returned unexpected error: %s", err)
	}

	_, err = NewArduPilotFromEnv()
	if err == nil {
		t.Fatal("NewArduPilotFromEnv did not return an error when it should have")
	}
}

func TestFunctionalArduPilot(t *testing.T) {
	system, err := NewArduPilotFromEnv()
	if err != nil {
		t.Fatalf("NewArduPilotFromEnv returned an unexpected error: %s", err)
	}

	ardupilot, ok := system.(*ArduPilot)
	if !ok {
		t.Fatal("Error: NewArduPilotFromEnv did not return an instance of ArduPilot")
	}

	err = ardupilot.Start()
	if err != nil {
		t.Fatalf("Start returned an unexpected error: %s", err)
	}

	time.Sleep(time.Second * 30)

	startTime := time.Now()
	isRunning := false
	for !isRunning && time.Now().Sub(startTime) < time.Second*10 {
		isRunning, _ = util.IsRunning("arducopter")
		time.Sleep(time.Millisecond * 250)
	}
	if !isRunning {
		t.Fatal("ArduCopter does not appear to have started")
	}

	ctx, cc := context.WithTimeout(context.Background(), time.Second*5)
	defer cc()
	err = ardupilot.Stop(ctx)
	if err != nil {
		t.Fatalf("ArduCopter could not stop: %s", err)
	}

	if !ardupilot.cmd.ProcessState.Exited() {
		t.Fatal("ArduCopter did not successfully stop")
	}
}
