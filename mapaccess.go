package mapaccess

import (
	"fmt"
	"reflect"
	"strconv"
)

// Get returns the corresponding key from the map
func Get(key string, data interface{}) (interface{}, error) {
	var err error
	state := data
	parser := parse(key)
	token := parser.nextItem()
	for token.typ != tokenEnd && token.typ != tokenError {
		if state, err = get(token, state); err != nil {
			return nil, err
		}
		token = parser.nextItem()
	}
	if token.typ == tokenError {
		return nil, fmt.Errorf(token.val)
	}
	return state, nil
}

func get(key token, data interface{}) (interface{}, error) {
	switch state := data.(type) {
	case map[string]interface{}:
		switch key.typ {
		case tokenIdentifier:
			if res, ok := state[key.val]; ok {
				return res, nil
			}
			return nil, fmt.Errorf("key <%s> not found", key.val)
		default:
			return nil, fmt.Errorf("expected key, but got %s", key.val)
		}
	case []interface{}:
		switch key.typ {
		case tokenArrayIndex:
			index, err := strconv.Atoi(key.val)
			if err != nil {
				return nil, fmt.Errorf("expected array index, but got %s", key.val)
			}
			if index < 0 || index >= len(state) {
				return nil, fmt.Errorf("index out of bounds %s", key.val)
			}
			return state[index], nil
		default:
			return nil, fmt.Errorf("expected index array, but got %s", key.val)
		}
	case nil:
		return nil, fmt.Errorf("key <%s> not found", key.val)
	default:
		return nil, fmt.Errorf("can't deal with this type %s", reflect.TypeOf(data))
	}
}
