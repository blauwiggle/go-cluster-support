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
	toolName := kubeconfig.PromptUser("Select tool:", []string{"octant", "k9s"})
	err = kubeconfig.ConnectToCluster(toolName, kubeconfigPath, clusterName)
	if err != nil {
		fmt.Println("Error connecting to cluster:", err)
		os.Exit(1)
	}
}
