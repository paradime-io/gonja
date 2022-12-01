package exec

import (
	"strings"

	"github.com/paradime-io/gonja/nodes"
	"github.com/pkg/errors"
	// "github.com/paradime-io/gonja/nodes"
)

// FilterFunction is the type filter functions must fulfil
type Macro func(params *VarArgs) *Value

type MacroSet map[string]Macro

// Exists returns true if the given filter is already registered
func (ms MacroSet) Exists(name string) bool {
	_, existing := ms[name]
	return existing
}

// Register registers a new filter. If there's already a filter with the same
// name, Register will panic. You usually want to call this
// function in the filter's init() function:
// http://golang.org/doc/effective_go.html#init
//
// See http://www.florian-schlachter.de/post/gonja/ for more about
// writing filters and tags.
func (ms *MacroSet) Register(name string, fn Macro) error {
	if ms.Exists(name) {
		return errors.Errorf("filter with name '%s' is already registered", name)
	}
	(*ms)[name] = fn
	return nil
}

// Replace replaces an already registered filter with a new implementation. Use this
// function with caution since it allows you to change existing filter behaviour.
func (ms *MacroSet) Replace(name string, fn Macro) error {
	if !ms.Exists(name) {
		return errors.Errorf("filter with name '%s' does not exist (therefore cannot be overridden)", name)
	}
	(*ms)[name] = fn
	return nil
}

// TransformToMapping transforms and validates VarArgs against an expected signature.
// It returns the mapping of argument name -> value and does some basic sanity checks based on the given signature.
func TransformToMapping(
	va *VarArgs,
	argNames []string,
	defaultKwargs []*KwArg,
	activateVarargs bool,
	activateKwargs bool,
) (map[string]*Value, error) {
	mapping := map[string]*Value{}

	varargs := []any{}
	kwargs := map[string]any{}

	allKnownArgNames := []string{}
	allKnownArgNames = append(allKnownArgNames, argNames...)

	// init with defaultKwargs
	for _, defaultValue := range defaultKwargs {
		mapping[defaultValue.Name] = AsValue(defaultValue.Default)
		allKnownArgNames = append(allKnownArgNames, defaultValue.Name)
	}

	// collect all known arg names as a map
	allKnownArgNamesAsMap := map[string]struct{}{}
	for _, argName := range allKnownArgNames {
		allKnownArgNamesAsMap[argName] = struct{}{}
	}

	// use the provided positional arguments to populate the mapping.
	// we go from the index of the argument to the arg map as defined.
	for i, arg := range va.Args {
		if len(va.Args) > len(allKnownArgNames) {
			if activateVarargs {
				varargs = append(varargs, arg.Interface())
			} else {
				return nil, errors.Errorf(`expected at most %d positional arguments, got %d`, len(allKnownArgNames), len(va.Args))
			}
		} else {
			mapping[allKnownArgNames[i]] = arg
		}
	}

	// use the provided keyword arguments to populate the mapping.
	for argName, value := range va.KwArgs {
		if _, isArgNameKnown := allKnownArgNamesAsMap[argName]; isArgNameKnown {
			mapping[argName] = value
		} else {
			if activateKwargs {
				kwargs[argName] = value.Interface()
			} else {
				return nil, errors.Errorf(`got an unexpected keyword argument: '%v'`, argName)
			}
		}
	}

	// sanity checks: are all arguments set?
	knownArgs := map[string]struct{}{} // (needed for step 2)
	for _, arg := range allKnownArgNames {
		if _, isValueSet := mapping[arg]; !isValueSet {
			return nil, errors.Errorf(`missing required positional arguments : '%v'`, arg)
		}
		knownArgs[arg] = struct{}{}
	}

	// add varargs and kwargs if required
	if activateVarargs {
		mapping["varargs"] = AsValue(varargs)
	}
	if activateKwargs {
		mapping["kwargs"] = AsValue(kwargs)
	}

	return mapping, nil
}

func MacroNodeToFunc(node *nodes.Macro, r *Renderer) (func(params *VarArgs) *Value, error) {
	// Compute default values once - lazily, only when requested
	defaultKwargs := []*KwArg{}
	activateVarargs := false
	activateKwargs := false

	alreadyInit := false
	init := func() error {
		if !alreadyInit {
			// https://stackoverflow.com/questions/13944751/args-kwargs-in-jinja2-macros
			for _, n := range node.Wrapper.Nodes {
				macroAsString := strings.ToLower(n.String())
				activateVarargs = activateVarargs || strings.Contains(macroAsString, "val=varargs")
				activateKwargs = activateKwargs || strings.Contains(macroAsString, "val=kwargs")
			}

			for _, pair := range node.Kwargs {
				key := r.Eval(pair.Key).String()
				value := r.Eval(pair.Value)
				if value.IsError() {
					return errors.Wrapf(value, `Unable to evaluate parameter %s=%s`, key, pair.Value)
				}
				defaultKwargs = append(defaultKwargs, &KwArg{Name: key, Default: value.Interface()})
			}
			alreadyInit = true
		}
		return nil
	}

	return func(params *VarArgs) *Value {
		if err := init(); err != nil {
			return AsValue(err)
		}

		var out strings.Builder
		sub := r.Inherit()
		sub.Out = &out

		mapping, mappingErr := TransformToMapping(params, node.Args, defaultKwargs, activateVarargs, activateKwargs)
		if mappingErr != nil {
			return AsValue(errors.Wrapf(mappingErr, `Wrong '%s' macro signature`, node.Name))
		}
		for argName, value := range mapping {
			sub.Ctx.Set(argName, value)
		}

		if err := sub.ExecuteWrapper(node.Wrapper); err != nil {
			return AsValue(errors.Wrapf(err, `Unable to execute macro '%s`, node.Name))
		}
		return AsSafeValue(out.String())
	}, nil
}
