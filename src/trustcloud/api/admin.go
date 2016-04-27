package api

func CheckLogin(user AdminUser) (bool, error) {
	result := (user.Password != "test")

	return result, nil
}
