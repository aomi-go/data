package mongoxentity

type DBRef struct {
	Ref string `bson:"$ref"`
	// 这里使用了StrObjectId类型来表示ObjectId的字符串形式，记得使用RegisterTypeEncoder\RegisterTypeDecoder
	ID StrObjectId `bson:"$id"`
}
