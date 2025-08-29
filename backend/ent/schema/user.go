package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().NotEmpty(),
		field.String("name"),
		field.String("email").NotEmpty(),
		field.String("profileURL").NotEmpty(),
		field.String("plan").Default("free"),
		field.Time("createdAt").
			Default(time.Now).Immutable(),
		field.Time("updatedAt").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.String("processingData").Optional().Nillable(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("planners", Planner.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("owned_chatrooms", Chatroom.Type),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name"),
		index.Fields("email"),
		index.Fields("profileURL"),
		index.Fields("plan"),
	}
}
