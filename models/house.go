package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// NewHouse creates a new base that points to house table
func NewHouse(db *sqlx.DB) *House {
	house := &House{}
	house.db = db
	house.table = "house"
	house.hasID = true

	return house
}

func NewItemInStorage(db *sqlx.DB) *ItemInStorage {
	storage := &ItemInStorage{}
	storage.db = db
	storage.table = "ITEM_IN_STORAGE"
	storage.hasID = true

	return storage
}

type ItemInStorage struct {
	Base
}

// House is a base
type House struct {
	Base
}

//HouseRowFromSQLResult returns the house that was last inserted to house table
func (h *House) HouseRowFromSQLResult(tx *sqlx.Tx, sqlResult sql.Result) (*HouseRow, error) {

	houseID, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return h.GetByID(tx, houseID)
}

// AllHouses returns every house in the housetable
func (h *House) AllHouses(tx *sqlx.Tx) ([]*HouseRow, error) {
	houses := []*HouseRow{}
	query := fmt.Sprintf("SELECT * FROM %v", h.table)
	err := h.db.Select(&houses, query)

	return houses, err
}

// GetByID returns a houseRow with the given id
func (h *House) GetByID(tx *sqlx.Tx, id int64) (*HouseRow, error) {
	house := &HouseRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=$1", h.table)
	err := h.db.Get(house, query, id)

	return house, err
}

// CreateHouse creates a house and a schedule for it
func (h *House) CreateHouse(tx *sqlx.Tx, name string, groceryDay string, household int64) (*HouseRow, error) {

	if name == "" {
		return nil, errors.New("House name cannot be blank")
	}

	data := make(map[string]interface{})
	data["name"] = name
	data["grocery_day"] = groceryDay
	data["household_number"] = household

	sqlResult, err := h.InsertIntoTable(tx, data)
	if err != nil {
		return nil, err
	}

	return h.HouseRowFromSQLResult(tx, sqlResult)
}

func (h *House) GetHouseUsers(tx *sqlx.Tx, houseID int64) ([]UserOwnTypeRow, error) {

	query := "SELECT U.ID, U.EMAIL, U.PASSWORD, U.USERNAME, O.OWN_TYPE, O.DESCRIPTION FROM USER_INFO U INNER JOIN MEMBER_OF M ON M.USER_ID = U.ID INNER JOIN OWNERSHIP O ON O.OWN_TYPE = M.OWN_TYPE WHERE M.HOUSE_ID = $1"

	data, err := h.GetCompoundModel(tx, query, houseID)

	users := createUserOwnTypeRows(data)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return users, err
}

func (h *House) GetHouseRecipes(tx *sqlx.Tx, houseID int64) ([]RecipeRow, error) {

	query := "SELECT R.ID, R.NAME, R.TYPE, R.SERVES_FOR FROM RECIPE R INNER JOIN HOUSE_RECIPE H ON R.ID = H.RECIPE_ID WHERE H.HOUSE_ID = $1"

	return h.GetRecipeForStruct(tx, query, houseID)
}

func (h *House) GetHouseStorage(tx *sqlx.Tx, houseID int64) ([]HouseStorageRow, error) {

	query := "SELECT S.HOUSE_ID, I.ID, S.AMOUNT, S.UNIT_ID, I.NAME, U.NAME FROM INGREDIENT I INNER JOIN ITEM_IN_STORAGE S ON I.ID = S.INGREDIENT_ID INNER JOIN UNIT U ON U.ID = S.UNIT_ID WHERE S.HOUSE_ID = $1"

	data, err := h.GetCompoundModel(tx, query, houseID)

	storage := createHouseStorageRows(data)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return storage, err
}

func (i *ItemInStorage) UpdateStorage(tx *sqlx.Tx, houseID int64, ingID int64, newAmt float64, newUnt int64) ([]HouseStorageRow, error) {

	var storage []HouseStorageRow
	var err error
	return storage, err
}

func (i *ItemInStorage) InsertToStorage(tx *sqlx.Tx, houseID int64, ingID int64, amount float64, unitID int64) (ItemInStorageRow, error) {

	data := make(map[string]interface{})
	data["house_id"] = houseID
	data["ingredient_id"] = ingID
	data["amount"] = amount
	data["unit_id"] = unitID

	_, err := i.InsertIntoMultiKeyTable(tx, data)

	if err != nil {
		fmt.Printf("%v", err)
	}

	query := fmt.Sprintf("SELECT * FROM ITEM_IN_STORAGE WHERE HOUSE_ID = %v AND INGREDIENT_ID = $1", houseID)

	res, err := i.GetCompoundModel(tx, query, ingID)

	storage := createItemInStorage(res)

	return storage, err
}

func (h *House) UpdateHouseHold(tx *sqlx.Tx, houseHold int64, houseID int64) (int64, error) {

	data := make(map[string]interface{})
	data["household_number"] = houseHold
	where := fmt.Sprintf("ID = %v", houseID)

	result, err := h.UpdateFromTable(tx, data, where)

	if err != nil {
		fmt.Println(err)
	}

	return h.AffectedRowsFromSqlResult(tx, result)

}

func (h *House) AffectedRowsFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (int64, error) {

	houseId, err := sqlResult.RowsAffected()
	if err != nil {

		fmt.Println(err)
	}

	return houseId, err
}
