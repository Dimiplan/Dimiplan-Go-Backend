//go:build ignore
// +build ignore

package main

import (
	"log"
	"strings"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/edge"
)

func main() {
	opts := []entc.Option{
		entc.Extensions(
			&EncodeExtension{},
		),
	}
	log.Println("Generating code...")
	err := entc.Generate("./schema", &gen.Config{}, opts...)
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}

// EncodeExtension is an implementation of entc.Extension that adds a MarshalJSON
// method to each generated type <T> and inlines the Edges field to the top level JSON.
type EncodeExtension struct {
	entc.DefaultExtension
}

// Templates of the extension.
func (e *EncodeExtension) Templates() []*gen.Template {
	return []*gen.Template{
		gen.MustParse(gen.NewTemplate("model/additional/jsonencode").
			Parse(`
{{ if $.Edges }}
    // MarshalJSON implements the json.Marshaler interface.
    func ({{ $.Receiver }} *{{ $.Name }}) MarshalJSON() ([]byte, error) {
        type Alias {{ $.Name }}
        return json.Marshal(&struct {
            *Alias
            {{ $.Name }}Edges
        }{
            Alias: (*Alias)({{ $.Receiver }}),
            {{ $.Name }}Edges: {{ $.Receiver }}.Edges,
        })
    }
{{ end }}
`)),
	}
}

// Hooks of the extension.
func (e *EncodeExtension) Hooks() []gen.Hook {
	return []gen.Hook{
		func(next gen.Generator) gen.Generator {
			return gen.GenerateFunc(func(g *gen.Graph) error {
				// Set edges to json:"-"
				edgeTag := edge.Annotation{StructTag: `json:"-"`}
				for _, n := range g.Nodes {
					n.Annotations.Set(edgeTag.Name(), edgeTag)

					// Remove omitempty from field JSON tags
					for _, f := range n.Fields {
						if f.StructTag != "" {
							// Replace omitempty in existing struct tags
							newTag := strings.Replace(f.StructTag, ",omitempty", "", -1)
							f.StructTag = newTag
						}
					}

					// Handle ID field specifically to remove omitempty
					if n.ID != nil && n.ID.StructTag != "" {
						newIDTag := strings.Replace(n.ID.StructTag, ",omitempty", "", -1)
						n.ID.StructTag = newIDTag
					}
				}
				return next.Generate(g)
			})
		},
	}
}
