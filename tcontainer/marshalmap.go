package tcontainer

import (
	"fmt"
	"strconv"
	"strings"
)

// MarshalMap is a wrapper type to attach converter methods to maps normally
// returned by marshalling methods, i.e. key/value parsers.
// All methods that do a conversion will return an error if the value stored
// behind key is not of the expected type or if the key is not existing in the
// map.
type MarshalMap map[string]interface{}

const (
	// MarshalMapSeparator defines the rune used for path separation
	MarshalMapSeparator = '/'
	// MarshalMapArrayBegin defines the rune starting array index notation
	MarshalMapArrayBegin = '['
	// MarshalMapArrayEnd defines the rune ending array index notation
	MarshalMapArrayEnd = ']'
)

// NewMarshalMap creates a new marshal map (string -> interface{})
func NewMarshalMap() MarshalMap {
	return make(map[string]interface{})
}

// ConvertToMarshalMap tries to convert a type compatible map to a marshal map.
// Compatible types are map[interface{}]interface{}, map[string]interface{} and of
// course MarshalMap. The underlying type must be a map[string]interface{}.
// Convertable child maps will be converted, too.
// You can pass a formatKey function to modify the map keys during conversion.
func ConvertToMarshalMap(value interface{}, formatKey func(string) string) (MarshalMap, error) {
	switch value.(type) {
	case map[interface{}]interface{}:
		mapValue := value.(map[interface{}]interface{})
		result := NewMarshalMap()

		for key, val := range mapValue {
			stringKey, isString := key.(string)
			if !isString {
				return nil, fmt.Errorf("Value is not a map[string]interface{}. Key is not a string")
			}

			if formatKey != nil {
				stringKey = formatKey(stringKey)
			}

			// Convert child nodes to marshal maps, too.
			switch val.(type) {
			case map[interface{}]interface{}, map[string]interface{}, MarshalMap:
				var err error
				if result[stringKey], err = ConvertToMarshalMap(val, formatKey); err != nil {
					return nil, err
				}
			default:
				result[stringKey] = val
			}
		}
		return result, nil

	case map[string]interface{}:
		mapValue := value.(map[string]interface{})
		// No formatting? pass through
		if formatKey == nil {
			return mapValue, nil
		}
		// Format keys
		result := NewMarshalMap()
		for key, val := range mapValue {
			result[formatKey(key)] = val
		}
		return result, nil

	case MarshalMap:
		mapValue := value.(MarshalMap)
		// No formatting? pass through
		if formatKey == nil {
			return mapValue, nil
		}
		// Format keys
		result := NewMarshalMap()
		for key, val := range mapValue {
			result[formatKey(key)] = val
		}
		return result, nil

	default:
		return nil, fmt.Errorf("Value is not compatible to map[string]interface{}")
	}
}

// Bool returns a value at key that is expected to be a boolean
func (mmap MarshalMap) Bool(key string) (bool, error) {
	val, exists := mmap.Value(key)
	if !exists {
		return false, fmt.Errorf(`"%s" is not set`, key)
	}

	boolValue, isBool := val.(bool)
	if !isBool {
		return false, fmt.Errorf(`"%s" is expected to be a boolean`, key)
	}
	return boolValue, nil
}

// Int returns a value at key that is expected to be an int
func (mmap MarshalMap) Int(key string) (int, error) {
	val, exists := mmap.Value(key)
	if !exists {
		return 0, fmt.Errorf(`"%s" is not set`, key)
	}

	intValue, isInt := val.(int)
	if !isInt {
		return 0, fmt.Errorf(`"%s" is expected to be an integer`, key)
	}
	return intValue, nil
}

// Uint64 returns a value at key that is expected to be an uint64
func (mmap MarshalMap) Uint64(key string) (uint64, error) {
	val, exists := mmap.Value(key)
	if !exists {
		return 0, fmt.Errorf(`"%s" is not set`, key)
	}

	intValue, isInt := val.(uint64)
	if !isInt {
		return 0, fmt.Errorf(`"%s" is expected to be an unsigned integer`, key)
	}
	return intValue, nil
}

// Int64 returns a value at key that is expected to be an int64
func (mmap MarshalMap) Int64(key string) (int64, error) {
	val, exists := mmap.Value(key)
	if !exists {
		return 0, fmt.Errorf(`"%s" is not set`, key)
	}

	intValue, isInt := val.(int64)
	if !isInt {
		return 0, fmt.Errorf(`"%s" is expected to be an unsigned integer`, key)
	}
	return intValue, nil
}

// Float64 returns a value at key that is expected to be a float64
func (mmap MarshalMap) Float64(key string) (float64, error) {
	val, exists := mmap.Value(key)
	if !exists {
		return 0, fmt.Errorf(`"%s" is not set`, key)
	}

	floatValue, isFloat := val.(float64)
	if !isFloat {
		return 0, fmt.Errorf(`"%s" is expected to be a float64`, key)
	}
	return floatValue, nil
}

