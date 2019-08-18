package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"

	"github.com/spf13/cobra"
)

const (
	zwS  = "\u200B" // Zero-width space
	zwNJ = "\u200C" // Zero-width non-joiner
	zwJ  = "\u2060" // Zero-width joiner
)

var stripper = regexp.MustCompile(fmt.Sprintf("[^%s|%s|%s]+", zwS, zwNJ, zwJ))

func main() {
	concealCmd := &cobra.Command{
		Use:   "conceal",
		Short: "Conceals a message into another string and copies the output to the clipboard",
		Run: func(_ *cobra.Command, args []string) {
			pub := getInput("Enter public message: ")
			pri := getInput("Enter private message: ")

			ct := encode(pub, pri)
			if err := clipboard.WriteAll(ct); err != nil {
				log.Fatalf("error copying to clipboard: %v", err)
			}
		},
	}

	revealCmd := &cobra.Command{
		Use:   "reveal",
		Short: "Reveals a hidden message and prints it to the console",
		Run: func(_ *cobra.Command, args []string) {
			pub := getInput("Enter public message: ")

			pt := decode(pub)
			fmt.Println(pt)
		},
	}

	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(concealCmd, revealCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("error executing command: %v", err)
	}
}

func getInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return text
}

func encode(pub, pri string) string {
	bin := strToBin(pri)
	pri = conceal(bin)

	return inject(pri, pub)
}

func decode(ct string) string {
	ct = stripper.ReplaceAllString(ct, "")
	return binToStr(reveal(ct))
}

func strToBin(s string) string {
	output := make([]string, len(s))

	for i, r := range s {
		output[i] = fmt.Sprintf("%b", r)
	}

	return strings.Join(output, " ")
}

func binToStr(s string) string {
	b := make([]byte, len(s))

	for i, r := range strings.Fields(s) {
		n, _ := strconv.ParseUint(r, 2, 8)
		b[i] = byte(n)
	}

	return string(b)
}

func inject(s, into string) string {
	output := into[:1]
	output += s
	output += into[1:]

	return output
}

func conceal(s string) string {
	s = strings.Replace(s, " ", zwS, -1)
	s = strings.Replace(s, "0", zwNJ, -1)
	s = strings.Replace(s, "1", zwJ, -1)

	return s
}

func reveal(s string) string {
	s = strings.Replace(s, zwS, " ", -1)
	s = strings.Replace(s, zwNJ, "0", -1)
	s = strings.Replace(s, zwJ, "1", -1)

	return s
}
