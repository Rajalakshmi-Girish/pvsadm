package cmd

import (
	goflag "flag"
	"os"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"k8s.io/klog/v2"

	"github.com/ppc64le-cloud/pvsadm/cmd/get"
	"github.com/ppc64le-cloud/pvsadm/cmd/image"
	"github.com/ppc64le-cloud/pvsadm/cmd/purge"
	"github.com/ppc64le-cloud/pvsadm/cmd/version"
	"github.com/ppc64le-cloud/pvsadm/pkg"
	"github.com/ppc64le-cloud/pvsadm/pkg/audit"
)

var rootCmd = &cobra.Command{
	Use:   "pvsadm",
	Short: "pvsadm is a command for managing powervs infra",
	Long: `Power Systems Virtual Server projects deliver flexible compute capacity for Power Systems workloads.
Integrated with the IBM Cloud platform for on-demand provisioning.

This is a tool built for the Power Systems Virtual Server helps managing and maintaining the resources easily`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if pkg.Options.APIKey == "" {
			if key := os.Getenv("IBMCLOUD_API_KEY"); key != "" {
				klog.Infof("Using an API key from IBMCLOUD_API_KEY environment variable")
				pkg.Options.APIKey = key
			}
		}
		return nil
	},
}

func init() {
	// Initilize the klog flags
	klog.InitFlags(nil)
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	rootCmd.AddCommand(purge.Cmd)
	rootCmd.AddCommand(get.Cmd)
	rootCmd.AddCommand(version.Cmd)
	rootCmd.AddCommand(image.Cmd)
	rootCmd.PersistentFlags().StringVarP(&pkg.Options.APIKey, "api-key", "k", "", "IBMCLOUD API Key(env name: IBMCLOUD_API_KEY)")
	rootCmd.PersistentFlags().BoolVar(&pkg.Options.Debug, "debug", false, "Enable PowerVS debug option(ATTENTION: dev only option, may print sensitive data from APIs)")
	rootCmd.PersistentFlags().StringVar(&pkg.Options.AuditFile, "audit-file", "pvsadm.log", "Audit logs for the tool")
	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false
	_ = rootCmd.Flags().MarkHidden("debug")

	// Hide the --audit-file for the image subcommand
	// TODO: Remove this after adding audit support to image subcommand
	origHelpFunc := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd.Name() == "image" || (cmd.Parent() != nil && cmd.Parent().Name() == "image") {
			cmd.Flags().MarkHidden("audit-file")
		}
		origHelpFunc(cmd, args)
	})

	audit.Logger = audit.New(pkg.Options.AuditFile)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		klog.Errorln(err)
		os.Exit(1)
	}
}
