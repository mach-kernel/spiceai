package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spiceai/spice/pkg/context"
	"github.com/spiceai/spice/pkg/pods"
	"github.com/spiceai/spice/pkg/spec"
	"github.com/spiceai/spice/pkg/util"
	"gopkg.in/yaml.v2"
)

var actionCmd = &cobra.Command{
	Use:   "action",
	Short: "Modify actions",
	Example: `
spice action add jump
`,
}

var actionAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an Action to the Pod manifest",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		context.SetContext(context.BareMetal)
		cmdActionName := args[0]

		podPath := pods.FindFirstManifestPath()
		if podPath == "" {
			fmt.Println("No pods detected!")
			return
		}

		pod, err := pods.LoadPodFromManifest(podPath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		actions := pod.Actions()

		if _, ok := actions[cmdActionName]; ok {
			fmt.Printf("Action %s already exists in %s. Overwrite? (y/n)\n", cmdActionName, pod.Name)
			var confirm string
			fmt.Scanf("%s", &confirm)
			if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
				return
			}
		}

		pod.PodSpec.Actions = append(pod.PodSpec.Actions, spec.PodActionSpec{Name: cmdActionName})

		marshalledPod, err := yaml.Marshal(pod.PodSpec)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = util.WriteToExistingFile(podPath, marshalledPod)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("Action '%s' added to pod %s.\n", cmdActionName, pod.Name)
	},
	Example: `
spice action add jump	
`,
}

func init() {
	actionCmd.AddCommand(actionAddCmd)
	actionCmd.Flags().BoolP("help", "h", false, "Print this help message")
	RootCmd.AddCommand(actionCmd)
}