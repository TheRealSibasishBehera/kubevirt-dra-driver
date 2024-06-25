package main

import (
	"fmt"
	nasclient "kubevirt.io/kubevirt-dra-driver/api/kubevirt.io/resource/pci/nas/v1alpha1/client"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"k8s.io/client-go/util/retry"

	nascrd "kubevirt.io/kubevirt-dra-driver/api/kubevirt.io/resource/pci/nas/v1alpha1"
	"kubevirt.io/kubevirt-dra-driver/pkg/flags"
)

func main() {
	if err := newApp().Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newApp() *cli.App {
	var (
		status string

		kubeClientConfig flags.KubeClientConfig
		loggingConfig    *flags.LoggingConfig
		nasConfig        flags.NasConfig
	)

	loggingConfig = flags.NewLoggingConfig()

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "status",
			Usage:    "The status to set [Ready | NotReady].",
			Required: true,
			Action: func(_ *cli.Context, value string) error {
				switch strings.ToLower(value) {
				case strings.ToLower(nascrd.NodeAllocationStateStatusReady):
					status = nascrd.NodeAllocationStateStatusReady
				case strings.ToLower(nascrd.NodeAllocationStateStatusNotReady):
					status = nascrd.NodeAllocationStateStatusNotReady
				default:
					return fmt.Errorf("unknown status: %s", value)
				}
				return nil
			},
			EnvVars: []string{"STATUS"},
		},
	}

	flags = append(flags, kubeClientConfig.Flags()...)
	flags = append(flags, loggingConfig.Flags()...)
	flags = append(flags, nasConfig.Flags()...)

	app := &cli.App{
		Name:            "set-nas-status",
		Usage:           "set-nas-status sets the status of the NodeAllocationState CRD managed by the DRA driver for GPUs.",
		ArgsUsage:       " ",
		HideHelpCommand: true,
		Flags:           flags,
		Before: func(ctx *cli.Context) error {
			if ctx.Args().Len() > 0 {
				return fmt.Errorf("arguments not supported: %v", ctx.Args().Slice())
			}
			return loggingConfig.Apply()
		},
		Action: func(c *cli.Context) error {
			ctx := c.Context
			clientSets, err := kubeClientConfig.NewClientSets()
			if err != nil {
				return fmt.Errorf("create client: %v", err)
			}

			nascr, err := nasConfig.NewNodeAllocationState(ctx, clientSets.Core)
			if err != nil {
				return fmt.Errorf("create NodeAllocationState CR: %v", err)
			}

			client := nasclient.New(nascr, clientSets.Example.NasV1alpha1())
			if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
				err := client.GetOrCreate(ctx)
				if err != nil {
					return err
				}
				return client.UpdateStatus(ctx, status)
			}); err != nil {
				return err
			}
			return nil
		},
	}

	return app
}
