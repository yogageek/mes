package db

import (
	"fmt"
	"log"
	"material-management/models"
	"material-management/util"

	_ "github.com/lib/pq" // postgres golang driver
)

// func recorderError(str string) error {
// 	errString := str[strings.LastIndex(str, ":")+1:]
// 	errString = "ops database: " + strings.TrimSpace(errString)
// 	return osb.HTTPStatusCodeError{
// 		StatusCode:   http.StatusServiceUnavailable,
// 		ErrorMessage: &errString,
// 	}
// }

// Get Latest Inventory
func QueryInventoryLastTime() (string, error) {
	db := CreatePGClient()
	defer db.Close()
	sqlStatement := `SELECT MAX(timestamp) FROM mes.inventory_list`

	var maxTimestamp string
	err := db.QueryRow(sqlStatement).Scan(&maxTimestamp)
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
	}
	return maxTimestamp, err
}

// 新增一筆物料
func InsertInventory(m models.Material) (models.Material, error) {

	db := CreatePGClient()
	defer db.Close()

	sqlStatement := `INSERT INTO mes.inventory_list (material_id, material_name, type_id, type_name, stock_location, total_stock, min_qty, timestamp) 
	VALUES ( $1, $2, $3, $4, $5, $6, $7, NOW()) 
	ON CONFLICT (material_id) DO UPDATE SET total_stock = mes.inventory_list.total_stock + $6;`
	util.LogOut(sqlStatement, m.MaterialID, m.MaterialName, m.TypeID, m.TypeName, m.StockLocation, m.TotalStock, m.MinQty)
	_, err := db.Exec(sqlStatement, m.MaterialID, m.MaterialName, m.TypeID, m.TypeName, m.StockLocation, m.TotalStock, m.MinQty)

	if err != nil {
		log.Println("Unable to execute the query, err:", err)
	}

	fmt.Println("Inserted a single record.")

	return m, err
}

// 修改指定id物料的數量
func UpdateInventory(m models.Material) error {

	db := CreatePGClient()
	defer db.Close()

	sqlStatement := `UPDATE mes.inventory_list SET total_stock = mes.inventory_list.total_stock - $2 WHERE material_id = $1;`
	res, err := db.Exec(sqlStatement, m.MaterialID, m.Quantity)
	util.LogOut(sqlStatement, m.MaterialID, m.Quantity)
	if err != nil {
		log.Println("Unable to execute the query, err:", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return err
}

// Get All Inventory
// func GetAllInventory() ([]models.Material, error) {
// 	var materials []models.Material

// 	db := CreatePGClient()
// 	defer db.Close()
// 	sqlStatement := `SELECT * FROM mes.inventory_list`

// 	rows, err := db.Query(sqlStatement)
// 	//defer不能放這 如果有err會錯
// 	if err != nil {
// 		log.Printf("Unable to execute the query. %v", err)
// 		return nil, err //需要return 不然接著rows.next會錯
// 	}
// 	defer rows.Close() //if row=nil, will cause this error

// 	for rows.Next() { // iterate over the rows
// 		var m models.Material
// 		err = rows.Scan(&m.MaterialID, &m.MaterialName, &m.TypeID, &m.TypeName, &m.TotalStock, &m.MinQty, &m.StockLocation, &m.Timestamp)
// 		if err != nil {
// 			log.Printf("Unable to scan the row. %v", err)
// 		}
// 		// get any error encountered during iteration
// 		err = rows.Err()
// 		if err != nil {
// 			log.Printf("err during iteration. %v", err)
// 		}
// 		materials = append(materials, m)
// 	}
// 	return materials, err
// }

// Get All Inventory 可帶條件 分頁
func QueryInventory(f *models.MaterialFilter) ([]models.Material, error) {
	var materials []models.Material

	//need a better way
	if f.Limit == 0 { //=size of one page
		f.Limit = 50
	}
	if f.Page == 0 {
		f.Page = 1
	}

	//如果有帶排序條件
	if f.SortBy != "" {
		if f.Desc == true {
			f.SortBy = "order by " + f.SortBy + " desc"
		} else {
			f.SortBy = "order by " + f.SortBy + " asc"
		}
	}

	db := CreatePGClient()
	defer db.Close()
	sqlStatement := fmt.Sprintf(`SELECT * FROM mes.inventory_list WHERE material_id LIKE '%%%s%%' 	
	AND material_name LIKE '%%%s%%'
	AND type_id LIKE '%%%s%%'
	AND type_name LIKE '%%%s%%'
	AND stock_location LIKE '%%%s%%' %s limit %d offset %d`,
		f.MaterialID, f.MaterialName, f.TypeID, f.TypeName, f.StockLocation, f.SortBy, f.Limit, (f.Page*f.Limit - f.Limit))
	util.LogOut(sqlStatement, f.MaterialID, f.MaterialName, f.TypeID, f.TypeName, f.StockLocation, f.SortBy, f.Limit, (f.Page*f.Limit - f.Limit))

	// offset := f.Page*f.Limit - f.Limit
	// offset2 := (f.Page - 1) * f.Limit
	// fmt.Println("offset", offset)
	// fmt.Println("offset2", offset2)

	rows, err := db.Query(sqlStatement)
	//defer不能放這 如果有err會錯
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
		return nil, err //需要return 不然接著rows.next會錯
	}
	defer rows.Close() //if row=nil, will cause this error

	for rows.Next() { // iterate over the rows
		var m models.Material
		err = rows.Scan(&m.MaterialID, &m.MaterialName, &m.TypeID, &m.TypeName, &m.TotalStock, &m.MinQty, &m.StockLocation, &m.Timestamp)
		if err != nil {
			log.Printf("Unable to scan the row. %v", err)
		}
		// get any error encountered during iteration
		err = rows.Err()
		if err != nil {
			log.Printf("err during iteration. %v", err)
		}
		materials = append(materials, m)
	}
	return materials, err
}

func QueryLogLastTime() (string, error) {
	db := CreatePGClient()
	defer db.Close()
	sqlStatement := `SELECT MAX(timestamp) FROM mes.inventory_log`

	var maxTimestamp string
	err := db.QueryRow(sqlStatement).Scan(&maxTimestamp)
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
	}
	return maxTimestamp, err
}

