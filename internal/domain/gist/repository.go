package gist

import "context"

type IStateTransition interface {
	CreateStateTransition(ctx context.Context, entity *StateTransition) (*StateTransition, error)
}
