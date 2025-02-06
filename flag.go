package flag

import (
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type configField struct {
	pointer     interface{}
	desired     interface{}
	key         string
	env         string
	longFlag    string
	description string
}

// FlagSet inherits flag.FlagSet and implements the struct-based flag generation and can map env values to the
// flag automatically.
// Additional list of methods:
// * AddOption: generate flags for a struct.
// * Add: generate a flag for a signal value.
// * Parse: parse the command-line flags from args.
type FlagSet struct { // nolint
	*flag.FlagSet
	output io.Writer
	fields map[string]*configField
	inited bool
}

// NewFlagSet returns a new, empty flag set with the given flag.FlagSet.
func NewFlagSet(fs *flag.FlagSet) *FlagSet {
	if fs == nil {
		fs = flag.CommandLine
	}

	return &FlagSet{
		FlagSet: fs,
		fields:  map[string]*configField{},
	}
}

// Parse parses the command-line flags from args.
// Must be called after all flags are defined and before flags
func (fs *FlagSet) Parse(args []string) error {
	if !fs.inited {
		fs.initFlags()
		fs.inited = true
	}
	return fs.FlagSet.Parse(args)
}

// AddOption will fill up options from flags/ENV after Parse.
//
// the option must be a pointer to struct.
//
// Here is an example:
//
//		type Option struct {
//		    FirstName string `desc:"Desc for First Name"`
//		    Age       uint `desc:"Desc for Age"`
//		    User struct{
//	           Name string `desc:"Desc for User Name"`
//			}
//		}
//
// The struct has two fields (with prefix example):
//
//	Field       Flag                   ENV
//	FirstName   --example-first-name   EXAMPLE_FIRST_NAME
//	Age         --example-age          EXAMPLE_AGE
//	User.Name   --example-user-name    EXAMPLE_USER_NAME
//
// When you execute command with `--help`, you can see the help doc of flags and
// descriptions (From field tag `desc`).
//
// The priority is:
//
//	Flag > ENV > The value you set in option
func (fs *FlagSet) AddOption(prefix string, options ...interface{}) {
	for _, opt := range options {
		val := reflect.ValueOf(opt)
		typ := val.Type()
		if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
			panic(fmt.Errorf("%s is not a pointer to struct", typ.String()))
		}
		if val.IsNil() {
			panic(fmt.Errorf("%s should not be nil", typ.String()))
		}
		val = val.Elem()
		typ = val.Type()
		walkthrough(prefix, typ, val, func(subPrefix string, field reflect.StructField, fieldVal reflect.Value) {
			ptr := fieldVal.Addr().Interface()
			if field.Type.Kind() == reflect.Ptr {
				if fieldVal.IsNil() {
					newVal := reflect.New(fieldVal.Type().Elem())
					fieldVal.Set(newVal)
				}
				ptr = fieldVal.Interface()
			}
			fs.Add(ptr, subPrefix+field.Name, field.Tag.Get("desc"))
		})
	}
}

// Add adds a field by key.
// If you don't have any struct to describe an option, you can use the method to
// add a single field into command.
// `pointer` must be a pointer to golang basic data type (e.g. *int, *string).
// `key` must be a unique value for each option.
// If you want a short flag for the field, you can only set a one-char string.
// `desc` describes the field.
func (fs *FlagSet) Add(pointer interface{}, key string, desc string) {
	if pointer == nil || reflect.ValueOf(pointer).IsNil() {
		panic(fmt.Errorf("pointer of %s should not be nil", key))
	}

	cf := &configField{
		key:         key,
		pointer:     pointer,
		description: desc,
	}

	// key: prefix+field.Name, e.g. EnV + OpAB
	// after split: [En V Op AB]
	keyParts := splitKey(key)

	// en-v-op-ab
	cf.longFlag = strings.ToLower(strings.Join(keyParts, "-"))
	// EN_V_OP_AB
	cf.env = strings.ToUpper(strings.Join(keyParts, "_"))

	if _, ok := fs.fields[cf.key]; ok {
		panic(fmt.Errorf("%s has been registered", cf.key))
	}
	fs.fields[cf.key] = cf
}

