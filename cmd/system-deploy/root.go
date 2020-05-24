package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ppacher/system-deploy/pkg/actions"
	"github.com/ppacher/system-deploy/pkg/deploy"
	"github.com/ppacher/system-deploy/pkg/runner"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	// plugin all builtin actions
	_ "github.com/ppacher/system-deploy/pkg/actions/builtin"
)

func init() {
	// Register all supported conditions from the
	// condition package.
	deploy.RegisterAllConditions()
}

func getRootCmd() *cobra.Command {
	var dropInSearchPaths []string

	var root = &cobra.Command{
		Use:   "system-deploy",
		Short: "Deploy and manage system configuration",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			var targets []deploy.Task
			for _, dir := range args {
				stat, err := os.Stat(dir)
				if err != nil {
					log.Fatal(err)
				}
				if !stat.IsDir() {
					continue
				}

				files, err := ioutil.ReadDir(dir)
				if err != nil {
					log.Fatal(err)
				}

				for _, fi := range files {
					// we skip directories for now.
					if fi.IsDir() {
						continue
					}

					if filepath.Ext(fi.Name()) != ".task" {
						continue
					}

					path := filepath.Join(dir, fi.Name())
					targets = append(targets, parseFile(path, dropInSearchPaths))
				}
			}

			if len(targets) == 0 {
				log.Fatal("no valid tasks found")
			}

			run, err := runner.NewRunner(actions.NewLogger(), targets)
			if err != nil {
				log.Fatal(err)
			}

			if err := run.Deploy(context.Background()); err != nil {
				log.Fatal(err)
			}
		},
	}

	defaultSearchPath := []string{
		".config", // inside the working directory
		"/etc/system-deploy",
	}
	root.Flags().StringSliceVarP(&dropInSearchPaths, "path", "p", defaultSearchPath, "Search paths for task drop-in files.")

	var logLevel string
	root.PersistentFlags().StringVarP(&logLevel, "log", "l", "info", "Log level")
	cobra.OnInitialize(func() {
		lvl, err := logrus.ParseLevel(logLevel)
		if err != nil {
			log.Fatal(err.Error())
		}

		logrus.SetLevel(lvl)
	})

	root.AddCommand(describe)
	root.AddCommand(runActionCommand)

	return root
}

func parseFile(filePath string, searchPaths []string) deploy.Task {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	target, err := deploy.Decode(filePath, f)
	if err != nil {
		log.Fatalf("Failed to decode target at %s: %s", filePath, err)
	}

	dropins, err := deploy.LoadDropIns(target.FileName, searchPaths)
	if err != nil {
		log.Fatalf("Failed to load drop-in files for unit %s: %s", target.FileName, err)
	}

	// TODO(ppacher): this is ugly, remove it and fix it in ApplyDropIns
	specs, err := actions.TaskSpec(target)
	if err != nil {
		log.Fatalf("Failed to apply dropins to %s: %s", target.FileName, err)
	}

	if err := deploy.ApplyDropIns(target, dropins, specs); err != nil {
		log.Fatalf("Failed to apply dropins to %s: %s", target.FileName, err)
	}

	if err = deploy.ApplyEnvironment(target); err != nil {
		log.Fatalf("Failed to apply environment to task %s: %s", target.FileName, err)
	}
	dump(target.FileName, *target)

	return *target
}

func dump(prefix string, x interface{}) {
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		log.Fatal(err.Error())
	}
	logrus.Debugf("dump %s: \n%s", prefix, string(b))
}
