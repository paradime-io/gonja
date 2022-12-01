package exec

type Context struct {
	data   map[string]interface{}
	parent *Context
}

func NewContext(data map[string]interface{}) *Context {
	return &Context{data: data}
}

func EmptyContext() *Context {
	return &Context{data: map[string]interface{}{}}
}

func (ctx *Context) Has(name string) bool {
	_, exists := ctx.data[name]
	if !exists && ctx.parent != nil {
		return ctx.parent.Has(name)
	}
	return exists
}

func (ctx *Context) Get(name string) interface{} {
	value, exists := ctx.data[name]
	if exists {
		return value
	} else if ctx.parent != nil {
		return ctx.parent.Get(name)
	} else {
		return nil
	}
}

func (ctx *Context) Overwrite(name string, value interface{}) bool {
	if _, hasKey := ctx.data[name]; hasKey {
		ctx.data[name] = value
		return true
	} else if ctx.parent != nil {
		return ctx.parent.Overwrite(name, value)
	} else {
		return false
	}
}

// Set first tries to overwrite the variable in a outer scope before adding it here.
// If that doesn't work, we'll set the variable locally.
func (ctx *Context) Set(name string, value interface{}) {
	if !ctx.Overwrite(name, value) {
		ctx.data[name] = value
	}
}

func (ctx *Context) Inherit() *Context {
	return &Context{
		data:   map[string]interface{}{},
		parent: ctx,
	}
}

// Update updates this context with the key/value pairs from a map.
func (ctx *Context) Update(other map[string]interface{}) *Context {
	for k, v := range other {
		ctx.data[k] = v
	}
	return ctx
}

// Merge updates this context with the key/value pairs from another context.
func (ctx *Context) Merge(other *Context) *Context {
	return ctx.Update(other.data)
}
