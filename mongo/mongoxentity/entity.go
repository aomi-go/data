package mongoxentity

type AbstractEntity struct {
	Id StrObjectId `bson:"_id,omitempty" json:"id,omitempty" yaml:"id,omitempty"`
}
