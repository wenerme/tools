package convert

type Context struct {
	o map[string]string
}

func (c *Context) With(o Option) *Context {
	c.o[string(o)] = "1"
	return c
}
func (c *Context) Has(o Option) bool {
	return len(c.o[string(o)]) > 0
}

type Option string
