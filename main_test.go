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

import "testing"

func Test_mark(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "normal comment",
			in:   "#foo",
			want: "#COMMENT_foo",
		},
		{
			name: "normal comment with space",
			in:   "# foo",
			want: "#COMMENT_ foo",
		},
		{
			name: "multiline",
			in:   `# foo
#bar`,
			want: `#COMMENT_ foo
#COMMENT_bar`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := markComments(tt.in); got != tt.want {
				t.Errorf("markComments(%q) got=%q, want=%q", tt.in, got, tt.want)
			}
		})
	}
}
