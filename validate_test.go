package validate

import (
	"fmt"
	"testing"
)

type Valid struct {
	Required  string   `json:"required" require:"+"`
	Number    string   `json:"number" goodChars:"0123456789"`
	NotNumber string   `json:"notnumber" badChars:"0123456789"`
	Match     string   `json:"match" match:"[0123456789]+"`
	Thing     []string `json:"thing" allowVals:"thing,thing2"`
}

func TestAll(t *testing.T) {
	if err := Validate(Valid{Required: "test", Number: "5564", NotNumber: "thing isn't", Match: "2135", Thing: []string{"thing", "thing2"}}); err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}

	if err := Validate(Valid{Number: "5564", NotNumber: "thing isn't", Match: "2135", Thing: []string{"thing", "thing2"}}); err != nil {
		fmt.Println(err.Error())
	} else {
		t.Fail()
	}

	if err := Validate(Valid{Required: "test", Number: "55asdg64", NotNumber: "thing isn't", Match: "2135", Thing: []string{"thing", "thing2"}}); err != nil {
		fmt.Println(err.Error())
	} else {
		t.Fail()
	}

	if err := Validate(Valid{Required: "test", Number: "5564", NotNumber: "th235ing isn't", Match: "2135", Thing: []string{"thing", "thing2"}}); err != nil {
		fmt.Println(err.Error())
	} else {
		t.Fail()
	}

	if err := Validate(Valid{Required: "test", Number: "5564", NotNumber: "thng isn't", Match: "testing ", Thing: []string{"thing", "other", "thing2"}}); err != nil {
		fmt.Println(err.Error())
	} else {
		t.Fail()
	}

	if err := Validate(Valid{Required: "test", Number: "5564", NotNumber: "thg isn't", Match: "235", Thing: []string{"thing", "other", "thing2"}}); err != nil {
		fmt.Println(err.Error())
	} else {
		t.Fail()
	}
}
