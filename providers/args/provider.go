package args

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type ArgProvider struct {
	args map[string]interface{}
}

func Provider() *ArgProvider {
	parsed := map[string]interface{}{}

	args := parseArgs(os.Args)
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		value := parts[1]
		v, ok := getNested(parsed, key)
		if !ok {
			setNested(parsed, key, value)
		} else {
			vt := reflect.TypeOf(v)
			switch vt.Kind() {
			case reflect.Slice:
				vs := v.([]interface{})
				vs = append(vs, value)
				setNested(parsed, key, vs)
				break
			default:
				vs := []interface{}{v, value}
				setNested(parsed, key, vs)
			}
		}
	}
	return &ArgProvider{
		args: parsed,
	}
}

func parseArgs(args []string) []string {
	var a []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !(strings.HasPrefix(arg, "--") || strings.HasPrefix(arg, "-")) {
			// skip args without - or --
			continue
		}
		key := ""
		val := ""

		if strings.Contains(arg, "=") {
			parts := strings.Split(arg, "=")
			if len(parts) == 1 {
				val = "true"
				key = strings.ToLower(strings.TrimLeft(arg, "-"))
			}
			if len(parts) == 2 {
				key = strings.ToLower(strings.TrimLeft(parts[0], "-"))
				val = parts[1]
			}
		} else {
			if i+1 < len(args) && !(strings.HasPrefix(args[i+1], "--") || strings.HasPrefix(args[i+1], "-")) {
				key = strings.ToLower(strings.TrimLeft(arg, "-"))
				val = args[i+1]
			}
		}

		if key != "" {
			a = append(a, fmt.Sprintf("%s=%s", key, val))
		}
	}

	return a
}

// setNested создаёт вложенное дерево для ключа вида "foo.bar.baz"
func setNested(data map[string]interface{}, key string, value interface{}) {
	parts := strings.Split(key, ".")
	if len(parts) == 1 {
		data[key] = value
		return
	}
	current := data

	// Проходим по частям ключа, создавая вложенные мапы
	for _, part := range parts[:len(parts)-1] {
		if _, ok := current[part]; !ok {
			current[part] = make(map[string]interface{})
		}
		next, ok := current[part].(map[string]interface{})
		if !ok {
			// Если значение не мапа, заменяем его новой мапой
			next = make(map[string]interface{})
			current[part] = next
		}
		current = next
	}

	// Последняя часть ключа
	lastPart := parts[len(parts)-1]
	if existing, ok := current[lastPart]; ok {
		// Если значение уже существует, делаем его списком
		if list, isList := existing.([]interface{}); isList {
			current[lastPart] = append(list, value)
		} else {
			current[lastPart] = []interface{}{existing, value}
		}
	} else {
		current[lastPart] = value
	}
}
func getNested(data map[string]interface{}, key string) (interface{}, bool) {
	parts := strings.Split(key, ".")
	if len(parts) == 1 {
		val, ok := data[parts[0]]
		return val, ok
	}
	current := data
	for i, part := range parts {
		next, ok := current[part].(map[string]interface{})
		if !ok {
			return nil, false
		}
		current = next
		if i == len(parts)-1 {
			return current, true
		}
	}
	return nil, false
}

func (p *ArgProvider) Read() (args map[string]interface{}, err error) {
	return p.args, nil
}

func (p *ArgProvider) ReadBytes() (raw []byte, err error) { return nil, nil }
