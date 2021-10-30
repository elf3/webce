package repository

type IQuery interface {
	Find(id int64) (map[string]interface{}, error)
	Select(query map[string]interface{}) map[string]interface{}
}

type IWriter interface {
	Update(query map[string]interface{}, data interface{}) (count int64, err error)
	Add(data map[string]interface{}) (int64, error)
}

type Repository interface {
	IQuery
	IWriter
}
