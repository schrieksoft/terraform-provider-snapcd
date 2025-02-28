package utils

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Generic function to convert struct to a map with UpperCamelCase keys
func PlanToJson(input interface{}, ignoreFields ...[]string) (map[string]interface{}, error) {

	ctx := context.Background()

	result := make(map[string]interface{})
	v := reflect.ValueOf(input)
	t := reflect.TypeOf(input)

	var ignored []string
	if len(ignoreFields) > 0 {
		ignored = ignoreFields[0]
	} else {
		ignored = []string{}
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Skip unexported fields
		if !value.CanInterface() {
			continue
		}

		// Skip if field is in list of ignored fields
		if contains(ignored, strings.ToLower(field.Name)) {
			continue
		}

		// Extract the Value field from types.String and types.Bool

		val := value.Interface()

		switch val.(type) {
		case types.String:
			castValue := value.Interface().(types.String)
			if castValue.IsUnknown() {
				continue
			} else if !(castValue.IsNull()) {
				result[field.Name] = castValue.ValueString()
			} else {
				result[field.Name] = nil
			}
		case types.Bool:
			castValue := value.Interface().(types.Bool)
			if castValue.IsUnknown() {
				continue
			} else if !(castValue.IsNull()) {
				result[field.Name] = castValue.ValueBool()
			} else {
				result[field.Name] = nil
			}
		case types.Number:
			castValue := value.Interface().(types.Number)
			if castValue.IsUnknown() {
				continue
			} else if !(castValue.IsNull()) {
				result[field.Name] = castValue.ValueBigFloat()
			} else {
				result[field.Name] = nil
			}
		case types.Float64:
			castValue := value.Interface().(types.Float64)
			if castValue.IsUnknown() {
				continue
			} else if !(castValue.IsNull()) {
				result[field.Name] = castValue.ValueFloat64()
			} else {
				result[field.Name] = nil
			}
		case types.Int64:
			castValue := value.Interface().(types.Int64)
			if castValue.IsUnknown() {
				continue
			} else if !(castValue.IsNull()) {
				result[field.Name] = castValue.ValueInt64()

			} else {
				result[field.Name] = nil
			}
		case types.Map:
			castValue := value.Interface().(types.Map)
			if castValue.IsUnknown() {
				continue
			} else if !(castValue.IsNull()) {
				v, _ := castValue.ToMapValue(context.Background())
				// TODO add error
				result[field.Name] = v

			} else {
				result[field.Name] = nil
			}
		case types.List:

			// TODO currently only list of strings is supported here!

			list := value.Interface().(types.List)

			if list.ElementType(ctx) != types.StringType {
				return nil, fmt.Errorf("expected list of strings but got: %s", list.ElementType(ctx).String())
			}

			// Extract the values from the ListValue
			var strArray []string
			diags := list.ElementsAs(ctx, &strArray, false)
			if diags.HasError() {
				return nil, fmt.Errorf("error extracting elements: %s", diags)
			}
			result[field.Name] = strArray

		case types.Set:
			castValue := value.Interface().(types.Set)
			if castValue.IsUnknown() {
				continue
			} else if !(castValue.IsNull()) {
				v, _ := castValue.ToSetValue(context.Background())
				// TODO add error
				result[field.Name] = v
			} else {
				result[field.Name] = nil
			}

		default:
			// Handle other types if necessary
			result[field.Name] = value.Interface()
		}
	}
	return result, nil
}

