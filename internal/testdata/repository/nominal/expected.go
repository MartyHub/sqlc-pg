package nominal

type MyRepository struct {
	db Database
}

func NewMyRepository(db Database) MyRepository {
	return MyRepository{db: db}
}
