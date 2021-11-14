package sub

import (
	"fmt"
	"os"

	"github.com/kooiot/robot/client"
	"github.com/kooiot/robot/client/config"
	"github.com/kooiot/robot/pkg/util/log"
	"github.com/kooiot/robot/pkg/util/version"

	"github.com/spf13/cobra"
)

var (
	cfgFile     string
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

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./config.yaml", "config file of iot-tun-client")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version of iot-tun-client")
}

func RegisterCommonFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&serverAddr, "server_addr", "s", "127.0.0.1:7000", "iot-tun server's address")
	cmd.PersistentFlags().StringVarP(&protocol, "protocol", "p", "tcp", "tcp or kcp")
	cmd.PersistentFlags().StringVarP(&user, "user", "u", "", "user")
	cmd.PersistentFlags().StringVarP(&token, "token", "t", "", "auth token")
	cmd.PersistentFlags().StringVarP(&logLevel, "log_level", "", "info", "log level")
	cmd.PersistentFlags().StringVarP(&logLink, "log_link", "", "latest_log", "latest log file link")
	cmd.PersistentFlags().StringVarP(&logDir, "log_dir", "", "log", "log file folder")
}

var rootCmd = &cobra.Command{
	Use:   "iot-tun-client",
	Short: "iot-tun-client is the client of iot-tun (https://github.com/kooiot/iot-tun)",
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
	cfg, err := config.ParseClientConfig(cfgFilePath)
	if err != nil {
		return err
	}
	return startService(cfg, cfgFilePath)
}

func startService(cfg config.ClientConf, cfgFile string) (err error) {
	log.InitLog(cfg.Log)

	svr := client.NewClient(&cfg)

	err = svr.Run()

	return
}
