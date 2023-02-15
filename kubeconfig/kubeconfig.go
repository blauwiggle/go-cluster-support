package kubeconfig

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func GetKubeConfig() (string, error, string) {
	stage, stageCorrected := selectStage()

	currentSubscription, err := exec.Command("az", "account", "show", "--query", "name", "-o", "tsv").Output()
	if err != nil {
		fmt.Println("Error getting current subscription:", err)
		os.Exit(1)
	}
	if strings.Contains(strings.ToLower(string(currentSubscription)), stageCorrected) {
		fmt.Printf("You are already in the %s subscription\n", strings.TrimSpace(string(currentSubscription)))
	} else {

		fmt.Println("You are not in the correct subscription. Please log in.")

		err := removeTokens()
		if err != nil {
			return "", err, ""
		}

		err = login()
		if err != nil {
			return "", err, ""
		}

		currentSubscription, err := getCurrentSubscription()
		if err != nil {
			return "", err, ""
		}

		if strings.Contains(strings.ToLower(currentSubscription), stageCorrected) {
			fmt.Printf("You are now in the %s subscription\n", strings.TrimSpace(currentSubscription))
		} else {
			fmt.Printf("You are not in the %s subscription\n", strings.TrimSpace(stage))
			os.Exit(1)
		}
	}

	color := PromptUser("Select color:", []string{"blue", "green"})
	resourceGroup := fmt.Sprintf("rg-cats-%s-aks-%s", stage, color)
	clusterName := fmt.Sprintf("aks-cats-westeurope-%s-%s", stage, color)
	kubeConfig := fmt.Sprintf("%s/.kube/cats-%s-%s.yml", os.Getenv("HOME"), stage, color)
	cmd := fmt.Sprintf("az aks get-credentials --resource-group %s --name %s --file %s", resourceGroup, clusterName, kubeConfig)
	fmt.Println("Executing command:", cmd)

	return kubeConfig, nil, clusterName
}

func selectStage() (string, string) {
	stage := PromptUser("Select stage:", []string{"dev", "prod"})
	stageCorrected := ""
	switch strings.ToLower(stage) {
	case "dev":
		stageCorrected = "dev"
	case "prod":
		stageCorrected = "prd"
	default:
		fmt.Println("Invalid stage:", stage)
		os.Exit(1)
	}
	return stage, stageCorrected
}

func executeCommand(cmd *exec.Cmd) ([]byte, error) {
	outPipe, _ := cmd.StdoutPipe()
	errPipe, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	output, _ := io.ReadAll(outPipe)
	errOutput, _ := io.ReadAll(errPipe)

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("%s\n%s", err, string(errOutput))
	}

	return output, nil
}

func removeTokens() error {
	cmd := exec.Command("kubelogin", "remove-tokens")
	output, err := executeCommand(cmd)
	if err != nil {
		return fmt.Errorf("error removing tokens: %v", string(output))
	}
	return nil
}

func login() error {
	cmd := exec.Command("az", "login", "--tenant", os.Getenv("tenantId"))
	output, err := executeCommand(cmd)
	if err != nil {
		return fmt.Errorf("Error while executing 'az login': %s\n%s", err, string(output))
	}
	fmt.Println("Output off 'az login': ", string(output))
	return nil
}

func getCurrentSubscription() (string, error) {
	output, err := exec.Command("az", "account", "show", "--query", "name", "-o", "tsv").Output()
	if err != nil {
		return "", fmt.Errorf("error getting current subscription: %s", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func PromptUser(message string, options []string) string {
	fmt.Println(message)
	for i, option := range options {
		fmt.Printf("%d. %s\n", i+1, option)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			os.Exit(1)
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		index, err := strconv.Atoi(input)
		if err != nil || index < 1 || index > len(options) {
			fmt.Println("Invalid choice")
			continue
		}
		return options[index-1]
	}
}

func ConnectToCluster(toolName string, kubeconfigPath string, clusterName string) error {
	// Set the context to the correct one
	if err := SetContext(clusterName, kubeconfigPath); err != nil {
		return fmt.Errorf("failed to set context: %v", err)
	}

	cmd := exec.Command(toolName, "--kubeconfig", kubeconfigPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func SetContext(context string, kubeconfigPath string) error {
	cmdArgs := []string{"config", "use-context", context}
	if kubeconfigPath != "" {
		cmdArgs = append(cmdArgs, "--kubeconfig", kubeconfigPath)
	}

	cmd := exec.Command("kubectl", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
