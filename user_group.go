package zabbix

// UserGroupGroup represent Zabbix usergroup object
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usergroup/object
type UserGroup struct {
	UserGroupID string               `json:"usrgrpid,omitempty"`
	Name        string               `json:"name"`
	DebugMode   int                  `json:"debug_mode,string"`
	GUIAccess   int                  `json:"gui_access,string"`
	Status      int                  `json:"users_status,string"`
	Permissions usergrouppermissions `json:"hostgroup_rights,omitempty"`
}

// UserGroups is an array of UserGroup
type UserGroups []UserGroup

// UserGroupID represent Zabbix UserGroupID
type UserGroupID struct {
	UserGroupID string `json:"usrgrpid"`
}

// usergroupids is an array of UserGroupId
type usergroupids []UserGroupID

// UserGroupPermission represents zabbix usergroup permission object
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usergroup/object
type UserGroupPermission struct {
	ID         string `json:"id"`
	Permission int    `json:"permission"`
}

// usergrouppermissions is an array of UserGroupPermission
type usergrouppermissions []UserGroupPermission

// UserGroupsGet Wrapper for usergroup.get
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usergroup/get
func (api *API) UserGroupsGet(params Params) (res UserGroups, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	err = api.CallWithErrorParse("usergroup.get", params, &res)
	return
}

// UserGroupGetByID Gets usergroup by Id only if there is exactly 1 matching usergroup.
func (api *API) UserGroupGetByID(id string) (res *UserGroup, err error) {
	groups, err := api.UserGroupsGet(Params{"usrgrpids": id})
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

// UserGroupsCreate Wrapper for usergroup.create
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usergroup/create
func (api *API) UserGroupsCreate(UserGroups UserGroups) (err error) {
	response, err := api.CallWithError("usergroup.create", UserGroups)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	usergroupids := result["usrgrpids"].([]interface{})
	for i, id := range usergroupids {
		UserGroups[i].UserGroupID = id.(string)
	}
	return
}

// UserGroupsUpdate Wrapper for usergroup.update
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usergroup/update
func (api *API) UserGroupsUpdate(UserGroups UserGroups) (err error) {
	_, err = api.CallWithError("usergroup.update", UserGroups)
	return
}

// UserGroupsDelete Wrapper for usergroup.delete
// Cleans UserGroupID in all UserGroups elements if call succeed.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usergroup/delete
func (api *API) UserGroupsDelete(UserGroups UserGroups) (err error) {
	ids := make([]string, len(UserGroups))
	for i, usergroup := range UserGroups {
		ids[i] = usergroup.UserGroupID
	}

	err = api.UserGroupsDeleteByIds(ids)
	if err == nil {
		for i := range UserGroups {
			UserGroups[i].UserGroupID = ""
		}
	}
	return
}

// UserGroupsDeleteByIds Wrapper for usergroup.delete
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usergroup/delete
func (api *API) UserGroupsDeleteByIds(ids []string) (err error) {
	response, err := api.CallWithError("usergroup.delete", ids)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	usergroupids := result["usrgrpids"].([]interface{})
	if len(ids) != len(usergroupids) {
		err = &ExpectedMore{len(ids), len(usergroupids)}
	}
	return
}
