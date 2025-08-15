package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Chat holds the schema definition for the Chat entity.
type Chat struct {
	ent.Schema
}

// Fields of the Chat.
func (Chat) Fields() []ent.Field {
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

// Edges of the Chat.
func (Chat) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("chatroom", ChatRoom.Type).
			Ref("chats").
			Unique().   // 각 Chat은 하나의 ChatRoom에만 속함 (N:1)
			Required(), // ChatRoom 없는 Chat 생성 금지 (FK NOT NULL)
	}
}
