package zabbix

// Proxy represent Zabbix proxy object
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxy/object#proxy
type Proxy struct {
	ProxyID        string `json:"proxyid,omitempty"`
	Host           string `json:"host"`
	Status         int    `json:"status,string"`
	Description    string `json:"description,omitempty"`
	TLSConnect     int    `json:"tls_connect,omitempty,string"`
	TLSAccept      int    `json:"tls_accept,omitempty,string"`
	TLSIssuer      string `json:"tls_issuer,omitempty"`
	TLSSubject     string `json:"tls_subject,omitempty"`
	TLSPSKIdentity string `json:"tls_psk_identity,omitempty"`
	TLSPSK         string `json:"tls_psk,omitempty"`
	ProxyAddress   string `json:"proxy_address,omitempty"`
}

// Proxies is an array of Proxy
type Proxies []Proxy

// ProxiesGet Wrapper for proxy.get
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxy/get
func (api *API) ProxiesGet(params Params) (res Proxies, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	err = api.CallWithErrorParse("proxy.get", params, &res)
	return
}

// ProxyGetByID Gets user by Id only if there is exactly 1 matching proxy.
func (api *API) ProxyGetByID(id string) (res *Proxy, err error) {
	groups, err := api.ProxiesGet(Params{"proxyids": id})
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

// ProxiesCreate Wrapper for proxy.create
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxy/create
func (api *API) ProxiesCreate(Proxies Proxies) (err error) {
	response, err := api.CallWithError("proxy.create", Proxies)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	proxyids := result["proxyids"].([]interface{})
	for i, id := range proxyids {
		Proxies[i].ProxyID = id.(string)
	}
	return
}

// ProxiesUpdate Wrapper for proxy.update
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxy/update
func (api *API) ProxiesUpdate(Proxies Proxies) (err error) {
	_, err = api.CallWithError("proxy.update", Proxies)
	return
}

// ProxiesDelete Wrapper for proxy.delete
// Cleans ProxyID in all Proxies elements if call succeed.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxy/delete
func (api *API) ProxiesDelete(Proxies Proxies) (err error) {
	ids := make([]string, len(Proxies))
	for i, proxy := range Proxies {
		ids[i] = proxy.ProxyID
	}

	err = api.ProxiesDeleteByIds(ids)
	if err == nil {
		for i := range Proxies {
			Proxies[i].ProxyID = ""
		}
	}
	return
}

// ProxiesDeleteByIds Wrapper for proxy.delete
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxy/delete
func (api *API) ProxiesDeleteByIds(ids []string) (err error) {
	response, err := api.CallWithError("proxy.delete", ids)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	proxyids := result["proxyids"].([]interface{})
	if len(ids) != len(proxyids) {
		err = &ExpectedMore{len(ids), len(proxyids)}
	}
	return
}
