package util

import "fmt"

func LogOut(args ...interface{}) {
	for _, arg := range args {
		switch arg.(type) {
		case *int:
			fmt.Print(*arg.(*int), " ")
		case *string:
			fmt.Print(*arg.(*string), " ")
		case string:
			fmt.Print(arg, " ")
		case int:
			fmt.Print(arg, " ")
		default:
			fmt.Print("unspport", " ")
		}
	}
	fmt.Println("")
}
