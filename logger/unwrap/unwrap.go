package unwrap

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"go.unistack.org/micro/v3/codec"
)

const sf = "0-+# "

var hexDigits = "0123456789abcdef"

var (
	panicBytes         = []byte("(PANIC=")
	plusBytes          = []byte("+")
	iBytes             = []byte("i")
	trueBytes          = []byte("true")
	falseBytes         = []byte("false")
	interfaceBytes     = []byte("(interface {})")
	openBraceBytes     = []byte("{")
	closeBraceBytes    = []byte("}")
	asteriskBytes      = []byte("*")
	ampBytes           = []byte("&")
	colonBytes         = []byte(":")
	openParenBytes     = []byte("(")
	closeParenBytes    = []byte(")")
	spaceBytes         = []byte(" ")
	commaBytes         = []byte(",")
	pointerChainBytes  = []byte("->")
	nilAngleBytes      = []byte("<nil>")
	circularShortBytes = []byte("<shown>")
	invalidAngleBytes  = []byte("<invalid>")
	openBracketBytes   = []byte("[")
	closeBracketBytes  = []byte("]")
	percentBytes       = []byte("%")
	precisionBytes     = []byte(".")
	openAngleBytes     = []byte("<")
	closeAngleBytes    = []byte(">")
	openMapBytes       = []byte("{")
	closeMapBytes      = []byte("}")
)

type unwrap struct {
	val            interface{}
	s              fmt.State
	depth          int
	pointers       map[uintptr]int
	opts           *Options
	ignoreNextType bool
}

type Options struct {
	Codec   codec.Codec
	Indent  string
	Methods bool
	Tagged  bool
}

