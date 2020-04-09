// Copyright 2020 Google Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

func main() {
	file := os.Args[1]
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	// decode
	var v yaml.Node
	if err := yaml.Unmarshal(b, &v); err != nil {
		panic(err)
	}

	// TODO investigate how are multiple docs coming to Node
	if len(v.Content) == 0 {
		panic("no yaml docs found")
	}

	// TODO support multiple documents
	content := v.Content[0]

	// build tree out of nodes due to lack of parent links in *yaml.Node
	scanNodeMappings(content, `$root`)

	colorizeKeys(content, "$root")

	colorizeComments(content)

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	enc.Encode(content)
	fmt.Print(render(buf))
}

// TODO do not make these global
var keyPathMap = make(map[*yaml.Node]string)
var kvMap = make(map[*yaml.Node]*yaml.Node)
var vkMap = make(map[*yaml.Node]*yaml.Node)

// scanNodeMappings visits the yaml document node recursively
// to scan keys into a map with keys like $root.foo.bar
func scanNodeMappings(doc *yaml.Node, rootPath string) {
	if doc.Kind == yaml.MappingNode {
		for i, ch := range doc.Content {
			if i%2 == 0 {
				// ch is a Key
				keyPathMap[ch] = fmt.Sprintf("%s.%s", rootPath, ch.Value)
				kvMap[ch] = doc.Content[i+1]
				vkMap[doc.Content[i+1]] = ch
			}
		}
	}
	for _, ch := range doc.Content {
		key := rootPath
		keyNode, ok := vkMap[ch]
		if ok {
			key = keyPathMap[keyNode]
		}
		scanNodeMappings(ch, key)
	}
}

func markComments(in string) string {
	re := regexp.MustCompile(`(?m)^#(.*)`)
	return re.ReplaceAllString(in, `#COMMENT_$1`)
}

func colorizeComments(node *yaml.Node) {
	for _, child := range node.Content {
		child.HeadComment = markComments(child.HeadComment)
		child.LineComment = markComments(child.LineComment)
		child.FootComment = markComments(child.FootComment)
		colorizeComments(child)
	}
}

func render(buf bytes.Buffer) string {
	s := buf.String()

	// render keys
	s = regexp.MustCompile(`(?m)(KEY_BLUE_)([^:]+)`).
		ReplaceAllString(s, color.New(color.FgBlue, color.Bold).Sprint(`$2$3`))
	s = regexp.MustCompile(`(?m)(KEY_YELLOW_)([^:]+)`).
		ReplaceAllString(s, color.New(color.FgYellow, color.Bold).Sprint(`$2$3`))
	s = regexp.MustCompile(`(?m)(KEY_RED_)([^:]+)`).
		ReplaceAllString(s, color.New(color.FgRed, color.Bold).Sprint(`$2$3`))
	s = regexp.MustCompile(`(?m)(KEY_GRAY_)([^:]+)`).
		ReplaceAllString(s, color.New(color.FgHiBlack, color.Bold).Sprint(`$2$3`))

	// render comments
	s = regexp.MustCompile(`(?m)#COMMENT_(.*)`).
		ReplaceAllString(s, color.New(color.FgHiBlack).Sprint(`#$1`))

	return s
}

// TODO remove this method
func colorizeKeys(node *yaml.Node, path string) {
	var prevKey string
	for i, child := range node.Content {
		if node.Kind == yaml.SequenceNode && child.Kind == yaml.ScalarNode {
			continue
		}

		if i%2 == 0 && child.Value != "" {
			// node is a field name
			keyPath := path + "." + child.Value
			prevKey = child.Value
			child.Value = "KEY_" + colorForKey(keyPath) + "_" + child.Value
		}

		subPath := path
		if node.Kind == yaml.MappingNode {
			subPath = path + "." + prevKey
		}
		colorizeKeys(child, subPath)
	}
}

func colorForKey(path string) string {
	redSuffixes := []string{"$root.apiVersion",
		"$root.kind",
		".spec",
		"$root.metadata.name",
		".containers.name",
		".containers.image"}
	for _, f := range redSuffixes {
		if strings.HasSuffix(path, f) {
			return "RED"
		}
	}

	if strings.HasPrefix(path, "$root.metadata") {
		return "YELLOW"
	}

	if strings.HasPrefix(path, "$root.spec") {
		return "BLUE"
	}

	if strings.HasPrefix(path, "$root.status") {
		return "GRAY"
	}
	return "UNKNOWN" // TODO choose a default key color
}
