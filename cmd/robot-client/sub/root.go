package sub

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/kooiot/robot/client"
	"github.com/kooiot/robot/client/config"
	"github.com/kooiot/robot/pkg/util/log"
	"github.com/kooiot/robot/pkg/util/version"

	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	workDir     string
	showVersion bool

	serverAddr string
	clientID   string
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
	cmd.PersistentFlags().StringVarP(&serverAddr, "server_addr", "s", "", "robot server's address")
	cmd.PersistentFlags().StringVarP(&clientID, "client_id", "", "", "robot client id")
	cmd.PersistentFlags().StringVarP(&logLevel, "log_level", "", "", "log level")
	cmd.PersistentFlags().StringVarP(&logLink, "log_link", "", "", "latest log file link")
	cmd.PersistentFlags().StringVarP(&logDir, "log_dir", "", "", "log file folder")
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
	cfg, err := config.ParseClientConfig(cfgFilePath)
	if err != nil {
		return err
	}
	if len(serverAddr) > 0 {
		ipStr, portStr, err := net.SplitHostPort(serverAddr)
		if err != nil {
			err = fmt.Errorf("invalid server_addr: %v", err)
			return err
		}

		cfg.Common.Addr = ipStr
		cfg.Common.Port, err = strconv.Atoi(portStr)
		if err != nil {
			err = fmt.Errorf("invalid server_addr: %v", err)
			return err
		}
	}
	if len(clientID) > 0 {
		cfg.Common.ClientID = clientID
	}
	if len(logDir) > 0 {
		cfg.Log.Dir = logDir
	}
	if len(logLink) > 0 {
		cfg.Log.Link = logLink
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

func startService(cfg config.ClientConf, cfgFile string) (err error) {
	log.InitLog(cfg.Log)

	svr := client.NewService(&cfg)

	err = svr.Run()

	return
}
