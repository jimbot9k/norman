package dbobjects

import "encoding/json"

type TriggerTiming string

const (
	TriggerTimingBefore    TriggerTiming = "BEFORE"
	TriggerTimingAfter     TriggerTiming = "AFTER"
	TriggerTimingInsteadOf TriggerTiming = "INSTEAD OF"
)

type TriggerEvent string

const (
	TriggerEventInsert TriggerEvent = "INSERT"
	TriggerEventUpdate TriggerEvent = "UPDATE"
	TriggerEventDelete TriggerEvent = "DELETE"
)

type Trigger struct {
	name       string
	table      *Table
	definition string
	timing     TriggerTiming
	events     []TriggerEvent
	function   *Function
	forEach    string // ROW or STATEMENT
}

func (t *Trigger) MarshalJSON() ([]byte, error) {
	var functionName string
	if t.function != nil {
		functionName = t.function.Name()
	}
	return json.Marshal(struct {
		Name       string         `json:"name"`
		Definition string         `json:"definition"`
		Timing     TriggerTiming  `json:"timing"`
		Events     []TriggerEvent `json:"events"`
		Function   string         `json:"function,omitempty"`
		ForEach    string         `json:"forEach"`
	}{
		Name:       t.name,
		Definition: t.definition,
		Timing:     t.timing,
		Events:     t.events,
		Function:   functionName,
		ForEach:    t.forEach,
	})
}

func NewTrigger(name string, definition string) *Trigger {
	return &Trigger{
		name:       name,
		definition: definition,
		events:     []TriggerEvent{},
		forEach:    "ROW",
	}
}

func (t *Trigger) Name() string {
	return t.name
}

func (t *Trigger) Table() *Table {
	return t.table
}

func (t *Trigger) SetTable(table *Table) {
	t.table = table
}

func (t *Trigger) Definition() string {
	return t.definition
}

func (t *Trigger) SetDefinition(definition string) {
	t.definition = definition
}

func (t *Trigger) Timing() TriggerTiming {
	return t.timing
}

func (t *Trigger) SetTiming(timing TriggerTiming) {
	t.timing = timing
}

func (t *Trigger) Events() []TriggerEvent {
	return t.events
}

func (t *Trigger) AddEvent(event TriggerEvent) {
	t.events = append(t.events, event)
}

func (t *Trigger) Function() *Function {
	return t.function
}

func (t *Trigger) SetFunction(function *Function) {
	t.function = function
}

func (t *Trigger) ForEach() string {
	return t.forEach
}

func (t *Trigger) SetForEach(forEach string) {
	t.forEach = forEach
}
