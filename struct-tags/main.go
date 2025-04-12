package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Test struct {
	Name string `json:"name_field" my_tag:"my_tag_value1"`
	Age  int    `json:"age_field" my_tag:"my_tag_value2"`
}

func main() {
	t := Test{
		Name: "John Doe",
		Age:  30,
	}

	v := reflect.ValueOf(t)

	for i := 0; i < v.NumField(); i++ {
		// Print the tag value
		field := v.Type().Field(i)
		jsonTag := field.Tag.Get("json")
		myTag := field.Tag.Get("my_tag")
		fmt.Printf("Field: %s, JSON Tag: %s, My Tag: %s\n", field.Name, jsonTag, myTag)
	}

	// Convert struct to pretty printed JSON
	jsonData, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	fmt.Print("JSON Data:\n", string(jsonData), "\n")

	// Convert JSON back to struct
	var newTest Test
	err = json.Unmarshal(jsonData, &newTest)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	fmt.Printf("Unmarshalled Struct: %+v\n", newTest)
}
