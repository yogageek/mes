package logic

import (
	"material-management/db"
	"material-management/models"
)

func AddMaterial(m models.Material) error {
	_, err := db.InsertInventory(m)
	if err != nil {
		return err
	}
	return nil
}

func UpdateMaterial(um models.Material) error {
	if um.MaterialID == "" {
		err := models.MyError{
			Message:     "empty value",
			Description: um.MaterialID,
		}
		return err
	}
	// 如果quantity不帶 會直接把quantity改為0
	// if um.Quantity == 0 {
	// 	err := models.MyError{
	// 		Message:     "illegal value on quantity",
	// 		Description: fmt.Sprint(um.Quantity),
	// 	}
	// 	return err
	// }

	err := db.UpdateInventory(um)
	if err != nil {
		return err
	}
	return nil
}
