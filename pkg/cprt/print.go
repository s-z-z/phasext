package cprt

import (
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
	"github.com/suzi1037/pcmd/pkg/symbol"
	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

func Green(format string, a ...interface{}) string {
	return color.HiGreenString(format, a...)
}
func Blue(format string, a ...interface{}) string {
	return color.HiBlueString(format, a...)
}
func Yellow(format string, a ...interface{}) string {
	return color.HiYellowString(format, a...)
}
func Red(format string, a ...interface{}) string {
	return color.HiRedString(format, a...)
}
func Magenta(format string, a ...interface{}) string {
	return color.HiMagentaString(format, a...)
}

func Ok(format string, a ...interface{}) {
	klog.Infoln(Green(format, a...))
}

func Debug(format string, a ...interface{}) {
	klog.Infoln(Magenta(format, a...))
}

func Info(format string, a ...interface{}) {
	klog.Infoln(Blue(format, a...))
}

func Warning(format string, a ...interface{}) {
	klog.Infoln(Yellow(format, a...))
}

func Error(format string, a ...interface{}) {
	klog.Infoln(Red(format, a...))
}

func PhaseTitle(title string, a ...interface{}) {
	prefix := fmt.Sprintf("%s", symbol.PHASE)
	title = fmt.Sprintf("[Phase] %s", title)
	title = fmt.Sprintf(title, a...)
	klog.Infof("%s %s\n", prefix, Blue(title))
}

func PhaseOK() {
	Ok("%s  ", symbol.OK)
}
func PhaseOKStr(s string) {
	Ok("%s %s  ", s, symbol.OK)
}

func PhaseWarning() {
	Warning("%s  ", symbol.WARN)
}

func PhaseError() {
	Error("%s  ", symbol.Error)
}

func PhaseEmoj(format string, a ...interface{}) {
	format = fmt.Sprintf("%s: %s", symbol.PHASE, format)
	Info(format, a...)
}

func PrettyJson(o interface{}) string {
	jsonBytes, _ := json.MarshalIndent(o, "", "  ")
	return string(jsonBytes)
}

func Yaml(o interface{}) string {
	bytes, _ := yaml.Marshal(o)
	return string(bytes)
}

func Spew(o ...interface{}) string {
	s := spew.ConfigState{
		Indent:                  "  ",
		DisableCapacities:       true,
		DisableMethods:          true,
		DisablePointerAddresses: true,
		DisablePointerMethods:   true,
	}
	return s.Sdump(o...)
}

func SpewInfo(o ...interface{}) {
	Info(Spew(o...))
}

func SpewWarning(o ...interface{}) {
	Warning(Spew(o...))
}
func SpewDebug(o ...interface{}) {
	Debug(Spew(o...))
}

func SpewError(o ...interface{}) {
	Error(Spew(o...))
}

func YamlInfo(o interface{}) {
	Info(Yaml(o))
}

func YamlWarning(o interface{}) {
	Warning(Yaml(o))
}
func YamlDebug(o interface{}) {
	Debug(Yaml(o))
}

func YamlError(o interface{}) {
	Error(Yaml(o))
}
