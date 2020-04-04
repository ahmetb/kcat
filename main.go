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

	colorizeKeys(content)
	colorizeComments(content)

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	enc.Encode(content)
	fmt.Print(render(buf))
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
	re := regexp.MustCompile(`(?m)(KEY_)([^:]+)`)
	s = re.ReplaceAllString(s, color.RedString(`$2$3`))

	// render comments
	s = regexp.MustCompile(`(?m)#COMMENT_(.*)`).
		ReplaceAllString(s, color.HiBlackString(`#$1`))

	return s
}

func colorizeKeys(node *yaml.Node) {
	for i, child := range node.Content {
		if node.Kind == yaml.SequenceNode && child.Kind == yaml.ScalarNode {
			continue
		}
		// key'ler 2'ye bolunuyor
		if i%2 == 0 && child.Value != "" {
			child.Value = "KEY_" + child.Value
		}
		colorizeKeys(child)
	}
}
