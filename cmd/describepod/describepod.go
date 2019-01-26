package describepod

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arunvelsriram/kube-fzf/cmd"
	"github.com/arunvelsriram/kube-fzf/pkg/kubectl"
	"github.com/arunvelsriram/kube-fzf/pkg/kubernetes"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var allNamespaces bool
var namespaceName string

const multiSelect = false

var rootCmd = &cobra.Command{
	Use:   "describepod [pod-query]",
	Short: "Describe a pod interactively",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var podName string
		if len(args) == 1 {
			podName = strings.TrimSpace(args[0])
		}

		kubeconfig := viper.GetString("kubeconfig")

		client, err := kubernetes.NewClient(kubeconfig)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if allNamespaces {
			pods, err := client.GetAllPods()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			filteredPod, err := pods.FilterOne(podName, multiSelect)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			kubectl.DescribePod(kubeconfig, filteredPod)
		} else {
			namespaces, err := client.GetNamespaces()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			filteredNamespace, err := namespaces.FilterOne(namespaceName)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			pods, err := client.GetPods(filteredNamespace)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			filteredPod, err := pods.FilterOne(podName, multiSelect)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			kubectl.DescribePod(kubeconfig, filteredPod)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initKubeconfig() {
	if !viper.IsSet("kubeconfig") || viper.GetString("kubeconfig") == "" {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.SetDefault("kubeconfig", filepath.Join(home, ".kube", "config"))
	}
}

func init() {
	cobra.OnInitialize(initKubeconfig)
	rootCmd.AddCommand(cmd.VersionCmd)
	rootCmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "a", false, "consider all namespaces")
	rootCmd.Flags().StringVarP(&namespaceName, "namespace", "n", "default", "namespace query")
	rootCmd.Flags().StringP("kubeconfig", "", "", "path to kubeconfig file (default is $HOME/.kube/config)")
	viper.BindPFlag("kubeconfig", rootCmd.Flags().Lookup("kubeconfig"))
	viper.AutomaticEnv()
}