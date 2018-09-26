package sigs

import (
	"github.com/calebamiles/keps/pkg/sigs/internal/generated"
)

func All() []string {
	return generated.SIGList
}

func Exists(s string) bool {
	return generated.SIGSet[s]
}

func ExtractNameFromPath(p string) string {

}

func ExtractSubprojectNameFromPath(p string) string {

}
