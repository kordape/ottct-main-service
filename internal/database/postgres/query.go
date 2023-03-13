package postgres

func (db *DB) GetTwitterEntities() (entities []TwitterEntity, err error) {
	result := db.db.Find(&entities)
	if result.Error != nil {
		db.log.Error("Error performing query: %s", result.Error)
	}

	db.log.Debug("Number of entites returned: %d", result.RowsAffected)

	return
}
