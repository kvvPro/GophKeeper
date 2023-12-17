package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	login    string
	password string
)

func init() {
	// reg
	regCmd.Flags().StringVarP(&login, "username", "u", "", "username to log in")
	regCmd.MarkFlagRequired("username")
	viper.BindPFlag("username", regCmd.Flags().Lookup("username"))

	regCmd.Flags().StringVarP(&password, "password", "p", "", "password to log in")
	regCmd.MarkFlagRequired("password")
	viper.BindPFlag("password", regCmd.Flags().Lookup("password"))

	rootCmd.AddCommand(regCmd)

	// auth
	authCmd.Flags().StringVarP(&login, "username", "u", "", "username to log in")
	authCmd.MarkFlagRequired("username")
	viper.BindPFlag("username", authCmd.Flags().Lookup("username"))

	authCmd.Flags().StringVarP(&password, "password", "p", "", "password to log in")
	authCmd.MarkFlagRequired("password")
	viper.BindPFlag("password", authCmd.Flags().Lookup("password"))

	rootCmd.AddCommand(authCmd)
}

var regCmd = &cobra.Command{
	Use:   "register",
	Short: "add new user",
	Long:  "",
	Run:   regOnServer,
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "pass authentication",
	Long:  "",
	Run:   authOnServer,
}

func regOnServer(cmd *cobra.Command, args []string) {
	client.AddUser(cmd.Context(), login, password)
}

func authOnServer(cmd *cobra.Command, args []string) {
	client.Auth(cmd.Context(), login, password)
}