// Generic function to convert a map with UpperCamelCase keys to a struct
func JsonToPlan(data map[string]interface{}, output interface{}, ignoreFields ...[]string) error {

	v := reflect.ValueOf(output)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("output must be a pointer to a struct")
	}
	v = v.Elem()
	t := v.Type()

	var ignored []string
	if len(ignoreFields) > 0 {
		ignored = ignoreFields[0]
		for i, field := range ignored {
			ignored[i] = strings.ToLower(field)
		}
	} else {
		ignored = []string{}
	}

	normalizedData := make(map[string]interface{})
	for key, value := range data {
		normalizedData[strings.ToLower(key)] = value
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Skip unexported fields
		if !fieldValue.CanSet() {
			continue
		}
		// Get the value from the map
		value, exists := normalizedData[strings.ToLower(field.Name)]
		if !exists {
			continue
		}

		if value == nil {
			continue
		}

		// Skip if field is in list of ignored fields
		if contains(ignored, strings.ToLower(field.Name)) {
			continue
		}

		// Set the value based on the field type
		switch fieldValue.Interface().(type) {
		case types.String:
			fieldValue.Set(reflect.ValueOf(types.StringValue(value.(string))))

		case types.Float64:
			floatValue := new(big.Float).SetFloat64(value.(float64))
			fieldValue.Set(reflect.ValueOf(types.NumberValue(floatValue)))

		case types.Int64:
			// Ensure value is an int64 (or perform type assertion accordingly)
			var intValue int64
			// First try to cast as int64
			if v, ok := value.(int64); ok {
				intValue = v
			} else if v, ok := value.(float64); ok {
				// If it's not an int64, try to convert from float64 to int64
				intValue = int64(v)
			} else {
				// Handle the case where it's neither int64 nor float64
				return fmt.Errorf("Unable to cast value %v of type %T as integer", value, value)
			}
			fieldValue.Set(reflect.ValueOf(types.Int64Value(intValue)))

		case types.Bool:
			fieldValue.Set(reflect.ValueOf(types.BoolValue(value.(bool))))

		case types.List:
			interfaceSlice, ok := value.([]interface{})
			if !ok {
				return fmt.Errorf("expected []interface{}, got %T for field %s", value, field.Name)
			}

			// Handle empty list
			if len(interfaceSlice) == 0 {
				emptyList, _ := types.ListValue(types.StringType, nil) // Default to empty string list
				fieldValue.Set(reflect.ValueOf(emptyList))
				continue
			}

			// Determine the element type dynamically based on the first element
			var elementType attr.Type
			var converter func(interface{}) (attr.Value, error)

			switch firstElem := interfaceSlice[0].(type) { //TODO, this may not be appropriate for a mix of integers and floats
			case string:
				elementType = types.StringType
				converter = func(val interface{}) (attr.Value, error) {
					str, ok := val.(string)
					if !ok {
						return nil, fmt.Errorf("expected string, got %T", val)
					}
					return types.StringValue(str), nil
				}
			case float64:
				elementType = types.Float64Type
				converter = func(val interface{}) (attr.Value, error) {
					num, ok := val.(float64)
					if !ok {
						return nil, fmt.Errorf("expected float64, got %T", val)
					}
					return types.Float64Value(num), nil
				}
			case int:
				elementType = types.Int64Type
				converter = func(val interface{}) (attr.Value, error) {
					integer, ok := val.(int)
					if !ok {
						return nil, fmt.Errorf("expected int, got %T", val)
					}
					return types.Int64Value(int64(integer)), nil
				}
			case int64:
				elementType = types.Int64Type
				converter = func(val interface{}) (attr.Value, error) {
					integer, ok := val.(int64)
					if !ok {
						return nil, fmt.Errorf("expected int64, got %T", val)
					}
					return types.Int64Value(integer), nil
				}
			case bool:
				elementType = types.BoolType
				converter = func(val interface{}) (attr.Value, error) {
					boolean, ok := val.(bool)
					if !ok {
						return nil, fmt.Errorf("expected bool, got %T", val)
					}
					return types.BoolValue(boolean), nil
				}
			default:
				return fmt.Errorf("unsupported list element type %T in field %s", firstElem, field.Name)
			}

			// Convert each element in the slice to the appropriate Terraform value
			attrValues := make([]attr.Value, len(interfaceSlice))
			for i, elem := range interfaceSlice {
				converted, err := converter(elem)
				if err != nil {
					return fmt.Errorf("error converting list element at index %d: %w", i, err)
				}
				attrValues[i] = converted
			}

			// Create the types.List value and handle diagnostics
			listValue, diags := types.ListValue(elementType, attrValues)
			if diags.HasError() {
				errorMessages := []string{}
				for _, diagnostic := range diags {
					errorMessages = append(errorMessages, diagnostic.Detail())
				}
				return fmt.Errorf("failed to construct ListValue for field %s: %s", field.Name, strings.Join(errorMessages, "; "))
			}

			// Set the constructed list value
			fieldValue.Set(reflect.ValueOf(listValue))

		default:
			fieldValue.Set(reflect.ValueOf(value))
		}
	}
	return nil
}

func GetEnvBool(key string, defaultValue bool) (bool, error) {
	s := os.Getenv(key)
	if len(s) == 0 {
		return defaultValue, nil
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}
