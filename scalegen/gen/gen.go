package gen

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"html/template"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
)

const (
	header = `// Code generated by github.com/spacemeshos/go-scale/scalegen. DO NOT EDIT.

	package {{ .Package }}
	
	import (
		"github.com/spacemeshos/go-scale"
		{{ range $pkg, $short := .Imported }}"{{ $pkg }}"
        {{ end }}
	)
	`
)

type temp struct {
	encode string
	decode string
}

const (
	encode = iota
	decode
)

var (
	start = temp{
		encode: `func (t *{{ .Name }}) EncodeScale(enc *scale.Encoder) (total int, err error) {
		`,
		decode: `func (t *{{ .Name }}) DecodeScale(dec *scale.Decoder) (total int,  err error) {
		`,
	}
	generic = temp{
		encode: `if n, err := scale.Encode{{ .ScaleType }}(enc, t.{{ .Name }}); err != nil {
			return total, err
		} else {
			total += n
		}
		`,
		decode: `if field, n, err := scale.Decode{{ .ScaleType }}{{ .TypeInfo }}(dec); err != nil {
			return total, err
		} else {
			total += n
			t.{{ .Name }} = field
		}
		`,
	}
	array = temp{
		encode: `if n, err := scale.Encode{{ .ScaleType }}(enc, t.{{ .Name }}[:]); err != nil {
			return total, err
		} else {
			total += n
		}
		`,
		decode: `if n, err := scale.Decode{{ .ScaleType }}(dec, t.{{ .Name }}[:]); err != nil {
			return total, err
		} else {
			total += n
		}
		`,
	}
	object = temp{
		encode: `if n, err := t.{{ .Name }}.EncodeScale(enc); err != nil {
			return total, err
		} else {
			total += n
		}
		`,
		decode: `if n, err := t.{{ .Name }}.DecodeScale(dec); err != nil {
			return total, err
		} else {
			total += n
		}
		`,
	}
)

func getAction(tm temp, action int) string {
	switch action {
	case encode:
		return tm.encode
	case decode:
		return tm.decode
	}
	panic("unreachable")
}

func Generate(pkg string, filepath string, objs ...interface{}) error {
	buf := bytes.NewBuffer(nil)
	ctx := &genContext{Package: pkg, Imported: generateImports(objs...)}

	err := executeTemplate(buf, header, ctx)
	if err != nil {
		return err
	}

	for _, obj := range objs {
		if err := generateType(buf, ctx, obj); err != nil {
			return err
		}
	}
	data := buf.Bytes()
	data, err = format.Source(data)
	if err != nil {
		return fmt.Errorf("can't format: \ndata: %s\n err:%w", buf.Bytes(), err)
	}
	return os.WriteFile(filepath, data, 0o664)
}

func generateImports(objs ...interface{}) map[string]struct{} {
	rst := map[string]struct{}{}
	for _, obj := range objs {
		typ := reflect.TypeOf(obj)
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if sameModule(field.Type, typ) {
				continue
			}
			if private(field) {
				continue
			}
			if builtin(field.Type) {
				continue
			}
			if needsImport(field.Type) {
				rst[canonicalPath(field.Type)] = struct{}{}
			}
		}
	}
	return rst
}

func canonicalPath(typ reflect.Type) string {
	if typ.Kind() == reflect.Ptr {
		return typ.Elem().PkgPath()
	}
	return typ.PkgPath()
}

func private(field reflect.StructField) bool {
	return strings.ToUpper(field.Name[:1]) != field.Name[:1]
}

func sameModule(a, b reflect.Type) bool {
	if b.Kind() == reflect.Ptr {
		b = b.Elem()
	}
	if a.Kind() == reflect.Ptr {
		a = a.Elem()
	}
	return a.PkgPath() == b.PkgPath()
}

func builtin(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr {
		return typ.Elem().PkgPath() == ""
	}
	return typ.PkgPath() == ""
}

