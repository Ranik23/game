package models



// TODO: это надо automigrate
// теперь вопрос, может быть не писать все эти миграции SQL, а просто создать структуры и автоматически создать таблицы
type LoginInfo struct {
	Login string
	Hash []byte
}