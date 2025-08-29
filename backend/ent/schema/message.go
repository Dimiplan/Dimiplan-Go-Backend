package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.String("sender").NotEmpty(),
		field.Text("message"),
		field.Time("createdAt").
			Default(time.Now).Immutable().StructTag(`json:"-"`),
		field.Time("updatedAt").
			Default(time.Now).
			UpdateDefault(time.Now).StructTag(`json:"-"`),
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("chatroom", Chatroom.Type).Ref("messages").Unique().Required(),
	}
}

// Indexes of the Message.
func (Message) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("sender"),
		index.Fields("message"),
		index.Edges("chatroom"),
	}
}
