package process

import (
	"fmt"
	"strings"

	"github.com/kam1sh/overprocessed/pkg/process/api"
)

func DumpTree(proc api.Process) string {
	return strings.Join(walkNode(proc, 0), "\n")
}

func walkNode(proc interface{}, level int) []string {
	if proc == nil {
		return []string{}
	}
	out := make([]string, 0)
	switch node := proc.(type) {
	case *Redirect:
		out = append(out, "Redirect{")
		if node.From != nil {
			out = append(out, "\tproc:")
			out = append(out, walkNode(node.From, level+1)...)
		}
		if node.Stdin != nil {
			out = append(out, "\tstdin:")
			out = append(out, walkNode(node.Stdin, level+1)...)
		}
		if node.Stdout != nil {
			out = append(out, "\tstdout:")
			out = append(out, walkNode(node.Stdout, level+1)...)
		}
		if node.Stderr != nil {
			out = append(out, "\tstderr:")
			out = append(out, walkNode(node.Stderr, level+1)...)
		}
		out = append(out, "}")
	case *Process:
		cmd := node.Cmd()
		out = append(out,
			"Process{",
			fmt.Sprint("\tcmd: ", cmd.Args),
		)
		stdin, ok := cmd.Stdin.(fmt.Stringer)
		if ok {
			out = append(out, fmt.Sprint("\tstdin:", stdin.String()))
		}
		stdout, ok := cmd.Stdout.(fmt.Stringer)
		if ok {
			out = append(out, fmt.Sprint("\tstdout:", stdout.String()))
		}
		stderr, ok := cmd.Stdout.(fmt.Stringer)
		if ok {
			out = append(out, fmt.Sprint("\tstderr:", stderr.String()))
		}
		out = append(out, "}")
	case fmt.Stringer:
		out = append(out, node.String())
	default:
		out = append(out, "[unknown]")
	}
	for i := range out {
		out[i] = fmt.Sprint(strings.Repeat("\t\t", level), out[i])
	}
	return out
}
