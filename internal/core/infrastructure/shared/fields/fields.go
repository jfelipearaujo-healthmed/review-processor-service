package fields

import "reflect"

var (
	ANY_CHAR string = "%"
)

func GetNonEmptyFields(input interface{}, prefix *string, suffix *string) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)

	// Check if input is a pointer and get the element
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// Iterate over the fields of the struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Check if the field is zero value
		if !isZeroValue(field) {
			fieldValue := field.Interface()
			// Add prefix and suffix if the field is of type string
			switch field.Kind() {
			case reflect.String:
				if prefix != nil {
					fieldValue = *prefix + fieldValue.(string)
				}
				if suffix != nil {
					fieldValue = fieldValue.(string) + *suffix
				}
			default:
			}

			// Get the JSON tag, use field name if no tag is present
			jsonTag := fieldType.Tag.Get("json")
			if jsonTag == "" {
				jsonTag = fieldType.Name
			}

			// Add to result map
			result[jsonTag] = fieldValue
		}
	}

	return result
}

func isZeroValue(v reflect.Value) bool {
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
