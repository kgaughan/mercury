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

// Compile compiles the filter's when and transform expressions. An error is
// returned if either expression could not be compiled.
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

// Run applies the filter to the given entry. It returns the transformed entry
// or nil if the entry should be excluded. An error is returned if the filter
// could not be applied.
func (f *Filter) Run(entry *gofeed.Item) (*gofeed.Item, error) {
	env := map[string]any{
		"entry": entry,
	}
	keep, err := expr.Run(f.compiledWhen, expr.Env(env))
	if err != nil {
		return nil, fmt.Errorf("cannot execute filter %q: %w", f.When, err)
	}
	if !keep.(bool) {
		return nil, nil
	}
	if f.compiledTransform == nil {
		return entry, nil
	}
	result, err := expr.Run(f.compiledTransform, expr.Env(env))
	if err != nil {
		return nil, fmt.Errorf("cannot execute filter transform %q: %w", f.Transform, err)
	}
	return result.(*gofeed.Item), nil
}
