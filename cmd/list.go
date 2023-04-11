/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"os"
	"path/filepath"
	"sort"
)

var Target string
var m map[string]*ExtensionInfo

type ExtensionInfo struct {
	name     string
	fileNum  int
	fileSize int64
}

type By func(ext1, ext2 *ExtensionInfo) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(exts []*ExtensionInfo) {
	ps := &ExtensionSorter{
		exts: exts,
		by:   by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// planetSorter joins a By function and a slice of Planets to be sorted.
type ExtensionSorter struct {
	exts []*ExtensionInfo
	by   func(pext1, ext2 *ExtensionInfo) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *ExtensionSorter) Len() int {
	return len(s.exts)
}

// Swap is part of sort.Interface.
func (s *ExtensionSorter) Swap(i, j int) {
	s.exts[i], s.exts[j] = s.exts[j], s.exts[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *ExtensionSorter) Less(i, j int) bool {
	return s.by(s.exts[i], s.exts[j])
}

func addFile(m map[string]*ExtensionInfo, ext string, fileInfo os.FileInfo) {
	if m[ext] == nil {
		m[ext] = &ExtensionInfo{ext, 1, fileInfo.Size()}
	} else {
		m[ext].fileNum += 1
		m[ext].fileSize += fileInfo.Size()
	}
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		size := func(ext1, ext2 *ExtensionInfo) bool {
			return ext1.fileSize > ext2.fileSize
		}

		filepath.Walk(Target,
			func(path string, fileInfo os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if fileInfo.IsDir() {

				} else {
					var ext = filepath.Ext(path)
					addFile(m, ext, fileInfo)
				}
				return nil
			})

		var exts []*ExtensionInfo
		for _, v := range m {
			exts = append(exts, v)
		}
		By(size).Sort(exts)
		for _, val := range exts {
			fmt.Printf("%-20s%-20d%-20d\t\n", val.name, val.fileNum, val.fileSize)
		}
	},
}

func init() {
	m = make(map[string]*ExtensionInfo)
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&Target, "target", "t", "", "Target directory to list")
	listCmd.MarkFlagRequired("target")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
