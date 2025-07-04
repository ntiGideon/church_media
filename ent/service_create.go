// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/ogidi/church-media/ent/attendancerecord"
	"github.com/ogidi/church-media/ent/service"
)

// ServiceCreate is the builder for creating a Service entity.
type ServiceCreate struct {
	config
	mutation *ServiceMutation
	hooks    []Hook
}

// SetDate sets the "date" field.
func (sc *ServiceCreate) SetDate(t time.Time) *ServiceCreate {
	sc.mutation.SetDate(t)
	return sc
}

// SetName sets the "name" field.
func (sc *ServiceCreate) SetName(s string) *ServiceCreate {
	sc.mutation.SetName(s)
	return sc
}

// SetNillableName sets the "name" field if the given value is not nil.
func (sc *ServiceCreate) SetNillableName(s *string) *ServiceCreate {
	if s != nil {
		sc.SetName(*s)
	}
	return sc
}

// SetType sets the "type" field.
func (sc *ServiceCreate) SetType(s service.Type) *ServiceCreate {
	sc.mutation.SetType(s)
	return sc
}

// SetID sets the "id" field.
func (sc *ServiceCreate) SetID(i int) *ServiceCreate {
	sc.mutation.SetID(i)
	return sc
}

// AddAttendanceRecordIDs adds the "attendance_records" edge to the AttendanceRecord entity by IDs.
func (sc *ServiceCreate) AddAttendanceRecordIDs(ids ...int) *ServiceCreate {
	sc.mutation.AddAttendanceRecordIDs(ids...)
	return sc
}

// AddAttendanceRecords adds the "attendance_records" edges to the AttendanceRecord entity.
func (sc *ServiceCreate) AddAttendanceRecords(a ...*AttendanceRecord) *ServiceCreate {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return sc.AddAttendanceRecordIDs(ids...)
}

// Mutation returns the ServiceMutation object of the builder.
func (sc *ServiceCreate) Mutation() *ServiceMutation {
	return sc.mutation
}

// Save creates the Service in the database.
func (sc *ServiceCreate) Save(ctx context.Context) (*Service, error) {
	sc.defaults()
	return withHooks(ctx, sc.sqlSave, sc.mutation, sc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (sc *ServiceCreate) SaveX(ctx context.Context) *Service {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sc *ServiceCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sc *ServiceCreate) ExecX(ctx context.Context) {
	if err := sc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (sc *ServiceCreate) defaults() {
	if _, ok := sc.mutation.Name(); !ok {
		v := service.DefaultName
		sc.mutation.SetName(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sc *ServiceCreate) check() error {
	if _, ok := sc.mutation.Date(); !ok {
		return &ValidationError{Name: "date", err: errors.New(`ent: missing required field "Service.date"`)}
	}
	if _, ok := sc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Service.name"`)}
	}
	if _, ok := sc.mutation.GetType(); !ok {
		return &ValidationError{Name: "type", err: errors.New(`ent: missing required field "Service.type"`)}
	}
	if v, ok := sc.mutation.GetType(); ok {
		if err := service.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "Service.type": %w`, err)}
		}
	}
	return nil
}

func (sc *ServiceCreate) sqlSave(ctx context.Context) (*Service, error) {
	if err := sc.check(); err != nil {
		return nil, err
	}
	_node, _spec := sc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = int(id)
	}
	sc.mutation.id = &_node.ID
	sc.mutation.done = true
	return _node, nil
}

func (sc *ServiceCreate) createSpec() (*Service, *sqlgraph.CreateSpec) {
	var (
		_node = &Service{config: sc.config}
		_spec = sqlgraph.NewCreateSpec(service.Table, sqlgraph.NewFieldSpec(service.FieldID, field.TypeInt))
	)
	if id, ok := sc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := sc.mutation.Date(); ok {
		_spec.SetField(service.FieldDate, field.TypeTime, value)
		_node.Date = value
	}
	if value, ok := sc.mutation.Name(); ok {
		_spec.SetField(service.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := sc.mutation.GetType(); ok {
		_spec.SetField(service.FieldType, field.TypeEnum, value)
		_node.Type = value
	}
	if nodes := sc.mutation.AttendanceRecordsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.AttendanceRecordsTable,
			Columns: []string{service.AttendanceRecordsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(attendancerecord.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// ServiceCreateBulk is the builder for creating many Service entities in bulk.
type ServiceCreateBulk struct {
	config
	err      error
	builders []*ServiceCreate
}

// Save creates the Service entities in the database.
func (scb *ServiceCreateBulk) Save(ctx context.Context) ([]*Service, error) {
	if scb.err != nil {
		return nil, scb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(scb.builders))
	nodes := make([]*Service, len(scb.builders))
	mutators := make([]Mutator, len(scb.builders))
	for i := range scb.builders {
		func(i int, root context.Context) {
			builder := scb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ServiceMutation)
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
					_, err = mutators[i+1].Mutate(root, scb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, scb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, scb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (scb *ServiceCreateBulk) SaveX(ctx context.Context) []*Service {
	v, err := scb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (scb *ServiceCreateBulk) Exec(ctx context.Context) error {
	_, err := scb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scb *ServiceCreateBulk) ExecX(ctx context.Context) {
	if err := scb.Exec(ctx); err != nil {
		panic(err)
	}
}
