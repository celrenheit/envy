package envy

import "os"

type EnvVarParser interface {
	Parse() (string, bool)
}

type Modifier interface {
	Modify(input string) (out string)
}

type Fusioner interface {
	Fusion(vars ...string) string
}

type FusionerFunc func(vars ...string) string

func (fn FusionerFunc) Fusion(vars ...string) string {
	return fn(vars...)
}

type ModifierFunc func(input string) (out string)

func (fn ModifierFunc) Modify(input string) string {
	return fn(input)
}

type EnvVar struct {
	name     string
	modifier Modifier
}

func (ev *EnvVar) Parse() (string, bool) {
	if val := os.Getenv(ev.name); val != "" {
		if ev.modifier != nil {
			return ev.modifier.Modify(val), true
		}
		return val, true
	}

	return "", false
}

type FusionVar struct {
	names  []string
	merger Fusioner
}

func (fv *FusionVar) Parse() (string, bool) {
	vals := []string{}
	for _, name := range fv.names {
		val := os.Getenv(name)
		if val == "" {
			return "", false
		}
		vals = append(vals, val)
	}

	if val := fv.merger.Fusion(vals...); val != "" {
		return val, true
	}

	return "", false
}

type Envy struct {
	vars         []EnvVarParser
	defaultValue string
}

func New() *Envy {
	return &Envy{
		vars: []EnvVarParser{},
	}
}

func (e *Envy) Add(ev string) *Envy {
	e.vars = append(e.vars, &EnvVar{
		name: ev,
	})
	return e
}

func (e *Envy) Default(val string) *Envy {
	e.defaultValue = val
	return e
}

func (e *Envy) AddWithModifier(ev string, mod Modifier) *Envy {
	e.vars = append(e.vars, &EnvVar{
		name:     ev,
		modifier: mod,
	})
	return e
}

func (e *Envy) AddWithModifierFunc(ev string, modFn ModifierFunc) *Envy {
	e.vars = append(e.vars, &EnvVar{
		name:     ev,
		modifier: ModifierFunc(modFn),
	})
	return e
}

func (e *Envy) Merge(merger Fusioner, names ...string) *Envy {
	e.vars = append(e.vars, &FusionVar{
		names:  names,
		merger: merger,
	})
	return e
}

func (e *Envy) MergeFunc(mergerFn FusionerFunc, names ...string) *Envy {
	e.vars = append(e.vars, &FusionVar{
		names:  names,
		merger: FusionerFunc(mergerFn),
	})

	return e
}

func (e *Envy) Getenv() string {
	val, _ := e.GetenvOk()
	return val
}

func (e *Envy) GetenvOk() (string, bool) {
	for _, v := range e.vars {
		if val, ok := v.Parse(); ok {
			return val, ok
		}
	}

	if e.defaultValue != "" {
		return e.defaultValue, true
	}

	return "", false
}
