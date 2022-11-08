package validator

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic(fmt.Errorf("could not get binding validator engine"))
	}

	v.RegisterValidation("objectid", objectid)
}

var objectid validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	if _, err := primitive.ObjectIDFromHex(value); err != nil {
		return false
	}

	return true
}
