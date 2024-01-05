package brokerTags

// Action - Custom type to hold value for broker action
type Action int

const (
	createdAtTagKey = "Created at"
	updatedAtTagKey = "Updated at"
)

const (
	Create Action = iota // EnumIndex = 0
	Update               // EnumIndex = 1
)

func (a Action) getTagKey() string {
	return [...]string{createdAtTagKey, updatedAtTagKey}[a]
}
