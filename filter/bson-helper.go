package filter

import "go.mongodb.org/mongo-driver/bson"

type BsonHelper struct {
}

func (f *BsonHelper) Filter(
	filters map[string][]interface{}) interface{} {

	if filters == nil {
		return nil
	}

	m := []bson.D{}
	for key, values := range filters {
		if len(values) == 0 {
			continue
		}
		vals := bson.A{}
		for _, v := range values {
			vals = append(vals, v)
		}
		m = append(m, bson.D{{key,
			bson.D{{"$in", vals}},
		}})
	}
	return bson.M{"$and": m}
}
