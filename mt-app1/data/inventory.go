package data

import (
	"database/sql"
	"mt-app1/entity"
)

func GetAllInventoriesById(itemId int) ([]entity.Inventory, error) {
	querySql := `SELECT
	t.id,
	t.item_id,
	t.batch,
	t.count
FROM
	inventory t
WHERE
	t.item_id = ?`

	return query(func(r *sql.Rows) entity.Inventory {
		var i entity.Inventory
		r.Scan(&i.Id, &i.ItemId, &i.Batch, &i.Count)
		return i
	}, querySql, itemId)
}

func GetInventory(itemId int, batch string) (entity.Inventory, error) {
	querySql := `SELECT
	t.id,
	t.item_id,
	t.batch,
	t.count
FROM
	inventory t
WHERE
	t.item_id = ?
	AND t.batch = ?`

	return queryRow(func(r *sql.Row) (entity.Inventory, error) {
		var i entity.Inventory
		err := r.Scan(&i.Id, &i.ItemId, &i.Batch, &i.Count)
		return i, err
	}, querySql, itemId, batch)
}

func GetOutBoundInventory(itemId int, count int) (entity.Inventory, error) {
	querySql := `SELECT
	t.id,
	t.item_id,
	t.batch,
	t.count
FROM
	inventory t
WHERE
	t.item_id = ?
	AND t.count > ?
ORDER BY
	t.batch
LIMIT
	1`
	return queryRow(func(r *sql.Row) (entity.Inventory, error) {
		var i entity.Inventory
		err := r.Scan(&i.Id, &i.ItemId, &i.Batch, &i.Count)
		return i, err
	}, querySql, itemId, count)
}

func AddInventoryCount(exec executor, id int, count int) (sql.Result, error) {
	updateSql := `UPDATE
	inventory
SET
	COUNT = COUNT + ?
WHERE
	id = ?
	AND COUNT + ? >= 0`
	return exec.Exec(updateSql, count, id, count)
}
