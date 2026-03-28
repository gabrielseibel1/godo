package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gabrielseibel1/godo/data"
)

// jsonActivityTest mirrors the JSON structure for assertion purposes.
type jsonActivityTest struct {
	ID          string              `json:"id"`
	Description string              `json:"description"`
	Duration    time.Duration       `json:"duration"`
	Done        bool                `json:"done"`
	Tags        map[string]struct{} `json:"tags"`
	Periods     []jsonPeriodTest    `json:"periods,omitempty"`
}

type jsonPeriodTest struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// testEnv sets up a temporary directory with .godo/godo.json and configures the repo.
// Returns the temp dir path and a cleanup function.
func testEnv(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	godoDir := filepath.Join(dir, data.Dir)
	if err := os.MkdirAll(godoDir, 0o755); err != nil {
		t.Fatal(err)
	}
	jsonPath := filepath.Join(godoDir, data.JSONFile)
	if err := os.WriteFile(jsonPath, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Set the package-level repo to use this temp directory
	repo = data.NewJSONRepository(
		data.FileReader(jsonPath),
		data.FileWriter(jsonPath),
		data.JSONDecode,
		data.JSONEncode,
		data.Compare,
	)

	return dir
}

// readJSON reads and parses the godo.json file from the test environment.
func readJSON(t *testing.T, dir string) map[string]jsonActivityTest {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(dir, data.Dir, data.JSONFile))
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]jsonActivityTest
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatal(err)
	}
	return result
}

// runCmd executes a godo cobra command with the given args.
func runCmd(t *testing.T, args ...string) error {
	t.Helper()
	cmd := NewRootCmd()
	cmd.SetArgs(args)
	return cmd.Execute()
}

func TestCreate(t *testing.T) {
	dir := testEnv(t)

	err := runCmd(t, "create", "my-task", "do something")
	if err != nil {
		t.Fatal(err)
	}

	activities := readJSON(t, dir)
	a, ok := activities["my-task"]
	if !ok {
		t.Fatal("expected activity 'my-task' in godo.json")
	}
	if a.ID != "my-task" {
		t.Errorf("expected id 'my-task', got %q", a.ID)
	}
	if a.Description != "do something" {
		t.Errorf("expected description 'do something', got %q", a.Description)
	}
	if a.Done {
		t.Error("expected done=false")
	}
}

func TestCreateWithDescription(t *testing.T) {
	dir := testEnv(t)

	err := runCmd(t, "create", "with-desc", "my description")
	if err != nil {
		t.Fatal(err)
	}

	activities := readJSON(t, dir)
	a := activities["with-desc"]
	if a.Description != "my description" {
		t.Errorf("expected description 'my description', got %q", a.Description)
	}
}

func TestList(t *testing.T) {
	testEnv(t)

	if err := runCmd(t, "create", "task-a", "desc a"); err != nil {
		t.Fatal(err)
	}
	if err := runCmd(t, "create", "task-b", "desc b"); err != nil {
		t.Fatal(err)
	}

	// list should not error
	err := runCmd(t, "list")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	testEnv(t)

	if err := runCmd(t, "create", "my-item", "description"); err != nil {
		t.Fatal(err)
	}

	err := runCmd(t, "get", "my-item")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetNotFound(t *testing.T) {
	testEnv(t)

	err := runCmd(t, "get", "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing item")
	}
}

func TestDelete(t *testing.T) {
	dir := testEnv(t)

	if err := runCmd(t, "create", "to-delete", "will be deleted"); err != nil {
		t.Fatal(err)
	}

	err := runCmd(t, "delete", "to-delete")
	if err != nil {
		t.Fatal(err)
	}

	activities := readJSON(t, dir)
	if _, ok := activities["to-delete"]; ok {
		t.Error("expected activity to be deleted")
	}
}

func TestDo(t *testing.T) {
	dir := testEnv(t)

	if err := runCmd(t, "create", "do-me", "mark done"); err != nil {
		t.Fatal(err)
	}

	err := runCmd(t, "do", "do-me")
	if err != nil {
		t.Fatal(err)
	}

	activities := readJSON(t, dir)
	if !activities["do-me"].Done {
		t.Error("expected done=true")
	}
}

