package filters

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type Filter struct {
	Use          string `toml:"use"`
	When         string `toml:"when"`
	compiledWhen *vm.Program
	Transform    string `toml:"transform"`
	Action       string `toml:"action"`
}

func (f *Filter) IsMatch(env any) (bool, error) {
	var err error
	if f.compiledWhen == nil {
		f.compiledWhen, err = expr.Compile(f.When, expr.AsBool())
		if err != nil {
			return false, fmt.Errorf("could not compile 'use' expression: %w", err)
		}
	}
	result, err := expr.Run(f.compiledWhen, env)
	if err != nil {
		return false, fmt.Errorf("could not execute 'use' expression: %w", err)
	}
	return result.(bool), nil
}
