// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/ogidi/church-media/ent/message"
	"github.com/ogidi/church-media/ent/response"
	"github.com/ogidi/church-media/ent/user"
)

// ResponseCreate is the builder for creating a Response entity.
type ResponseCreate struct {
	config
	mutation *ResponseMutation
	hooks    []Hook
}

// SetEmail sets the "email" field.
func (rc *ResponseCreate) SetEmail(s string) *ResponseCreate {
	rc.mutation.SetEmail(s)
	return rc
}

// SetSubject sets the "subject" field.
func (rc *ResponseCreate) SetSubject(r response.Subject) *ResponseCreate {
	rc.mutation.SetSubject(r)
	return rc
}

// SetNillableSubject sets the "subject" field if the given value is not nil.
func (rc *ResponseCreate) SetNillableSubject(r *response.Subject) *ResponseCreate {
	if r != nil {
		rc.SetSubject(*r)
	}
	return rc
}

// SetUserID sets the "user_id" field.
func (rc *ResponseCreate) SetUserID(i int) *ResponseCreate {
	rc.mutation.SetUserID(i)
	return rc
}

// SetDescription sets the "description" field.
func (rc *ResponseCreate) SetDescription(s string) *ResponseCreate {
	rc.mutation.SetDescription(s)
	return rc
}

// SetCreatedAt sets the "created_at" field.
func (rc *ResponseCreate) SetCreatedAt(t time.Time) *ResponseCreate {
	rc.mutation.SetCreatedAt(t)
	return rc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (rc *ResponseCreate) SetNillableCreatedAt(t *time.Time) *ResponseCreate {
	if t != nil {
		rc.SetCreatedAt(*t)
	}
	return rc
}

// SetID sets the "id" field.
func (rc *ResponseCreate) SetID(i int) *ResponseCreate {
	rc.mutation.SetID(i)
	return rc
}

// SetMessageID sets the "message" edge to the Message entity by ID.
func (rc *ResponseCreate) SetMessageID(id int) *ResponseCreate {
	rc.mutation.SetMessageID(id)
	return rc
}

// SetMessage sets the "message" edge to the Message entity.
func (rc *ResponseCreate) SetMessage(m *Message) *ResponseCreate {
	return rc.SetMessageID(m.ID)
}

// SetUser sets the "user" edge to the User entity.
func (rc *ResponseCreate) SetUser(u *User) *ResponseCreate {
	return rc.SetUserID(u.ID)
}

// Mutation returns the ResponseMutation object of the builder.
func (rc *ResponseCreate) Mutation() *ResponseMutation {
	return rc.mutation
}

// Save creates the Response in the database.
func (rc *ResponseCreate) Save(ctx context.Context) (*Response, error) {
	rc.defaults()
	return withHooks(ctx, rc.sqlSave, rc.mutation, rc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (rc *ResponseCreate) SaveX(ctx context.Context) *Response {
	v, err := rc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rc *ResponseCreate) Exec(ctx context.Context) error {
	_, err := rc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rc *ResponseCreate) ExecX(ctx context.Context) {
	if err := rc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (rc *ResponseCreate) defaults() {
	if _, ok := rc.mutation.Subject(); !ok {
		v := response.DefaultSubject
		rc.mutation.SetSubject(v)
	}
	if _, ok := rc.mutation.CreatedAt(); !ok {
		v := response.DefaultCreatedAt()
		rc.mutation.SetCreatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (rc *ResponseCreate) check() error {
	if _, ok := rc.mutation.Email(); !ok {
		return &ValidationError{Name: "email", err: errors.New(`ent: missing required field "Response.email"`)}
	}
	if _, ok := rc.mutation.Subject(); !ok {
		return &ValidationError{Name: "subject", err: errors.New(`ent: missing required field "Response.subject"`)}
	}
	if v, ok := rc.mutation.Subject(); ok {
		if err := response.SubjectValidator(v); err != nil {
			return &ValidationError{Name: "subject", err: fmt.Errorf(`ent: validator failed for field "Response.subject": %w`, err)}
		}
	}
	if _, ok := rc.mutation.UserID(); !ok {
		return &ValidationError{Name: "user_id", err: errors.New(`ent: missing required field "Response.user_id"`)}
	}
	if _, ok := rc.mutation.Description(); !ok {
		return &ValidationError{Name: "description", err: errors.New(`ent: missing required field "Response.description"`)}
	}
	if _, ok := rc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Response.created_at"`)}
	}
	if len(rc.mutation.MessageIDs()) == 0 {
		return &ValidationError{Name: "message", err: errors.New(`ent: missing required edge "Response.message"`)}
	}
	if len(rc.mutation.UserIDs()) == 0 {
		return &ValidationError{Name: "user", err: errors.New(`ent: missing required edge "Response.user"`)}
	}
	return nil
}

func (rc *ResponseCreate) sqlSave(ctx context.Context) (*Response, error) {
	if err := rc.check(); err != nil {
		return nil, err
	}
	_node, _spec := rc.createSpec()
	if err := sqlgraph.CreateNode(ctx, rc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = int(id)
	}
	rc.mutation.id = &_node.ID
	rc.mutation.done = true
	return _node, nil
}

func (rc *ResponseCreate) createSpec() (*Response, *sqlgraph.CreateSpec) {
	var (
		_node = &Response{config: rc.config}
		_spec = sqlgraph.NewCreateSpec(response.Table, sqlgraph.NewFieldSpec(response.FieldID, field.TypeInt))
	)
	if id, ok := rc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := rc.mutation.Email(); ok {
		_spec.SetField(response.FieldEmail, field.TypeString, value)
		_node.Email = value
	}
	if value, ok := rc.mutation.Subject(); ok {
		_spec.SetField(response.FieldSubject, field.TypeEnum, value)
		_node.Subject = value
	}
	if value, ok := rc.mutation.Description(); ok {
		_spec.SetField(response.FieldDescription, field.TypeString, value)
		_node.Description = value
	}
	if value, ok := rc.mutation.CreatedAt(); ok {
		_spec.SetField(response.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if nodes := rc.mutation.MessageIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   response.MessageTable,
			Columns: []string{response.MessageColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(message.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.UserID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rc.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   response.UserTable,
			Columns: []string{response.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.UserID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// ResponseCreateBulk is the builder for creating many Response entities in bulk.
type ResponseCreateBulk struct {
	config
	err      error
	builders []*ResponseCreate
}

// Save creates the Response entities in the database.
func (rcb *ResponseCreateBulk) Save(ctx context.Context) ([]*Response, error) {
	if rcb.err != nil {
		return nil, rcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(rcb.builders))
	nodes := make([]*Response, len(rcb.builders))
	mutators := make([]Mutator, len(rcb.builders))
	for i := range rcb.builders {
		func(i int, root context.Context) {
			builder := rcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ResponseMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, rcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, rcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil && nodes[i].ID == 0 {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, rcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (rcb *ResponseCreateBulk) SaveX(ctx context.Context) []*Response {
	v, err := rcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rcb *ResponseCreateBulk) Exec(ctx context.Context) error {
	_, err := rcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rcb *ResponseCreateBulk) ExecX(ctx context.Context) {
	if err := rcb.Exec(ctx); err != nil {
		panic(err)
	}
}
