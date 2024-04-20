package wyze

type WyzeRefreshTokenRequest struct {
  Name string     `json:"email"`
  Password string `json:"password"`
}

type WyzeRefreshTokenResponse struct {
  RefreshToken string     `json:"refresh_token"`
}

type WyzeAccessTokenRequest struct {
  AppVer string           `json:"app_ver"`
  PhoneId string          `json:"phone_id"`
  RefreshToken string      `json:"refresh_token"`
  SC string               `json:"sc"`
  SV string               `json:"sv"`
  TS string               `json:"ts"`
}

type WyzeAccessTokenData struct {
  AccessToken  string`json:"access_token"`
}

type WyzeAccessTokenResponse struct {
  Code string              `json:"code"`
  Data WyzeAccessTokenData `json:"data"`
}

type WyzeDevice struct {
  MAC string                   `json:"mac"`
  Nickname string              `json:"nickname"`
  Model string                 `json:"product_model"`
  ProductType string           `json:"product_type"`
  DeviceMac string             `json:"device_mac"`
  Properties map[string]string `json:"properties,omitempty"`
}

type WyzeDeviceGroup struct {
  Name string             `json:"group_name"`
  DeviceList []WyzeDevice `json:"device_list"`
  PoweredOn bool          `json:"powered_on"`
}

type WyzeDeviceData struct {
  DeviceList []WyzeDevice     `json:"device_list"`
  GroupList []WyzeDeviceGroup `json:"device_group_list"`
}

type WyzeDeviceListResponse struct {
  Timestamp int       `json:"ts"`
  Code string         `json:"code"`
  Message string      `json:"msg"`
  Data WyzeDeviceData `json:"data"`
}

type WyzeDeviceListRequest struct {
  AppVer string           `json:"app_ver"`
  PhoneId string          `json:"phone_id"`
  AccessToken string      `json:"access_token"`
  SC string               `json:"sc"`
  SV string               `json:"sv"`
  TS string               `json:"ts"`
}

type WyzeActionProperty struct {
  Pid string    `json:"pid"`
  Pvalue string `json:"pvalue"`
}

type WyzeActionParamEntry struct {
  MAC string                 `json:"mac"`
  PList []WyzeActionProperty `json:"plist"`
}

type WyzeActionParams struct {
  List []WyzeActionParamEntry `json:"list"`
}

type WyzeActionList struct {
  ActionKey string        `json:"action_key"`
  InstanceId string       `json:"instance_id"`
  ProviderKey string      `json:"provider_key"`
  Params WyzeActionParams `json:"action_params"`
}

type WyzeRunActionListRequest struct {
  AppVer string               `json:"app_ver"`
  PhoneId string              `json:"phone_id"`
  AccessToken string          `json:"access_token"`
  SC string                   `json:"sc"`
  SV string                   `json:"sv"`
  TS string                   `json:"ts"`
  ActionList []WyzeActionList `json:"action_list"`
}

type WyzeDevicePropertyRequest struct {
  AppVer string               `json:"app_ver"`
  PhoneId string              `json:"phone_id"`
  AccessToken string          `json:"access_token"`
  SC string                   `json:"sc"`
  SV string                   `json:"sv"`
  TS string                   `json:"ts"`
  DeviceList []string         `json:"device_list"`
  TargetPropertyList []string `json:"target_pid_list"`
}

type WyzeDevicePropertyEntry struct {
  Pid string   `json:"pid"`
  Value string `json:"value"`
}

type WyzeDeviceProperties struct {
  DeviceMac string                     `json:"device_mac"`
  Properties []WyzeDevicePropertyEntry `json:"device_property_list"`
  PropertyMap map[string]string
}

type WyzeDevicePropertyData struct {
  DeviceList []WyzeDeviceProperties `json:"device_list"`
}

type WyzeDevicePropertyResponse struct {
  Code string                 `json:"code"`
  Data WyzeDevicePropertyData `json:"data"`
}

type WyzeEventRequest struct {
  AppVer string           `json:"app_ver"`
  PhoneId string          `json:"phone_id"`
  AccessToken string      `json:"access_token"`
  SC string               `json:"sc"`
  SV string               `json:"sv"`
  Devices []string        `json:"device_mac_list"`
  Count int               `json:"count"`
  OrderBy string          `json:"order_by"`
  PhoneSystemType string  `json:"phone_system_type"`
  BeginTime string        `json:"begin_time"`
  EndTime string          `json:"end_time"`
  Tags []int              `json:"event_tag_list"`
  TS string               `json:"ts"`
}

type WyzeDownloadedFile struct {
  Path string
  Url string
  Mac string
  Timestamp int64
}

type WyzeFile struct {
  URL string `json:"url"`
}

type WyzeCamEvent struct {
  DeviceMac string    `json:"device_mac"`
  EventTime int       `json:"event_ts"`
  FileList []WyzeFile `json:"file_list"`
}

type WyzeEventResponseData struct {
  EventList []WyzeCamEvent `json:"event_list"`
}

type WyzeEventResponse struct {
  Timestamp int              `json:"ts"`
  Code string                `json:"code"`
  Message string             `json:"msg"`
  Data WyzeEventResponseData `json:"data"`
}