func NewOptions(opts ...Option) Options {
	options := Options{
		Indent:  " ",
		Methods: false,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

type Option func(*Options)

func Indent(f string) Option {
	return func(o *Options) {
		o.Indent = f
	}
}

func Methods(b bool) Option {
	return func(o *Options) {
		o.Methods = b
	}
}

func Codec(c codec.Codec) Option {
	return func(o *Options) {
		o.Codec = c
	}
}

func Tagged(b bool) Option {
	return func(o *Options) {
		o.Tagged = b
	}
}

func Unwrap(val interface{}, opts ...Option) *unwrap {
	options := NewOptions(opts...)
	return &unwrap{val: val, opts: &options, pointers: make(map[uintptr]int)}
}

func (f *unwrap) unpackValue(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface {
		f.ignoreNextType = false
		if !v.IsNil() {
			v = v.Elem()
		}
	}
	return v
}

// formatPtr handles formatting of pointers by indirecting them as necessary.
func (f *unwrap) formatPtr(v reflect.Value) {
	// Display nil if top level pointer is nil.
	showTypes := f.s.Flag('#')
	if v.IsNil() && (!showTypes || f.ignoreNextType) {
		_, _ = f.s.Write(nilAngleBytes)
		return
	}

	// Remove pointers at or below the current depth from map used to detect
	// circular refs.
	for k, depth := range f.pointers {
		if depth >= f.depth {
			delete(f.pointers, k)
		}
	}

	// Keep list of all dereferenced pointers to possibly show later.
	pointerChain := make([]uintptr, 0)

	// Figure out how many levels of indirection there are by derferencing
	// pointers and unpacking interfaces down the chain while detecting circular
	// references.
	nilFound := false
	cycleFound := false
	indirects := 0
	ve := v
	for ve.Kind() == reflect.Ptr {
		if ve.IsNil() {
			nilFound = true
			break
		}
		indirects++
		addr := ve.Pointer()
		pointerChain = append(pointerChain, addr)
		if pd, ok := f.pointers[addr]; ok && pd < f.depth {
			cycleFound = true
			indirects--
			break
		}
		f.pointers[addr] = f.depth

		ve = ve.Elem()
		if ve.Kind() == reflect.Interface {
			if ve.IsNil() {
				nilFound = true
				break
			}
			ve = ve.Elem()
		}
	}

	// Display type or indirection level depending on flags.
	if showTypes && !f.ignoreNextType {
		if f.depth > 0 {
			_, _ = f.s.Write(openParenBytes)
		}
		if f.depth > 0 {
			_, _ = f.s.Write(bytes.Repeat(asteriskBytes, indirects))
		} else {
			_, _ = f.s.Write(bytes.Repeat(ampBytes, indirects))
		}
		_, _ = f.s.Write([]byte(ve.Type().String()))
		if f.depth > 0 {
			_, _ = f.s.Write(closeParenBytes)
		}
	} else {
		if nilFound || cycleFound {
			indirects += strings.Count(ve.Type().String(), "*")
		}
		_, _ = f.s.Write(openAngleBytes)
		_, _ = f.s.Write([]byte(strings.Repeat("*", indirects)))
		_, _ = f.s.Write(closeAngleBytes)
	}

	// Display pointer information depending on flags.
	if f.s.Flag('+') && (len(pointerChain) > 0) {
		_, _ = f.s.Write(openParenBytes)
		for i, addr := range pointerChain {
			if i > 0 {
				_, _ = f.s.Write(pointerChainBytes)
			}
			getHexPtr(f.s, addr)
		}
		_, _ = f.s.Write(closeParenBytes)
	}

	// Display dereferenced value.
	switch {
	case nilFound:
		_, _ = f.s.Write(nilAngleBytes)

	case cycleFound:
		_, _ = f.s.Write(circularShortBytes)

	default:
		f.ignoreNextType = true
		f.format(ve)
	}
}

// format is the main workhorse for providing the Formatter interface.  It
// uses the passed reflect value to figure out what kind of object we are
// dealing with and formats it appropriately.  It is a recursive function,
// however circular data structures are detected and handled properly.
func (f *unwrap) format(v reflect.Value) {
	if f.opts.Codec != nil {
		buf, err := f.opts.Codec.Marshal(v.Interface())
		if err != nil {
			_, _ = f.s.Write(invalidAngleBytes)
			return
		}
		_, _ = f.s.Write(buf)
		return
	}
	// Handle invalid reflect values immediately.
	kind := v.Kind()
	if kind == reflect.Invalid {
		_, _ = f.s.Write(invalidAngleBytes)
		return
	}

	if (kind == reflect.Ptr) && (!reflect.Indirect(v).IsValid()) {
		f.formatPtr(v)
		return
	}

	// Handle pointers specially.
	if kind == reflect.Ptr {
		fmt.Printf("AAAA %s\n", reflect.Indirect(v).Type().String())
		switch reflect.Indirect(v).Type().String() {
		case "sql.NullBool":
			if eva := reflect.Indirect(v).FieldByName("Valid"); eva.IsValid() && eva.Bool() {
				v = reflect.Indirect(v).FieldByName("Bool")
				kind = v.Kind()
			}
		case "sql.NullByte":
			if eva := reflect.Indirect(v).FieldByName("Valid"); eva.IsValid() && eva.Bool() {
				v = reflect.Indirect(v).FieldByName("Byte")
				kind = v.Kind()
			}
		case "sql.NullFloat64":
			if eva := reflect.Indirect(v).FieldByName("Valid"); eva.IsValid() && eva.Bool() {
				v = reflect.Indirect(v).FieldByName("Float64")
				kind = v.Kind()
			}
		case "sql.NullInt16":
			if eva := reflect.Indirect(v).FieldByName("Valid"); eva.IsValid() && eva.Bool() {
				v = reflect.Indirect(v).FieldByName("Int16")
				kind = v.Kind()
			}
		case "sql.NullInt32":
			if eva := reflect.Indirect(v).FieldByName("Valid"); eva.IsValid() && eva.Bool() {
				v = reflect.Indirect(v).FieldByName("Int32")
				kind = v.Kind()
			}
		case "sql.NullInt64":
			if eva := reflect.Indirect(v).FieldByName("Valid"); eva.IsValid() && eva.Bool() {
				v = reflect.Indirect(v).FieldByName("Int64")
				kind = v.Kind()
			}
		case "sql.NullString":
			fmt.Printf("AAAAAAAAAA")
			if eva := reflect.Indirect(v).FieldByName("Valid"); eva.IsValid() && eva.Bool() {
				v = reflect.Indirect(v).FieldByName("String")
				kind = v.Kind()
			}
		case "sql.NullTime":
			if eva := reflect.Indirect(v).FieldByName("Valid"); eva.IsValid() && eva.Bool() {
				v = reflect.Indirect(v).FieldByName("Time")
				kind = v.Kind()
			}
		case "wrapperspb.BoolValue", "wrapperspb.BytesValue",
			"wrapperspb.DoubleValue", "wrapperspb.FloatValue",
			"wrapperspb.Int32Value", "wrapperspb.Int64Value",
			"wrapperspb.UInt32Value", "wrapperspb.UInt64Value",
			"wrapperspb.StringValue":
			if eva := reflect.Indirect(v).FieldByName("Value"); eva.IsValid() {
				v = eva
				kind = v.Kind()
			}
		default:
			f.formatPtr(v)
			return
		}
	}

	// get type information unless already handled elsewhere.
	if !f.ignoreNextType && f.s.Flag('#') {
		if v.Type().Kind() != reflect.Map &&
			v.Type().Kind() != reflect.String &&
			v.Type().Kind() != reflect.Array &&
			v.Type().Kind() != reflect.Slice {
			_, _ = f.s.Write(openParenBytes)
		}
		if v.Kind() != reflect.String {
			_, _ = f.s.Write([]byte(v.Type().String()))
		}
		if v.Type().Kind() != reflect.Map &&
			v.Type().Kind() != reflect.String &&
			v.Type().Kind() != reflect.Array &&
			v.Type().Kind() != reflect.Slice {
			_, _ = f.s.Write(closeParenBytes)
		}
	}
	f.ignoreNextType = false

	// Call Stringer/error interfaces if they exist and the handle methods
	// flag is enabled.
	if !f.opts.Methods {
		if (kind != reflect.Invalid) && (kind != reflect.Interface) {
			if handled := handleMethods(f.opts, f.s, v); handled {
				return
			}
		}
	}

	switch kind {
	case reflect.Invalid:
		_, _ = f.s.Write(invalidAngleBytes)
	case reflect.Bool:
		getBool(f.s, v.Bool())
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		getInt(f.s, v.Int(), 10)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		getUint(f.s, v.Uint(), 10)
	case reflect.Float32:
		getFloat(f.s, v.Float(), 32)
	case reflect.Float64:
		getFloat(f.s, v.Float(), 64)
	case reflect.Complex64:
		getComplex(f.s, v.Complex(), 32)
	case reflect.Complex128:
		getComplex(f.s, v.Complex(), 64)
	case reflect.Slice:
		if v.IsNil() {
			_, _ = f.s.Write(nilAngleBytes)
			break
		}
		fallthrough
	case reflect.Array:
		_, _ = f.s.Write(openBraceBytes)
		f.depth++
		numEntries := v.Len()
		for i := 0; i < numEntries; i++ {
			if i > 0 {
				_, _ = f.s.Write(commaBytes)
				_, _ = f.s.Write(spaceBytes)
			}
			f.ignoreNextType = true
			f.format(f.unpackValue(v.Index(i)))
		}
		f.depth--
		_, _ = f.s.Write(closeBraceBytes)
	case reflect.String:
		_, _ = f.s.Write([]byte(`"` + v.String() + `"`))
	case reflect.Interface:
		// The only time we should get here is for nil interfaces due to
		// unpackValue calls.
		if v.IsNil() {
			_, _ = f.s.Write(nilAngleBytes)
		}
	case reflect.Ptr:
		// Do nothing.  We should never get here since pointers have already
		// been handled above.
	case reflect.Map:
		// nil maps should be indicated as different than empty maps
		if v.IsNil() {
			_, _ = f.s.Write(nilAngleBytes)
			break
		}
		_, _ = f.s.Write(openMapBytes)
		f.depth++
		keys := v.MapKeys()
		for i, key := range keys {
			if i > 0 {
				_, _ = f.s.Write(spaceBytes)
			}
			f.ignoreNextType = true
			f.format(f.unpackValue(key))
			_, _ = f.s.Write(colonBytes)
			f.ignoreNextType = true
			f.format(f.unpackValue(v.MapIndex(key)))
		}
		f.depth--
		_, _ = f.s.Write(closeMapBytes)
	case reflect.Struct:
		numFields := v.NumField()
		_, _ = f.s.Write(openBraceBytes)
		f.depth++
		vt := v.Type()
		prevSkip := false
		for i := 0; i < numFields; i++ {
			sv, ok := vt.Field(i).Tag.Lookup("logger")
			if ok {
				if sv == "omit" {
					prevSkip = true
					continue
				}
			} else if f.opts.Tagged {
				prevSkip = true
				continue
			}
			if i > 0 && !prevSkip {
				_, _ = f.s.Write(commaBytes)
				_, _ = f.s.Write(spaceBytes)
			}
			if prevSkip {
				prevSkip = false
			}
			vtf := vt.Field(i)
			if f.s.Flag('+') || f.s.Flag('#') {
				_, _ = f.s.Write([]byte(vtf.Name))
				_, _ = f.s.Write(colonBytes)
			}
			f.format(f.unpackValue(v.Field(i)))
		}
		f.depth--
		_, _ = f.s.Write(closeBraceBytes)
	case reflect.Uintptr:
		getHexPtr(f.s, uintptr(v.Uint()))
	case reflect.UnsafePointer, reflect.Chan, reflect.Func:
		getHexPtr(f.s, v.Pointer())
	// There were not any other types at the time this code was written, but
	// fall back to letting the default fmt package handle it if any get added.
	default:
		format := f.buildDefaultFormat()
		if v.CanInterface() {
			_, _ = fmt.Fprintf(f.s, format, v.Interface())
		} else {
			_, _ = fmt.Fprintf(f.s, format, v.String())
		}
	}
}

func (f *unwrap) Format(s fmt.State, verb rune) {
	f.s = s

	// Use standard formatting for verbs that are not v.
	if verb != 'v' {
		format := f.constructOrigFormat(verb)
		_, _ = fmt.Fprintf(s, format, f.val)
		return
	}

	if f.val == nil {
		if s.Flag('#') {
			_, _ = s.Write(interfaceBytes)
		}
		_, _ = s.Write(nilAngleBytes)
		return
	}

	f.format(reflect.ValueOf(f.val))
}

// handle special methods like error.Error() or fmt.Stringer interface
func handleMethods(_ *Options, w io.Writer, v reflect.Value) (handled bool) {
	if !v.CanInterface() {
		// not our case
		return false
	}

	if !v.CanAddr() {
		// not our case
		return false
	}

	if v.CanAddr() {
		v = v.Addr()
	}

	// Is it an error or Stringer?
	switch iface := v.Interface().(type) {
	case error:
		defer catchPanic(w, v)
		_, _ = w.Write([]byte(iface.Error()))
		return true
	case fmt.Stringer:
		defer catchPanic(w, v)
		_, _ = w.Write([]byte(iface.String()))
		return true
	}

	return false
}

// getBool outputs a boolean value as true or false to Writer w.
func getBool(w io.Writer, val bool) {
	if val {
		_, _ = w.Write(trueBytes)
	} else {
		_, _ = w.Write(falseBytes)
	}
}

// getInt outputs a signed integer value to Writer w.
func getInt(w io.Writer, val int64, base int) {
	_, _ = w.Write([]byte(strconv.FormatInt(val, base)))
}

// getUint outputs an unsigned integer value to Writer w.
func getUint(w io.Writer, val uint64, base int) {
	_, _ = w.Write([]byte(strconv.FormatUint(val, base)))
}

// getFloat outputs a floating point value using the specified precision,
// which is expected to be 32 or 64bit, to Writer w.
func getFloat(w io.Writer, val float64, precision int) {
	_, _ = w.Write([]byte(strconv.FormatFloat(val, 'g', -1, precision)))
}

// getComplex outputs a complex value using the specified float precision
// for the real and imaginary parts to Writer w.
func getComplex(w io.Writer, c complex128, floatPrecision int) {
	r := real(c)
	_, _ = w.Write(openParenBytes)
	_, _ = w.Write([]byte(strconv.FormatFloat(r, 'g', -1, floatPrecision)))
	i := imag(c)
	if i >= 0 {
		_, _ = w.Write(plusBytes)
	}
	_, _ = w.Write([]byte(strconv.FormatFloat(i, 'g', -1, floatPrecision)))
	_, _ = w.Write(iBytes)
	_, _ = w.Write(closeParenBytes)
}

// getHexPtr outputs a uintptr formatted as hexadecimal with a leading '0x'
// prefix to Writer w.
func getHexPtr(w io.Writer, p uintptr) {
	// Null pointer.
	num := uint64(p)
	if num == 0 {
		_, _ = w.Write(nilAngleBytes)
		return
	}

	// Max uint64 is 16 bytes in hex + 2 bytes for '0x' prefix
	buf := make([]byte, 18)

	// It's simpler to construct the hex string right to left.
	base := uint64(16)
	i := len(buf) - 1
	for num >= base {
		buf[i] = hexDigits[num%base]
		num /= base
		i--
	}
	buf[i] = hexDigits[num]

	// Add '0x' prefix.
	i--
	buf[i] = 'x'
	i--
	buf[i] = '0'

	// Strip unused leading bytes.
	buf = buf[i:]
	_, _ = w.Write(buf)
}

func catchPanic(w io.Writer, _ reflect.Value) {
	if err := recover(); err != nil {
		_, _ = w.Write(panicBytes)
		_, _ = fmt.Fprintf(w, "%v", err)
		_, _ = w.Write(closeParenBytes)
	}
}

func (f *unwrap) buildDefaultFormat() (format string) {
	buf := bytes.NewBuffer(percentBytes)

	for _, flag := range sf {
		if f.s.Flag(int(flag)) {
			_, _ = buf.WriteRune(flag)
		}
	}

	_, _ = buf.WriteRune('v')

	format = buf.String()
	return format
}

func (f *unwrap) constructOrigFormat(verb rune) (format string) {
	buf := bytes.NewBuffer(percentBytes)

	for _, flag := range sf {
		if f.s.Flag(int(flag)) {
			_, _ = buf.WriteRune(flag)
		}
	}

	if width, ok := f.s.Width(); ok {
		_, _ = buf.WriteString(strconv.Itoa(width))
	}

	if precision, ok := f.s.Precision(); ok {
		_, _ = buf.Write(precisionBytes)
		_, _ = buf.WriteString(strconv.Itoa(precision))
	}

	_, _ = buf.WriteRune(verb)

	format = buf.String()
	return format
}
