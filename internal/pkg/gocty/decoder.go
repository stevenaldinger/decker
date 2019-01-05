package gocty

import (
	"fmt"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
	"strconv"
)

// Decoder contains helpers for decoding go-cty values to native go types.
type Decoder struct{}

// GetInt gets an integer from a cty value.
func (*Decoder) GetInt(ctyVal cty.Value) int {
	var val int
	err := gocty.FromCtyValue(ctyVal, &val)

	if err != nil {
		fmt.Println("Error decoding int:", err)
	}

	return val
}

// GetString gets a string from a cty value.
func (*Decoder) GetString(ctyVal cty.Value) string {
	var val string
	err := gocty.FromCtyValue(ctyVal, &val)

	if err != nil {
		fmt.Println("CTY VAL:", ctyVal)
		fmt.Println("Error decoding string:", err)
	}

	return val
}

// GetMap will get a map of strings from a cty value map. Booleans will be
// converted to strings to keep things simple right now and no other types are
// supported yet.
func (*Decoder) GetMap(ctyVal cty.Value) *map[string]string {
	var decodedMapVal string
	var decodedBool bool
	var decodedMap = map[string]string{}

	// if this errors out use default?
	ctyValMap := ctyVal.AsValueMap()

	for key, val := range ctyValMap {
		err := gocty.FromCtyValue(val, &decodedMapVal)
		if err != nil {
			boolErr := gocty.FromCtyValue(val, &decodedBool)
			if boolErr != nil {
				fmt.Println("Error trying to decode cty val for string and bool in map:", boolErr)
				// exit the program here?
			} else {
				decodedMap[key] = strconv.FormatBool(decodedBool)
			}
		} else {
			decodedMap[key] = decodedMapVal
		}
	}

	// if err != nil {
	//   fmt.Println("CTY VAL:", ctyVal)
	//   fmt.Println("Error decoding string:", err)
	// }

	return &decodedMap
}

// GetBool gets a boolean from a cty value.
func (*Decoder) GetBool(ctyVal cty.Value) bool {
	var val bool

	err := gocty.FromCtyValue(ctyVal, &val)

	if err != nil {
		fmt.Println("CTY VAL:", ctyVal)
		fmt.Println("Error decoding boolean:", err)
	}

	return val
}

// GetStringOrBool will get a string from a cty value (of types string or boolean).
func (*Decoder) GetStringOrBool(ctyVal cty.Value) string {
	var decodedVal string
	var decodedBool bool
	err := gocty.FromCtyValue(ctyVal, &decodedVal)

	if err != nil {
		boolErr := gocty.FromCtyValue(ctyVal, &decodedBool)
		if boolErr != nil {
			fmt.Println("decoder.go - Decoding error for string and bool and default value was an empty string:", boolErr)
		} else {
			decodedVal = strconv.FormatBool(decodedBool)
		}
	}

	return decodedVal
}
