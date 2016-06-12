package packets

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"strconv"
	"unsafe"
)

// Read decodes a packet structure from binary data. Data must be a pointer to a 
// structure of either a dynamic or fixed amount of memory. Bytes read from r are 
// decoded using NetDragon's byte order and written to successive fields of the 
// data. When reading into structs, the field data for fields with blank field 
// names are skipped (i.e., used for padding between values). All non-blank fields
// must be exported.
func Read(r io.Reader, data interface{}) error {
	e := reflect.ValueOf(data).Elem()
	if e.Kind() != reflect.Struct {
		return errors.New("packets.Read: invalid type " + e.Kind().String())
	}
	
	// Iterate through each field in the structure and read data. Switch on 
	// the field's assigned data type for appropriate reading.
	count := e.NumField()
	for index := 0; index < count; index++ {
		err := readfield(r, e.Field(index), e.Type().Field(index))
		if err != nil { return err }
	}
	return nil
}

// readfield is the recursive function body for reading in values from a field. 
// The Read function calls this function to read in a structure's body. It can 
// also be called by itself to read in nested structures and arrays.
func readfield(r io.Reader, f reflect.Value, t reflect.StructField) error {
	
	// Determine if this field is a nested structure or collection that requires a
	// different type of handling.
	switch f.Kind() {
	case reflect.Struct:
		for i := 0; i < f.NumField(); i++ {
			err := readfield(r, f.Field(i), f.Type().Field(i))
			if err != nil { return err }
		}
		return nil
		
	case reflect.Array:
		for i := 0; i < f.Len(); i++ {
			err := readfield(r, f.Index(i), t)
			if err != nil { return err }
		}
		return nil
	
	case reflect.Slice:
		b, err := readbytes(r, 1)
		if err != nil { return err }
		length := int(b[0])
		
		f.Set(reflect.MakeSlice(f.Type(), length, length))
		for i := 0; i < length; i++ {
			err := readfield(r, f.Index(i), t)
			if err != nil { return err }
		}
		return nil	
	}
	
	// Determine the length of the read required to decode data from the binary 
	// packet structure. Assume that if the kind isn't a special case, that the 
	// field has a fixed size.
	var length int
	switch f.Kind() {
	case reflect.String: 
		if t.Tag != "" { // Fixed string.
			fixedlength, err := strconv.Atoi(t.Tag.Get("len"))
			if err != nil { return err }
			length = fixedlength
			
		} else { // Dynamic string.
			buffer, err := readbytes(r, 1)
			if err != nil { return err }
			length = int(buffer[0])
		}

	default: 
		length = int(f.Type().Size())
	}
	
	// Read the length of bytes into the buffer and switch according to the data
	// type of the field.  Assume correctness and process write according to type.
	b, err := readbytes(r, length)
	if err != nil { return err }
	if f.CanSet() {
		switch f.Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32,
			reflect.Uint64, reflect.Uintptr, reflect.Uint:
			f.SetUint(*(*uint64)(unsafe.Pointer(&b[0])))
		
		case reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Int64, reflect.Int:
			f.SetInt(*(*int64)(unsafe.Pointer(&b[0])))
		
		case reflect.Bool:
			f.SetBool(*(*bool)(unsafe.Pointer(&b[0])))
		
		case reflect.String: 
			f.SetString(strings.TrimRight(string(b), "\x00"))
			
		default:
			return errors.New("packets.read: Invalid kind detected: " +
				f.Kind().String())
		}
	}
	return nil
}

// readbytes is an internal function is called by a read function to read bytes of
// a specified length from the Reader interface. If the read is successful, it will
// return the byte array created for parsing; else, it will only return an error 
// message.
func readbytes(r io.Reader, length int) ([]byte, error) {
	buffer := make([]byte, length)
	if count, err := r.Read(buffer); err != nil || count < length {
		return nil, errors.New("packets.readbytes: failed")
	} else { return buffer, nil }
}
