package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/mitchellh/go-wordwrap"
	"github.com/ppacher/system-deploy/pkg/actions"
	"github.com/ppacher/system-deploy/pkg/deploy"
	"github.com/spf13/cobra"
)

var describe = &cobra.Command{
	Use:   "describe",
	Short: "Display documentation for an action",
	Run: func(_ *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Printf(" - %s\n", strings.Join(actions.ListActions(), "\n - "))
			return
		}

		if len(args) > 1 {
			log.Fatal("only one parameter expected")
		}

		plg, ok := actions.GetPlugin(args[0])
		if !ok {
			log.Fatalf("Action %s does not exist", args[0])
		}

		header := color.New(color.Bold, color.Underline)
		fmt.Printf("[ %s ]\n", header.Sprint(plg.Name))

		if plg.Description != "" {
			fmt.Printf("\n%s\n", wrap(plg.Description, ""))
		}

		for _, section := range plg.Help {
			if section.Title != "" {
				fmt.Printf("\n%s\n", header.Sprint(strings.ToUpper(section.Title)))
			}

			if section.Description != "" {
				fmt.Printf("\n%s\n", wrap(section.Description, ""))
			}
		}

		if !deploy.IsAllowAny(plg.Options) {
			fmt.Printf("\n%s\n\n", header.Sprint("OPTIONS"))

			for _, opt := range plg.Options {
				required := ""
				defaultValue := ""

				if opt.Required {
					required = " (required)"
				}

				if opt.Default != "" {
					defaultValue = fmt.Sprintf(" (Default: %q)", opt.Default)
				}

				fmt.Printf("   %s (%s)\n      %s\n\n",
					color.New(color.Bold).Sprint(opt.Name),
					opt.Type.String(),
					wrap(opt.Description+required+defaultValue, "      "),
				)
			}
		} else {
			fmt.Println("Any options allowed")
		}

		if plg.Example != "" {
			fmt.Printf("\n%s\n", header.Sprint("EXAMPLE"))
			fmt.Printf("\n%s\n", plg.Example)
		}

		if plg.Author != "" || plg.Website != "" {
			fmt.Printf("\n%s\n", header.Sprint("CONTACT"))
			fmt.Printf("\n%s", color.New(color.Underline).Sprint(plg.Author))
			fmt.Printf("\n%s", plg.Website)

			fmt.Println()
		}
	},
}

// wrap ensures that text is no longer thatn 80 characters per line.
// It automatically breaks text into multiple lines that fit into a
// 80 character (including indention) limit.
func wrap(text string, indent string) string {
	lines := strings.Split(wordwrap.WrapString(text, uint(80-len(indent))), "\n")
	return strings.Join(lines, "\n"+indent)
}