func TestUndo(t *testing.T) {
	dir := testEnv(t)

	if err := runCmd(t, "create", "undo-me", "mark undone"); err != nil {
		t.Fatal(err)
	}
	if err := runCmd(t, "do", "undo-me"); err != nil {
		t.Fatal(err)
	}

	err := runCmd(t, "undo", "undo-me")
	if err != nil {
		t.Fatal(err)
	}

	activities := readJSON(t, dir)
	if activities["undo-me"].Done {
		t.Error("expected done=false")
	}
}

func TestWorkDuration(t *testing.T) {
	dir := testEnv(t)

	if err := runCmd(t, "create", "work-task", "work on it"); err != nil {
		t.Fatal(err)
	}

	err := runCmd(t, "work", "work-task", "2h30m")
	if err != nil {
		t.Fatal(err)
	}

	activities := readJSON(t, dir)
	expected := time.Duration(2*time.Hour + 30*time.Minute)
	if activities["work-task"].Duration != expected {
		t.Errorf("expected duration %v, got %v", expected, activities["work-task"].Duration)
	}
}

func TestWorkTimeRange(t *testing.T) {
	dir := testEnv(t)

	if err := runCmd(t, "create", "range-2026-03-02", "time range"); err != nil {
		t.Fatal(err)
	}

	err := runCmd(t, "work", "range-2026-03-02", "18:00-21:00")
	if err != nil {
		t.Fatal(err)
	}

	activities := readJSON(t, dir)
	a := activities["range-2026-03-02"]
	if len(a.Periods) != 1 {
		t.Fatalf("expected 1 period, got %d", len(a.Periods))
	}
	if a.Periods[0].Start.Hour() != 18 {
		t.Errorf("expected start hour 18, got %d", a.Periods[0].Start.Hour())
	}
	if a.Periods[0].End.Hour() != 21 {
		t.Errorf("expected end hour 21, got %d", a.Periods[0].End.Hour())
	}
	expectedDur := 3 * time.Hour
	if a.Duration != expectedDur {
		t.Errorf("expected duration %v, got %v", expectedDur, a.Duration)
	}
}

func TestWorkTimeRangeCrossesMidnight(t *testing.T) {
	dir := testEnv(t)

	if err := runCmd(t, "create", "midnight-2026-03-13", "crosses midnight"); err != nil {
		t.Fatal(err)
	}

	err := runCmd(t, "work", "midnight-2026-03-13", "22:00-02:00")
	if err != nil {
		t.Fatal(err)
	}

	activities := readJSON(t, dir)
	a := activities["midnight-2026-03-13"]
	if len(a.Periods) != 1 {
		t.Fatalf("expected 1 period, got %d", len(a.Periods))
	}
	expectedDur := 4 * time.Hour
	if a.Duration != expectedDur {
		t.Errorf("expected duration %v, got %v", expectedDur, a.Duration)
	}
	// End should be on the next day
	if a.Periods[0].End.Day() != a.Periods[0].Start.Day()+1 {
		t.Error("expected end to be on the next day")
	}
}

func TestWorkTimeRange2400(t *testing.T) {
	dir := testEnv(t)

	if err := runCmd(t, "create", "midnight24-2026-03-24", "24:00 edge case"); err != nil {
		t.Fatal(err)
	}

	err := runCmd(t, "work", "midnight24-2026-03-24", "18:00-24:00")
	if err != nil {
		t.Fatal(err)
	}

	activities := readJSON(t, dir)
	a := activities["midnight24-2026-03-24"]
	expectedDur := 6 * time.Hour
	if a.Duration != expectedDur {
		t.Errorf("expected duration %v, got %v", expectedDur, a.Duration)
	}
}

