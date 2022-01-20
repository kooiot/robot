package sub

import (
	"fmt"
	"os"

	"github.com/kooiot/robot/pkg/util/log"
	"github.com/kooiot/robot/pkg/util/version"
	"github.com/kooiot/robot/server"
	"github.com/kooiot/robot/server/config"

	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	workDir     string
	showVersion bool

	serverAddr string
	user       string
	protocol   string
	token      string
	logLevel   string
	logLink    string
	logDir     string
)

func init() {
	RegisterCommonFlags(rootCmd)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./config.yaml", "config file of robot-client")
	rootCmd.PersistentFlags().StringVarP(&workDir, "work_dir", "d", "", "work directory")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version of robot-client")
}

func RegisterCommonFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&serverAddr, "server_addr", "s", "127.0.0.1:7000", "robot server's address")
	cmd.PersistentFlags().StringVarP(&protocol, "protocol", "p", "tcp", "tcp or kcp")
	cmd.PersistentFlags().StringVarP(&user, "user", "u", "", "user")
	cmd.PersistentFlags().StringVarP(&token, "token", "t", "", "auth token")
	cmd.PersistentFlags().StringVarP(&logLevel, "log_level", "", "info", "log level")
	cmd.PersistentFlags().StringVarP(&logLink, "log_link", "", "latest_log", "latest log file link")
	cmd.PersistentFlags().StringVarP(&logDir, "log_dir", "", "log", "log file folder")
}

var rootCmd = &cobra.Command{
	Use:   "robot-client",
	Short: "robot-client is the client of robot (https://github.com/kooiot/robot)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(version.Full())
			return nil
		}
		if len(workDir) > 0 {
			os.Chdir(workDir)
		}

		// Do not show command usage here.
		err := runServer(cfgFile)
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

func runServer(cfgFilePath string) error {
	cfg, err := config.ParseServerConfig(cfgFilePath)
	if err != nil {
		return err
	}
	return startService(cfg, cfgFilePath)
}

func startService(cfg config.ServerConf, cfgFile string) (err error) {
	log.InitLog(cfg.Log)

	svr := server.NewServer(&cfg, cfgFile)

	err = svr.Init()
	if err != nil {
		return
	}

	err = svr.Run()

	return
}
