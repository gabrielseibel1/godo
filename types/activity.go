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
	AddTag(tag ID)
	RemoveTag(tag ID)
	Tags() map[ID]struct{}
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
	tags        map[ID]struct{}
}

func NewActivity(id ID, description string) *Activity {
	return &Activity{id: id, description: description, tags: make(map[ID]struct{})}
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

func (a *Activity) AddTag(tag ID) {
	a.tags[tag] = struct{}{}
}

func (a *Activity) RemoveTag(tag ID) {
	delete(a.tags, tag)
}

func (a Activity) Tags() map[ID]struct{} {
	return a.tags
}

func IDToString(id ID) string {
	return string(id)
}