// Array returns a value at key that is expected to be a string
func (mmap MarshalMap) String(key string) (string, error) {
	val, exists := mmap.Value(key)
	if !exists {
		return "", fmt.Errorf(`"%s" is not set`, key)
	}

	strValue, isString := val.(string)
	if !isString {
		return "", fmt.Errorf(`"%s" is expected to be a string`, key)
	}
	return strValue, nil
}

// Array returns a value at key that is expected to be a []interface{}
func (mmap MarshalMap) Array(key string) ([]interface{}, error) {
	val, exists := mmap.Value(key)
	if !exists {
		return nil, fmt.Errorf(`"%s" is not set`, key)
	}

	arrayValue, isArray := val.([]interface{})
	if !isArray {
		return nil, fmt.Errorf(`"%s" is expected to be an array`, key)
	}
	return arrayValue, nil
}

// Map returns a value at key that is expected to be a
// map[interface{}]interface{}.
func (mmap MarshalMap) Map(key string) (map[interface{}]interface{}, error) {
	val, exists := mmap.Value(key)
	if !exists {
		return nil, fmt.Errorf(`"%s" is not set`, key)
	}

	mapValue, isMap := val.(map[interface{}]interface{})
	if !isMap {
		return nil, fmt.Errorf(`"%s" is expected to be a map`, key)
	}
	return mapValue, nil
}

func castToStringArray(key string, value interface{}) ([]string, error) {
	switch value.(type) {
	case string:
		return []string{value.(string)}, nil

	case []interface{}:
		arrayVal := value.([]interface{})
		stringArray := make([]string, 0, len(arrayVal))

		for _, val := range arrayVal {
			strValue, isString := val.(string)
			if !isString {
				return nil, fmt.Errorf(`"%s" does not contain string keys`, key)
			}
			stringArray = append(stringArray, strValue)
		}
		return stringArray, nil

	case []string:
		return value.([]string), nil

	default:
		return nil, fmt.Errorf(`"%s" is not a valid string array type`, key)
	}
}

// StringArray returns a value at key that is expected to be a []string
// This function supports conversion (by copy) from
//  * []interface{}
func (mmap MarshalMap) StringArray(key string) ([]string, error) {
	val, exists := mmap.Value(key)
	if !exists {
		return nil, fmt.Errorf(`"%s" is not set`, key)
	}

	return castToStringArray(key, val)
}

// StringMap returns a value at key that is expected to be a map[string]string.
// This function supports conversion (by copy) from
//  * map[interface{}]interface{}
//  * map[string]interface{}
func (mmap MarshalMap) StringMap(key string) (map[string]string, error) {
	val, exists := mmap.Value(key)
	if !exists {
		return nil, fmt.Errorf(`"%s" is not set`, key)
	}

	switch val.(type) {
	case map[interface{}]interface{}:
		result := make(map[string]string)
		interfaceMap, _ := val.(map[interface{}]interface{})

		for key, value := range interfaceMap {
			stringKey, isString := key.(string)
			if !isString {
				return nil, fmt.Errorf(`"%s" is expected to be a map[string]string. Key is not a string`, key)
			}
			stringValue, isString := value.(string)
			if !isString {
				return nil, fmt.Errorf(`"%s" is expected to be a map[string]string. Value is not a string`, key)
			}
			result[stringKey] = stringValue
		}
		return result, nil

	case map[string]interface{}:
		result := make(map[string]string)
		interfaceMap, _ := val.(map[string]interface{})

		for key, value := range interfaceMap {
			stringValue, isString := value.(string)
			if !isString {
				return nil, fmt.Errorf(`"%s" is expected to be a map[string]string. Value is not a string`, key)
			}
			result[key] = stringValue
		}
		return result, nil

	case map[string]string:
		return val.(map[string]string), nil

	default:
		return nil, fmt.Errorf(`"%s" is expected to be a map[string]string`, key)
	}
}

