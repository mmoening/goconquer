package packets

import (
	"errors"
	"io"
	"reflect"
	"strconv"
	"unsafe"
)

// Write encodes a packet structure into binary data. Data must be a pointer to a
// structure of either a dynamic or fixed amount of memory. Bytes written to w are
// encoded using NetDragon's byte order and read from to successive fields of the 
// data. When reading from structs, the field data for fields with blank field 
// names are skipped (i.e., used for padding between values). All non-blank fields
// must be exported.
func Write(w io.Writer, data interface{}) error {
	e := reflect.ValueOf(data).Elem()
	if e.Kind() != reflect.Struct {
		return errors.New("packets.Write: invalid type " + e.Kind().String())
	}
	
	// Iterate through each field in the structure and write data. Switch on 
	// the field's assigned data type for appropriate writing.
	count := e.NumField()
	for index := 0; index < count; index++ {
		err := writefield(w, e.Field(index), e.Type().Field(index))
		if err != nil { return err }
	}
	return nil
}

// writefield is the recursive function body for writing values from a field. The
// Write function calls this function to write a structure's body out to a binary 
// writer. It can also be called by itself to write out nested structures and 
// arrays.
func writefield(w io.Writer, f reflect.Value, t reflect.StructField) error {

	// Determine if this field is a nested structure or collection that 
	// requires a different type of handling.
	switch f.Kind() {
	case reflect.Struct:
		for i := 0; i < f.NumField(); i++ {
			err := writefield(w, f.Field(i), f.Type().Field(i))
			if err != nil { return err }
		}
		return nil
	
	case reflect.Array:
		for i := 0; i < f.Len(); i++ {
			err := writefield(w, f.Index(i), t)
			if err != nil { return err }
		}
		return nil
		
	case reflect.Slice:
		err := writebytes(w, []byte { byte(f.Len()) })
		if err != nil { return err }
		for i := 0; i < f.Len(); i++ {
			err := writefield(w, f.Index(i), t)
			if err != nil { return err }
		}
		return nil	
	}
	
	// Determine the length of the write required to encode data to the binary 
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
			length = f.Len()
			err := writebytes(w, []byte { byte(length) })
			if err != nil { return err }
		}

	default: 
		length = int(f.Type().Size())
	}
	
	// Make a slice for the buffer being sent to the writebytes function according
	// to the type of data being written. Assume correctness and process write 
	// according to type.
	b := make([]byte, length)
	if b == nil { return errors.New("packets.write: make failed") }
	if f.CanSet() {
		switch f.Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32,
			reflect.Uint64, reflect.Uintptr, reflect.Uint:
			*(*uint64)(unsafe.Pointer(&b[0])) = f.Uint()
			
		case reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Int64, reflect.Int:
			*(*int64)(unsafe.Pointer(&b[0])) = f.Int()
			
		case reflect.Bool:
			*(*bool)(unsafe.Pointer(&b[0])) = f.Bool()
			
		case reflect.String: 
			copy(b, []byte(f.String()))
			
		default:
			return errors.New("packets.write: invalid kind: " +
				f.Kind().String())
		}
	}
	
	// Write bytes to the Writer interface.
	err := writebytes(w, b)
	return err	
}

// writebytes is an internal function is called by a write function to write bytes
// of a specified length to the Writer interface. If the write is unsuccessful, it
// will return an error message. 
func writebytes(w io.Writer, buffer []byte) error {
	if count, err := w.Write(buffer); err != nil || count < len(buffer) {
		return errors.New("packets.writebytes: failed")
	} else { return nil }
}
