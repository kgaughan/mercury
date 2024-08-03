package filters

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/mmcdole/gofeed"
)

type Filter struct {
	When              string `toml:"when"`
	compiledWhen      *vm.Program
	Transform         string `toml:"transform"`
	compiledTransform *vm.Program
}

func (f *Filter) Compile() error {
	var err error
	if f.When == "" {
		f.When = "true"
	}
	if f.compiledWhen, err = expr.Compile(f.When, expr.AsBool()); err != nil {
		return fmt.Errorf("cannot compile filter when expression %q: %w", f.When, err)
	}
	if f.Transform != "" {
		if f.compiledTransform, err = expr.Compile(f.Transform); err != nil {
			return fmt.Errorf("cannot compile filter transform expression %q: %w", f.Transform, err)
		}
	}
	return nil
}

func (f *Filter) Run(entry *gofeed.Item) (bool, *gofeed.Item, error) {
	env := map[string]any{
		"entry": entry,
	}
	keep, err := expr.Run(f.compiledWhen, expr.Env(env))
	if err != nil {
		return false, nil, fmt.Errorf("cannot execute filter %q: %w", f.When, err)
	}
	if keep.(bool) {
		if f.compiledTransform == nil {
			return true, entry, nil
		}
		result, err := expr.Run(f.compiledTransform, expr.Env(env))
		if err != nil {
			return false, entry, fmt.Errorf("cannot execute filter transform %q: %w", f.Transform, err)
		}
		return true, result.(*gofeed.Item), nil
	}
	return false, entry, nil
}
