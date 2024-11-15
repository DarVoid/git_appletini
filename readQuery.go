package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

func loadTrackingConfig() {
	file_contents, err := os.ReadFile(TRACKING_CONFIG_FILE)
	if err != nil {
		fmt.Println()
	}
	// fmt.Println(string(file_contents))
	err = json.Unmarshal(file_contents, &TrackingConfig)
	ehp(err)
	fmt.Println(TrackingConfig)
}

func createQuery() {
	templatePath := "./queries/prMultiRepoByLabel.gql"
	tpl, err := os.ReadFile(templatePath)
	if err != nil {
		panic(fmt.Sprintf("cannot read template path: %v", templatePath))
	}
	templateName := "graphqlQueryLabeledRepos"
	loadedTemplate, err := template.New(templateName).Funcs(funcMap).Parse(string(tpl))
	if err != nil {
		panic(fmt.Sprintf("cannot load template funcmap"))
	}
	PrQuery = new(bytes.Buffer)

	err = loadedTemplate.Execute(PrQuery, TrackingConfig)
	if err != nil {
		panic(fmt.Sprintf("cannot create query"))
	}

	// query created from config
	SavedPRQuerry = fmt.Sprint(PrQuery)
	// fmt.Println(SavedPRQuerry)

	// fileOutputPath := "./finalQuery.gql"
	// outputFile, err := os.Create(fileOutputPath)
	// if err != nil {
	// 	panic(fmt.Sprintf("cannot create file: %v", fileOutputPath))
	// }
	// // write string to file
	// _, err = outputFile.Write([]byte(SavedPRQuerry))
	// if err != nil {
	// 	panic(fmt.Sprintf("cannot write to file: %v", fileOutputPath))
	// }
}

var funcMap = template.FuncMap{
	"ToSnake": strcase.ToSnake,
	"ToKebab": strcase.ToKebab,
	"ToLower": strings.ToLower,
	"ToUpper": strings.ToUpper,
}
