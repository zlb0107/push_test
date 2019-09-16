package log

type Kit interface {
	// Bussiness log
	B() Logger
	// Gen log
	G() Logger
	// Acceess log
	A() Logger

	// Slow log
	S() Logger
}

type kit struct {
	b, g, a, s Logger
}

func NewKit(b, g, a, s Logger) Kit {
	return kit{
		b: b,
		g: g,
		a: a,
		s: s,
	}
}

func (c kit) B() Logger {
	return c.b
}

func (c kit) G() Logger {
	return c.g
}

func (c kit) A() Logger {
	return c.a
}

func (c kit) S() Logger {
	return c.s
}