// initFlags will generate the flags and mapping the value.
// you need call this function before call the flag.Parse.
func (fs *FlagSet) initFlags() {
	for _, f := range fs.fields {
		var envValue interface{}
		env := os.Getenv(f.env)

		switch v := f.pointer.(type) {
		case *uint:
			if env != "" {
				value, err := strconv.ParseUint(env, 10, 64)
				if err != nil {
					panic(fmt.Errorf("env %s is not a valid uint", f.env))
				}
				envValue = uint(value)
			}
			f.desired = chooseValue(envValue, *v)
			fs.UintVar(v, f.longFlag, f.desired.(uint), f.description)
		case *uint64:
			if env != "" {
				value, err := strconv.ParseUint(env, 10, 64)
				if err != nil {
					panic(fmt.Errorf("env %s is not a valid uint64", f.env))
				}
				envValue = value
			}
			f.desired = chooseValue(envValue, *v)
			fs.Uint64Var(v, f.longFlag, f.desired.(uint64), f.description)
		case *int:
			if env != "" {
				value, err := strconv.ParseInt(env, 10, 64)
				if err != nil {
					panic(fmt.Errorf("env %s is not a valid int", f.env))
				}
				envValue = int(value)
			}
			f.desired = chooseValue(envValue, *v)
			fs.IntVar(v, f.longFlag, f.desired.(int), f.description)
		case *int64:
			if env != "" {
				value, err := strconv.ParseInt(env, 10, 64)
				if err != nil {
					panic(fmt.Errorf("env %s is not a valid int64", f.env))
				}
				envValue = value
			}
			f.desired = chooseValue(envValue, *v)
			fs.Int64Var(v, f.longFlag, f.desired.(int64), f.description)
		case *float64:
			if env != "" {
				value, err := strconv.ParseFloat(env, 64)
				if err != nil {
					panic(fmt.Errorf("env %s is not a valid float64", f.env))
				}
				envValue = value
			}
			f.desired = chooseValue(envValue, *v)
			fs.Float64Var(v, f.longFlag, f.desired.(float64), f.description)
		case *string:
			if env != "" {
				envValue = env
			}
			f.desired = chooseValue(envValue, *v)
			fs.StringVar(v, f.longFlag, f.desired.(string), f.description)
		case *bool:
			if env != "" {
				value, err := strconv.ParseBool(env)
				if err != nil {
					panic(fmt.Errorf("env %s is not a valid bool", f.env))
				}
				envValue = value
			}
			f.desired = chooseValue(envValue, *v)
			fs.BoolVar(v, f.longFlag, f.desired.(bool), f.description)
		case *time.Duration:
			if env != "" {
				value, err := time.ParseDuration(env)
				if err != nil {
					panic(fmt.Errorf("env %s is not a valid time.Duration", f.env))
				}
				envValue = value
			}
			f.desired = chooseValue(envValue, *v)
			fs.DurationVar(v, f.longFlag, f.desired.(time.Duration), f.description)
		default:
			panic(fmt.Errorf("unrecognized type %s for %s", reflect.TypeOf(f.pointer).String(), f.key))
		}
	}

	p := NewTable(30)
	rows := make([][]interface{}, 0, len(fs.fields))
	for _, f := range fs.fields {
		rows = append(rows, []interface{}{0, f.env, "--" + f.longFlag, f.desired})
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][1].(string) < rows[j][1].(string)
	})
	p.AddRow("", "ENV", "Flag", "Current Value")
	for i, row := range rows {
		row[0] = i + 1
		p.AddRow(row...)
	}
	fs.Usage = func() {
		fs.PrintDefaults()
		fmt.Fprintf(fs.out(), "\nENV-Flag Mapping Table\n%s\n", p.String())
	}
}

// SetOutput sets the destination for usage and error messages.
func (fs *FlagSet) SetOutput(output io.Writer) {
	fs.output = output
	fs.FlagSet.SetOutput(output)
}

func (f *FlagSet) out() io.Writer {
	if f.output == nil {
		return os.Stderr
	}
	return f.output
}

func walkthrough(prefix string, typ reflect.Type, val reflect.Value, f func(subPrefix string, field reflect.StructField, fieldVal reflect.Value)) {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)
		if field.Anonymous { // must be struct
			if field.Type.Kind() == reflect.Ptr {
				if fieldVal.IsNil() {
					newVal := reflect.New(fieldVal.Type().Elem())
					fieldVal.Set(newVal)
				}
				fieldVal = fieldVal.Elem()
			}
			walkthrough(prefix, fieldVal.Type(), fieldVal, f)
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			if fieldVal.IsNil() {
				newVal := reflect.New(fieldVal.Type().Elem())
				fieldVal.Set(newVal)
			}
			fieldVal = fieldVal.Elem()
			walkthrough(prefix+field.Name, fieldVal.Type(), fieldVal, f)
		} else if field.Type.Kind() == reflect.Struct {
			walkthrough(prefix+field.Name, fieldVal.Type(), fieldVal, f)
		} else {
			if field.Name != "" && field.Name[0] >= 'A' && field.Name[0] <= 'Z' {
				f(prefix, field, fieldVal)
			}
		}
	}
}

// splitKey splits key to parts.
func splitKey(key string) []string {
	parts := make([]string, 0)

	lastIsCapital := true
	lastIndex := 0
	for i, char := range key {
		if char >= '0' && char <= '9' {
			// Numbers inherit last char.
			continue
		}
		currentIsCapital := char >= 'A' && char <= 'Z'
		if i > 0 && lastIsCapital != currentIsCapital {
			end := 0
			if currentIsCapital {
				end = i
			} else {
				end = i - 1
			}
			if end > lastIndex {
				parts = append(parts, key[lastIndex:end])
				lastIndex = end
			}
		}
		lastIsCapital = currentIsCapital
	}
	if lastIndex < len(key) {
		parts = append(parts, key[lastIndex:])
	}
	return parts
}

// chooseValue chooses expected value form multiple configurations.
// Priority: ENV > Default
func chooseValue(envValue interface{}, defaultValue interface{}) interface{} {
	if envValue != nil {
		return envValue
	}
	return defaultValue
}
