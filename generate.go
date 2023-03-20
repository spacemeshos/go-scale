package scale

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"html/template"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	header = `// Code generated by github.com/spacemeshos/go-scale/scalegen. DO NOT EDIT.

	// nolint
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
		encode: `{
		n, err := scale.Encode{{ .ScaleType }}(enc, t.{{ .Name }}{{.ScaleTypeArgs}})
		if err != nil {
			return total, err
		}
		total += n
		}
		`,
		decode: `{
		field, n, err := scale.Decode{{ .ScaleType }}{{ .TypeInfo }}(dec{{.ScaleTypeArgs}})
		if err != nil {
			return total, err
		}
		total += n
		t.{{ .Name }} = field
		}
		`,
	}
	genericTyped = temp{
		encode: `{
		n, err := scale.Encode{{ .ScaleType }}(enc, {{.EncodeModifier}}(t.{{ .Name }}){{.ScaleTypeArgs}})
		if err != nil {
			return total, err
		}
		total += n
		}
		`,
		decode: `{
		field, n, err := scale.Decode{{ .ScaleType }}{{ .TypeInfo }}(dec{{.ScaleTypeArgs}})
		if err != nil {
			return total, err
		}
		total += n
		t.{{ .Name }} = {{.DecodeModifier}}(field)
		}
		`,
	}
	array = temp{
		encode: `{
		n, err := scale.Encode{{ .ScaleType }}(enc, t.{{ .Name }}[:]{{.ScaleTypeArgs}})
		if err != nil {
			return total, err
		}
		total += n
		}
		`,
		decode: `{
		n, err := scale.Decode{{ .ScaleType }}(dec, t.{{ .Name }}[:]{{.ScaleTypeArgs}})
		if err != nil {
			return total, err
		}
		total += n
		}
		`,
	}
	object = temp{
		encode: `{
	n, err := t.{{ .Name }}.EncodeScale(enc)
		if err != nil {
			return total, err
		}
		total += n
		}
		`,
		decode: `{
		n, err := t.{{ .Name }}.DecodeScale(dec)
		if err != nil {
			return total, err
		}
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
			if !skipPackageImport(field.Type) {
				rst[canonicalPath(field.Type)] = struct{}{}
			}
		}
	}
	return rst
}

func skipPackageImport(typ reflect.Type) bool {
	if typ.Kind() == reflect.Struct {
		return true
	}
	if typ.Kind() == reflect.Slice {
		if typ.Elem().Kind() == reflect.Slice && typ.Elem().Elem().Kind() == reflect.Uint8 {
			return true
		}
		if typ.Elem().Kind() == reflect.String {
			return true
		}
		if typ.Elem().Kind() == reflect.Uint8 {
			return true
		}
		return false
	}
	if typ.Kind() == reflect.Array {
		return true
	}
	return false
}

func canonicalPath(typ reflect.Type) string {
	if typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Array || typ.Kind() == reflect.Slice {
		return typ.Elem().PkgPath()
	}
	return typ.PkgPath()
}

func private(field reflect.StructField) bool {
	return strings.ToUpper(field.Name[:1]) != field.Name[:1]
}

func sameModule(a, b reflect.Type) bool {
	if b.Kind() == reflect.Ptr || b.Kind() == reflect.Slice || b.Kind() == reflect.Array {
		b = b.Elem()
	}
	if a.Kind() == reflect.Ptr || a.Kind() == reflect.Slice || a.Kind() == reflect.Array {
		a = a.Elem()
	}
	return a.PkgPath() == b.PkgPath()
}

func builtin(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Array || typ.Kind() == reflect.Slice {
		return typ.Elem().PkgPath() == ""
	}
	return typ.PkgPath() == ""
}

func generateType(w io.Writer, gc *genContext, obj interface{}) error {
	typ := reflect.TypeOf(obj)
	tc := &typeContext{
		Name:          typ.Name(),
		Type:          typ,
		ParentPackage: typ.PkgPath(),
		TypeName:      typ.Name(),
	}
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
	Name           string
	ScaleType      string
	ScaleTypeArgs  string
	EncodeModifier string
	DecodeModifier string
	TypeName       string
	TypeInfo       string
	Type           reflect.Type
	ParentPackage  string
}

type scaleType struct {
	Name           string
	Args           string
	EncodeModifier string
	DecodeModifier string
}

func getDecodeModifier(parentType reflect.Type, field reflect.StructField) string {
	if sameModule(field.Type, parentType) {
		parts := strings.Split(field.Type.String(), ".")
		return parts[len(parts)-1]
	}
	return field.Type.String()
}

