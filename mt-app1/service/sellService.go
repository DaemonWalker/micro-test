package service

import (
	"mt-app1/data"
	"mt-app1/model"
	"mt-app1/optl"
	optlattr "mt-app1/optl/attr"

	"github.com/labstack/echo/v4"
)

func BuyItem(c echo.Context, model *model.BuyItemModel) error {
	span := optl.StartTrace(c.Request().Context(), "Sell Service", optlattr.MethodName("BuyItem"))
	defer span.End()

	user, err := data.GetUserByEMail(model.User)
	if err != nil {
		return err
	}

	item, err := GetItemInfo(c, model.ItemId)
	if err != nil {
		return err
	}

	err = OutBound(c, item.Id, model.Count)
	if err != nil {
		return err
	}

	err = MinusMoney(c, user.Id, item.Price)
	if err != nil {
		return err
	}

	return nil
}
