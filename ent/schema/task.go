package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
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
		field.Int("priority").Default(1), // 완료했다면 -1
		field.Time("createdAt").
			Default(time.Now).Immutable(),
		field.Time("updatedAt").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Task.
func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("planner", Planner.Type).
			Ref("tasks").
			Unique().   // 각 Task는 하나의 Planner에만 속함 (N:1)
			Required(), // Planner 없는 Task 생성 금지 (FK NOT NULL)
	}
}
