package vecty

import (
	"reflect"
	"strings"

	"github.com/gopherjs/gopherwasm/js"
)

const (
	structTag = "js"

	structTagOptionIncludeEmpty = "includeEmpty"
)

// Value Returns the js value of a type
func Value(p interface{}) js.Value {
	t := reflect.TypeOf(p)
	rv := reflect.ValueOf(p)

	switch t.Kind() {
	case reflect.Struct:
		v := js.Global().Get("Object").New()
		structValue(v, p)
		return v
	case reflect.Ptr:
		return Value(rv.Addr().Interface())
	default:
		return js.ValueOf(p)
	}
}

func structValue(v js.Value, p interface{}) {
	t := reflect.TypeOf(p)
	rv := reflect.ValueOf(p)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fv := rv.Field(i)

		fn := field.Name

		tag := strings.Split(field.Tag.Get(structTag), ",")

		if len(tag[0]) > 0 {
			fn = tag[0]
		}

		if field.Anonymous {
			structValue(v, fv.Interface())
			continue
		}

		includeEmpty := len(tag) > 1 && tag[1] == structTagOptionIncludeEmpty

		if !includeEmpty && fv.Interface() == reflect.Zero(field.Type).Interface() {
			continue
		}

		v.Set(fn, Value(fv.Interface()))
	}
}
