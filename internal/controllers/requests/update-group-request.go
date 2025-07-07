package requests

type UpdateGroupRequest struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	TeacherId int    `json:"teacher_id"`
}
