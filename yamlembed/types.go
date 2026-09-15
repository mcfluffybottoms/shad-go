package yamlembed

import (
	"strings"
)

type Foo struct {
	A string `yaml:"aa"`
	p int64
}

type Bar struct {
	I      int64    `yaml:"-"`
	B      string   `yaml:"b"`
	UpperB string   `yaml:"-"`
	OI     []string `yaml:"-"`
	F      []any    `yaml:"f,flow"`
}

func (b *Bar) UnmarshalYAML(unmarshal func(any) error) error {
	var v struct {
		B  string   `yaml:"b"`
		OI []string `yaml:"oi"`
		F  []any    `yaml:"f"`
	}
	if err := unmarshal(&v); err != nil {
		return err
	}
	b.B = v.B
	b.UpperB = strings.ToUpper(v.B)
	b.OI = v.OI
	b.F = v.F
	return nil
}

type Baz struct {
	Foo `yaml:",inline"`
	Bar `yaml:",inline"`
}

func (b *Baz) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v struct {
		A  string   `yaml:"aa"`
		B  string   `yaml:"b"`
		OI []string `yaml:"oi"`
		F  []any    `yaml:"f"`
	}

	if err := unmarshal(&v); err != nil {
		return err
	}

	b.Foo.A = v.A

	b.Bar.B = v.B
	b.Bar.UpperB = strings.ToUpper(v.B)
	b.Bar.OI = v.OI
	b.Bar.F = v.F

	return nil
}