func TestWorkWithDateFlag(t *testing.T) {
	dir := testEnv(t)

	if err := runCmd(t, "create", "date-flag", "with date"); err != nil {
		t.Fatal(err)
	}

	err := runCmd(t, "work", "date-flag", "14:00-18:00", "--date", "2026-05-15")
	if err != nil {
		t.Fatal(err)
	}

	activities := readJSON(t, dir)
	a := activities["date-flag"]
	if a.Periods[0].Start.Month() != time.May {
		t.Errorf("expected May, got %v", a.Periods[0].Start.Month())
	}
	if a.Periods[0].Start.Day() != 15 {
		t.Errorf("expected day 15, got %d", a.Periods[0].Start.Day())
	}
}

func TestTag(t *testing.T) {
	dir := testEnv(t)

	if err := runCmd(t, "create", "tag-me", "taggable"); err != nil {
		t.Fatal(err)
	}

	err := runCmd(t, "tag", "tag-me", "urgent")
	if err != nil {
		t.Fatal(err)
	}

	activities := readJSON(t, dir)
	if _, ok := activities["tag-me"].Tags["urgent"]; !ok {
		t.Error("expected tag 'urgent'")
	}
}

func TestUntag(t *testing.T) {
	dir := testEnv(t)

	if err := runCmd(t, "create", "untag-me", "untaggable"); err != nil {
		t.Fatal(err)
	}
	if err := runCmd(t, "tag", "untag-me", "remove-this"); err != nil {
		t.Fatal(err)
	}

	err := runCmd(t, "untag", "untag-me", "remove-this")
	if err != nil {
		t.Fatal(err)
	}

	activities := readJSON(t, dir)
	if _, ok := activities["untag-me"].Tags["remove-this"]; ok {
		t.Error("expected tag to be removed")
	}
}

func TestTagMultiple(t *testing.T) {
	dir := testEnv(t)

	if err := runCmd(t, "create", "a", "task a"); err != nil {
		t.Fatal(err)
	}
	if err := runCmd(t, "create", "b", "task b"); err != nil {
		t.Fatal(err)
	}

	err := runCmd(t, "tag", "a", "b", "shared-tag")
	if err != nil {
		t.Fatal(err)
	}

	activities := readJSON(t, dir)
	if _, ok := activities["a"].Tags["shared-tag"]; !ok {
		t.Error("expected tag on activity 'a'")
	}
	if _, ok := activities["b"].Tags["shared-tag"]; !ok {
		t.Error("expected tag on activity 'b'")
	}
}

func TestAutoWork(t *testing.T) {
	dir := testEnv(t)

	err := runCmd(t, "auto-work", "18:00-21:00", "--date", "2026-03-02")
	if err != nil {
		t.Fatal(err)
	}

	activities := readJSON(t, dir)
	a, ok := activities["work-2026-03-02"]
	if !ok {
		t.Fatal("expected auto-created activity")
	}
	if a.Description != "" {
		t.Errorf("expected empty description, got %q", a.Description)
	}
	if len(a.Periods) != 1 {
		t.Fatalf("expected 1 period, got %d", len(a.Periods))
	}
	expectedDur := 3 * time.Hour
	if a.Duration != expectedDur {
		t.Errorf("expected duration %v, got %v", expectedDur, a.Duration)
	}
}

func TestAutoWorkAppendsToExisting(t *testing.T) {
	dir := testEnv(t)

	if err := runCmd(t, "auto-work", "14:00-18:00", "--date", "2026-03-13"); err != nil {
		t.Fatal(err)
	}
	if err := runCmd(t, "auto-work", "19:00-01:00", "--date", "2026-03-13"); err != nil {
		t.Fatal(err)
	}

	activities := readJSON(t, dir)
	a := activities["work-2026-03-13"]
	if len(a.Periods) != 2 {
		t.Fatalf("expected 2 periods, got %d", len(a.Periods))
	}
	expectedDur := 10 * time.Hour
	if a.Duration != expectedDur {
		t.Errorf("expected duration %v, got %v", expectedDur, a.Duration)
	}
}

func TestAutoWorkYesterday(t *testing.T) {
	dir := testEnv(t)

	err := runCmd(t, "auto-work", "20:00-23:00", "--yesterday")
	if err != nil {
		t.Fatal(err)
	}

	yesterday := time.Now().AddDate(0, 0, -1)
	expectedID := "work-" + yesterday.Format("2006-01-02")

	activities := readJSON(t, dir)
	if _, ok := activities[expectedID]; !ok {
		t.Fatalf("expected activity %q", expectedID)
	}
}

