package docker

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func IsInstalled() bool {
	cmd := exec.Command("docker", "--version")
	return cmd.Run() == nil
}

func IsComposeAvailable() bool {
	cmd := exec.Command("docker", "compose", "version")
	return cmd.Run() == nil
}

func Install() error {
	switch runtime.GOOS {
	case "linux":
		return installLinux()
	case "darwin":
		return fmt.Errorf("please install Docker Desktop from https://www.docker.com/products/docker-desktop")
	case "windows":
		return fmt.Errorf("please install Docker Desktop from https://www.docker.com/products/docker-desktop")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func installLinux() error {
	if _, err := os.Stat("/etc/debian_version"); err == nil {
		return installDockerDebian()
	}
	if _, err := os.Stat("/etc/redhat-release"); err == nil {
		return installDockerRedHat()
	}
	if _, err := os.Stat("/etc/arch-release"); err == nil {
		return installDockerArch()
	}

	return fmt.Errorf("unsupported Linux distribution. Please install Docker manually")
}

func installDockerDebian() error {
	commands := [][]string{
		{"sudo", "apt-get", "update"},
		{"sudo", "apt-get", "install", "-y", "ca-certificates", "curl"},
		{"sudo", "install", "-m", "0755", "-d", "/etc/apt/keyrings"},
		{"sudo", "curl", "-fsSL", "https://download.docker.com/linux/ubuntu/gpg", "-o", "/etc/apt/keyrings/docker.asc"},
		{"sudo", "chmod", "a+r", "/etc/apt/keyrings/docker.asc"},
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run %v: %w", cmdArgs, err)
		}
	}

	addRepoCmd := `echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null`
	cmd := exec.Command("bash", "-c", addRepoCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	installCmd := exec.Command("sudo", "apt-get", "update")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		return err
	}

	installCmd = exec.Command("sudo", "apt-get", "install", "-y", "docker-ce", "docker-ce-cli", "containerd.io", "docker-compose-plugin")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	return installCmd.Run()
}

func installDockerRedHat() error {
	cmd := exec.Command("sudo", "yum", "install", "-y", "docker")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	startCmd := exec.Command("sudo", "systemctl", "start", "docker")
	startCmd.Stdout = os.Stdout
	startCmd.Stderr = os.Stderr
	if err := startCmd.Run(); err != nil {
		return err
	}

	enableCmd := exec.Command("sudo", "systemctl", "enable", "docker")
	enableCmd.Stdout = os.Stdout
	enableCmd.Stderr = os.Stderr
	return enableCmd.Run()
}

func installDockerArch() error {
	cmd := exec.Command("sudo", "pacman", "-S", "--noconfirm", "docker", "docker-compose")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	startCmd := exec.Command("sudo", "systemctl", "start", "docker")
	startCmd.Stdout = os.Stdout
	startCmd.Stderr = os.Stderr
	if err := startCmd.Run(); err != nil {
		return err
	}

	enableCmd := exec.Command("sudo", "systemctl", "enable", "docker")
	enableCmd.Stdout = os.Stdout
	enableCmd.Stderr = os.Stderr
	return enableCmd.Run()
}

func ComposePull() error {
	cmd := exec.Command("docker", "compose", "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ComposeUp() error {
	cmd := exec.Command("docker", "compose", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ComposeStop() error {
	cmd := exec.Command("docker", "compose", "stop")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ComposeRestart() error {
	cmd := exec.Command("docker", "compose", "restart")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ComposeDown(removeVolumes bool) error {
	args := []string{"compose", "down"}
	if removeVolumes {
		args = append(args, "-v")
	}
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RemoveImages() error {
	images := []string{
		"ghcr.io/usekaneo/api:latest",
		"ghcr.io/usekaneo/web:latest",
		"postgres:16-alpine",
		"caddy:2-alpine",
	}

	var errors []string
	for _, image := range images {
		cmd := exec.Command("docker", "rmi", image)
		if err := cmd.Run(); err != nil {
			errors = append(errors, fmt.Sprintf("failed to remove %s: %v", image, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "; "))
	}

	return nil
}
