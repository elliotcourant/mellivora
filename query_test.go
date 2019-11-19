package mellivora

import (
	"testing"
)

func TestQuery_Where(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		q := &Query{}
		q.Where(func(criteria *Criteria) *Criteria {
			return criteria.And("tableName", 0, "products")
		}).LeftJoin("nodes", "nodes", nil)
	})
}
