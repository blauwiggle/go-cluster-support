package main

import (
	"flag"
	"fmt"
	"github.com/blauwiggle/go-cluster-support/kubeconfig"
	"github.com/blauwiggle/go-cluster-support/tools"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var toolsFlag bool

func init() {

	err := tools.CheckPrerequisites()
	if err != nil {
		fmt.Println(err)
		toolName := tools.GetToolName()
		tools.InstallTool(toolName)
	}

	flag.BoolVar(&toolsFlag, "tools", false, "Use this flag to install tools")
	flag.Parse()
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	loadEnv()
	if toolsFlag {
		toolName := tools.GetToolName()
		tools.InstallTool(toolName)
		return
	}

	kubeconfigPath, err, clusterName := kubeconfig.GetKubeConfig()

	availableTools := []string{"kubectl"}

	if tools.ShowK9sOption() {
		availableTools = append(availableTools, "k9s")
	} else {
		fmt.Println("k9s is not installed. You can install it with the --tools flag, e.g. ./bin/go-cluster-support --tools")
	}

	if tools.ShowOctantOption() {
		availableTools = append(availableTools, "octant")
	} else {
		fmt.Println("octant is not installed. You can install it with the --tools flag, e.g. ./bin/go-cluster-support --tools")
	}

	toolName := kubeconfig.PromptUser("Select tool:", availableTools)

	err = kubeconfig.ConnectToCluster(toolName, kubeconfigPath, clusterName)
	if err != nil {
		fmt.Println("Error connecting to cluster:", err)
		os.Exit(1)
	}
}