func needsImport(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Ptr:
		return true
	case reflect.Slice:
		return true
	}
	return false
}

func generateType(w io.Writer, gc *genContext, obj interface{}) error {
	typ := reflect.TypeOf(obj)
	tc := &typeContext{
		Name:          typ.Name(),
		Type:          typ,
		ParentPackage: typ.PkgPath(),
		TypeName:      typ.Name(),
	}
	log.Printf("generating codec for type %+v", tc)
	if err := executeAction(encode, w, gc, tc); err != nil {
		return err
	}
	if err := executeAction(decode, w, gc, tc); err != nil {
		return err
	}
	return nil
}

type genContext struct {
	Package  string
	Imported map[string]struct{} // full path to shortname
}

type typeContext struct {
	Name          string
	ScaleType     string
	TypeName      string
	TypeInfo      string
	Type          reflect.Type
	ParentPackage string
}

func getScaleType(t reflect.Type) (string, error) {
	switch t.Kind() {
	case reflect.Bool:
		return "Bool", nil
	case reflect.Uint8:
		return "Compact8", nil
	case reflect.Uint16:
		return "Compact16", nil
	case reflect.Uint32:
		return "Compact32", nil
	case reflect.Uint64:
		return "Compact64", nil
	case reflect.Struct:
		return "Object", nil
	case reflect.Ptr:
		switch t.Elem().Kind() {
		case reflect.Array:
			return "", errors.New("ptr to array is not supported")
		case reflect.Slice:
			return "", errors.New("ptr to slice is not supported")
		}
		return "Option", nil
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return "ByteSlice", nil
		}
		return "StructSlice", nil
	case reflect.Array:
		if t.Elem().Kind() == reflect.Uint8 {
			return "ByteArray", nil
		}
		return "StructArray", nil
	}
	return "", fmt.Errorf("type %v is not supported", t.Kind())
}

func getTemplate(stype string) temp {
	switch stype {
	case "StructArray", "ByteArray":
		return array
	case "Object":
		return object
	default:
		return generic
	}
}

func executeAction(action int, w io.Writer, gc *genContext, tc *typeContext) error {
	typ := tc.Type
	if err := executeTemplate(w, getAction(start, action), tc); err != nil {
		return err
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if private(field) {
			continue
		}

		stype, err := getScaleType(field.Type)
		if err != nil {
			return err
		}

		tctx := &typeContext{
			Name:          field.Name,
			Type:          field.Type,
			TypeName:      fullTypeName(gc, tc, field.Type),
			ScaleType:     stype,
			ParentPackage: tc.ParentPackage,
		}

		if stype == "StructSlice" {
			tctx.TypeInfo = "[" + field.Type.Elem().Name() + "]"
		} else if strings.Contains(stype, "Struct") || strings.Contains(stype, "Option") {
			tctx.TypeInfo = fmt.Sprintf("[%v]", tctx.TypeName)
		}
		log.Printf("type context %+v", tctx)
		fmt.Fprintf(w, "// field %v (%d)\n", field.Name, i)
		if err := executeTemplate(w, getAction(getTemplate(stype), action), tctx); err != nil {
			return err
		}
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w, "return total, nil")
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
	return nil
}

func fullTypeName(gc *genContext, tc *typeContext, typ reflect.Type) string {
	pkg := typ.PkgPath()
	name := typ.Name()
	str := typ.String()
	if typ.Kind() == reflect.Ptr {
		pkg = typ.Elem().PkgPath()
		name = typ.Elem().Name()
		str = typ.Elem().String()
	}

	if typ.Kind() == reflect.Slice {
		if typ.Elem().PkgPath() == tc.ParentPackage {
			return "[]" + typ.Elem().Name()
		}
	}
	if tc.ParentPackage == pkg {
		return name
	}
	return str
}

func executeTemplate(w io.Writer, text string, ctx interface{}) error {
	tpl, err := template.New("").Parse(text)
	if err != nil {
		return err
	}
	return tpl.Execute(w, ctx)
}
