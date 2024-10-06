package cmd

import (
	"github.com/hashicorp/memberlist"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	utils "github.com/cedi/gossip-experiment/pkg/utils"
)

var (
	defaultConfigType string
	nodeName          string
	noWait            bool
	nodePort          int
)

func init() {
	memberlistCreateCmd.Flags().StringVarP(&defaultConfigType, "config", "c", "local", "default config type for memberlist")
	memberlistCreateCmd.RegisterFlagCompletionFunc("config", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"local", "lan", "wan"}, cobra.ShellCompDirectiveDefault
	})

	memberlistCreateCmd.Flags().BoolVar(&noWait, "nowait", false, "exit without waiting")
	memberlistCreateCmd.Flags().StringVarP(&nodeName, "name", "n", "", "Node Name")
	memberlistCreateCmd.MarkFlagRequired("name")

	memberlistJoinCmd.Flags().StringVarP(&defaultConfigType, "config", "c", "local", "default config type for memberlist")
	memberlistJoinCmd.RegisterFlagCompletionFunc("config", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"local", "lan", "wan"}, cobra.ShellCompDirectiveDefault
	})

	memberlistJoinCmd.Flags().BoolVar(&noWait, "nowait", false, "exit without waiting")
	memberlistJoinCmd.Flags().StringVarP(&nodeName, "name", "n", "", "Node Name")
	memberlistJoinCmd.MarkFlagRequired("name")

	memberlistJoinCmd.Flags().IntVar(&nodePort, "port", 7947, "Port to bind to")
	memberlistJoinCmd.MarkFlagRequired("port")

	memberlistCmd.AddCommand(memberlistJoinCmd)
	memberlistCmd.AddCommand(memberlistCreateCmd)
	rootCmd.AddCommand(memberlistCmd)
}

var memberlistCmd = &cobra.Command{
	Use:   "memberlist",
	Short: "A CLI to play around with hashicorp memberlist",
}

var memberlistCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a memberlist cluster",
	Example: "gossip memberlist create",
	Args:    cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := utils.GetMemberlistConfig(defaultConfigType)
		if err != nil {
			log.WithError(err).Fatal("Could not aquire config")
		}

		config.Name = nodeName

		list, err := memberlist.Create(config)
		if err != nil {
			log.WithError(err).Fatal("Could not create memerblist cluster")
		}

		local := list.LocalNode()

		log.WithFields(log.Fields{
			"node_name": nodeName,
			"node_port": local.Port,
			"node_ip":   local.Addr.To4().String(),
		}).Printf("Node running - waiting for other nodes to join")

		if !noWait {
			utils.WaitSignal()
		}

		return nil
	},
}

var memberlistJoinCmd = &cobra.Command{
	Use:     "join [nodeIp]",
	Short:   "Join a memberlist cluster",
	Example: "gossip memberlist join 0.0.0.0:8081 ",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		joinNode := ""
		if len(args) == 1 {
			joinNode = args[0]
		}

		config, err := utils.GetMemberlistConfig(defaultConfigType)
		if err != nil {
			log.WithError(err).Fatal("Could not aquire config")
		}

		config.Name = nodeName
		config.BindPort = nodePort // avoid port confliction
		config.AdvertisePort = config.BindPort

		list, err := memberlist.Create(config)
		if err != nil {
			log.WithError(err).Fatal("Could not create memerblist cluster")
		}

		local := list.LocalNode()

		log.WithFields(log.Fields{
			"node_name": nodeName,
			"node_port": local.Port,
			"node_ip":   local.Addr.To4().String(),
		}).Printf("Node up - attempting to reach other nodes")

		if _, err := list.Join([]string{joinNode}); err != nil {
			log.WithError(err).Fatal("Failed to join existing node")
		}

		for _, member := range list.Members() {
			log.WithFields(log.Fields{
				"local_ip":    local.Addr.To4().String(),
				"local_name":  nodeName,
				"local_port":  local.Port,
				"member_name": member.Name,
				"member_ip":   member.Addr.To4().String(),
				"member_port": member.Port,
			}).Printf("Node up - found other member %s!", member.Name)
		}

		if !noWait {
			utils.WaitSignal()
		}

		return nil
	},
}
