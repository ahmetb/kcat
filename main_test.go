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
