package dto

type FileReference struct {
	ID       string `json:"id"`
	Filename []byte `json:"filename"`
}
