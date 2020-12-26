package support

// Field represents a field that is apart of a ticket
type Field struct {
	Label string      `json:"label" bson:"label"`
	Kind  string      `json:"kind" bson:"kind"`
	Value interface{} `json:"value" bson:"value"`
}
