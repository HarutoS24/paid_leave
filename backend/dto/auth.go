package dto

type LoginRequest struct {
	EmployeeID *string `json:"employeeId"`
	Password   *string `json:"password"`
}

type LoginResponse = BaseReponse
