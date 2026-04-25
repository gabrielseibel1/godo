package types

import (
	"testing"
	"time"
)

func tp(h, m int) time.Time {
	return time.Date(2026, 3, 15, h, m, 0, 0, time.Local)
}

func TestMergePeriodsEmpty(t *testing.T) {
	got := MergePeriods(nil)
	if len(got) != 0 {
		t.Errorf("expected 0 periods, got %d", len(got))
	}
}

func TestMergePeriodsSingle(t *testing.T) {
	periods := []Period{{Start: tp(18, 0), End: tp(21, 0)}}
	got := MergePeriods(periods)
	if len(got) != 1 {
		t.Fatalf("expected 1 period, got %d", len(got))
	}
	if got[0].Start != tp(18, 0) || got[0].End != tp(21, 0) {
		t.Errorf("unexpected period: %v-%v", got[0].Start, got[0].End)
	}
}

func TestMergePeriodsNoOverlap(t *testing.T) {
	periods := []Period{
		{Start: tp(14, 0), End: tp(16, 0)},
		{Start: tp(18, 0), End: tp(20, 0)},
	}
	got := MergePeriods(periods)
	if len(got) != 2 {
		t.Fatalf("expected 2 periods, got %d", len(got))
	}
}

func TestMergePeriodsPartialOverlap(t *testing.T) {
	periods := []Period{
		{Start: tp(18, 0), End: tp(21, 0)},
		{Start: tp(20, 0), End: tp(23, 0)},
	}
	got := MergePeriods(periods)
	if len(got) != 1 {
		t.Fatalf("expected 1 period, got %d", len(got))
	}
	if got[0].Start != tp(18, 0) || got[0].End != tp(23, 0) {
		t.Errorf("expected 18:00-23:00, got %v-%v", got[0].Start.Format("15:04"), got[0].End.Format("15:04"))
	}
}

func TestMergePeriodsAdjacent(t *testing.T) {
	periods := []Period{
		{Start: tp(18, 0), End: tp(20, 0)},
		{Start: tp(20, 0), End: tp(22, 0)},
	}
	got := MergePeriods(periods)
	if len(got) != 1 {
		t.Fatalf("expected 1 period, got %d", len(got))
	}
	if got[0].Start != tp(18, 0) || got[0].End != tp(22, 0) {
		t.Errorf("expected 18:00-22:00, got %v-%v", got[0].Start.Format("15:04"), got[0].End.Format("15:04"))
	}
}

func TestMergePeriodsFullOverlap(t *testing.T) {
	periods := []Period{
		{Start: tp(18, 0), End: tp(23, 0)},
		{Start: tp(19, 0), End: tp(21, 0)},
	}
	got := MergePeriods(periods)
	if len(got) != 1 {
		t.Fatalf("expected 1 period, got %d", len(got))
	}
	if got[0].Start != tp(18, 0) || got[0].End != tp(23, 0) {
		t.Errorf("expected 18:00-23:00, got %v-%v", got[0].Start.Format("15:04"), got[0].End.Format("15:04"))
	}
}

func TestMergePeriodsChain(t *testing.T) {
	periods := []Period{
		{Start: tp(22, 0), End: tp(23, 0)},
		{Start: tp(18, 0), End: tp(20, 0)},
		{Start: tp(19, 0), End: tp(22, 30)},
	}
	got := MergePeriods(periods)
	if len(got) != 1 {
		t.Fatalf("expected 1 period, got %d", len(got))
	}
	if got[0].Start != tp(18, 0) || got[0].End != tp(23, 0) {
		t.Errorf("expected 18:00-23:00, got %v-%v", got[0].Start.Format("15:04"), got[0].End.Format("15:04"))
	}
}

func TestMergePeriodsUnsorted(t *testing.T) {
	periods := []Period{
		{Start: tp(20, 0), End: tp(22, 0)},
		{Start: tp(14, 0), End: tp(16, 0)},
		{Start: tp(18, 0), End: tp(21, 0)},
	}
	got := MergePeriods(periods)
	if len(got) != 2 {
		t.Fatalf("expected 2 periods, got %d", len(got))
	}
	if got[0].Start != tp(14, 0) || got[0].End != tp(16, 0) {
		t.Errorf("first period: expected 14:00-16:00, got %v-%v", got[0].Start.Format("15:04"), got[0].End.Format("15:04"))
	}
	if got[1].Start != tp(18, 0) || got[1].End != tp(22, 0) {
		t.Errorf("second period: expected 18:00-22:00, got %v-%v", got[1].Start.Format("15:04"), got[1].End.Format("15:04"))
	}
}
