package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hvs-fasya/dd/internal/copier"
)

var (
	inFile  string
	outFile string
	Params  copier.CopyParams

	ErrNotRegularFile = errors.New("source file is not regular file")
	ErrLimitOffset    = errors.New("source file size is less than offset + limit")

	rootCmd = &cobra.Command{
		Use:   "./go-dd --from=in_file --to=out_file",
		Short: `./go-dd - copy a file`,
		Long:  `./go-dd - copy a file; show progress`,
		RunE: func(cmd *cobra.Command, args []string) error {
			src, err := os.Open(inFile)
			if err != nil {
				return err
			}
			defer src.Close()
			if mode, _ := src.Stat(); !mode.Mode().IsRegular() {
				return ErrNotRegularFile
			}
			srcFi, _ := src.Stat() //source file info
			if srcFi.Size() < Params.Offset+Params.Limit {
				return ErrLimitOffset
			}
			if Params.Limit == 0 {
				Params.Limit = srcFi.Size()
			}

			dst, err := os.Create(outFile)
			if err != nil {
				return err
			}
			defer dst.Close()
			src.Seek(Params.Offset, 0)
			bytesCopied, err := copier.CopyBuffered(src, dst, Params)
			if err != nil {
				fmt.Printf("TOTAL bytes copied: %d\n", bytesCopied)
				return fmt.Errorf("error while copying: %s", err)
			}
			fmt.Printf("TOTAL bytes copied: %d\n", bytesCopied)
			return nil
		},
	}
)

//todo: -- -> -
func main() {
	rootCmd.Flags().StringVar(&inFile, "from", "", "File that should be copied")
	rootCmd.MarkFlagRequired("from")
	rootCmd.Flags().StringVar(&outFile, "to", "", "File copy destination")
	rootCmd.MarkFlagRequired("to")
	rootCmd.Flags().Int64Var(&Params.Offset, "offset", 0, "Source file bytes count which should be skipped while copying")
	rootCmd.Flags().Int64Var(&Params.Limit, "limit", 0, "Bytes count to be copied")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
