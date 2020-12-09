package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"./pkg/shamir"
	"./pkg/shares"
)

func main() {
	os.Exit(run(os.Args, os.Stdin))
}

func run(osArgs []string, input io.Reader) int {
	splitCmd := flag.NewFlagSet("split", flag.ContinueOnError)
	combineCmd := flag.NewFlagSet("combine", flag.ContinueOnError)

	splitParts := splitCmd.Int("parts", 0, "number of shares to create")
	splitThreshold := splitCmd.Int("threshold", 0, "number of shares required to reconstruct the secret")

	combineParts := combineCmd.Int("parts", 0, "number of shares you have to recombine")

	helpMsg := func() {
		fmt.Println("ERR: `split` or `combine` subcommand is required")
		fmt.Println()
		fmt.Println("split:")
		splitCmd.PrintDefaults()
		fmt.Println()
		fmt.Println("combine:")
		combineCmd.PrintDefaults()
	}

	if len(osArgs) < 2 {
		helpMsg()
		return 1
	}

	scanner := bufio.NewScanner(input)

	switch osArgs[1] {
	case "split":
		if err := splitCmd.Parse(osArgs[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}

		parts := *splitParts
		threshold := *splitThreshold
		if err := shamir.ValidateSplit(parts, threshold); err != nil {
			fmt.Fprintln(os.Stderr, err)
			splitCmd.PrintDefaults()

			return 1
		}

		fmt.Fprintf(os.Stderr, "Secret: ")
		scanner.Scan()
		secret := []byte(scanner.Text())

		byteShares, err := shamir.Split(secret, parts, threshold)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		fmt.Printf("# parts = %d, threshold = %d\n", parts, threshold)
		fmt.Printf("# %s combine --parts=%d\n", os.Args[0], threshold)
		for _, byteShare := range byteShares {
			fmt.Println(shares.Encode(byteShare))
		}
	case "combine":
		if err := combineCmd.Parse(osArgs[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}

		parts := *combineParts
		if parts < 2 {
			fmt.Fprintln(os.Stderr, "less than two parts cannot be used to reconstruct the secret")
			combineCmd.PrintDefaults()
			return 1
		}

		byteShares := make([][]byte, 0, parts)
		for i := 0; i < parts; i++ {
		readshare:
			fmt.Fprintf(os.Stderr, "Share %d: ", i)
			scanner.Scan()
			share := scanner.Text()
			byteShare, err := shares.Decode(share)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erroring reading share, try again (ctrl-c to exit): %v\n", err)
				goto readshare
			}

			byteShares = append(byteShares, byteShare)
		}

		secret, err := shamir.Combine(byteShares)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Println("# SECRET")
		fmt.Println(string(secret))
	default:
		helpMsg()
		return 1
	}

	return 0
}
