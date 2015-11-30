package msgspec
import "testing"

func TestValidateStructAcceptsStructs(t *testing.T){
	var x = struct {
		A string
	}{ "Hello" }

	defer func(){
		if r := recover(); r != nil {
			t.Error("Should not have panicked")
		}
	}()

	ValidateStruct(x)
}

func TestValidateStructDoesNotAcceptOtherTypes(t *testing.T){

	defer func(){
		if r := recover(); r != nil {
			t.Log("yay, panicked")
		}
	}()

	ValidateStruct("lkjkj")

	t.Error("Should have panicked.")
}
