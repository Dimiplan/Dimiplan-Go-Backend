package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
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
			Default(time.Now).Immutable(),
		field.Time("updatedAt").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("chatroom", Chatroom.Type).
			Ref("messages").
			Unique().   // 각 Message는 하나의 ChatRoom에만 속함 (N:1)
			Required(), // ChatRoom 없는 Message 생성 금지 (FK NOT NULL)
	}
}
