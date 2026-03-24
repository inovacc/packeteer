package cmd

import (

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "packeteer",
	Short: "MCP server for Wireshark CLI tools — packet capture and network analysis for AI",
	Long: `Packeteer is an MCP (Model Context Protocol) server that wraps Wireshark's CLI tools
(tshark, capinfos, editcap, mergecap) giving AI assistants the ability to capture,
analyze, and dissect network traffic.

Run 'packeteer serve' to start the MCP server on stdio transport.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {


	rootCmd.Version = GetVersionJSON()
	rootCmd.CompletionOptions.DisableDefaultCmd = true


	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

}

