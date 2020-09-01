package platforms

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/obicons/rmck/util"
)

type PX4 struct {
	srcPath string
	cmd     *exec.Cmd
}

func NewPX4FromEnv() (System, error) {
	px4Path := os.Getenv("PX4_PATH")
	if px4Path == "" {
		return nil, fmt.Errorf("error: NewPX4FromEnv(): set PX4_PATH")
	} else if _, err := os.Stat(px4Path); err != nil {
		return nil, fmt.Errorf("error: NewPX4FromEnv(): %s", err)
	}

	px4 := PX4{
		srcPath: px4Path,
		cmd:     nil,
	}

	return &px4, nil
}

// implements System
func (px4 *PX4) Start() error {
	binaryPath := path.Join(px4.srcPath, "bin/px4")
	romfsPath := path.Join(px4.srcPath, "ROMFS/px4fmu_common")
	rcPath := path.Join(px4.srcPath, "etc/init.d-posix/rcS")
	testDataPath := path.Join(px4.srcPath, "test_data")
	if _, err := os.Stat(binaryPath); err != nil {
		return fmt.Errorf("error: Start(): build px4")
	}

	cmd := exec.Command(
		binaryPath,
		"-d", // disable user input
		romfsPath,
		"-s", // set startup path
		rcPath,
		"-t", // set test data
		testDataPath,
	)
	cmd.Dir = px4.srcPath
	cmd.Env = px4Environ()

	// do we need cmd.Stdin to be set?
	logging, err := util.GetLogger("px4")
	if err != nil {
		return err
	}

	err = util.LogProcess(cmd, logging)
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err == nil {
		px4.cmd = cmd
	}

	return err
}

// implements System
func (px4 *PX4) Stop(ctx context.Context) error {
	return util.GracefulStop(px4.cmd, ctx)
}

func px4Environ() []string {
	env := os.Environ()
	env = append(
		env,
		"HEADLESS=1",
		"PX4_HOME_LAT=-35.363261",
		"PX4_HOME_LON=149.165230",
		"PX4_HOME_ALT=584",
		"DISPLAY=:0",
	)
	return env
}
