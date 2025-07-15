package docker

import (
	"log"
	"os"
	"os/exec"
	"strconv"
)

func MakeDockerImage(hackathonId int, imageName string, dir string) error {
	image := imageName + ":" + strconv.Itoa(hackathonId) + ":latest"
	dockerfileDir := "."

	cmd := exec.Command("docker", "build", "-t", image, dockerfileDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("error creating docker image: %v\n", err)
		return err
	}

	return nil
}
