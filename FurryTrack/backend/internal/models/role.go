package models


// Константы для присваевания роли пользователю
type Role string

const (
    RoleAdmin Role = "ADMIN"
    RoleUser  Role = "USER"
    RoleVet   Role = "VET"
)

func (r Role) IsValid() bool {
    switch r {
    case RoleAdmin, RoleUser, RoleVet:
        return true
    }
    return false
}