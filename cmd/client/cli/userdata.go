package cli

import (
	"strings"

	"github.com/kvvPro/gophkeeper/cmd/client/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	dataLogin    string
	dataPassword string
	dataMeta     []string
	meta         map[string]string
)

func init() {
	// put userdata
	putUserData.Flags().StringVarP(&dataLogin, "userdata-login", "u", "", "login")
	viper.BindPFlag("userdata-login", putUserData.Flags().Lookup("userdata-login"))

	putUserData.Flags().StringVarP(&dataPassword, "userdata-password", "p", "", "password")
	viper.BindPFlag("userdata-password", putUserData.Flags().Lookup("userdata-password"))

	putUserData.Flags().StringArrayVarP(&dataMeta, "meta", "m", nil, "meta info")

	meta = make(map[string]string, 0)
	for _, el := range dataMeta {
		a := strings.Split(el, "=")
		if len(a) == 2 {
			meta[a[0]] = a[1]
		} else {
			app.Sugar.Fatalf("invalid metadata - %v", el)
			return
		}
	}

	putUserData.MarkFlagsRequiredTogether("userdata-login", "userdata-password")

	writeCmd.AddCommand(putUserData)

	// get userdata
	writeCmd.AddCommand(getUserData)
}

var putUserData = &cobra.Command{
	Use:   "userdata [key] [options]",
	Short: "update userdata",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	Run:   updateUserDataOnServer,
}

var getUserData = &cobra.Command{
	Use:   "userdata [key]",
	Short: "get userdata by key",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	Run:   getUserDataOnServer,
}

func updateUserDataOnServer(cmd *cobra.Command, args []string) {
	client.WriteUserData(cmd.Context(), args[0], dataLogin, dataPassword, meta)
}

func getUserDataOnServer(cmd *cobra.Command, args []string) {
	var err error
	dataLogin, dataPassword, meta, err = client.GetUserData(cmd.Context(), args[0])
	if err != nil {
		// output
		app.Sugar.Infof("\nKey=%v", args[0])
		app.Sugar.Infof("\nLogin=%v", dataLogin)
		app.Sugar.Infof("\nPassword=%v", dataPassword)
		app.Sugar.Infof("\nMeta info=%v", meta)
	}
}
