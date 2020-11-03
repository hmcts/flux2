/*
Copyright 2020 The Flux authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"

	"github.com/fluxcd/flux2/internal/utils"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/types"
)

var deleteSourceHelmCmd = &cobra.Command{
	Use:   "helm [name]",
	Short: "Delete a HelmRepository source",
	Long:  "The delete source helm command deletes the given HelmRepository from the cluster.",
	Example: `  # Delete a Helm repository
  flux delete source helm podinfo
`,
	RunE: deleteSourceHelmCmdRun,
}

func init() {
	deleteSourceCmd.AddCommand(deleteSourceHelmCmd)
}

func deleteSourceHelmCmdRun(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("name is required")
	}
	name := args[0]

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	kubeClient, err := utils.KubeClient(kubeconfig, kubecontext)
	if err != nil {
		return err
	}

	namespacedName := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}

	var helmRepository sourcev1.HelmRepository
	err = kubeClient.Get(ctx, namespacedName, &helmRepository)
	if err != nil {
		return err
	}

	if !deleteSilent {
		prompt := promptui.Prompt{
			Label:     "Are you sure you want to delete this source",
			IsConfirm: true,
		}
		if _, err := prompt.Run(); err != nil {
			return fmt.Errorf("aborting")
		}
	}

	logger.Actionf("deleting source %s in %s namespace", name, namespace)
	err = kubeClient.Delete(ctx, &helmRepository)
	if err != nil {
		return err
	}
	logger.Successf("source deleted")

	return nil
}