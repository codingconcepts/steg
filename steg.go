package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/atotto/clipboard"
	"golang.org/x/crypto/sha3"
	"golang.org/x/term"

	"github.com/spf13/cobra"
)

const (
	zwS  = "\u200B" // Zero-width space
	zwNJ = "\u200C" // Zero-width non-joiner
	zwJ  = "\u2060" // Zero-width joiner
)

var (
	encrypt  bool
	stripper = regexp.MustCompile(fmt.Sprintf("[^%s|%s|%s]+", zwS, zwNJ, zwJ))
)

func main() {
	concealCmd := &cobra.Command{
		Use:   "conceal",
		Short: "Conceals a message into another string and copies the output to the clipboard",
		Run:   runConceal,
	}

	revealCmd := &cobra.Command{
		Use:   "reveal",
		Short: "Reveals a hidden message and prints it to the console",
		Run:   runReveal,
	}

	rootCmd := &cobra.Command{}
	rootCmd.PersistentFlags().BoolVarP(&encrypt, "encrypt", "e", false, "encrypt data before hiding it")

	rootCmd.AddCommand(concealCmd, revealCmd)

	fmt.Print("\033[H\033[2J")
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("error executing command: %v", err)
	}
}

func runConceal(_ *cobra.Command, args []string) {
	pub := getInput("Enter public message: ")
	pri := getInput("Enter private message: ")

	if encrypt {
		key := getInputSecret("Enter password: ")
		pri = encryptData(pri, key)
	}

	ct := encode(pub, pri)
	if err := clipboard.WriteAll(ct); err != nil {
		log.Fatalf("error copying to clipboard: %v", err)
	}
}

func runReveal(_ *cobra.Command, args []string) {
	pub := getInput("Enter public message: ")

	pt := decode(pub)

	if encrypt {
		key := getInputSecret("Enter password: ")
		ptBytes := decryptData(pt, key)
		pt = string(ptBytes)
	}

	fmt.Print(pt)
}

func encode(pub, pri string) string {
	bin := strToBin(pri)
	pri = conceal(bin)

	return strings.Trim(inject(pri, pub), "\n")
}

func decode(ct string) string {
	ct = stripper.ReplaceAllString(ct, "")
	return strings.Trim(binToStr(reveal(ct)), "\n\x00")
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

func encryptData(data string, key []byte) string {
	hash := make([]byte, len(data))
	sha3.ShakeSum256(hash, key)

	dataBytes := []byte(data)
	for i := 0; i < len(dataBytes); i++ {
		dataBytes[i] ^= hash[i]
	}

	return base64.StdEncoding.EncodeToString(dataBytes)
}

func decryptData(data string, key []byte) string {
	hash := make([]byte, len(data))
	sha3.ShakeSum256(hash, key)

	dataBytes, _ := base64.StdEncoding.DecodeString(data)
	for i := 0; i < len(dataBytes); i++ {
		dataBytes[i] ^= hash[i]
	}

	return string(dataBytes)
}

func getInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')

	clear()
	return text
}

func getInputSecret(prompt string) []byte {
	fmt.Print(prompt)
	password, _ := term.ReadPassword(int(syscall.Stdin))

	clear()
	return bytes.TrimSpace(password)
}

func clear() {
	fmt.Print("\033[H\033[2J")
}
