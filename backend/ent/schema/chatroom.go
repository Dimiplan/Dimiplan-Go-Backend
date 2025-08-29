package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Chatroom holds the schema definition for the ChatRoom entity.
type Chatroom struct {
	ent.Schema
}

// Fields of the Chatroom.
func (Chatroom) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.Time("createdAt").
			Default(time.Now).Immutable().StructTag(`json:"-"`),
		field.Time("updatedAt").
			Default(time.Now).
			UpdateDefault(time.Now).StructTag(`json:"-"`),
	}
}

// Edges of the Chatroom.
func (Chatroom) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("owned_chatrooms").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("messages", Message.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (Chatroom) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("owner"),
		index.Fields("name"),
	}
}