func TestAutoWorkCrossesMidnight(t *testing.T) {
	dir := testEnv(t)

	err := runCmd(t, "auto-work", "22:00-02:00", "--date", "2026-03-15")
	if err != nil {
		t.Fatal(err)
	}

	activities := readJSON(t, dir)
	a := activities["work-2026-03-15"]
	expectedDur := 4 * time.Hour
	if a.Duration != expectedDur {
		t.Errorf("expected duration %v, got %v", expectedDur, a.Duration)
	}
}

func TestAutoWorkDotDot(t *testing.T) {
	dir := testEnv(t)

	now := time.Now()
	startHour := now.Add(-2 * time.Hour).Format("15:04")
	input := startHour + ".."

	err := runCmd(t, "auto-work", input)
	if err != nil {
		t.Fatal(err)
	}

	date := now
	if now.Hour() < 5 {
		date = now.AddDate(0, 0, -1)
	}
	expectedID := "work-" + date.Format("2006-01-02")

	activities := readJSON(t, dir)
	a, ok := activities[expectedID]
	if !ok {
		t.Fatalf("expected activity %q", expectedID)
	}
	if len(a.Periods) != 1 {
		t.Fatalf("expected 1 period, got %d", len(a.Periods))
	}
	// Duration should be roughly 2h (within 1 minute tolerance)
	if a.Duration < 2*time.Hour-time.Minute || a.Duration > 2*time.Hour+time.Minute {
		t.Errorf("expected ~2h duration, got %v", a.Duration)
	}
}

func TestAutoWorkNow(t *testing.T) {
	dir := testEnv(t)

	now := time.Now()
	startHour := now.Add(-1 * time.Hour).Format("15:04")
	input := startHour + "-now"

	err := runCmd(t, "auto-work", input)
	if err != nil {
		t.Fatal(err)
	}

	date := now
	if now.Hour() < 5 {
		date = now.AddDate(0, 0, -1)
	}
	expectedID := "work-" + date.Format("2006-01-02")

	activities := readJSON(t, dir)
	a, ok := activities[expectedID]
	if !ok {
		t.Fatalf("expected activity %q", expectedID)
	}
	if len(a.Periods) != 1 {
		t.Fatalf("expected 1 period, got %d", len(a.Periods))
	}
	// Duration should be roughly 1h (within 1 minute tolerance)
	if a.Duration < 1*time.Hour-time.Minute || a.Duration > 1*time.Hour+time.Minute {
		t.Errorf("expected ~1h duration, got %v", a.Duration)
	}
}

