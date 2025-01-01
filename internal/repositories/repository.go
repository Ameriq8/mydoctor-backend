package repositories

// Repository defines the basic CRUD operations.
type Repository[T any] interface {
	Find(id int64) (*T, error)
	FindMany(filter map[string]interface{}) ([]T, error)
	Create(entity *T) (*T, error)
	CreateMany(entities []T) ([]T, error)
	Update(id int64, updates map[string]interface{}) (*T, error)
	UpdateMany(filter map[string]interface{}, updates map[string]interface{}) (int64, error)
	Delete(id int64) (*T, error)
	DeleteMany(filter map[string]interface{}) ([]T, error)
}
