package gen

var (
	encodeMethodStart = `func (t *{{ .Name }}) EncodeScale(enc *scale.Encoder) (int,  error) {
		var total int
	`

	encodeTemplate = `
	if n, err := scale.Encode{{ .ScaleType }}(enc, t.{{ .Name }}); err != nil {
		return err
	} else {
		total += n
	}
	`
)
