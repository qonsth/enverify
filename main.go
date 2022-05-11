package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	var err error
	defer func() { handleFatal(err) }()

	uniqKeys := map[string]*struct{}{}
	envFilesKeys := map[string]map[string]*struct{}{}
	for _, file := range os.Args[1:] {
		var envMap map[string]string
		if envMap, err = godotenv.Read(file); err != nil {
			return
		}

		envFilesKeys[file] = map[string]*struct{}{}
		for k := range envMap {
			envFilesKeys[file][k] = &struct{}{}
		}

		for k := range envMap {
			uniqKeys[k] = &struct{}{}
		}
	}

	filesNotFound := map[string]map[string]*struct{}{}
	for key := range uniqKeys {
		for file, keysMap := range envFilesKeys {
			if keysMap[key] == nil {
				if filesNotFound[file] == nil {
					filesNotFound[file] = map[string]*struct{}{}
				}
				filesNotFound[file][key] = &struct{}{}
			}
		}
	}

	for file := range filesNotFound {
		if err = handleCommented(file, filesNotFound[file]); err != nil {
			return
		}
	}

	code := 0
	for file, notFound := range filesNotFound {
		notFoundMap := map[string]string{}
		for key := range notFound {
			code = 1
			notFoundMap[key] = ""
		}

		if len(notFoundMap) > 0 {
			var envInfo string
			if envInfo, err = godotenv.Marshal(notFoundMap); err != nil {
				return
			}

			fmt.Printf("\u001B[33mKeys not found in \u001B[34m%s\u001B[33m:\u001B[0m\n", file)
			fmt.Printf("\u001B[31m%v\u001B[0m\r\n\n", envInfo)
		}
	}

	if code == 0 {
		fmt.Printf("\u001B[32mEnverification succeed\u001B[0m\r\n")
	} else {
		fmt.Printf("\u001B[31mEnverification failed\u001B[0m\r\n")
	}
	os.Exit(code)
}

func handleCommented(file string, notFound map[string]*struct{}) (err error) {
	if len(notFound) > 0 {
		var b []byte
		if b, err = ioutil.ReadFile(file); err != nil {
			return
		}
		body := string(b)

		for k := range notFound {
			if strings.Contains(body, "#"+k+"=") {
				delete(notFound, k)
			}
		}
	}

	return
}

func handleFatal(err error) {
	if r := recover(); r != nil {
		switch rt := r.(type) {
		case error:
			err = fmt.Errorf("panic: %w", rt)
		default:
			err = fmt.Errorf("panic: %v", rt)
		}
	}

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\u001B[31mError: %s\u001B[0m\r\n", err)
		os.Exit(1)
	}
}
