package reflect

import "reflect"

// Function to copy fields from one struct to another using reflection
func CopyFields(src interface{}, dst interface{}) error {
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst).Elem()

	// Check if the source is a pointer
	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}

	// Iterate over the fields of the source
	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		dstField := dstVal.FieldByName(srcVal.Type().Field(i).Name)

		// Check if the destination has the field and if it can be set
		if dstField.IsValid() && dstField.CanSet() {
			dstField.Set(srcField)
		}
	}

	return nil
}
