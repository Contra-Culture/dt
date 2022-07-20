package dt

import (
	"fmt"
	"strings"
)

type Template []interface{}
type inj struct{}

var safeTextReplacer = strings.NewReplacer("<", "&lt;", ">", "&gt;", "\"", "&quot", "'", "&quot")

func Safe(s string) string {
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
		case inj:
			fragments = append(fragments, f)
			wasString = false
		}
	}
	t := Template(fragments)
	return &t
}
func I(_ ...string) interface{} {
	return inj{}
}
func (t *Template) Render(data ...string) (string, error) {
	var sb strings.Builder
	var injIdx = -1
	for _, _f := range *t {
		switch f := _f.(type) {
		case string:
			sb.WriteString(f)
		case inj:
			injIdx++
			if len(data) <= injIdx {
				return "", fmt.Errorf("*Template.Render(): injection [%d] not provided, got: \"%#v\"", injIdx, data)
			}
			sb.WriteString(data[injIdx])
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
func (s *Stylesheet) Append(cs ...string) interface{} {
	for _, c := range cs {
		s.sb.WriteString(c)
	}
	return nil
}
func R(b string, ss ...string) string {
	return fmt.Sprintf("%s %s\n\n", strings.Join(ss, ",\n"), b)
}
func (s *Stylesheet) S(n string) *Styling {
	st := &Styling{}
	if _, exists := s.templates[n]; exists {
		panic(fmt.Errorf("*Stylesheet.S(): styling \"%s\" already specified", n))
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
func (s *Stylesheet) SC(cn, n string, inj ...string) string {
	st, exists := s.templates[n]
	if !exists {
		panic(fmt.Errorf("*Stylesheet.SC(): can't add styling use case: styling \"%s\" is not specified", n))
	}
	inj = append([]string{"." + cn}, inj...)
	for _, rt := range st.ruleTemplates {
		selector, err := rt.selectorTemplate.Render(inj...)
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
