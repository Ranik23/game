package models




type Player struct {
	ID   int
	Name string
	Age  int
	Projects []*Project
}