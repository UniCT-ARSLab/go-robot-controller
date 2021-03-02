package utilities

import (
	"github.com/fatih/color"
)

/*
PrintStatus prints the result of a CheckService result and returns the same error, if there was, otherwise nil will be returned.
*/

func CreateColorString(str string, clr color.Attribute) string {
	return color.New(clr).SprintFunc()(str)
}