func getScaleType(parentType reflect.Type, field reflect.StructField) (scaleType, error) {
	decodeModifier := getDecodeModifier(parentType, field)

	switch field.Type.Kind() {
	case reflect.Bool:
		return scaleType{Name: "Bool"}, nil
	case reflect.String:
		maxElements, err := getMaxElements(field.Tag)
		if err != nil {
			return scaleType{}, fmt.Errorf("scale tag has incorrect max value: %w", err)
		}
		if maxElements == 0 {
			return scaleType{}, fmt.Errorf("strings must have max scale tag")
		}
		return scaleType{
			Name:           "StringWithLimit",
			Args:           fmt.Sprintf(", %d", maxElements),
			EncodeModifier: "string",
			DecodeModifier: decodeModifier,
		}, nil
	case reflect.Uint8:
		return scaleType{Name: "Compact8", EncodeModifier: "uint8", DecodeModifier: decodeModifier}, nil
	case reflect.Uint16:
		return scaleType{Name: "Compact16", EncodeModifier: "uint16", DecodeModifier: decodeModifier}, nil
	case reflect.Uint32:
		return scaleType{Name: "Compact32", EncodeModifier: "uint32", DecodeModifier: decodeModifier}, nil
	case reflect.Uint64:
		return scaleType{Name: "Compact64", EncodeModifier: "uint64", DecodeModifier: decodeModifier}, nil
	case reflect.Struct:
		return scaleType{Name: "Object"}, nil
	case reflect.Ptr:
		return scaleType{Name: "Option"}, nil
	case reflect.Slice:
		if field.Type.Elem().Kind() == reflect.Slice && field.Type.Elem().Elem().Kind() == reflect.Uint8 {
			// [][]byte
			return scaleType{}, fmt.Errorf("nested slices are not supported")
		}
		if field.Type.Elem().Kind() == reflect.String {
			// []string
			return scaleType{}, fmt.Errorf("string slices are not supported")
		}
		maxElements, err := getMaxElements(field.Tag)
		if err != nil {
			return scaleType{}, fmt.Errorf("scale tag has incorrect max value: %w", err)
		}
		if maxElements == 0 {
			return scaleType{}, fmt.Errorf("slices must have max scale tag")
		}
		if field.Type.Elem().Kind() == reflect.Uint8 {
			return scaleType{Name: "ByteSliceWithLimit", Args: fmt.Sprintf(", %d", maxElements)}, nil
		}
		return scaleType{Name: "StructSliceWithLimit", Args: fmt.Sprintf(", %d", maxElements)}, nil
	case reflect.Array:
		if field.Type.Elem().Kind() == reflect.Uint8 {
			return scaleType{Name: "ByteArray"}, nil
		}
		return scaleType{Name: "StructArray"}, nil
	}
	return scaleType{}, fmt.Errorf("type %v is not supported", field.Type.Kind())
}

func getMaxElements(tag reflect.StructTag) (uint32, error) {
	scaleTagValue, exists := tag.Lookup("scale")
	if !exists || scaleTagValue == "" {
		return 0, nil
	}
	if scaleTagValue == "" {
		return 0, errors.New("scale tag is not defined")
	}
	pairs := strings.Split(scaleTagValue, ",")
	if len(pairs) == 0 {
		return 0, errors.New("no max value found in scale tag")
	}
	var maxElementsStr string
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		data := strings.Split(pair, "=")
		if len(data) < 2 {
			continue
		}
		if data[0] != "max" {
			continue
		}
		maxElementsStr = strings.TrimSpace(data[1])
		break
	}
	if maxElementsStr == "" {
		return 0, errors.New("no max value found in scale tag")
	}
	maxElements, err := strconv.Atoi(maxElementsStr)
	if err != nil {
		return 0, fmt.Errorf("parsing max value: %w", err)
	}
	return uint32(maxElements), nil
}

func getTemplate(stype scaleType) temp {
	switch {
	case stype.Name == "StructArray":
		return array
	case stype.Name == "ByteArray":
		return array
	case stype.Name == "Object":
		return object
	case stype.EncodeModifier != "":
		return genericTyped
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

		scaleType, err := getScaleType(typ, field)
		if err != nil {
			return fmt.Errorf("getting scale type for %s: %w", typ, err)
		}

		tctx := &typeContext{
			Name:           field.Name,
			Type:           field.Type,
			TypeName:       fullTypeName(gc, tc, field.Type),
			ScaleType:      scaleType.Name,
			ScaleTypeArgs:  scaleType.Args,
			EncodeModifier: scaleType.EncodeModifier,
			DecodeModifier: scaleType.DecodeModifier,
			ParentPackage:  tc.ParentPackage,
		}

		if strings.HasPrefix(scaleType.Name, "StructSlice") {
			tctx.TypeInfo = "[" + fullTypeName(gc, tc, field.Type.Elem()) + "]"
		} else if strings.Contains(scaleType.Name, "Struct") || strings.Contains(scaleType.Name, "Option") {
			tctx.TypeInfo = fmt.Sprintf("[%v]", tctx.TypeName)
		}
		if err := executeTemplate(w, getAction(getTemplate(scaleType), action), tctx); err != nil {
			return err
		}
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
