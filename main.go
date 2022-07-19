package main

import (
	"errors"
	"fmt"
	"github.com/stretchr/stew/objects"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

func main() {
	myInput := os.Getenv("INPUT_MYINPUT")

	output := fmt.Sprintf("Hello %s", myInput)

	valueFile := os.Getenv("INPUT_VALUEFILE")
	propertyPath := os.Getenv("INPUT_PROPERTYPATH")
	key := os.Getenv("INPUT_KEY")
	value := os.Getenv("INPUT_VALUE")
	action := os.Getenv("INPUT_ACTION")

	err := yamlEdit(valueFile, propertyPath, key, value, action)
	if err != nil {
		return
	}

	fmt.Println(fmt.Sprintf(`::set-output name=myOutput::%s`, output))
}

func yamlEdit(valueFile, propertyPath, key, value, action string) error {

	yamlFile, err := ioutil.ReadFile(valueFile)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return err
	}

	var yamlConfig map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
		return err
	}

	if action == "insert" {
		obj := objects.Map(yamlConfig).Get(propertyPath)
		objMap, ok := obj.(map[string]interface{})
		if ok { // easy, just add another key-value pair
			objMap[key] = value
			objects.Map(yamlConfig).Set(propertyPath, objMap)
		} else { // this is just a key-value pair, so I need to delete the value and make it an interface

		}
	} else if action == "update" {
		objects.Map(yamlConfig).Set(propertyPath, value)
	} else if action == "delete" {
		if propertyPath != "" { // I need to access the parent map to delete the key
			obj := objects.Map(yamlConfig).Get(propertyPath)
			objMap := obj.(map[string]interface{})
			delete(objMap, key)
			objects.Map(yamlConfig).Set(propertyPath, objMap)
		} else { // easy, key is on the root so I can delete it right away
			delete(yamlConfig, key)
		}
	} else {
		fmt.Printf("Wrong action selected: %s\n", action)
		return errors.New("wrong action")
	}

	m, err := yaml.Marshal(yamlConfig)
	if err != nil {
		fmt.Printf("Error marshaling YAML file: %s\n", err)
		return err
	}

	err = ioutil.WriteFile(valueFile, m, 0644)
	if err != nil {
		fmt.Printf("Error writing YAML file: %s\n", err)
		return err
	}

	return nil
}
