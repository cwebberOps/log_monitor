package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cwebberOps/log_monitor/pkg"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "log_monitor <LOGFILE>",
	Short: "Monitor W3C Common Formated Log File",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting Monitoring of:", args[0])
		config := pkg.Config{
			IntervalDuration:   viper.GetString("interval"),
			TrafficThreshold:   float64(viper.GetInt("threshold")),
			RollingAvgDuration: viper.GetString("average"),
			DbPath:             viper.GetString("dbpath"),
			TopCount:           viper.GetInt64("count"),
		}
		pkg.Run(args[0], config)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().StringP("interval", "i", "10s", "Interval between reports. (ex: 10s, 1m, 15s)")
	viper.BindPFlag("interval", rootCmd.Flags().Lookup("interval"))

	rootCmd.Flags().StringP("average", "a", "2m", "Amount of time used to calclulate the rolling average. (ex: 2m, 1m, 30s)")
	viper.BindPFlag("average", rootCmd.Flags().Lookup("average"))

	rootCmd.Flags().String("dbpath", "/tmp/log.db", "Path to sqlite database used for managing state.")
	viper.BindPFlag("dbpath", rootCmd.Flags().Lookup("dbpath"))

	rootCmd.Flags().IntP("threshold", "t", 5, "Threshold the rolling average is compared against for alerting.")
	viper.BindPFlag("threshold", rootCmd.Flags().Lookup("threshold"))

	rootCmd.Flags().IntP("count", "c", 5, "Display the top N results in last interval traffic.")
	viper.BindPFlag("count", rootCmd.Flags().Lookup("count"))
}
