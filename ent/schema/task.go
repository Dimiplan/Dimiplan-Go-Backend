package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Task holds the schema definition for the Task entity.
type Task struct {
	ent.Schema
}

// Fields of the Task.
func (Task) Fields() []ent.Field {
	return []ent.Field{
		field.Time("deadline").Optional(),
		field.String("title").NotEmpty(),
		field.Int("priority").Default(1),
		field.Time("createdAt").
			Default(time.Now).Immutable().StructTag(`json:"-"`),
		field.Time("updatedAt").
			Default(time.Now).
			UpdateDefault(time.Now).StructTag(`json:"-"`),
	}
}

// Edges of the Task.
func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("planner", Planner.Type).Ref("tasks").Unique().Required(),
	}
}

func (Task) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("title"),
		index.Fields("priority"),
		index.Edges("planner"),
	}
}
