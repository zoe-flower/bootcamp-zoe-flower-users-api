package env

import (
	"os/exec"
	"strings"
)

// IsDockerRunning is a check to see if docker is running on the system
func IsDockerRunning() bool {
	cmd := exec.Command("docker", "info")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	if strings.Contains(string(out), "Is the docker daemon running?") {
		return false
	}

	return true
}

func IsGoAWSRunning() bool {
	return isContainerRunning("flyt-microservices_goaws") && isContainerRunning("flyt-microservices_stepfn")
}

// IsS3Running will tell us if the S3 container is running
func IsS3Running() bool {
	return isContainerRunning("flyt-microservices_s3")
}

// IsDynamoRunning will tell us if the dynamodb container is running
func IsDynamoRunning() bool {
	return isContainerRunning("flyt-microservices_dynamodb")
}

func isContainerRunning(name string) bool {
	if !IsDockerRunning() {
		return false
	}
	cmd := exec.Command("docker", "ps", "--format='{{.Names}}'")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	if !strings.Contains(string(out), name) {
		return false
	}

	return true
}
