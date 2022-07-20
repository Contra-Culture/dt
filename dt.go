package dt

import (
	"fmt"
	"strings"
)

type Template []interface{}
type inj struct {
	unsafe bool
	key    string
}

var safeTextReplacer = strings.NewReplacer("<", "&lt;", ">", "&gt;", "\"", "&quot", "'", "&quot")

func Protected(s string) string {
	return safeTextReplacer.Replace(s)
}
func T(_fragments ...interface{}) *Template {
	fragments := []interface{}{}
	wasString := false
	for _, _f := range _fragments {
		switch f := _f.(type) {
		case nil:
			// do nothing
		case string:
			if wasString {
				fragments[len(fragments)-1] = fragments[len(fragments)-1].(string) + f
				continue
			}
			wasString = true
			fragments = append(fragments, f)
		case *inj:
			fragments = append(fragments, f)
			wasString = false
		}
	}
	t := Template(fragments)
	return &t
}
func I(k string) interface{} {
	return &inj{
		unsafe: false,
		key:    k,
	}
}
func UI(k string) interface{} {
	return &inj{
		unsafe: true,
		key:    k,
	}
}
func (t *Template) Render(data map[string]string) (string, error) {
	var sb strings.Builder
	for _, _f := range *t {
		switch f := _f.(type) {
		case string:
			sb.WriteString(f)
		case *inj:
			inj, exists := data[f.key]
			if !exists {
				return "", fmt.Errorf("injection \"%s\" not provided", f.key)
			}
			if !f.unsafe {
				inj = safeTextReplacer.Replace(inj)
			}
			sb.WriteString(inj)
		}
	}
	return sb.String(), nil
}

type Stylesheet struct {
	name      string
	sb        strings.Builder
	templates map[string]*Styling
}
type Styling struct {
	ruleTemplates []*ruleTemplate
}
type ruleTemplate struct {
	selectorTemplate *Template
	selectors        []string
	block            string
}
type RuleTemplateNesting struct {
	styling          *Styling
	selectorTemplate *Template
}

func S(n, c string) *Stylesheet {
	s := &Stylesheet{
		name:      n,
		templates: map[string]*Styling{},
	}
	s.sb.WriteString(c)
	return s
}
func D(n string, values ...string) string {
	return fmt.Sprintf("\t%s: %s;", n, strings.Join(values, ", "))
}
func B(ds ...string) string {
	return fmt.Sprintf("{\n%s\n}", strings.Join(ds, "\n"))
}
func (s *Stylesheet) C(c string) interface{} {
	s.sb.WriteString(fmt.Sprintf("\n/* %s */\n", c))
	return nil
}
func (s *Stylesheet) Append(c string) interface{} {
	s.sb.WriteString(c)
	return nil
}
func R(b string, ss ...string) string {
	return fmt.Sprintf("%s %s\n\n", strings.Join(ss, ",\n"), b)
}
func (s *Stylesheet) S(n string) *Styling {
	st := &Styling{}
	if _, exists := s.templates[n]; exists {
		panic(fmt.Errorf("styling \"%s\" already specified", n))
	}
	s.templates[n] = st
	return st
}
func (s *Styling) RT(b string, st *Template) *RuleTemplateNesting {
	rt := &ruleTemplate{
		selectorTemplate: st,
		block:            b,
	}
	s.ruleTemplates = append(s.ruleTemplates, rt)
	return &RuleTemplateNesting{
		styling:          s,
		selectorTemplate: st,
	}
}
func Join(ts ...*Template) *Template {
	fragments := []interface{}{}
	wasString := false
	for _, t := range ts {
		for _, _f := range []interface{}(*t) {
			switch f := _f.(type) {
			case string:
				wasString = true
				if wasString {
					fragments[len(fragments)-1] = fragments[len(fragments)-1].(string) + f
					continue
				}
				wasString = true
				fragments = append(fragments, f)
			default:
				wasString = false
				fragments = append(fragments, f)
			}
		}
	}
	t := Template(fragments)
	return &t
}
func (rtn *RuleTemplateNesting) RT(b string, st *Template) *RuleTemplateNesting {
	return rtn.styling.RT(b, Join(rtn.selectorTemplate, st))
}

const self = "self"

func Self() interface{} {
	return I(self)
}
func (s *Stylesheet) SC(cn, n string, inj map[string]string) string {
	st, exists := s.templates[n]
	if !exists {
		panic(fmt.Errorf("*Stylesheet.SC(): can't add styling use case: styling \"%s\" is not specified", n))
	}
	if inj == nil {
		inj = map[string]string{self: "." + cn}
	} else if _, exists := inj[self]; exists {
		panic(fmt.Errorf("*Stylesheet.SC(): \"self\" key is reserved"))
	} else {
		inj[self] = "." + cn
	}
	for _, rt := range st.ruleTemplates {
		selector, err := rt.selectorTemplate.Render(inj)
		if err != nil {
			panic(fmt.Errorf("*Stylesheet.SC(): can't add styling use case: %w", err))
		}
		rt.selectors = append(rt.selectors, selector)
	}
	return cn
}
func (s *Stylesheet) Compile() string {
	for n, st := range s.templates {
		s.sb.WriteString(fmt.Sprintf("/* styling: %s */\n", n))
		for _, rt := range st.ruleTemplates {
			s.sb.WriteString(fmt.Sprintf("%s %s\n\n", strings.Join(rt.selectors, ",\n"), rt.block))
		}
	}
	return s.sb.String()
}
