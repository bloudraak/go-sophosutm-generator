package main

import "testing"

func Test_removeCommonPrefix(t *testing.T) {
    type args struct {
        left  string
        right string
    }
    tests := []struct {
        name string
        args args
        want string
    }{
        {
            name: "empty",
            args: args{
                left:  "",
                right: "",
            },
            want: "",
        },
        {
            name: "empty left",
            args: args{
                left:  "",
                right: "abc",
            },
            want: "",
        },
        {
            name: "empty right",
            args: args{
                left:  "abc",
                right: "",
            },
            want: "abc",
        },
        {
            name: "no common prefix",
            args: args{
                left:  "abc",
                right: "def",
            },
            want: "abc",
        },

        {
            name: "common prefix abc/abcd",
            args: args{
                left:  "abc",
                right: "abcd",
            },
            want: "",
        },
        {
            name: "common prefix abcd/abc",
            args: args{
                left:  "abcd",
                right: "abc",
            },
            want: "d",
        },
        {
            name: "common prefix abc/abc",
            args: args{
                left:  "abc",
                right: "abc",
            },
            want: "",
        },
        {
            name: "common prefix abcdef/abcghi",
            args: args{
                left:  "abcdef",
                right: "abcghi",
            },
            want: "def",
        },

    }
    for _, tt := range tests {
        t.Run(
            tt.name, func(t *testing.T) {
                if got := removeCommonPrefix(tt.args.left, tt.args.right); got != tt.want {
                    t.Errorf("removeCommonPrefix() = %v, want %v", got, tt.want)
                }
            },
        )
    }
}