func TestAutoListToday(t *testing.T) {
	testEnv(t)

	today := time.Now().Format("2006-01-02")
	if err := runCmd(t, "auto-work", "14:00-18:00", "--date", today); err != nil {
		t.Fatal(err)
	}

	// auto-list with no args should show today
	err := runCmd(t, "auto-list")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAutoListSince(t *testing.T) {
	testEnv(t)

	if err := runCmd(t, "auto-work", "14:00-18:00", "--date", "2026-03-10"); err != nil {
		t.Fatal(err)
	}
	if err := runCmd(t, "auto-work", "18:00-20:00", "--date", "2026-03-12"); err != nil {
		t.Fatal(err)
	}

	// auto-list since 2026-03-10 should show both
	err := runCmd(t, "auto-list", "2026-03-10")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAutoListNoResults(t *testing.T) {
	testEnv(t)

	// No auto-work entries, should print "No extra hours today"
	err := runCmd(t, "auto-list")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAutoListMonth(t *testing.T) {
	testEnv(t)

	now := time.Now()
	first := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	mid := time.Date(now.Year(), now.Month(), 15, 0, 0, 0, 0, time.Local).Format("2006-01-02")

	if err := runCmd(t, "auto-work", "14:00-18:00", "--date", first); err != nil {
		t.Fatal(err)
	}
	if err := runCmd(t, "auto-work", "18:00-20:00", "--date", mid); err != nil {
		t.Fatal(err)
	}

	// auto-list month should show both entries
	err := runCmd(t, "auto-list", "month")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAutoListInvalidDate(t *testing.T) {
	testEnv(t)

	err := runCmd(t, "auto-list", "not-a-date")
	if err == nil {
		t.Fatal("expected error for invalid date")
	}
	if !strings.Contains(err.Error(), "invalid date") {
		t.Errorf("expected 'invalid date' message, got: %s", err.Error())
	}
}

func TestInit(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(dir)

	err := runCmd(t, "init")
	if err != nil {
		t.Fatal(err)
	}

	jsonPath := filepath.Join(data.Dir, data.JSONFile)
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		t.Error("expected godo.json to exist after init")
	}
}

func TestVersion(t *testing.T) {
	testEnv(t)

	err := runCmd(t, "version")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMultipleWorkAccumulates(t *testing.T) {
	dir := testEnv(t)

	if err := runCmd(t, "create", "accum", "accumulate"); err != nil {
		t.Fatal(err)
	}

	if err := runCmd(t, "work", "accum", "1h"); err != nil {
		t.Fatal(err)
	}
	if err := runCmd(t, "work", "accum", "30m"); err != nil {
		t.Fatal(err)
	}

	activities := readJSON(t, dir)
	expected := time.Duration(1*time.Hour + 30*time.Minute)
	if activities["accum"].Duration != expected {
		t.Errorf("expected duration %v, got %v", expected, activities["accum"].Duration)
	}
}

func TestSublist(t *testing.T) {
	testEnv(t)

	if err := runCmd(t, "create", "tagged-1", "first"); err != nil {
		t.Fatal(err)
	}
	if err := runCmd(t, "create", "untagged", "second"); err != nil {
		t.Fatal(err)
	}
	if err := runCmd(t, "tag", "tagged-1", "mytag"); err != nil {
		t.Fatal(err)
	}

	// sublist should not error
	err := runCmd(t, "sublist", "mytag")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCat(t *testing.T) {
	testEnv(t)

	if err := runCmd(t, "create", "cat-item", "for cat"); err != nil {
		t.Fatal(err)
	}
	if err := runCmd(t, "tag", "cat-item", "alpha"); err != nil {
		t.Fatal(err)
	}

	err := runCmd(t, "cat")
	if err != nil {
		t.Fatal(err)
	}
}

func TestWorkDurationThenPeriodPreservesBoth(t *testing.T) {
	dir := testEnv(t)

	if err := runCmd(t, "create", "mixed-2026-03-10", "mixed work"); err != nil {
		t.Fatal(err)
	}

	// Log 2h via plain duration
	if err := runCmd(t, "work", "mixed-2026-03-10", "2h"); err != nil {
		t.Fatal(err)
	}

	activities := readJSON(t, dir)
	if activities["mixed-2026-03-10"].Duration != 2*time.Hour {
		t.Fatalf("expected 2h after duration work, got %v", activities["mixed-2026-03-10"].Duration)
	}

	// Now log 1h via time range
	if err := runCmd(t, "work", "mixed-2026-03-10", "17:00-18:00"); err != nil {
		t.Fatal(err)
	}

	activities = readJSON(t, dir)
	a := activities["mixed-2026-03-10"]

	// Duration should be 2h + 1h = 3h (not just 1h from the period)
	if a.Duration != 3*time.Hour {
		t.Errorf("expected 3h total duration, got %v", a.Duration)
	}

	// Should have exactly 1 period
	if len(a.Periods) != 1 {
		t.Errorf("expected 1 period, got %d", len(a.Periods))
	}
}

// --- Error scenario tests ---

func TestCreateOptionalDescription(t *testing.T) {
	dir := testEnv(t)

	err := runCmd(t, "create", "no-desc")
	if err != nil {
		t.Fatal(err)
	}

	activities := readJSON(t, dir)
	a, ok := activities["no-desc"]
	if !ok {
		t.Fatal("expected activity 'no-desc' in godo.json")
	}
	if a.Description != "" {
		t.Errorf("expected empty description, got %q", a.Description)
	}
}

func TestCreateNoArgs(t *testing.T) {
	testEnv(t)

	err := runCmd(t, "create")
	if err == nil {
		t.Fatal("expected error for no args")
	}
	if !strings.Contains(err.Error(), "missing required argument") {
		t.Errorf("expected 'missing required argument' message, got: %s", err.Error())
	}
	if !strings.Contains(err.Error(), "godo create <id>") {
		t.Errorf("expected usage in error, got: %s", err.Error())
	}
}

func TestDeleteNotFound(t *testing.T) {
	testEnv(t)

	err := runCmd(t, "delete", "nonexistent")
	if err == nil {
		t.Fatal("expected error for deleting nonexistent item")
	}
	if !strings.Contains(err.Error(), `"nonexistent" not found`) {
		t.Errorf("expected not found message, got: %s", err.Error())
	}
}

func TestDoNotFound(t *testing.T) {
	testEnv(t)

	err := runCmd(t, "do", "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), `"nonexistent" not found`) {
		t.Errorf("expected not found message, got: %s", err.Error())
	}
}

func TestUndoNotFound(t *testing.T) {
	testEnv(t)

	err := runCmd(t, "undo", "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), `"nonexistent" not found`) {
		t.Errorf("expected not found message, got: %s", err.Error())
	}
}

func TestWorkNotFound(t *testing.T) {
	testEnv(t)

	err := runCmd(t, "work", "nonexistent", "1h")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), `"nonexistent" not found`) {
		t.Errorf("expected not found message, got: %s", err.Error())
	}
}

func TestGetNotFoundMessage(t *testing.T) {
	testEnv(t)

	err := runCmd(t, "get", "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), `"nonexistent" not found`) {
		t.Errorf("expected not found message, got: %s", err.Error())
	}
}

func TestWorkInvalidDuration(t *testing.T) {
	testEnv(t)

	if err := runCmd(t, "create", "bad-dur", "test"); err != nil {
		t.Fatal(err)
	}
	err := runCmd(t, "work", "bad-dur", "xyz")
	if err == nil {
		t.Fatal("expected error for invalid duration")
	}
}

func TestWorkInvalidTimeRange(t *testing.T) {
	testEnv(t)

	if err := runCmd(t, "create", "bad-range-2026-01-01", "test"); err != nil {
		t.Fatal(err)
	}
	err := runCmd(t, "work", "bad-range-2026-01-01", "25:00-26:00")
	if err == nil {
		t.Fatal("expected error for invalid time range")
	}
}

func TestMissingArgsShowUsage(t *testing.T) {
	testEnv(t)

	tests := []struct {
		name    string
		args    []string
		usage   string
	}{
		{"create no args", []string{"create"}, "godo create <id>"},
		{"delete no args", []string{"delete"}, "godo delete <id>"},
		{"get no args", []string{"get"}, "godo get <id>"},
		{"do no args", []string{"do"}, "godo do <id>"},
		{"undo no args", []string{"undo"}, "godo undo <id>"},
		{"work no args", []string{"work"}, "godo work <id>"},
		{"work one arg", []string{"work", "something"}, "godo work <id>"},
		{"auto-work no args", []string{"auto-work"}, "godo auto-work <HH:MM-HH:MM>"},
		{"tag no args", []string{"tag"}, "godo tag <id>"},
		{"tag one arg", []string{"tag", "only-id"}, "godo tag <id>"},
		{"untag no args", []string{"untag"}, "godo untag <id>"},
		{"untag one arg", []string{"untag", "only-id"}, "godo untag <id>"},
		{"sublist no args", []string{"sublist"}, "godo sublist <tag>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runCmd(t, tt.args...)
			if err == nil {
				t.Fatal("expected error for missing args")
			}
			if !strings.Contains(err.Error(), "missing required argument") {
				t.Errorf("expected 'missing required argument' message, got: %s", err.Error())
			}
			if !strings.Contains(err.Error(), tt.usage) {
				t.Errorf("expected usage %q in error, got: %s", tt.usage, err.Error())
			}
		})
	}
}
