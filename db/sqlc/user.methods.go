package db

import (
	"restapi/utils"
	"strconv"
)

func (u *User) CheckPassword(password string) (bool, error) {

	res, err := utils.ComparePasswordAndHash(password, u.Password)

	if err != nil {
		return false, err
	}

	return res, nil
}

func (u *User) CreateTokens() (string, string, error) {
	jwtAccess, errAccess := utils.GenerateJWT(strconv.Itoa(int(u.ID)), "access")
	jwtRefresh, errRefresh := utils.GenerateJWT(strconv.Itoa(int(u.ID)), "refresh")

	if errAccess != nil {

		return "", "", errAccess
	}

	if errRefresh != nil {

		return "", "", errRefresh
	}

	return jwtAccess, jwtRefresh, nil
}
