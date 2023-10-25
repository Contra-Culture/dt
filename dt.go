package dt

import (
	"fmt"
	"strings"
)

// Represents Template entity.
type Template []interface{}

// Represents injection entity (place for content injection).
type inj struct{}

// We use this replacer to make HTML-safe strings out of user's input.
var safeTextReplacer = strings.NewReplacer("<", "&lt;", ">", "&gt;", "\"", "&quot", "'", "&quot")

// Safe() generated HTML safe string out of its input.
func Safe(s string) string {
	return safeTextReplacer.Replace(s)
}

// T() is a template constructor function.
// All the templates MUST be made with this constructor.
// This function can take three types of arguments:
// - nil - does nothing,
// - string - appends to the previous string fragment or creates new,
// - inj - adds place for code injection.
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

// I() is an injection (place for injection) constructor function.
// You can name injections, like:
//
//	I("title")
//
// ... or provide explicitly an order:
//
// I("title", 1)
//
// ... for commenting purpose, but only the order which injections following in the template definition is important. That means that all the I()'s arguments will be ignored.
func I(_ ...interface{}) interface{} {
	return inj{}
}

// *Template.Render() renders template with provided content injections.
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

// *Template.MustRender() the same as *Template.Render() except it panics instead of returning error.
func (t *Template) MustRender(data ...string) string {
	r, err := t.Render(data...)
	if err != nil {
		panic(err)
	}
	return r
}

// *Template.RenderCollection() - renders the template N times,
// where N is a length of slice with injections for every template render.
func (t *Template) RenderCollection(cdata ...[]string) (string, error) {
	var sb strings.Builder
	for _, data := range cdata {
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
	}
	return sb.String(), nil
}

// *Template.MustRenderCollection() is the same as *Template.RenderCollection() except it panics instead of returning error.
func (t *Template) MustRenderCollection(cdata ...[]string) string {
	r, err := t.RenderCollection(cdata...)
	if err != nil {
		panic(err)
	}
	return r
}

// Stylesheet type represents a CSS stylesheet
type Stylesheet struct {
	name      string
	sb        strings.Builder
	templates map[string]*Styling
}

// Styling type represents a reusable CSS styling trait.
type Styling struct {
	ruleTemplates []*ruleTemplate
}

// ruleTemplate type represents a single CSS rule template
type ruleTemplate struct {
	selectorTemplate *Template
	selectors        []string
	block            string
}
type RuleTemplateNesting struct {
	styling          *Styling
	selectorTemplate *Template
}

// S() is a stylesheet constructor function.
func S(n, c string) *Stylesheet {
	s := &Stylesheet{
		name:      n,
		templates: map[string]*Styling{},
	}
	s.sb.WriteString(c)
	return s
}

// D() is CSS rule's declaration constructor.
func D(n string, values ...string) string {
	return fmt.Sprintf("\t%s: %s;", n, strings.Join(values, ", "))
}

// B() is CSS rule's declarations block constructor.
func B(ds ...string) string {
	return fmt.Sprintf("{\n%s\n}", strings.Join(ds, "\n"))
}

// C() is a CSS comment constructor.
func (s *Stylesheet) C(c string) interface{} {
	s.sb.WriteString(fmt.Sprintf("\n/* %s */\n", c))
	return nil
}

// *Stylesheet.Append() allows to add content to the Stylesheet on which it was called.
// Use *Stylesheet.Append() in Template constructors to generate CSS from HTML.
// For the purpose of compatibility with Template constructor (T()) it returns interface{}(nil)).
// It should be used in template constructors, but adds nothing to templates.
func (s *Stylesheet) Append(cs ...string) interface{} {
	for _, c := range cs {
		s.sb.WriteString(c)
	}
	return nil
}

// R() is a CSS rule constructor. It receives body and selectors.
// Examples:
//
//	R(B(D("color", "red")), ".danger")
//	R(B("color: red;"), ".danger")
//	R("{color: red", ".danger")
func R(b string, ss ...string) string {
	return fmt.Sprintf("%s %s\n\n", strings.Join(ss, ",\n"), b)
}

// *Stylesheet.S() is a CSS Styling constructor.
// You can use Styling for creation reusable traits that later can be applied by using use cases.
func (s *Stylesheet) S(n string) *Styling {
	st := &Styling{}
	if _, exists := s.templates[n]; exists {
		panic(fmt.Errorf("*Stylesheet.S(): styling \"%s\" already specified", n))
	}
	s.templates[n] = st
	return st
}

// *Styling.RT() defines CSS Rule template for styling.
// *Styling.RT() uses templates for selector templating.
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

// Join() allows to join several templates into a single one.
func Join(ts ...*Template) *Template {
	fragments := []interface{}{}
	wasString := false
	for _, t := range ts {
		for _, _f := range []interface{}(*t) {
			switch f := _f.(type) {
			case string:
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

// *RuleTemplateNesting.RT() is the same as *Styling.RT() except it alows to add Styling/traits for nested selectors.
func (rtn *RuleTemplateNesting) RT(b string, st *Template) *RuleTemplateNesting {
	return rtn.styling.RT(b, Join(rtn.selectorTemplate, st))
}

// *Stylesheet.SC() adds styling use case for particular CSS class and returns that class name.
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

// *Stylesheet.Compile() generates CSS code for the stylesheet.
func (s *Stylesheet) Compile() string {
	for n, st := range s.templates {
		s.sb.WriteString(fmt.Sprintf("/* styling: %s */\n", n))
		for _, rt := range st.ruleTemplates {
			s.sb.WriteString(fmt.Sprintf("%s %s\n\n", strings.Join(rt.selectors, ",\n"), rt.block))
		}
	}
	return s.sb.String()
}
