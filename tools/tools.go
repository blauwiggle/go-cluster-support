package tools

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// GetToolName returns the name of the tool selected by the user
func GetToolName() string {
	toolName := *flag.String("tool", "", "Specify the tool you want to use")
	flag.Parse()

	if toolName != "" {
		return toolName
	}

	options := []string{"kubectl", "k9s", "octant", "Azure/kubelogin/kubelogin"}
	fmt.Println("Choose a tool:")
	for i, option := range options {
		fmt.Printf("%d. %s\n", i+1, option)
	}

	choice, err := readChoice()
	if err != nil {
		fmt.Println("Invalid choice:", err)
		os.Exit(1)
	}

	index, err := strconv.Atoi(strings.TrimSpace(choice))
	if err != nil || index < 1 || index > len(options) {
		fmt.Println("Invalid choice")
		os.Exit(1)
	}

	return options[index-1]
}

// InstallTool installs the tool with the given name
func InstallTool(toolName string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("brew", "install", toolName)
	case "linux":
		cmd = exec.Command("sudo", "apt", "install", "-y", toolName)
	case "windows":
		cmd = exec.Command("powershell", "choco", "install", "-y", toolName)
	default:
		fmt.Println("Unsupported operating system")
		os.Exit(1)
	}

	fmt.Printf("Installing %s...\n", toolName)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error installing tool:", err)
		os.Exit(1)
	}
	fmt.Printf("%s successfully installed.\n", toolName)
}

// readChoice reads the user's choice from the command line
func readChoice() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return choice, nil
}
