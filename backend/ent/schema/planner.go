package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Planner holds the schema definition for the Planner entity.
type Planner struct {
	ent.Schema
}

// Fields of the Planner.
func (Planner) Fields() []ent.Field {
	return []ent.Field{
		field.String("type"),
		field.String("name").NotEmpty(),
		field.Time("createdAt").
			Default(time.Now).Immutable().StructTag(`json:"-"`),
		field.Time("updatedAt").
			Default(time.Now).
			UpdateDefault(time.Now).StructTag(`json:"-"`),
	}
}

// Edges of the Planner.
func (Planner) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("planners").Unique().Required(),
		edge.To("tasks", Task.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (Planner) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name"),
		index.Fields("type"),
		index.Edges("user"),
	}
}