func QueryLogs() ([]models.MaterialLog, error) {
	var logs []models.MaterialLog

	db := CreatePGClient()
	defer db.Close()
	sqlStatement := `SELECT * FROM mes.inventory_log order by timestamp desc`

	rows, err := db.Query(sqlStatement)
	//defer不能放這 如果有err會錯
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
		return nil, err //需要return 不然接著rows.next會錯
	}
	defer rows.Close() //if row=nil, will cause this error

	for rows.Next() { // iterate over the rows
		var l models.MaterialLog
		err = rows.Scan(&l.MaterialID, &l.MaterialName, &l.TypeID, &l.TypeName, &l.InOrOut, &l.Timestamp, &l.StaffID, &l.StaffName, &l.Quantity, &l.LogSn)
		if err != nil {
			log.Printf("Unable to scan the row. %v", err)
		}
		// get any error encountered during iteration
		err = rows.Err()
		if err != nil {
			log.Printf("err during iteration. %v", err)
		}
		logs = append(logs, l)
	}
	return logs, err
}

// 新增一筆物料
func InsertLog(l models.MaterialLog) error {

	db := CreatePGClient()
	defer db.Close()

	sqlStatement := `INSERT INTO mes.inventory_log (material_id, material_name, type_id, type_name, in_or_out, quantity, staff_id, timestamp) 
	VALUES ( $1, $2, $3, $4, $5, $6, $7, NOW())`
	util.LogOut(sqlStatement, l.MaterialID, l.MaterialName, l.TypeID, l.TypeName, l.InOrOut, l.Quantity, l.StaffID)

	_, err := db.Exec(sqlStatement, l.MaterialID, l.MaterialName, l.TypeID, l.TypeName, l.InOrOut, l.Quantity, l.StaffID)

	if err != nil {
		log.Println("Unable to execute the query, err:", err)
	}

	fmt.Println("Inserted a single record.")

	return err
}
