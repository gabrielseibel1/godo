package types

import (
	"slices"
	"time"
)

type ID string

type Identifiable interface {
	Identify(ID)
	Identity() ID
}

type Describable interface {
	Describe(string)
	Description() string
}

type Workable interface {
	Work(time.Duration)
	WorkPeriod(start, end time.Time)
	AddPeriod(start, end time.Time)
	Worked() time.Duration
	Periods() []Period
}

type Doable interface {
	Do()
	Undo()
	Done() bool
}

type Taggable interface {
	AddTag(ID)
	RemoveTag(ID)
	Tags() map[ID]struct{}
}

type Actionable interface {
	Identifiable
	Describable
	Workable
	Doable
	Taggable
}

type Period struct {
	Start time.Time
	End   time.Time
}

type Activity struct {
	id          ID
	description string
	duration    time.Duration
	done        bool
	tags        map[ID]struct{}
	periods     []Period
}

func NewActivity(id ID, description string) *Activity {
	return &Activity{id: id, description: description, tags: make(map[ID]struct{})}
}

func (a *Activity) Identify(id ID) {
	a.id = id
}

func (a Activity) Identity() ID {
	return a.id
}

func (a *Activity) Describe(description string) {
	a.description = description
}

func (a Activity) Description() string {
	return a.description
}

func (a *Activity) Work(duration time.Duration) {
	a.duration += duration
}

func (a *Activity) WorkPeriod(start, end time.Time) {
	oldPeriodDuration := sumPeriods(a.periods)
	a.periods = append(a.periods, Period{Start: start, End: end})
	a.periods = MergePeriods(a.periods)
	newPeriodDuration := sumPeriods(a.periods)
	a.duration += newPeriodDuration - oldPeriodDuration
}

func MergePeriods(periods []Period) []Period {
	if len(periods) <= 1 {
		return periods
	}
	sorted := make([]Period, len(periods))
	copy(sorted, periods)
	slices.SortFunc(sorted, func(a, b Period) int { return a.Start.Compare(b.Start) })

	merged := []Period{sorted[0]}
	for _, p := range sorted[1:] {
		last := &merged[len(merged)-1]
		if !p.Start.After(last.End) {
			if p.End.After(last.End) {
				last.End = p.End
			}
		} else {
			merged = append(merged, p)
		}
	}
	return merged
}

func sumPeriods(periods []Period) time.Duration {
	var d time.Duration
	for _, p := range periods {
		d += p.End.Sub(p.Start)
	}
	return d
}

func (a *Activity) AddPeriod(start, end time.Time) {
	a.periods = append(a.periods, Period{Start: start, End: end})
}

func (a Activity) Worked() time.Duration {
	return a.duration
}

func (a Activity) Periods() []Period {
	return a.periods
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
