package controllers

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

const testServerURL = "http://localhost:8081"

func startDockerCompose() error {
	cmd := exec.Command("docker", "compose", "-f", "../../../docker-compose.test.yml", "up", "--build", "-d")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker compose up failed: %v, output: %s", err, string(output))
	}
	return nil
}

func waitForServer() error {
	for i := 0; i < 15; i++ {
		resp, err := http.Get(testServerURL + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			return nil
		}
		fmt.Printf("Waiting for server... Attempt %d/15 (Error: %v)\n", i+1, err)
		time.Sleep(3 * time.Second)
	}
	return fmt.Errorf("server did not start in time")
}

func SetupDocker() {
	_ = startDockerCompose()
	_ = waitForServer()
}

func StopDockerCompose() {
	cmd := exec.Command("docker", "compose", "-f", "../../../docker-compose.test.yml", "down")
	cmd.Run()
}
