/*
Copyright © 2020 Guillaume Pancak <gpkfr@imelbox.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/ecs"
	gtdAWS "github.com/gpkfr/gtd-ng/aws"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show statut of ECS service(s) described in 'gtd/env.yaml' ",
	Long: `show the current status of ECS services. For example:

gtd-ng show -e sleep360.`,
	Run: func(cmd *cobra.Command, args []string) {
		show()
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	showCmd.Flags().StringVarP(&gtdenv, "env", "e", "sleep360", "Environment to show")
	showCmd.Flags().StringSliceVarP(&selectedService, "service", "s", []string{}, "Selected service")
}

func show() {

	var services gtdAWS.Services

	err := gtdAWS.LoadService(&services, &gtdenv)
	if err != nil {
		log.Fatal(err)
	}

	//Da AWS Stuff
	awsSession, err := gtdAWS.NewAWSSession(&services.ECSRegion, &awsProfile)
	if err != nil {
		log.Fatal(err)
	}

	svc := ecs.New(awsSession)

	if len(selectedService) <= 0 {
		err = gtdAWS.GetServiceTask(services, svc)
		if err != nil {
			log.Fatal(err)
		}
	} else {

		err = gtdAWS.GetServiceTask(services, svc, selectedService...)
		if err != nil {
			log.Fatal(err)
		}
	}

	for i, s := range services.Services {
		if s.TaskARN != "" {
			currentTask, err := gtdAWS.GetCurrentServiceTaskDefinition(svc, s.TaskARN)
			if err != nil {
				log.Println(err)
			}
			services.Services[i].TaskDefinition = currentTask.TaskDefinition
		}
	}
	fmt.Printf("Current Services Status\n")
	// loop under services
	for _, aService := range services.Services {
		if aService.TaskDefinition != nil {
			fmt.Printf("Services: %s - TaskARN: %s\nFamily: %s\n", aService.Name, aService.TaskARN, *aService.TaskDefinition.Family)
			fmt.Printf("Revision: %d\nActual Image: %s\nStatus: %s\nRunning Count: %d\n",
				*aService.TaskDefinition.Revision,
				*aService.TaskDefinition.ContainerDefinitions[0].Image,
				aService.Status,
				aService.RunningCount)
		}
	}
}
