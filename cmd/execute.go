package cmd

import (
	"os"
	"strings"

	"github.com/dx-oss/wait-for/pkg/utils"
	"github.com/dx-oss/wait-for/pkg/wfi"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// executeCmd represents the execute command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Wait for service(s)",
	Long: `Check if service(s) is active. When it is true the exit code is 0.	`,
	Run: func(cmd *cobra.Command, args []string) {
		services, _ := cmd.Flags().GetStringArray("services")
		wait, _ := cmd.Flags().GetInt("wait")
		timeout, _ := cmd.Flags().GetInt("timeout")
		sleep, _ := cmd.Flags().GetInt("sleep")

		svc, err := wfi.New(services, wait, timeout, sleep)
		if err != nil {
			log.Debugln(err)
			os.Exit(1)
		}

		err = svc.Execute()
		if err != nil {
			log.Debugln(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)

	executeCmd.Flags().IntP("timeout", "t", utils.GetOsEnvNumberOrDefault("TIMEOUT", wfi.DefaultTimeoutInSeconds), "Total max running time in seconds (TIMEOUT)")
	executeCmd.Flags().StringArrayP("services", "s", getDefaultServices(), "Service(s) (SERVICES - use , to multiple services)")
	executeCmd.Flags().IntP("wait", "w", utils.GetOsEnvNumberOrDefault("WAIT", wfi.DefaultWaitForBeforeStart), "Wait before starting in seconds (WAIT)")
	executeCmd.Flags().Int("sleep", utils.GetOsEnvNumberOrDefault("SLEEP", wfi.DefaultSleepBetweenChecks), "Sleep between execution in seconds (SLEEP)")

	executeCmd.Example = strings.Join([]string{
		"wfi execute -s tcp://127.0.0.1:8080 -s http://google.com -w 30 -t 120 --sleep 5",
		"",
		"Supported envs (TIMEOUT, SERVICES, WAIT, SLEEP)",
	}, "\n")

	//executeCmd.MarkFlagRequired("services")

}

func getDefaultServices() []string {
	rawServices := os.Getenv("SERVICES")
	if len(rawServices) > 0 && strings.Contains(rawServices, ",") {
		return strings.Split(rawServices, ",")
	}

	return []string{rawServices}
}
