# steg
A simple steganographic CLI that uses zero-width characters to hide messages.

Inspired by https://neatnik.net/steganographr

## Installation

```
$ go get -u github.com/codingconcepts/steg
```

## Usage

Help:

```
$ steg

Usage:
   [command]

Available Commands:
  conceal     Conceals a message into another string and copies the output to the clipboard
  help        Help about any command
  reveal      Reveals a hidden message and prints it to the console

Flags:
  -h, --help   help for this command

Use " [command] --help" for more information about a command.
```

### Conceal

```
$ steg conceal
Enter public message: hello
Enter private message: goodbye
```

### Reveal

```
$ steg conceal
Enter public message: h‌‌​​‌‌‌⁠‌‌​‌‌‌‌⁠‌‌​‌‌‌‌⁠‌‌​​‌​​⁠‌‌​​​‌​⁠‌‌‌‌​​‌⁠‌‌​​‌​‌⁠‌​‌​ello
goodbye

```
