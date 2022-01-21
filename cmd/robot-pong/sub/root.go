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
	workDir     string
	showVersion bool

	logLevel string
	logName  string
	logDir   string
)

func init() {
	RegisterCommonFlags(rootCmd)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./config.yaml", "config file of robot-client")
	rootCmd.PersistentFlags().StringVarP(&workDir, "work_dir", "d", "", "work directory")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version of robot-client")
}

func RegisterCommonFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&logLevel, "log_level", "", "", "log level")
	cmd.PersistentFlags().StringVarP(&logName, "log_name", "", "", "log file name")
	cmd.PersistentFlags().StringVarP(&logDir, "log_dir", "", "", "log file folder")
}

var rootCmd = &cobra.Command{
	Use:   "robot-pong",
	Short: "robot-pong is the pong server for serial (https://github.com/kooiot/robot)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(version.Full())
			return nil
		}
		if len(workDir) > 0 {
			os.Chdir(workDir)
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
	if len(logDir) > 0 {
		cfg.Log.Dir = logDir
	}
	if len(logName) > 0 {
		cfg.Log.Filename = logName
	}
	if len(logLevel) > 0 {
		cfg.Log.Level = logLevel
	}

	cfg.Complete()
	if err = cfg.Validate(); err != nil {
		err = fmt.Errorf("parse config error: %v", err)
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
