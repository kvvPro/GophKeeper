package cli

// keeper
//	 register -u=[login] -p=[password]
//	 auth -u=[login] -p=[password]
//   get
//     userdata [key] -o=[output-to-file]
//     text [key] -o=[output-to-file]
//     bin [key] -o=[output-to-file]
//     card [key] -o=[output-to-file]
//	   -a -o=[output-to-file]
//   write
//     userdata [key] -u=[login] -p=[password] -m="[key]=[value]" -m="[key]=[value]"
//     text [key] -v=[value] -m="[key]=[value]" -m="[key]=[value]"
//     bin [key] -v=[value] -f=[path-to-file] -m="[key]=[value]" -m="[key]=[value]"
//     card [key] -n=[number] --pin=[pin] --cvc=[cvc] -m="[key]=[value]" -m="[key]=[value]"
//   version -v

import (
	"strings"

	"github.com/kvvPro/gophkeeper/cmd/client/app"
	"github.com/kvvPro/gophkeeper/cmd/client/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gophkeeper",
		Short: "A manager of data",
		Long:  ``,
	}

	client *app.Client
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	app.Sugar = *logger.Sugar()

	cobra.OnInitialize(func() {
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
		postInitCommands(rootCmd.Commands())
	})

	agentFlags := initConfigs()
	client, err = app.NewClient(agentFlags)
	if err != nil {
		app.Sugar.Fatalw(err.Error())
	}
}

func initConfigs() *config.ClientFlags {
	viper.AutomaticEnv()
	agentFlags, err := config.ReadConfig()
	if err != nil {
		app.Sugar.Fatalw(err.Error(), "event", "read config")
	}
	config.Initialize(agentFlags)
	if err != nil {
		app.Sugar.Fatalw(err.Error(), "event", "read config")
	}
	return agentFlags
}

func postInitCommands(commands []*cobra.Command) {
	for _, cmd := range commands {
		presetRequiredFlags(cmd)
		if cmd.HasSubCommands() {
			postInitCommands(cmd.Commands())
		}
	}
}

func presetRequiredFlags(cmd *cobra.Command) {
	viper.BindPFlags(cmd.Flags())
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			cmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}
