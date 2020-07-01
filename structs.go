package main

import "bufio"

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

type connectionXdsl struct {
	Success bool `json:"success"`
	Result  struct {
		Status struct {
			Status     string `json:"status"`
			Modulation string `json:"modulation"`
			Protocol   string `json:"protocol"`
			Uptime     int    `json:"uptime"`
		} `json:"status"`
		Down struct {
			Attn       int  `json:"attn"`
			Attn10     int  `json:"attn_10"`
			Crc        int  `json:"crc"`
			Es         int  `json:"es"`
			Fec        int  `json:"fec"`
			Ginp       bool `json:"ginp"`
			Hec        int  `json:"hec"`
			Maxrate    int  `json:"maxrate"`
			Nitro      bool `json:"nitro"`
			Phyr       bool `json:"phyr"`
			Rate       int  `json:"rate"`
			RtxC       int  `json:"rtx_c,omitempty"`
			RtxTx      int  `json:"rtx_tx,omitempty"`
			RtxUc      int  `json:"rtx_uc,omitempty"`
			Rxmt       int  `json:"rxmt"`
			RxmtCorr   int  `json:"rxmt_corr"`
			RxmtUncorr int  `json:"rxmt_uncorr"`
			Ses        int  `json:"ses"`
			Snr        int  `json:"snr"`
			Snr10      int  `json:"snr_10"`
		} `json:"down"`
		Up struct {
			Attn       int  `json:"attn"`
			Attn10     int  `json:"attn_10"`
			Crc        int  `json:"crc"`
			Es         int  `json:"es"`
			Fec        int  `json:"fec"`
			Ginp       bool `json:"ginp"`
			Hec        int  `json:"hec"`
			Maxrate    int  `json:"maxrate"`
			Nitro      bool `json:"nitro"`
			Phyr       bool `json:"phyr"`
			Rate       int  `json:"rate"`
			RtxC       int  `json:"rtx_c,omitempty"`
			RtxTx      int  `json:"rtx_tx,omitempty"`
			RtxUc      int  `json:"rtx_uc,omitempty"`
			Rxmt       int  `json:"rxmt"`
			RxmtCorr   int  `json:"rxmt_corr"`
			RxmtUncorr int  `json:"rxmt_uncorr"`
			Ses        int  `json:"ses"`
			Snr        int  `json:"snr"`
			Snr10      int  `json:"snr_10"`
		} `json:"up"`
	}
}

type database struct {
	DB        string   `json:"db"`
	DateStart int      `json:"date_start,omitempty"`
	DateEnd   int      `json:"date_end,omitempty"`
	Precision int      `json:"precision,omitempty"`
	Fields    []string `json:"fields"`
}

type freeplug struct {
	Success bool              `json:"success"`
	Result  []freeplugNetwork `json:"result"`
}

type freeplugNetwork struct {
	ID      string           `json:"id"`
	Members []freeplugMember `json:"members"`
}

type freeplugMember struct {
	ID            string `json:"id"`
	Local         bool   `json:"local"`
	NetRole       string `json:"net_role"`
	EthPortStatus string `json:"eth_port_status"`
	EthFullDuplex bool   `json:"eth_full_duplex"`
	HasNetwork    bool   `json:"has_network"`
	EthSpeed      int    `json:"eth_speed"`
	Inative       int    `json:"inactive"`
	NetId         string `json:"net_id"`
	RxRate        int    `json:"rx_rate"`
	TxRate        int    `json:"tx_rate"`
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

type system struct {
	Success bool `json:"success"`
	Result  struct {
		Mac              string `json:"mac,omitempty"`
		FanRPM           int    `json:"fan_rpm,omitempty"`
		BoxFlavor        string `json:"box_flavor,omitempty"`
		TempCpub         int    `json:"temp_cpub,omitempty"`
		TempCpum         int    `json:"temp_cpum,omitempty"`
		DiskStatus       string `json:"disk_status,omitempty"`
		TempHDD          int    `json:"temp_hdd,omitempty"`
		BoardName        string `json:"board_name,omitempty"`
		TempSW           int    `json:"temp_sw,omitempty"`
		Uptime           string `json:"uptime,omitempty"`
		UptimeVal        int    `json:"uptime_val,omitempty"`
		UserMainStorage  string `json:"user_main_storage,omitempty"`
		BoxAuthenticated bool   `json:"box_authenticated,omitempty"`
		Serial           string `json:"serial,omitempty"`
		FirmwareVersion  string `json:"firmware_version,omitempty"`
	}
}

type wifiAccessPoint struct {
	Name string `json:"name,omitempty"`
	ID   int    `json:"id,omitempty"`
}

type wifi struct {
	Success bool              `json:"success"`
	Result  []wifiAccessPoint `json:"result,omitempty"`
}

type wifiStation struct {
	Hostname           string `json:"hostname,omitempty"`
	MAC                string `json:"mac,omitempty"`
	State              string `json:"state,omitempty"`
	Inactive           int    `json:"inactive,omitempty"`
	RXBytes            int    `json:"rx_bytes,omitempty"`
	TXBytes            int    `json:"tx_bytes,omitempty"`
	ConnectionDuration int    `json:"conn_duration,omitempty"`
	TXRate             int    `json:"tx_rate,omitempty"`
	RXRate             int    `json:"rx_rate,omitempty"`
	Signal             int    `json:"signal,omitempty"`
}

type wifiStations struct {
	Success bool          `json:"success"`
	Result  []wifiStation `json:"result,omitempty"`
}

type app struct {
	AppID      string `json:"app_id"`
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
	DeviceName string `json:"device_name"`
}

type api struct {
	authz        string
	login        string
	loginSession string
}

type store struct {
	location string
}

type authInfo struct {
	myApp    app
	myAPI    api
	myStore  store
	myReader *bufio.Reader
}

type postRequest struct {
	method, url, header string
}
