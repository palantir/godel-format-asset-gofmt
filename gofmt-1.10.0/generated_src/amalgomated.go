// generated by amalgomate; DO NOT EDIT
package amalgomated

import (
	"fmt"
	"sort"

	gofmt "github.com/palantir/godel-format-asset-gofmt/gofmt-1.10.0/generated_src/internal/github.com/golang/go/src/cmd/gofmt"
)

var programs = map[string]func(){"gofmt": func() {
	gofmt.AmalgomatedMain()
},
}

func Instance() Amalgomated {
	return &amalgomated{}
}

type Amalgomated interface {
	Run(cmd string)
	Cmds() []string
}

type amalgomated struct{}

func (a *amalgomated) Run(cmd string) {
	if _, ok := programs[cmd]; !ok {
		panic(fmt.Sprintf("Unknown command: \"%v\". Valid values: %v", cmd, a.Cmds()))
	}
	programs[cmd]()
}

func (a *amalgomated) Cmds() []string {
	var cmds []string
	for key := range programs {
		cmds = append(cmds, key)
	}
	sort.Strings(cmds)
	return cmds
}
