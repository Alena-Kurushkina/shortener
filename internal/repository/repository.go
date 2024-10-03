package repository

type Repository map[string]string

func NewRepository() *Repository{
	db:=make(Repository)
	return &db
}

func (r Repository) Insert(key, value string) {
	r[key]=value
}

func (r Repository) Select(key string) string {
	return r[key]
}


