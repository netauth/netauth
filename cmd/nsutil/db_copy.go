package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/netauth/netauth/internal/db"
	_ "github.com/netauth/netauth/internal/db/bitcask"
	_ "github.com/netauth/netauth/internal/db/filesystem"

	"github.com/netauth/netauth/internal/startup"
)

var (
	dbCopyCmd = &cobra.Command{
		Use:   "copy <source> <target>",
		Short: "Copy an existing datastore to a new one",
		Long:  dbCopyCmdLongDocs,
		Run:   dbCopyCmdRun,
		Args:  cobra.ExactArgs(2),
	}

	dbCopyCmdLongDocs = `
The copy command allows you to copy all data from one backend to
another.  This is useful for when you want to migrate from a
host-local storage option to something distributed, or migrate into or
out of a format that can be passed to other external tools.
`

	dbCopyCmdNoDryRun bool
	dbCopyCmdTruncate bool
)

func init() {
	dbCopyCmd.Flags().BoolVar(&dbCopyCmdNoDryRun, "no-dry-run", false, "Make changes, potentially destructive.")
	dbCopyCmd.Flags().BoolVar(&dbCopyCmdTruncate, "truncate", false, "Truncate target storage.")

	rootCmd.AddCommand(dbCopyCmd)
}

func dbCopyCmdRun(c *cobra.Command, args []string) {
	startup.DoCallbacks()

	source, err := db.NewKV(args[0], nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing source: %s\n", err)
		os.Exit(1)
	}
	source.SetEventFunc(func(db.Event) {})
	sourceKeys, err := source.Keys("/*/*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving keys from source: %s\n", err)
	}
	fmt.Printf("Source initialized, contains %d objects.\n", len(sourceKeys))

	target, err := db.NewKV(args[1], nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing target: %s\n", err)
		os.Exit(1)
	}
	target.SetEventFunc(func(db.Event) {})
	targetKeys, err := target.Keys("/*/*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving keys from target: %s\n", err)
	}
	fmt.Printf("Target initialized, contains %d objects.\n", len(targetKeys))

	// Bail out at this point if we're in a dry-run, otherwise continue
	fmt.Printf("All data will be copied from %s to %s.\n", args[0], args[1])
	if dbCopyCmdTruncate {
		fmt.Println("Target database will be truncated.")
	}
	if !dbCopyCmdNoDryRun {
		fmt.Println("You are in dry-run mode, pass --no-dry-run to make changes described above.")
		os.Exit(0)
	}

	// If we're truncating then go ahead and do that.
	if dbCopyCmdTruncate {
		for _, k := range targetKeys {
			if err := target.Del(k); err != nil {
				fmt.Fprintf(os.Stderr, "Error truncating target: %s\n", err)
			}
		}
	}

	// Copy  over all the  keys, maintain a  copy of the  count so
	// that if there are errors an operator can see that.
	count := 0
	for _, k := range sourceKeys {
		v, err := source.Get(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error copying object (%s): %s\n", k, err)
			continue
		}
		if err := target.Put(k, v); err != nil {
			fmt.Fprintf(os.Stderr, "Error putting object (%s): %s\n", k, err)
			continue
		}
		count++
	}
	fmt.Println("Copy complete; %d objects were copied.  If this number does not match the number of source objects you have encountered copy errors and should review them above!", count)
}
