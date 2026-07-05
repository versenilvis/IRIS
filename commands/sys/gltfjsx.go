package sys

import (
	"github.com/versenilvis/iris/spec"
)

func init() {
	spec.Register(&spec.Spec{
		Name:        "gltfjsx",
		Description: "GLTF to JSX converter",
		Options: []spec.Option{
			{Name: "-t", Description: "Add Typescript definitions"},
			{Name: "-v", Description: "Verbose output w/ names and empty groups"},
			{Name: "-m", Description: "Include metadata (as userData)"},
			{Name: "-s", Description: "Let meshes cast and receive shadows"},
			{Name: "-w", Description: "Prettier printWidth (default: 120)"},
			{Name: "-p", Description: "Number of fractional digits (default: 2)"},
			{Name: "-d", Description: "Draco binary path"},
			{Name: "-r", Description: "Sets directory from which .gltf file is served"},
			{Name: "-D", Description: "Debug output"},
		},
	})
}
