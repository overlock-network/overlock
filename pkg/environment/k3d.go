package environment

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (e *Environment) CreateK3dEnvironment(logger *zap.SugaredLogger) (string, error) {

	args := []string{
		"cluster", "create", e.name,
	}

	if e.mountPath != "" {
		args = append(args, "-v", e.mountPath+":/storage")
	}

	cmd := exec.Command("k3d", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", errors.Wrap(err, "error creating k3d cluster")
	}

	logger.Info("k3d cluster created successfully")
	return e.K3dContextName(), nil
}

func (e *Environment) DeleteK3dEnvironment(logger *zap.SugaredLogger) error {
	cmd := exec.Command("k3d", "cluster", "delete", e.name)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start k3d command: %w", err)
	}

	stderrScanner := bufio.NewScanner(stderr)
	for stderrScanner.Scan() {
		logger.Info(stderrScanner.Text())
	}
	return nil
}

func (e *Environment) K3dContextName() string {
	return "k3d-" + e.name
}
