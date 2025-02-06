package flag

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type NestedStruct struct {
	NestedString string `desc:"Nested string field"`
	NestedInt    int    `desc:"Nested int field"`
}

type EmbeddedStruct struct {
	EmbeddedString string `desc:"Embedded string field"`
	EmbeddedInt    int    `desc:"Embedded int field"`
}

type EmbeddedStructPtr struct {
	EmbeddedPtrString string `desc:"Embedded string field"`
}

type ComplexOption struct {
	String    string        `desc:"String field"`
	Int       int           `desc:"Int field"`
	Uint      uint          `desc:"Uint field"`
	Float64   float64       `desc:"Float64 field"`
	Bool      bool          `desc:"Bool field"`
	Duration  time.Duration `desc:"Duration field"`
	Nested    NestedStruct  `desc:"Nested struct field"`
	NestedPtr *NestedStruct `desc:"Nested struct pointer field"`
	PtrField  *string       `desc:"Pointer to string field"`
	// embedded 结构体
	EmbeddedStruct `desc:"Embedded struct field"`
	// embedded 结构体指针
	*EmbeddedStructPtr `desc:"Embedded struct pointer field"`
	//匿名结构体
	AnonymousStruct struct {
		AnonymousString string `desc:"Anonymous string field"`
		AnonymousInt    int    `desc:"Anonymous int field"`
	} `desc:"Anonymous struct field"`
	// 匿名结构体指针
	AnonymousStructPtr *struct {
		AnonymousString string `desc:"Anonymous string field"`
		AnonymousInt    int    `desc:"Anonymous int field"`
	} `desc:"Anonymous struct pointer field"`
	unexported string
}

