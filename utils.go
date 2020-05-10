package migrations

func getMigrationIdx(migration Migration, list []Migration) int {
	for i, m := range list {
		if m == migration {
			return i
		}
	}
	return -1
}
