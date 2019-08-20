package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const chunk = 1024

type source struct {
	path   string
	reader io.Reader
	writer io.Writer
}

func main() {
	args := os.Args[1:]
	sources := []source{}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		random := false
		if arg == "-k" {
			i++
			if len(args) == i {
				fmt.Fprintln(os.Stderr, "-k must be not last argument")
				os.Exit(1)
			}
			arg = args[i]
			random = true
		}

		path, err := filepath.Abs(arg)
		if err != nil {
			fmt.Fprintln(os.Stderr, "cannot convert relative path", path, "to absolute:", err)
			os.Exit(1)
		}

		var reader io.Reader

		if random {
			reader = rand.Reader
		} else {
			file, err := os.Open(path)
			if err != nil {
				fmt.Fprintln(os.Stderr, "cannot open file", path, "for reading:", err)
				os.Exit(1)
			}
			defer file.Close()
			reader = file
		}

		var writer io.Writer

		if random {
			file, err := os.Create(path)
			if err != nil {
				fmt.Fprintln(os.Stderr, "cannot open file", path, "for write:", err)
				os.Exit(1)
			}
			defer file.Close()
			writer = file
			fmt.Fprintln(os.Stderr, "created key file", path)
		}

		sources = append(sources, source{
			path,
			reader,
			writer,
		})
	}

	if len(sources) <= 1 || os.Args[0] == "-n" || os.Args[0] == "--n" || os.Args[0] == "-help" || os.Args[0] == "--help" {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "<file> <files...> > <output>")
		fmt.Fprintln(os.Stderr, "you can pass argument -k before file name to generate XOR key and write it to specified file")
		fmt.Fprintln(os.Stderr, "examples:")
		fmt.Fprintln(os.Stderr, ">", os.Args[0], "text.txt key key2 > encrypted1")
		fmt.Fprintln(os.Stderr, ">", os.Args[0], "encrypted1 -k key3 > encrypted2")
		fmt.Fprintln(os.Stderr, ">", os.Args[0], "encrypted2 key key2 key3 > decrypted")
		// TODO: fix it
		fmt.Fprintln(os.Stderr, "please, don't pipe output of", os.Args[0], "to one of source files. it can cause bugs")
		os.Exit(1)
	}

	bytes := make([][]byte, len(sources))
	for i := range bytes {
		bytes[i] = make([]byte, chunk)
	}

	for {
		isEnd := false
		xorBytes := 1024

		for i, source := range sources {
			read, err := source.reader.Read(bytes[i])
			if err != nil && err != io.EOF {
				fmt.Fprintln(os.Stderr, "cannot read data from file ", source.path, ": ", err)
				os.Exit(1)
			}

			if err == io.EOF {
				isEnd = true
			}

			if read < xorBytes {
				xorBytes = read
			}
		}

		for i, source := range sources {
			if source.writer != nil {
				_, err := source.writer.Write(bytes[i][:xorBytes])
				if err != nil {
					fmt.Fprintln(os.Stderr, "cannot write data to key file file", source.path, ":", err)
					os.Exit(1)
				}

			}
		}

		result := bytes[0]
		for s := 1; s < len(sources); s++ {
			for i := 0; i < xorBytes; i++ {
				result[i] = result[i] ^ bytes[s][i]
			}
		}

		_, err := os.Stdout.Write(result[:xorBytes])
		if err != nil {
			fmt.Fprintln(os.Stderr, "cannot write to stdout:", err)
			os.Exit(1)
		}

		if isEnd {
			os.Exit(0)
		}
	}
}
