package main

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"k8s.io/klog/v2"

	"github.com/suzi1037/pcmd/pkg/pcmd"

	"github.com/suzi1037/pcmd/app/demo/cmd"
	"github.com/suzi1037/pcmd/pkg/cprt"
)

func main() {
	if err := cmd.Run(); err != nil {

		var e validator.ValidationErrors
		if errors.As(err, &e) {
			cprt.Error("Error:field validation, causes:")
			for idx, e := range e {
				cprt.Error("\t%d: %s", idx+1, e)
			}
			return
		}

		if errors.Is(err, pcmd.ErrUserAbort) {
			cprt.Error("User abort")
			return
		}

		cprt.Error("LOG FILE: ./demo.log")
		klog.Fatalf("execute error: %+v", err)
	}
}