// StringArrayMap returns a value at key that is expected to be a
// map[string][]string. This function supports conversion (by copy) from
//  * map[interface{}][]interface{}
//  * map[interface{}]interface{}
//  * map[string]interface{}
func (mmap MarshalMap) StringArrayMap(key string) (map[string][]string, error) {
	val, exists := mmap.Value(key)
	if !exists {
		return nil, fmt.Errorf(`"%s" is not set`, key)
	}

	switch val.(type) {
	case map[interface{}][]interface{}:
		interfaceMap := val.(map[interface{}][]interface{})
		result := make(map[string][]string)

		for key, value := range interfaceMap {
			strKey, isString := key.(string)
			if !isString {
				return nil, fmt.Errorf(`"%s" is expected to be a map[string][]string. Key is not a string`, key)
			}
			arrayValue, err := castToStringArray(strKey, value)
			if err != nil {
				return nil, fmt.Errorf(`"%s" is expected to be a map[string][]string. Value is not a []string`, key)
			}
			result[strKey] = arrayValue
		}
		return result, nil

	case map[interface{}]interface{}:
		interfaceMap := val.(map[interface{}]interface{})
		result := make(map[string][]string)

		for key, value := range interfaceMap {
			strKey, isString := key.(string)
			if !isString {
				return nil, fmt.Errorf(`"%s" is expected to be a map[string][]string. Key is not a string`, key)
			}
			arrayValue, err := castToStringArray(strKey, value)
			if err != nil {
				return nil, fmt.Errorf(`"%s" is expected to be a map[string][]string. Value is not a []string`, key)
			}
			result[strKey] = arrayValue
		}
		return result, nil

	case map[string]interface{}:
		interfaceMap := val.(map[string]interface{})
		result := make(map[string][]string)

		for key, value := range interfaceMap {
			arrayValue, err := castToStringArray(key, value)
			if err != nil {
				return nil, fmt.Errorf(`"%s" is expected to be a map[string][]string. Value is not a []string`, key)
			}
			result[key] = arrayValue
		}
		return result, nil

	case map[string][]string:
		return val.(map[string][]string), nil

	default:
		return nil, fmt.Errorf(`"%s" is expected to be a map[string][]string`, key)
	}
}

// MarshalMap returns a value at key that is expected to be another MarshalMap
// This function supports conversion (by copy) from
//  * map[interface{}]interface{}
func (mmap MarshalMap) MarshalMap(key string) (MarshalMap, error) {
	val, exists := mmap.Value(key)
	if !exists {
		return nil, fmt.Errorf(`"%s" is not set`, key)
	}

	return ConvertToMarshalMap(val, nil)
}

// Value returns a value from a given value path.
// Fields can be accessed by their name. Nested fields can be accessed by using
// "/" as a separator. Arrays can be addressed using the standard array
// notation "[<index>]".
// Examples:
// "key"         -> mmap["key"]              single value
// "key1/key2"   -> mmap["key1"]["key2"]     nested map
// "key1[0]"     -> mmap["key1"][0]          nested array
// "key1[0]key2" -> mmap["key1"][0]["key2"]  nested array, nested map
func (mmap MarshalMap) Value(key string) (interface{}, bool) {
	return mmap.resolvePath(key, mmap)
}

func (mmap MarshalMap) resolvePathKey(key string) (int, int) {
	keyEnd := len(key)
	nextKeyStart := keyEnd
	pathIdx := strings.IndexRune(key, MarshalMapSeparator)
	arrayIdx := strings.IndexRune(key, MarshalMapArrayBegin)

	if pathIdx > -1 && pathIdx < keyEnd {
		keyEnd = pathIdx
		nextKeyStart = pathIdx + 1 // don't include slash
	}
	if arrayIdx > -1 && arrayIdx < keyEnd {
		keyEnd = arrayIdx
		nextKeyStart = arrayIdx // include bracket because of multidimensional arrays
	}

	// a       -> key: "a", remain: ""       -- value
	// a/b/c   -> key: "a", remain: "b/c"    -- nested map
	// a[1]b/c -> key: "a", remain: "[1]b/c" -- nested array

	return keyEnd, nextKeyStart
}

func (mmap MarshalMap) resolvePath(key string, value interface{}) (interface{}, bool) {
	if len(key) == 0 {
		return value, true // ### return, found requested value ###
	}

	// Expecting nested type
	switch value.(type) {
	case []interface{}:
		startIdx := strings.IndexRune(key, MarshalMapArrayBegin) // Must be first char, otherwise malformed
		endIdx := strings.IndexRune(key, MarshalMapArrayEnd)     // Must be > startIdx, otherwise malformed

		if startIdx == -1 || endIdx == -1 {
			return nil, false
		}

		if startIdx == 0 && endIdx > startIdx {
			arrayValue := value.([]interface{})
			index, err := strconv.Atoi(key[startIdx+1 : endIdx])

			// [1]    -> index: "1", remain: ""    -- value
			// [1]a/b -> index: "1", remain: "a/b" -- nested map
			// [1][2] -> index: "1", remain: "[2]" -- nested array

			if err == nil && index < len(arrayValue) {
				item := arrayValue[index]
				key := key[endIdx+1:]
				return mmap.resolvePath(key, item) // ### return, nested array ###
			}
		}

	case MarshalMap, map[string]interface{}:
		mapValue := value.(MarshalMap)
		if storedValue, exists := mapValue[key]; exists {
			return storedValue, true
		}

		keyEnd, nextKeyStart := mmap.resolvePathKey(key)
		if storedValue, exists := mapValue[key[:keyEnd]]; exists {
			remain := key[nextKeyStart:]
			return mmap.resolvePath(remain, storedValue) // ### return, nested map ###
		}
	}

	return nil, false
}
