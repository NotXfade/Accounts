package methods

import (
	"log"
	"strconv"
)

//ConvertID into int
func ConvertID(i interface{}) int {
	var id int
	switch v := i.(type) {
	case int:
		id = i.(int)
	case string:
		id, _ = strconv.Atoi(i.(string))
	case float64:
		id = int(i.(float64))
	case float32:
		id = int(i.(float32))
	default:
		log.Println("The type of interface is ", v)
		id = 0
	}
	return id
}
