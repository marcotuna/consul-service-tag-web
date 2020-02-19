package commands

import (
	"consul-service-tag-web/app"
	"consul-service-tag-web/model"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/mattermost/mattermost-server/mlog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:          "server",
	Short:        "Run the Web server",
	RunE:         serverCmdF,
	SilenceUsage: true,
}

func init() {
	RootCmd.PersistentFlags().StringP("server.address", "s", ":9999", "Server listening address")
	RootCmd.AddCommand(serverCmd)

	viper.BindPFlag("server.address", RootCmd.PersistentFlags().Lookup("server.address"))

	RootCmd.RunE = serverCmdF
}

func serverCmdF(command *cobra.Command, args []string) error {

	interruptChan := make(chan os.Signal, 1)

	return runServer(interruptChan)
}

func runServer(interruptChan chan os.Signal) error {

	// Load configurations to structure
	config := model.Config{
		HTTP: model.ConfigHTTP{
			Address: viper.GetString("server.address"),
		},
		Tag: model.ConfigTag{
			Filter: viper.GetString("tag.filter"),
		},
	}

	app, err := app.NewApp(&config)

	if err != nil {
		return err
	}

	defer app.Shutdown()

	notifyReady()

	// wait for kill signal before attempting to gracefully shutdown
	// the running service
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan

	return nil
}

func notifyReady() {
	// If the environment vars provide a systemd notification socket,
	// notify systemd that the server is ready.
	systemdSocket := os.Getenv("NOTIFY_SOCKET")
	if systemdSocket != "" {
		mlog.Info("Sending systemd READY notification.")

		err := sendSystemdReadyNotification(systemdSocket)
		if err != nil {
			mlog.Error(err.Error())
		}
	}
}

func sendSystemdReadyNotification(socketPath string) error {
	msg := "READY=1"
	addr := &net.UnixAddr{
		Name: socketPath,
		Net:  "unixgram",
	}
	conn, err := net.DialUnix(addr.Net, nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write([]byte(msg))
	return err
}
