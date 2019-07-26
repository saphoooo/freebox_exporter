package main

type track struct {
	Success bool `json:"success"`
	Result  struct {
		AppToken string `json:"app_token"`
		TrackID  int    `json:"track_id"`
	} `json:"result"`
}

type grant struct {
	Success bool `json:"success"`
	Result  struct {
		Status    string `json:"status"`
		Challenge string `json:"challenge"`
	} `json:"result"`
}

type challenge struct {
	Success bool `json:"success"`
	Result  struct {
		LoggedIN  bool   `json:"logged_in,omitempty"`
		Challenge string `json:"challenge"`
	} `json:"result"`
}

type session struct {
	AppID    string `json:"app_id"`
	Password string `json:"password"`
}

type sessionToken struct {
	Msg       string `json:"msg,omitempty"`
	Success   bool   `json:"success"`
	UID       string `json:"uid,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	Result    struct {
		SessionToken string `json:"session_token,omitempty"`
		Challenge    string `json:"challenge"`
		Permissions  struct {
			Settings   bool `json:"settings,omitempty"`
			Contacts   bool `json:"contacts,omitempty"`
			Calls      bool `json:"calls,omitempty"`
			Explorer   bool `json:"explorer,omitempty"`
			Downloader bool `json:"downloader,omitempty"`
			Parental   bool `json:"parental,omitempty"`
			Pvr        bool `json:"pvr,omitempty"`
			Home       bool `json:"home,omitempty"`
			Camera     bool `json:"camera,omitempty"`
		} `json:"permissions,omitempty"`
	} `json:"result"`
}

type rrd struct {
	UID     string `json:"uid,omitempty"`
	Success bool   `json:"success"`
	Msg     string `json:"msg,omitempty"`
	Result  struct {
		DateStart int              `json:"date_start,omitempty"`
		DateEnd   int              `json:"date_end,omitempty"`
		Data      []map[string]int `json:"data,omitempty"`
	} `json:"result"`
	ErrorCode string `json:"error_code"`
}

type database struct {
	DB        string   `json:"db"`
	DateStart int      `json:"date_start,omitempty"`
	DateEnd   int      `json:"date_end,omitempty"`
	Precision int      `json:"precision,omitempty"`
	Fields    []string `json:"fields"`
}

type lanHost struct {
	Reachable   bool   `json:"reachable,omitempty"`
	PrimaryName string `json:"primary_name,omitempty"`
}

type lan struct {
	Success   bool      `json:"success"`
	Result    []lanHost `json:"result"`
	ErrorCode string    `json:"error_code"`
}

type idNameValue struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Value int    `json:"value,omitempty"`
}

type systemR struct {
	Sensors []idNameValue `json:"sensors"`
	Fans    []idNameValue `json:"fans"`
}

type system struct {
	Success   bool    `json:"success"`
	Result    systemR `json:"result"`
	ErrorCode string  `json:"error_code"`
}

type app struct {
	AppID      string `json:"app_id"`
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
	DeviceName string `json:"device_name"`
}

type api struct {
	authz string
	login string
}

type store struct {
	location string
}

type authInfo struct {
	myApp   app
	myAPI   api
	myStore store
}

type postRequest struct {
	method, url, header string
}
