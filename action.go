package brokerTags

// Action - Custom type to hold value for broker action
type Action int

const (
	Create Action = iota // EnumIndex = 0
	Update               // EnumIndex = 1
)

func (a Action) String() string {
	return [...]string{"Created", "Updated"}[a]
}

func (a Action) getTagKey() string {
	return [...]string{CreatedAtTagKey, UpdatedAtTagKey}[a]
}
