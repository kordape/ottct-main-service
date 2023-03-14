package postgres

func (db *DB) GetTwitterEntities() (entities []TwitterEntity, err error) {
	result := db.db.Find(&entities)
	if result.Error != nil {
		db.log.Error("Error performing query: %s", result.Error)
	}

	db.log.Debug("Number of entites returned: %d", result.RowsAffected)

	return
}

func (db *DB) GetUserSupscriptions(userId string) (entities []TwitterEntity, err error) {
	var user User

	err = db.db.First(&user, "id = ?", userId).Error
	if err != nil {
		db.log.Error("Error performing query: %s", err)
	}

	err = db.db.Model(&user).Association("Subscriptions").Find(&entities)
	if err != nil {
		db.log.Error("Error performing query: %s", err)
	}

	return
}