func TestFlagSet_ComplexOption(t *testing.T) {
	fs := NewFlagSet(nil)
	opt := &ComplexOption{}
	fs.AddOption("test", opt)

	// 测试命令行参数解析
	args := []string{"--test-string=flag_string",
		"--test-int=42", "--test-uint=10", "--test-float64=3.14",
		"--test-bool=true", "--test-duration=10s",
		"--test-nested-nested-string=nested_string",
		"--test-nested-nested-int=42",
		"--test-nested-ptr-nested-string=nested_ptr_string",
		"--test-nested-ptr-nested-int=42",
		"--test-ptr-field=ptr_field", "--test-embedded-string=embedded_string",
		"--test-embedded-int=42",
		"--test-embedded-ptr-string=embedded_ptr_string",
		"--test-anonymous-struct-anonymous-string=anonymous_string",
		"--test-anonymous-struct-anonymous-int=42",
		"--test-anonymous-struct-ptr-anonymous-string=anonymous_string",
		"--test-anonymous-struct-ptr-anonymous-int=42",
	}
	err := fs.Parse(args)
	assert.NoError(t, err)
	assert.Equal(t, "flag_string", opt.String)
	assert.Equal(t, 42, opt.Int)
	assert.Equal(t, uint(10), opt.Uint)
	assert.Equal(t, 3.14, opt.Float64)
	assert.True(t, opt.Bool)
	assert.Equal(t, 10*time.Second, opt.Duration)
	assert.Equal(t, "nested_string", opt.Nested.NestedString)
	assert.Equal(t, 42, opt.Nested.NestedInt)
	assert.Equal(t, "nested_ptr_string", opt.NestedPtr.NestedString)
	assert.Equal(t, 42, opt.NestedPtr.NestedInt)
	assert.Equal(t, "ptr_field", *opt.PtrField)
	assert.Equal(t, "embedded_string", opt.EmbeddedString)
	assert.Equal(t, 42, opt.EmbeddedInt)
	assert.Equal(t, "embedded_ptr_string", opt.EmbeddedPtrString)
	assert.Equal(t, "anonymous_string", opt.AnonymousStruct.AnonymousString)
	assert.Equal(t, 42, opt.AnonymousStruct.AnonymousInt)
	assert.Equal(t, "anonymous_string", opt.AnonymousStructPtr.AnonymousString)
	assert.Equal(t, 42, opt.AnonymousStructPtr.AnonymousInt)

	// 测试环境变量解析
	os.Setenv("TEST_STRING", "env_string")
	os.Setenv("TEST_INT", "42")
	os.Setenv("TEST_UINT", "10")
	os.Setenv("TEST_FLOAT64", "3.14")
	os.Setenv("TEST_BOOL", "true")
	os.Setenv("TEST_DURATION", "10s")
	os.Setenv("TEST_NESTED_NESTED_STRING", "nested_env_string")
	os.Setenv("TEST_NESTED_NESTED_INT", "42")
	os.Setenv("TEST_NESTED_PTR_NESTED_STRING", "nested_ptr_env_string")
	os.Setenv("TEST_NESTED_PTR_NESTED_INT", "42")
	os.Setenv("TEST_PTR_FIELD", "ptr_env_field")
	os.Setenv("TEST_EMBEDDED_STRING", "embedded_env_string")
	os.Setenv("TEST_EMBEDDED_INT", "42")
	os.Setenv("TEST_EMBEDDED_PTR_STRING", "embedded_ptr_env_string")
	os.Setenv("TEST_ANONYMOUS_STRUCT_ANONYMOUS_STRING", "anonymous_env_string")
	os.Setenv("TEST_ANONYMOUS_STRUCT_ANONYMOUS_INT", "42")
	os.Setenv("TEST_ANONYMOUS_STRUCT_PTR_ANONYMOUS_STRING", "anonymous_env_string")
	os.Setenv("TEST_ANONYMOUS_STRUCT_PTR_ANONYMOUS_INT", "42")
	defer os.Unsetenv("TEST_STRING")
	defer os.Unsetenv("TEST_INT")
	defer os.Unsetenv("TEST_UINT")
	defer os.Unsetenv("TEST_FLOAT64")
	defer os.Unsetenv("TEST_BOOL")
	defer os.Unsetenv("TEST_DURATION")
	defer os.Unsetenv("TEST_NESTED_NESTED_STRING")
	defer os.Unsetenv("TEST_NESTED_NESTED_INT")
	defer os.Unsetenv("TEST_NESTED_PTR_NESTED_STRING")
	defer os.Unsetenv("TEST_NESTED_PTR_NESTED_INT")
	defer os.Unsetenv("TEST_PTR_FIELD")
	defer os.Unsetenv("TEST_EMBEDDED_STRING")
	defer os.Unsetenv("TEST_EMBEDDED_INT")
	defer os.Unsetenv("TEST_EMBEDDED_PTR_STRING")
	defer os.Unsetenv("TEST_ANONYMOUS_STRUCT_ANONYMOUS_STRING")
	defer os.Unsetenv("TEST_ANONYMOUS_STRUCT_ANONYMOUS_INT")
	defer os.Unsetenv("TEST_ANONYMOUS_STRUCT_PTR_ANONYMOUS_STRING")
	defer os.Unsetenv("TEST_ANONYMOUS_STRUCT_PTR_ANONYMOUS_INT")

	fs = NewFlagSet(flag.NewFlagSet("test", flag.ExitOnError))
	fs.AddOption("test", opt)
	err = fs.Parse([]string{})
	assert.NoError(t, err)
	assert.Equal(t, "env_string", opt.String)
	assert.Equal(t, 42, opt.Int)
	assert.Equal(t, uint(10), opt.Uint)
	assert.Equal(t, 3.14, opt.Float64)
	assert.True(t, opt.Bool)
	assert.Equal(t, 10*time.Second, opt.Duration)
	assert.Equal(t, "nested_env_string", opt.Nested.NestedString)
	assert.Equal(t, 42, opt.Nested.NestedInt)
	assert.Equal(t, "nested_ptr_env_string", opt.NestedPtr.NestedString)
	assert.Equal(t, 42, opt.NestedPtr.NestedInt)
	assert.Equal(t, "ptr_env_field", *opt.PtrField)
	assert.Equal(t, "embedded_env_string", opt.EmbeddedString)
	assert.Equal(t, 42, opt.EmbeddedInt)
	assert.Equal(t, "embedded_ptr_env_string", opt.EmbeddedPtrString)
	assert.Equal(t, "anonymous_env_string", opt.AnonymousStruct.AnonymousString)
	assert.Equal(t, 42, opt.AnonymousStruct.AnonymousInt)
	assert.Equal(t, "anonymous_env_string", opt.AnonymousStructPtr.AnonymousString)
	assert.Equal(t, 42, opt.AnonymousStructPtr.AnonymousInt)
}

func TestFlagSet_Unset(t *testing.T) {
	type UnsetOption struct {
		UnsetString string `desc:"Unset string field"`
		UnsetInt    int    `desc:"Unset int field"`
	}

	fs := NewFlagSet(nil)
	opt := &UnsetOption{
		UnsetString: "default_string",
		UnsetInt:    42,
	}
	fs.AddOption("test", opt)

	// 测试命令行参数解析
	err := fs.Parse([]string{})
	assert.NoError(t, err)
	assert.Equal(t, "default_string", opt.UnsetString)
	assert.Equal(t, 42, opt.UnsetInt)

}

func TestFlagSet_UnsetPtr(t *testing.T) {
	type UnsetPtrOption struct {
		UnsetPtr *string `desc:"Unset pointer to string field"`
	}

	fs := NewFlagSet(nil)
	opt := &UnsetPtrOption{
		UnsetPtr: new(string),
	}
	fs.AddOption("test", opt)

	// 测试命令行参数解析
	err := fs.Parse([]string{"--test-unset-ptr=flag_string"})
	assert.NoError(t, err)
	assert.Equal(t, "flag_string", *opt.UnsetPtr)
}
