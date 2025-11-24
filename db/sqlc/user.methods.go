package db

import "restapi/utils"

func (u *User) CheckPassword(password string) bool {

	res, err := utils.ComparePasswordAndHash(password, u.Password)

	if err != nil {
		return false
	}

	return res
}
