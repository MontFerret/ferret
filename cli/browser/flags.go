package browser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

type Flags map[string]interface{}

func (flags Flags) Get(arg string) (interface{}, error) {
	var values interface{}
	var err error

	if !flags.Has(arg) {
		err = errors.Errorf("The specified argument '%s' does not exist", arg)
	} else {
		values = flags[arg]
	}

	return values, err
}

func (flags Flags) GetString(arg string) (string, error) {
	found, err := flags.Get(arg)

	if err != nil {
		return "", err
	}

	str, ok := found.(string)

	if ok {
		return str, nil
	}

	return "", nil
}

func (flags Flags) GetInt(arg string) (int, error) {
	found, err := flags.Get(arg)

	if err != nil {
		return 0, err
	}

	num, ok := found.(int)

	if ok {
		return num, nil
	}

	return 0, nil
}

func (flags Flags) Has(arg string) bool {
	_, exists := flags[arg]

	return exists
}

func (flags Flags) List() []string {
	orderedFlags := make([]string, 0, 10)

	for arg := range flags {
		orderedFlags = append(orderedFlags, arg)
	}

	sort.Strings(orderedFlags)

	list := make([]string, len(orderedFlags))

	for i, arg := range orderedFlags {
		val, err := flags.Get(arg)

		if err != nil {
			continue
		}

		switch v := val.(type) {
		case int:
			arg = fmt.Sprintf("--%s=%d", arg, v)
		case string:
			arg = fmt.Sprintf("--%s=%s", arg, v)
		default:
			arg = fmt.Sprintf("--%s", arg)
		}

		list[i] = arg
	}

	return list
}

func (flags Flags) Set(arg string, value interface{}) (err error) {
	if value == nil {
		if _, ok := flags[arg]; !ok {
			flags[arg] = nil
		}
	}

	if value != nil {
		switch value.(type) {
		case int:
			flags[arg] = value
		case string:
			flags[arg] = value
		default:
			return errors.Errorf("Invalid data type '%T' for argument %s: %+v", value, arg, value)
		}
	}

	return nil
}

func (flags Flags) SetN(arg string) (err error) {
	return flags.Set(arg, nil)
}

func (flags Flags) String() string {
	return strings.Join(flags.List(), " ")
}
