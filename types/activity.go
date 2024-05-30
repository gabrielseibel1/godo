package types

import "time"

type ID string

type Identifiable interface {
	Identify() ID
}

type Describable interface {
	Describe() string
}

type Workable interface {
	Work(duration time.Duration)
	Worked() time.Duration
}

type Doable interface {
	Do()
	Undo()
	Done() bool
}

type Taggable interface {
	Tag(tag ID)
	Tags() []ID
}

type Actionable interface {
	Identifiable
	Describable
	Workable
	Doable
	Taggable
}

type Activity struct {
	id          ID
	description string
	duration    time.Duration
	done        bool
	tags        []ID
}

func NewActivity(id ID, description string) *Activity {
	return &Activity{id: id, description: description}
}

func (a Activity) Identify() ID {
	return a.id
}

func (a Activity) Describe() string {
	return a.description
}

func (a *Activity) Work(duration time.Duration) {
	a.duration += duration
}

func (a Activity) Worked() time.Duration {
	return a.duration
}

func (a *Activity) Do() {
	a.done = true
}

func (a *Activity) Undo() {
	a.done = false
}

func (a Activity) Done() bool {
	return a.done
}

func (a *Activity) Tag(tag ID) {
	a.tags = append(a.tags, tag)
}

func (a Activity) Tags() []ID {
	return a.tags
}
