package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

func main() {
	myInput := os.Getenv("INPUT_MYINPUT")

	output := fmt.Sprintf("Hello %s", myInput)

	valueFile := os.Getenv("INPUT_VALUEFILE")
	propertyPath := os.Getenv("INPUT_PROPERTYPATH")
	value := os.Getenv("INPUT_VALUE")
	action := os.Getenv("INPUT_ACTION")

	yamlFile, err := ioutil.ReadFile(valueFile)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
	}

	var yamlConfig map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	fmt.Println(fmt.Sprintf(`::set-output name=myOutput::%s,%s,%s,%s,%s`, output, valueFile, propertyPath, value, action))
}
