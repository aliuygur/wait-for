package cmd

import (
	"github.com/dx-oss/wait-for/pkg/wfi"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// executeCmd represents the execute command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Wait for service(s)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		services, _ := cmd.Flags().GetStringArray("services")
		timeout, _ := cmd.Flags().GetInt("timeout")
		repeat, _ := cmd.Flags().GetInt("repeat")

		svc, err := wfi.New(services, timeout, repeat)
		if err != nil {
			log.Fatalln(err)
		}

		err = svc.Execute()
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)

	executeCmd.Flags().IntP("timeout", "t", wfi.DefaultTimeoutInSeconds, "Timeout in seconds")
	executeCmd.Flags().StringArrayP("services", "s", []string{}, "Service(s)")
	executeCmd.Flags().IntP("repeat", "r", wfi.DefaultRepeat, "Retries")
	executeCmd.MarkFlagRequired("services")

}
