/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"math/big"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var capacityBytes int64
var word string
var metadataSizeBytes int64

const byteSizeInBits int64 = 8

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rc-mega",
	Short: "Silly application that calculates how many words will fit in a given size",
	Long: `Exceedingly silly application that calculates how many times a word
will fit in a given size, using a special compression algorithm.

This algorithm assumes the size required for N copies of word is:
1KiB (for metadata) + the size of the word + the number N encoded as binary

Another way of putting this is that, for a size S bits you can fit a word of
size W:
S - (8 * (1024 + W) )
times`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func (cmd *cobra.Command, args []string){
	// calculate number of bits in capacity
	var capacityBits = capacityBytes * byteSizeInBits
	// calculate size of word in bytes
	var wordSize int64 = int64(len(word))
	// calculate word size in bits
	var wordSizeBits = wordSize * byteSizeInBits
	// calculate metadata size in bits
	var metadataSizeBits = metadataSizeBytes * byteSizeInBits

	// calculate remaining bits after metadata and word
	var availableBits int64 = capacityBits - (wordSizeBits + metadataSizeBits)
  // big integer variant
  var availableBitsBig = big.NewInt(availableBits)

	// calculate max number that will fit in availableBits
	var numWords *big.Int = big.NewInt(2)
  numWords.Exp(numWords, availableBitsBig, nil)

	// print output
	fmt.Printf("Guess what! The number of times \"%v\" will fit into %v bytes is: %v\n", word, capacityBytes,
	numWords.Text(10))
},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().Int64VarP(&capacityBytes, "capacity", "c", 16 * 1024,
		"number of bytes to fit word into (default is 16 * 1024 for 16KiB)")
	rootCmd.PersistentFlags().StringVarP(&word, "word", "w", "sophie", "word to fit into capacity (default is sophie)")
	rootCmd.PersistentFlags().Int64VarP(&metadataSizeBytes, "metadata-size", "m", 1024,
		"metadata size in bytes (default is 1024 for 1KiB)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".rc-mega" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".rc-mega")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func Power(a, n uint64) uint64 {
	var i, result uint64
	result = 1
	for i = 0; i < n; i++ {
		result *= a
	}
	return result
}