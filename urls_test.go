package main

import "testing"

func Test_join(t *testing.T) {
    type args struct {
        paths []string
    }
    tests := []struct {
        name string
        args args
        want string
    }{
        {
            name: "empty",
            args: args{
                paths: []string{},
            },
            want: "",
        },
        {
            name: "one",
            args: args{
                paths: []string{"one"},
            },
            want: "one",
        },
        {
            name: "two",
            args: args{
                paths: []string{"one", "two"},
            },
            want: "one/two",
        },
        {
            name: "three",
            args: args{
                paths: []string{"one", "two", "three"},
            },

            want: "one/two/three",
        },
        {
            name: "empty middle",
            args: args{
                paths: []string{"one", "", "three"},
            },
            want: "one/three",
        },
        {
            name: "empty start",
            args: args{
                paths: []string{"", "two", "three"},
            },
            want: "two/three",
        },
        {
            name: "empty end",
            args: args{
                paths: []string{"one", "two", ""},
            },
            want: "one/two",

        },
        {
            name: "empty all",
            args: args{
                paths: []string{"", "", ""},
            },
            want: "",
        },
        {
            name: "slash",
            args: args{
                paths: []string{"one/", "two"},
            },
            want: "one/two",
        },
        {
            name: "standalone slash",
            args: args{
                paths: []string{"one", "/", "two"},
            },
            want: "one/two",
        },
        {
            name: "slash prefix",
            args: args{
                paths: []string{"one", "/two"},
            },
            want: "one/two",
        },
    }
    for _, tt := range tests {
        t.Run(
            tt.name, func(t *testing.T) {
                if got := join(tt.args.paths...); got != tt.want {
                    t.Errorf("join() = %v, want %v", got, tt.want)
                }
            },
        )
    }
}
