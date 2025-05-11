package util

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

type FlagKind int

const (
	Local FlagKind = iota + 1
	Persistent
)

type FieldProp struct {
	FieldName string
	FieldAddr any
	FlagName  string
	Usage     string
	Kind      reflect.Kind
}

type ExportProp struct {
	Export bool
	Usage  string
}

func getExportFields(o interface{}) ([]FieldProp, error) {
	var exportFields []FieldProp
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)

	if t.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("must be a pointer")
	}

	t = t.Elem()

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("provided value is not a struct type")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldAddr := v.Elem().Field(i).Addr().Interface()

		ep := parseExportTag(field.Tag.Get("export"))
		if ep.Export {
			usage := ep.Usage
			if usage == "" {
				usage = field.Name
			}
			exportFields = append(exportFields, FieldProp{
				FieldName: field.Name,
				FieldAddr: fieldAddr,
				FlagName:  getFieldUse(field),
				Usage:     usage,
				Kind:      field.Type.Kind(),
			})

		}
	}
	return exportFields, nil
}

func parseExportTag(s string) ExportProp {
	var ret ExportProp
	s = strings.TrimSpace(s)
	props := strings.Split(s, ",")
	propsMap := make(map[string]bool, len(props))
	for _, p := range props {
		propsMap[p] = true
	}
	if propsMap["true"] {
		delete(propsMap, "true")
		ret.Export = true
	}
	if len(propsMap) > 0 {
		for u := range propsMap {
			ret.Usage = u
		}
	}
	return ret
}

func getFieldUse(field reflect.StructField) string {
	var ret string

	tagStr, ok := field.Tag.Lookup("json")
	if ok {
		ret = getTag0String(tagStr)
		if ret != "" {
			return ret
		}
	}

	tagStr, ok = field.Tag.Lookup("yaml")
	if ok {
		ret = getTag0String(tagStr)
		if ret != "" {
			return ret
		}
	}

	ret = field.Name
	return strings.ToLower(ret)
}

func getTag0String(tag string) string {
	s := strings.Split(tag, ",")
	if len(s) > 0 {
		if s[0] != "" {
			return s[0]
		}
		return ""
	}
	return ""
}

func AddExportFlags(cmd *cobra.Command, o any, specIncludes []string, flagKind FlagKind, bindAddr bool) {
	exportFields, err := getExportFields(o)
	if err != nil {
		klog.Fatalf("Failed to get export fields: %v", err)
	}

	flagSet := cmd.Flags()
	if flagKind == Persistent {
		flagSet = cmd.PersistentFlags()
	}

	for _, f := range exportFields {
		if len(specIncludes) > 0 && !lo.Contains(specIncludes, f.FieldName) {
			continue
		}
		FlagSet(flagSet, f, bindAddr)
	}
}

func FlagSet(fs *pflag.FlagSet, f FieldProp, bindAddr bool) {
	switch f.Kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if bindAddr {
			fs.IntVar(f.FieldAddr.(*int), f.FlagName, 0, f.Usage)
		} else {
			fs.Int(f.FlagName, 0, f.Usage)
		}
	case reflect.Float32, reflect.Float64:
		if bindAddr {
			fs.Float64Var(f.FieldAddr.(*float64), f.FlagName, 0.0, f.Usage)
		} else {
			fs.Float64(f.FlagName, 0.0, f.Usage)
		}
	case reflect.String:
		if bindAddr {
			fs.StringVar(f.FieldAddr.(*string), f.FlagName, "", f.Usage)
		} else {
			fs.String(f.FlagName, "", f.Usage)
		}
	case reflect.Bool:
		if bindAddr {
			fs.BoolVar(f.FieldAddr.(*bool), f.FlagName, false, f.Usage)
		} else {
			fs.Bool(f.FlagName, false, f.Usage)
		}
	case reflect.Slice:
		if bindAddr {
			fs.StringSliceVar(f.FieldAddr.(*[]string), f.FlagName, nil, f.Usage)
		} else {
			fs.StringSlice(f.FlagName, nil, f.Usage)
		}
	default:
	}
}

// AddConfigFlag 添加config flag
func AddConfigFlag(cmd *cobra.Command, cFlag string, configPathPtr *string) {
	cmd.PersistentFlags().StringVar(configPathPtr, cFlag, *configPathPtr, "Path to config file")
	_ = cmd.MarkPersistentFlagRequired(cFlag)
}
