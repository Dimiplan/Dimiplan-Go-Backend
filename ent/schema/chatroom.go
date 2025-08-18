package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Chatroom holds the schema definition for the ChatRoom entity.
type Chatroom struct {
	ent.Schema
}

// Fields of the Chatroom.
func (Chatroom) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.Bool("isProcessing").Default(false),
		field.Time("createdAt").
			Default(time.Now).Immutable(),
		field.Time("updatedAt").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Chatroom.
func (Chatroom) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("chatrooms").
			Unique().   // 각 ChatRoom은 하나의 User에만 속함 (N:1)
			Required(), // User 없는 ChatRoom 생성 금지 (FK NOT NULL)

		edge.To("messages", Message.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
