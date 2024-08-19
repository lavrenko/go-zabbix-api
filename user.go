package zabbix

// User represent Zabbix user object
// https://www.zabbix.com/documentation/current/en/manual/api/reference/user/object
type User struct {
	UserID   string       `json:"userid,omitempty"`
	Username string       `json:"username"`
	Password string       `json:"passwd"`
	RoleID   string       `json:"roleid"`
	Name     string       `json:"name"`
	Surname  string       `json:"surname"`
	Groups   usergroupids `json:"usrgrps"`
}

// Users is an array of User
type Users []User

// UserID represent Zabbix UserID
type UserID struct {
	UserID string `json:"userid"`
}

// userids is an array of UserId
type userids []UserID

// UsersGet Wrapper for user.get
// https://www.zabbix.com/documentation/current/en/manual/api/reference/user/get
func (api *API) UsersGet(params Params) (res Users, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	err = api.CallWithErrorParse("user.get", params, &res)
	return
}

// UserGetByID Gets user by Id only if there is exactly 1 matching user.
func (api *API) UserGetByID(id string) (res *User, err error) {
	groups, err := api.UsersGet(Params{"userids": id})
	if err != nil {
		return
	}

	if len(groups) == 1 {
		res = &groups[0]
	} else {
		e := ExpectedOneResult(len(groups))
		err = &e
	}
	return
}

// UsersCreate Wrapper for user.create
// https://www.zabbix.com/documentation/current/en/manual/api/reference/user/create
func (api *API) UsersCreate(Users Users) (err error) {
	response, err := api.CallWithError("user.create", Users)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	userids := result["userids"].([]interface{})
	for i, id := range userids {
		Users[i].UserID = id.(string)
	}
	return
}

// UsersUpdate Wrapper for user.update
// https://www.zabbix.com/documentation/current/en/manual/api/reference/user/update
func (api *API) UsersUpdate(Users Users) (err error) {
	_, err = api.CallWithError("user.update", Users)
	return
}

// UsersDelete Wrapper for user.delete
// Cleans UserID in all Users elements if call succeed.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/user/delete
func (api *API) UsersDelete(Users Users) (err error) {
	ids := make([]string, len(Users))
	for i, user := range Users {
		ids[i] = user.UserID
	}

	err = api.UsersDeleteByIds(ids)
	if err == nil {
		for i := range Users {
			Users[i].UserID = ""
		}
	}
	return
}

// UsersDeleteByIds Wrapper for user.delete
// https://www.zabbix.com/documentation/current/en/manual/api/reference/user/delete
func (api *API) UsersDeleteByIds(ids []string) (err error) {
	response, err := api.CallWithError("user.delete", ids)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	userids := result["userids"].([]interface{})
	if len(ids) != len(userids) {
		err = &ExpectedMore{len(ids), len(userids)}
	}
	return
}
