package data

import (
	"database/sql"
	"mt-app1/entity"
)

func GetItemById(itemId int) (entity.Item, error) {
	querySql := `SELECT
	t.id,
	t.name,
	t.price
FROM
	item t
WHERE
	t.id = ?`
	return queryRow(func(r *sql.Row) (entity.Item, error) {
		var item entity.Item
		err := r.Scan(&item.Id, &item.Name, &item.Price)
		return item, err
	}, querySql, itemId)
}
