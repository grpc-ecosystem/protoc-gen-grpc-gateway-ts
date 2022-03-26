package data

// Message stores the rendering information about message
type Message struct {
	// Nested shows whether this message is a nested message and needs to be exported
	Nested bool
	// Name is the name of the Message
	Name string
	//FQType is the fully qualified type name for the message itself
	FQType string
	// Enums is a list of NestedEnums inside
	Enums []*NestedEnum
	// Fields is a list of fields to render
	Fields []*Field
	// NonOneOfFields contains a subset of fields that are not in the one-of groups
	NonOneOfFields []*Field
	// Message is the nested messages defined inside the message
	Messages []*Message
	// OneOfFieldsGroups is the grouped list of one of fields with same index. so that renderer can render the clearing of other fields on set.
	OneOfFieldsGroups map[int32][]*Field
	// OneOfFieldNames is the names of one of fields with same index. so that renderer can render the clearing of other fields on set.
	OneOfFieldsNames map[int32]string
	// Comment is the comment of the message.
	Comment string
}

// HasOneOfFields returns true when the message has a one of field.
func (m *Message) HasOneOfFields() bool {
	return len(m.OneOfFieldsGroups) > 0
}

// NewMessage initialises and return a Message
func NewMessage() *Message {
	return &Message{
		Nested:            false,
		Name:              "",
		Enums:             make([]*NestedEnum, 0),
		Fields:            make([]*Field, 0),
		Messages:          make([]*Message, 0),
		OneOfFieldsGroups: make(map[int32][]*Field),
		OneOfFieldsNames:  make(map[int32]string),
		Comment:           "",
	}
}

// NestedEnum stores the information of enums defined inside a message
type NestedEnum struct {
	// Name of the Enum inside the class, which will be identical to the name
	// defined inside the message
	Name string
	// Type will have two types of value, and the difference can be told by
	// IsExternal attribute.
	// For external one, because during analysis stage there might not be a full map
	// of the types inside Registry. So the actual translation of this will
	// be left in the render time
	// If it is only types inside the file, it will be filled with the unique type name defined
	// up at the top level
	Type string
}

// Field stores the information about a field inside message
type Field struct {
	Name string
	// Type will be similar to NestedEnum.Type. Where scalar type and types inside
	// the same file will be short type
	// external types will have fully-qualified name and translated during render time
	Type string
	// IsExternal tells whether the type of this field is an external dependency
	IsExternal bool
	// IsOneOfField tells whether this field is part of a one of field.
	// one of fields will have extra method clearXXX,
	// and the setter accessor will clear out other fields in the group on set
	IsOneOfField bool
	// Message is the reference back to the parent message
	Message *Message
	// OneOfIndex is the index in the one of fields
	OneOfIndex int32
	// IsRepeated indicates whether the field is a repeated field
	IsRepeated bool
	// JSONName is the name of JSON.
	JSONName string
	// Comment is the comment of the field
	Comment string
}

// GetType returns some information of the type to aid the rendering
func (f *Field) GetType() *TypeInfo {
	return &TypeInfo{
		Type:       f.Type,
		IsRepeated: f.IsRepeated,
		IsExternal: f.IsExternal,
	}
}

// SetExternal mutate the IsExternal attribute
func (f *Field) SetExternal(external bool) {
	f.IsExternal = external
}

// MapEntryType is the generic entry type for both key and value
type MapEntryType struct {
	// Type of the map entry
	Type string
	// IsExternal indicates the field typeis external to its own package
	IsExternal bool
}

// GetType returns the type information for the type entry
func (m *MapEntryType) GetType() *TypeInfo {
	return &TypeInfo{
		Type:       m.Type,
		IsRepeated: false,
		IsExternal: m.IsExternal,
	}
}

// SetExternal mutate the IsExternal attribute inside
func (m *MapEntryType) SetExternal(external bool) {
	m.IsExternal = external
}
