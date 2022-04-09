package sub

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/kooiot/robot/pkg/util/log"
	"github.com/kooiot/robot/pkg/util/version"
	"github.com/kooiot/robot/server"
	"github.com/kooiot/robot/server/api"
	"github.com/kooiot/robot/server/config"

	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	workDir     string
	showVersion bool

	bindAddr string
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
	cmd.PersistentFlags().StringVarP(&bindAddr, "bind_addr", "s", "", "server bind address")
	cmd.PersistentFlags().StringVarP(&logLevel, "log_level", "", "", "log level")
	cmd.PersistentFlags().StringVarP(&logName, "log_name", "", "", "log file name")
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
	if len(bindAddr) > 0 {
		ipStr, portStr, err := net.SplitHostPort(bindAddr)
		if err != nil {
			err = fmt.Errorf("invalid bind_addr: %v", err)
			return err
		}

		cfg.Common.Bind = ipStr
		cfg.Common.Port, err = strconv.Atoi(portStr)
		if err != nil {
			err = fmt.Errorf("invalid bind_addr: %v", err)
			return err
		}
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

func runRobotServer(cfg *config.ServerConf) error {
	svr := server.NewServer(cfg, cfgFile)

	err := svr.Init()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	err = svr.Run()
	log.Error(err.Error())
	return err
}

func startService(cfg config.ServerConf, cfgFile string) (err error) {
	log.InitLog(cfg.Log)

	if !cfg.Api.Enable {
		return runRobotServer(&cfg)
	} else {
		go func() {
			err := runRobotServer(&cfg)
			log.Error(err.Error())
		}()

		// Run HTTP API Server
		err = api.RunServer(&cfg.Api)
		log.Error(err.Error())
		return err
	}
}
