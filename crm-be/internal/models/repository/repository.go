package global_repository

import "skripsi-be/internal/config/database"

var db = database.GetDB()

func CreateData(entity interface{}) (interface{}, error) {
	tx := db.Create(&entity)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return entity, nil
}
func GetAllData(entity interface{}, condition interface{}) (interface{}, error) {
	tx := db.Find(&entity)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return entity, nil

}
func GetByConditionData(id string, entity interface{}, condition any, by any) (interface{}, error) {
	tx := db.First(&entity, condition, by)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return entity, nil

}
func UpdateData(id string, entity interface{}) (interface{}, error) {
	tx := db.Save(&entity)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return entity, nil

}
func DeleteData(id string, entity interface{}) (interface{}, error) {
	tx := db.Delete(&entity)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return entity, nil
}
