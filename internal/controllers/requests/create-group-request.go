package requests

type CreateGroupRequest struct {
	Name      string `json:"name"`
	TeacherId int    `json:"teacher_id"`
}
