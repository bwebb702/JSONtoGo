package JSONtoGo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

var keywords [25]string = [...]string{"break", "default", "func", "interface", "select", "case", "defer",
	"go", "map", "struct", "chan", "else", "goto", "package", "switch", "const", "fallthrough",
	"if", "range", "type", "continue", "for", "import", "return", "var"}

var (
	fileName string
	codeMap  map[string]interface{} // struct name:struct type
)

func CreateStruct(file *os.File, name string) map[string]interface{} {
	codeMap = make(map[string]interface{})

	scanner := bufio.NewScanner(file)
	generateGoCode(scanner, name)

	return codeMap
}

func generateGoCode(scanner *bufio.Scanner, name string) {
	name = contains(name)

	for scanner.Scan() {
		var m map[string]interface{}

		if err := json.Unmarshal(scanner.Bytes(), &m); err != nil {
			log.Fatal(err)
		}

		getKeyValuePairs(m)
		if len(codeMap) == 0 {
			codeMap = m
		}

		filterKeyValuePairs(codeMap, m)
		// fmt.Println(m)
	}
	// fmt.Println(codeMap)
	createCode(name)
}

// This function is used to create only a single instance of each key:value pair
func getKeyValuePairs(m map[string]interface{}) {
	for k, v := range m {
		if fmt.Sprint(reflect.TypeOf(v)) == "map[string]interface {}" {
			getKeyValuePairs(v.(map[string]interface{}))
		} else if fmt.Sprint(reflect.TypeOf(v)) == "[]interface {}" {
			var t string
			for _, k := range v.([]interface{}) {
				t = fmt.Sprintf("[]%v", reflect.TypeOf(k))
			}
			m[k] = t
		} else {
			m[k] = fmt.Sprintf("%v", reflect.TypeOf(v))
		}
	}
}

func filterKeyValuePairs(cm, m map[string]interface{}) {
	for k, v := range m {
		for j, l := range cm {
			if j == k {
				if fmt.Sprint(reflect.TypeOf(v)) == "map[string]interface {}" && fmt.Sprint(reflect.TypeOf(l)) == "map[string]interface {}" {
					filterKeyValuePairs(l.(map[string]interface{}), v.(map[string]interface{}))
				}
			} else if _, ok := cm[k]; !ok {
				cm[k] = v
			}
		}
	}
}

// creates a Go code struct format that can be pasted into code
func createCode(name string) {
	println(fmt.Sprintf("type %s struct {", name))
	recursion(codeMap)
	println("}")
}

// loops through objects using recursion until the lowest level object is reached
func recursion(m map[string]interface{}) {
	for k, v := range m {
		s := strings.Title(k)
		if fmt.Sprint(reflect.TypeOf(v)) == "map[string]interface {}" {
			fmt.Printf("%s struct {\n", s)
			recursion(v.(map[string]interface{}))
			fmt.Println("}")
			// not sure if needed
			// } else if fmt.Sprint(reflect.TypeOf(v)) == "interface {}" {
			// 	fmt.Printf("%s %v `json:%q`\n", s, reflect.TypeOf(v), k)
		} else if fmt.Sprint(reflect.TypeOf(v)) == "[]interface {}" {
			var t string
			for _, k := range v.([]interface{}) {
				t = fmt.Sprintf("%v", reflect.TypeOf(k))
			}
			fmt.Printf("%s []%s `json:%q`\n", s, t, k)
		} else {
			fmt.Printf("%s %v `json:%q`\n", s, v, k)
		}
	}
}

func getFileName() {
	out, err := exec.Command("/python/python", "fileopendialog.py").Output()

	if err != nil {
		panic(err)
	} else {
		fileName = string(out[:])
	}
}

// ensures that none of the struct variable names are Go keywords
func contains(s string) string {
	for i := 0; i < len(keywords); i++ {
		if keywords[i] == s {
			return fmt.Sprintf("%s_", s)
		}
	}

	return s
}
