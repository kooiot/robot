package sub

import (
	"fmt"
	"os"

	"github.com/kooiot/robot/cmd/robot-pong/pong"
	"github.com/kooiot/robot/pkg/util/log"
	"github.com/kooiot/robot/pkg/util/version"

	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	showVersion bool

	logLevel string
	logLink  string
	logDir   string
)

func init() {
	RegisterCommonFlags(rootCmd)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./config.yaml", "config file of robot-client")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version of robot-client")
}

func RegisterCommonFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&logLevel, "log_level", "", "info", "log level")
	cmd.PersistentFlags().StringVarP(&logLink, "log_link", "", "latest_log", "latest log file link")
	cmd.PersistentFlags().StringVarP(&logDir, "log_dir", "", "log", "log file folder")
}

var rootCmd = &cobra.Command{
	Use:   "robot-pong",
	Short: "robot-pong is the pong server for serial (https://github.com/kooiot/robot)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(version.Full())
			return nil
		}

		// Do not show command usage here.
		err := runClient(cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runClient(cfgFilePath string) error {
	cfg, err := pong.ParseServerConfig(cfgFilePath)
	if err != nil {
		return err
	}
	return startService(cfg, cfgFilePath)
}

func startService(cfg pong.PongServerConf, cfg_file string) (err error) {
	log.InitLog(cfg.Log)

	svr := pong.NewService(&cfg)

	err = svr.Run()

	return
}
