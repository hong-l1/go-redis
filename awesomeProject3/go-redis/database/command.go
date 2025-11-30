package database

import "strings"

var CmdMap = map[string]*Command{}

type Command struct {
	exec  ExecFunc
	arity int
}

func RegisterCommand(name string, exec ExecFunc, arity int) {
	name = strings.ToLower(name)
	CmdMap[name] = &Command{exec, arity}
}
