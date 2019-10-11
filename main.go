package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/template"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func replaceEnvs(match []byte) []byte {
	var varContent = string(match)

	// remove wrapping chars and split
	var varParts = strings.Split(varContent[2:len(varContent)-1], ":-")
	// do we have the env?
	var envValue = os.Getenv(varParts[0])

	if len(varParts) == 2 && envValue == "" {
		envValue = varParts[1]
	}

	return []byte(envValue)
}

func main() {
	programArgs := os.Args[1:]

	var templateString = ""
	if len(programArgs) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		buffer := bytes.Buffer{}
		for scanner.Scan() {
			buffer.WriteString(scanner.Text() + "\n")
		}
		templateString = buffer.String()
	} else if len(programArgs) == 1 {
		fileContent, err := ioutil.ReadFile(programArgs[0])
		check(err)
		templateString = string(fileContent)
	} else {
		panic(errors.New("you need to either pass a filename as first arg or pass content through stdin"))
	}

	// first, execute go templateString
	var envVars = make(map[string]string)
	for _, e := range os.Environ() {
		envVarPair := strings.Split(e, "=")
		envVars[envVarPair[0]] = envVarPair[1]
	}

	t, err := template.New("tmpl").Parse(templateString)
	if err != nil {
		panic(err)
	}

	var templateBytes bytes.Buffer
	terr := t.Execute(&templateBytes, envVars)
	if terr != nil {
		panic(terr)
	}

	envRegex, _ := regexp.Compile("\\$\\{.*\\}")
	envReplaced := envRegex.ReplaceAllFunc(templateBytes.Bytes(), replaceEnvs)

	fmt.Println(string(envReplaced))
}
