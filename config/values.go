package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/jsii-runtime-go"
)

func prefixEnv(key string) string {
	return fmt.Sprintf("ANYWHERE_%s", key)
}

func envStrs(key string) []string {
	if value, ok := os.LookupEnv(prefixEnv(key)); ok {
		return strings.Split(value, ",")
	}

	return nil
}

func envStr(key string) *string {
	if value, ok := os.LookupEnv(prefixEnv(key)); ok {
		return &value
	}

	return nil
}

func envBool(key string) *bool {
	if value, ok := os.LookupEnv(prefixEnv(key)); ok {
		return jsii.Bool(value == "true")
	}

	return nil
}

func assertStrs(name string, value []string, defaultValue ...[]string) []string {
	if value != nil && len(value) > 0 {
		return value
	}

	if len(defaultValue) == 1 && len(defaultValue[0]) > 0 {
		return defaultValue[0]
	}

	panic(fmt.Errorf("%s: could not resolve value", name))
}

func assertStr(name string, value *string, defaultValue ...string) string {
	if value != nil && *value != "" {
		return *value
	}

	if len(defaultValue) == 1 && len(defaultValue[0]) > 0 {
		return defaultValue[0]
	}

	panic(fmt.Errorf("%s: could not resolve value", name))
}

func assertBool(name string, value *bool, defaultValue ...bool) bool {
	if value != nil {
		return *value

	}

	if len(defaultValue) == 1 {
		return defaultValue[0]
	}

	panic(fmt.Errorf("%s: could not resolve value", name))
}
