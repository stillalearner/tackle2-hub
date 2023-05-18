package filter

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

const (
	QueryParam = "filter"
)

// New filter.
func New(ctx *gin.Context, assertions []Assert) (f Filter, err error) {
	p := Parser{}
	q := strings.Join(
		ctx.QueryArray(QueryParam),
		string(COMMA))
	f, err = p.Filter(q)
	if err != nil {
		return
	}
	err = f.Validate(assertions)
	return
}

// Filter is a collection of predicates.
type Filter struct {
	predicates []Predicate
}

// Validate -
func (f *Filter) Validate(assertions []Assert) (err error) {
	if len(assertions) == 0 {
		return
	}
	find := func(name string) (assert *Assert, found bool) {
		name = strings.ToLower(name)
		for i := range assertions {
			assert = &assertions[i]
			if strings.ToLower(assert.Field) == name {
				found = true
				break
			}
		}
		return
	}
	for _, p := range f.predicates {
		name := p.Field.Value
		v, found := find(name)
		if !found {
			err = Errorf("'%s' not supported.", name)
			return
		}
		err = v.assert(&p)
		if err != nil {
			break
		}
	}

	return
}

// Field returns a field.
func (f *Filter) Field(name string) (field Field, found bool) {
	fields := f.Fields(name)
	if len(fields) > 0 {
		field = fields[0]
		found = true
	}
	return
}

// Fields returns fields.
func (f *Filter) Fields(name string) (fields []Field) {
	name = strings.ToLower(name)
	for _, p := range f.predicates {
		if strings.ToLower(p.Field.Value) == name {
			f := Field{p}
			fields = append(fields, f)
		}
	}
	return
}

// Resource returns a filter scoped to resource.
func (f *Filter) Resource(r string) (filter Filter) {
	r = strings.ToLower(r)
	var predicates []Predicate
	for _, p := range f.predicates {
		field := Field{p}
		fr := field.Resource()
		fr = strings.ToLower(fr)
		if fr == r {
			p.Field.Value = field.Name()
			predicates = append(predicates, p)
		}
	}
	filter.predicates = predicates
	return
}

// Where applies (root) fields to the where clause.
func (f *Filter) Where(in *gorm.DB) (out *gorm.DB) {
	out = in
	for _, p := range f.predicates {
		field := Field{p}
		fr := field.Resource()
		if fr == "" {
			out = out.Where(field.SQL())
		}
	}
	return
}

// Empty returns true when the filter has no predicates.
func (f *Filter) Empty() bool {
	return len(f.predicates) == 0
}

// Field predicate.
type Field struct {
	Predicate
}

// Name returns the field name.
func (f *Field) Name() (s string) {
	_, s = f.split()
	return
}

// As returns the renamed field.
func (f *Field) As(s string) (named Field) {
	named = Field{f.Predicate}
	named.Field.Value = s
	return
}

// Resource returns the field resource.
func (f *Field) Resource() (s string) {
	s, _ = f.split()
	return
}

// SQL builds SQL.
// Returns statement and value (for ?).
func (f *Field) SQL() (s string, v interface{}) {
	name := f.Name()
	switch len(f.Value) {
	case 0:
	case 1:
		switch f.Operator.Value {
		case string(LIKE):
			v = strings.Replace(f.Value[0].Value, "*", "%", -1)
			s = strings.Join(
				[]string{
					name,
					f.operator(),
					"?",
				},
				" ")
		default:
			v = AsValue(f.Value[0])
			s = strings.Join(
				[]string{
					name,
					f.operator(),
					"?",
				},
				" ")
		}
	default:
		if f.Value.Operator(AND) {
			// not supported.
			break
		}
		values := f.Value.ByKind(LITERAL, STRING)
		collection := []interface{}{}
		for i := range values {
			v := AsValue(values[i])
			collection = append(collection, v)
		}
		v = collection
		s = strings.Join(
			[]string{
				name,
				f.operator(),
				"?",
			},
			" ")
	}
	return
}

// split field name.
// format: resource.name
// The resource may be "" (anonymous).
func (f *Field) split() (relation string, name string) {
	s := f.Field.Value
	mark := strings.Index(s, ".")
	if mark == -1 {
		name = s
		return
	}
	relation = s[:mark]
	name = s[mark+1:]
	return
}

// operator returns SQL operator.
func (f *Field) operator() (s string) {
	switch len(f.Value) {
	case 1:
		s = f.Operator.Value
		switch s {
		case string(COLON):
			s = "="
		case string(LIKE):
			s = "LIKE"
		}
	default:
		s = "IN"
	}

	return
}

// Assert -
type Assert struct {
	Field    string
	Kind     byte
	Relation bool
}

// assert validation.
func (r *Assert) assert(p *Predicate) (err error) {
	name := p.Field.Value
	switch r.Kind {
	case LITERAL:
		switch p.Operator.Value {
		case string(LIKE):
			err = Errorf("'~' cannot be used with '%s'", name)
			return
		}
	}
	if !r.Relation {
		if (&Field{*p}).Value.Operator(AND) {
			err = Errorf("(,,) cannot be used with '%s'.", name)
			return
		}
	}
	return
}

// AsValue returns the real value.
func AsValue(t Token) (object interface{}) {
	v := t.Value
	object = v
	switch t.Kind {
	case LITERAL:
		n, err := strconv.Atoi(v)
		if err == nil {
			object = n
			break
		}
		b, err := strconv.ParseBool(v)
		if err == nil {
			object = b
			break
		}
	default:
	}
	return
}
