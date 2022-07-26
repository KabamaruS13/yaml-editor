package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/stew/objects"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

func main() {

	valueFile := "test.yaml"                                      //os.Getenv("INPUT_VALUEFILE")
	parentPath := "key2"                                          //os.Getenv("INPUT_PARENTPATH")
	key := "key7"                                                 //os.Getenv("INPUT_KEY")
	value := "[{\"key71\": \"value75\", \"key72\": \"value76\"}]" //os.Getenv("INPUT_VALUE")
	action := "upsert"                                            //os.Getenv("INPUT_ACTION")

	output, err := yamlEdit(valueFile, parentPath, key, value, action)
	if err != nil {
		return
	}

	fmt.Println(fmt.Sprintf(`::set-output name=yamlContent::%s`, output))
}

func yamlEdit(valueFile, parentPath, key, value, action string) (string, error) {

	// read the file
	yamlFile, err := ioutil.ReadFile(valueFile)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return "", err
	}

	// unmarshal it to a dynamic struct
	var yamlConfig map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
		return "", err
	}

	// check if the value is json or a string
	var jv interface{}
	err = json.Unmarshal([]byte(value), &jv)
	if err != nil { // our value is a string, and we are going to use it as it is
		jv = value
	}

	if action == "upsert" {
		upsertAction(parentPath, key, jv, yamlConfig)
	} else if action == "delete" {
		deleteAction(parentPath, key, yamlConfig)
	} else {
		fmt.Printf("Wrong action selected: %s\n", action)
		return "", errors.New("wrong action")
	}

	m, err := yaml.Marshal(yamlConfig)
	if err != nil {
		fmt.Printf("Error marshaling YAML file: %s\n", err)
		return "", err
	}

	err = ioutil.WriteFile(valueFile, m, 0644)
	if err != nil {
		fmt.Printf("Error writing YAML file: %s\n", err)
		return "", err
	}

	return string(m), nil
}

func upsertAction(parentPath, key string, value interface{}, yamlConfig map[string]interface{}) {
	path := key
	typ := "interface"
	switch value.(type) {
	case []interface{}:
		typ = "array"
	}
	if parentPath != "" {
		path = parentPath + "." + key
		obj := objects.Map(yamlConfig).Get(parentPath)
		_, ok := obj.(map[string]interface{})
		if !ok { // this is a key-value pair, and I need to turn it into a new node, so I need to delete the value first
			keys := strings.Split(parentPath, ".")
			keyD := keys[len(keys)-1]
			keys = keys[:len(keys)-1]
			deleteAction(strings.Join(keys, "."), keyD, yamlConfig)
		}
	}
	if typ == "array" {
		obj := objects.Map(yamlConfig).Get(path)
		if obj != nil {
			var arr []interface{}
			for _, o := range obj.([]interface{}) {
				arr = append(arr, o.(interface{}))
			}
			for _, v := range value.([]interface{}) {
				arr = append(arr, v.(interface{}))
			}
			objects.Map(yamlConfig).Set(path, arr)
		} else {
			objects.Map(yamlConfig).Set(path, value.([]interface{}))
		}
	} else {
		objects.Map(yamlConfig).Set(path, value)
	}
}

func deleteAction(parentPath, key string, yamlConfig map[string]interface{}) {
	if parentPath != "" { // I need to access the parent map to delete the key
		obj := objects.Map(yamlConfig).Get(parentPath)
		objMap := obj.(map[string]interface{})
		delete(objMap, key)
		objects.Map(yamlConfig).Set(parentPath, objMap)
	} else { // easy, key is on the root, so I can delete it right away
		delete(yamlConfig, key)
	}
}
