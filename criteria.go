package mellivora

type Criteria struct {
	query *Query
}

func (c *Criteria) And(fieldName string, operator int, value interface{}) *Criteria {
	return c
}

func (c *Criteria) AndEx(modelName, fieldName string, operator int, value interface{}) *Criteria {
	return c
}
