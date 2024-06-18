package service

import (
	"mt-app1/data"
	"mt-app1/entity"
	"mt-app1/optl"
	optlattr "mt-app1/optl/attr"

	"github.com/labstack/echo/v4"
)

const (
	itemServiceSpanName = "Item Service"
)

func GetItemInfo(c echo.Context, itemId int) (entity.Item, error) {
	span := optl.StartTrace(c.Request().Context(), itemServiceSpanName,
		optlattr.MethodName("GetItemInfo"), optlattr.IntParam("itemId", itemId))
	defer span.End()
	item, error := data.GetItemById(itemId)
	return item, error
}

func OutBound(c echo.Context, itemId int, count int) error {
	span := optl.StartTrace(c.Request().Context(), itemServiceSpanName,
		optlattr.MethodName("OutBound"),
		optlattr.IntParam("itemId", itemId), optlattr.IntParam("count", count))
	defer span.End()

	inventory, err := data.GetOutBoundInventory(itemId, count)
	if err != nil {
		return err
	}

	tx, err := getTransaction(c)
	if err != nil {
		return err
	}
	_, err = checkExecResult(data.AddInventoryCount(tx, inventory.Id, -count))
	return err
}
