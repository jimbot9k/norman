package dbobjects

import "encoding/json"

type Sequence struct {
	name       string
	schema     *Schema
	startValue int64
	increment  int64
	minValue   int64
	maxValue   int64
	cache      int64
	cycle      bool
}

func (s *Sequence) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name       string `json:"name"`
		StartValue int64  `json:"startValue"`
		Increment  int64  `json:"increment"`
		MinValue   int64  `json:"minValue"`
		MaxValue   int64  `json:"maxValue"`
		Cache      int64  `json:"cache"`
		Cycle      bool   `json:"cycle"`
	}{
		Name:       s.name,
		StartValue: s.startValue,
		Increment:  s.increment,
		MinValue:   s.minValue,
		MaxValue:   s.maxValue,
		Cache:      s.cache,
		Cycle:      s.cycle,
	})
}

func NewSequence(name string, startValue int64, increment int64) *Sequence {
	return &Sequence{
		name:       name,
		startValue: startValue,
		increment:  increment,
		minValue:   1,
		maxValue:   9223372036854775807, // max int64
		cache:      1,
		cycle:      false,
	}
}

func (s *Sequence) Name() string {
	return s.name
}

func (s *Sequence) Schema() *Schema {
	return s.schema
}

func (s *Sequence) SetSchema(schema *Schema) {
	s.schema = schema
}

func (s *Sequence) StartValue() int64 {
	return s.startValue
}

func (s *Sequence) SetStartValue(value int64) {
	s.startValue = value
}

func (s *Sequence) Increment() int64 {
	return s.increment
}

func (s *Sequence) SetIncrement(increment int64) {
	s.increment = increment
}

func (s *Sequence) MinValue() int64 {
	return s.minValue
}

func (s *Sequence) SetMinValue(minValue int64) {
	s.minValue = minValue
}

func (s *Sequence) MaxValue() int64 {
	return s.maxValue
}

func (s *Sequence) SetMaxValue(maxValue int64) {
	s.maxValue = maxValue
}

func (s *Sequence) Cache() int64 {
	return s.cache
}

func (s *Sequence) SetCache(cache int64) {
	s.cache = cache
}

func (s *Sequence) Cycle() bool {
	return s.cycle
}

func (s *Sequence) SetCycle(cycle bool) {
	s.cycle = cycle
}

// FullyQualifiedName returns schema.sequence format if schema is set
func (s *Sequence) FullyQualifiedName() string {
	if s.schema != nil {
		return s.schema.Name() + "." + s.name
	}
	return s.name
}
