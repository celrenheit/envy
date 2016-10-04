package envy

import (
	"os"
	"strings"
)

type EnvVarParser interface {
	Parse() (string, bool)
}

type Modifier func(input string) (out string)

type Fusioner func(vars ...string) string

type EnvVar struct {
	name     string
	modifier Modifier
}

func (ev *EnvVar) Parse() (string, bool) {
	if val := os.Getenv(ev.name); val != "" {
		if ev.modifier != nil {
			return ev.modifier(val), true
		}
		return val, true
	}

	return "", false
}

type FusionVar struct {
	names  []string
	fusion Fusioner
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

	if val := fv.fusion(vals...); val != "" {
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

func (e *Envy) Merge(fusionner Fusioner, names ...string) *Envy {
	e.vars = append(e.vars, &FusionVar{
		names:  names,
		fusion: fusionner,
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

func PrependWith(str string) Modifier {
	return func(in string) string {
		return str + in
	}
}

func AppendWith(str string) Modifier {
	return func(in string) string {
		return in + str
	}
}

func Join(s string) Fusioner {
	return func(vals ...string) string {
		return strings.Join(vals, s)
	}
}
