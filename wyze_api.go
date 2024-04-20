package wyze

import (
  "github.com/go-resty/resty/v2"
  "github.com/golang-jwt/jwt/v5"
  "fmt"
  "log"
  "strconv"
  "strings"
  "net/http"
  "time"
  "io"
  "os"
)

func isRefreshTokenValid(refreshToken string) bool {
  if (len(refreshToken) == 0) {
    return false
  }

  token, _, tokenErr := new(jwt.Parser).ParseUnverified(refreshToken, jwt.MapClaims{})

  if (tokenErr != nil) {
    fmt.Println("Token Error: ", tokenErr)
    return false;
  }

  expiration, expErr := token.Claims.GetExpirationTime()

  if (expErr != nil) {
    fmt.Println("Expiration Claim Error: ", expErr)
    return false;
  }

  return (time.Now().Before(expiration.Time))
}

func GetWyzeRefreshToken(client *resty.Client, username string, password string, keyId string, apiKey string) string {
  var refreshTokenResponse WyzeRefreshTokenResponse
  
  payload := WyzeRefreshTokenRequest{
    Name: username,
    Password: password,
  }

  _, err := client.R().
    SetHeader("Content-Type", wyzeContentType).
    SetHeader("Host", wyzeAuthHost).
    SetHeader("Keyid", keyId).
    SetHeader("Apikey", apiKey).
    SetBody(&payload).
    SetResult(&refreshTokenResponse).
    Post(wyzeAuthEndpoint)

  if err != nil {
    fmt.Println(err)
    return ""
  } else {
    return refreshTokenResponse.RefreshToken
  }
}


func GetWyzeAccessToken(client *resty.Client, env map[string]string) string {
  filePath := env["WYZE_HOME"] + "refresh_token.txt"
  refreshToken := readRefreshToken(filePath)
  if (!isRefreshTokenValid(refreshToken)) {
    fmt.Println("Create new refresh token")
    refreshToken = GetWyzeRefreshToken(client, env["WYZE_USERNAME"], env["WYZE_PASSWORD_HASH"], env["WYZE_KEY_ID"], env["WYZE_API_KEY"])
    writeRefreshToken(refreshToken, filePath)
  }

  var accessTokenResponse WyzeAccessTokenResponse

  payload := WyzeAccessTokenRequest{
    AppVer: wyzeDeveloperApi,
    PhoneId: wyzeDeveloperApi,
    RefreshToken: refreshToken,
    SC: wyzeDeveloperApi,
    SV: wyzeDeveloperApi,
    TS: wyzeRequestTimestamp,
  }

  _, err := client.R().
    SetHeader("Content-Type", wyzeContentType).
    SetHeader("Host", wyzeApiHost).
    SetBody(&payload).
    SetResult(&accessTokenResponse).
    Post(wyzeAccessTokenEndpoint)

    if err != nil {
      fmt.Println(err)
      return ""
    } else {
      return accessTokenResponse.Data.AccessToken
    }
}

func GetWyzeCamThumbnails(client *resty.Client,
    downloadDirectory string,
    accessToken string,
    count int,
    devices []string,
    tags []int,
    begin_time time.Time,
    end_time time.Time) []WyzeDownloadedFile {
  var thumbnailPaths []WyzeDownloadedFile
  var eventResponse WyzeEventResponse

  payload := WyzeEventRequest{
    AppVer: wyzeDeveloperApi,
    PhoneId: wyzeDeveloperApi,
    AccessToken: accessToken,
    SC: wyzeDeveloperApi,
    SV: wyzeDeveloperApi,
    Devices: devices,
    Count: count,
    OrderBy: "1",
    PhoneSystemType: "1",
    BeginTime: strconv.FormatInt(begin_time.UnixMilli(), 10),
    EndTime: strconv.FormatInt(end_time.UnixMilli(), 10),
    Tags: tags,
    TS: wyzeRequestTimestamp,
  }

  _, err := client.R().
    SetHeader("Content-Type", wyzeContentType).
    SetHeader("Host", wyzeApiHost).
    SetBody(&payload).
    SetResult(&eventResponse).
    Post(wyzeGetEventListEndpoint)

  if err != nil {
    fmt.Println(err)
  } else {
    for _, event := range eventResponse.Data.EventList {
      for _, file := range event.FileList {
        fileName := fmt.Sprintf("%s%d.jpg", downloadDirectory, event.EventTime)

        _,err := os.Stat(fileName)

        if err == nil {
          log.Println("File", fileName, "already exists")
        } else {
          saveFile(file.URL, fileName)
          log.Println("File", fileName, "downloaded")
          downloadedFile := WyzeDownloadedFile{
            Path: fileName,
            Url: file.URL,
            Mac: event.DeviceMac,
            Timestamp: getTimestampFromFile(fileName),
          }
          thumbnailPaths = append(thumbnailPaths, downloadedFile)
        }
      }
    }
  }

  return thumbnailPaths
}

func getTimestampFromFile(filePath string) int64 {
  ts := strings.Split(filePath, ".")[0]
  slice := strings.Split(ts, "/")
  ts = strings.Split(ts, "/")[len(slice) - 1]
  value,_ := strconv.ParseInt(ts, 10, 64)

  return value
}

func GetWyzeDeviceList(client *resty.Client,
    accessToken string) []WyzeDevice {
  var deviceListResponse WyzeDeviceListResponse

  payload := WyzeDeviceListRequest{
    AppVer: wyzeDeveloperApi,
    PhoneId: wyzeDeveloperApi,
    AccessToken: accessToken,
    SC: wyzeDeveloperApi,
    SV: wyzeDeveloperApi,
    TS: wyzeRequestTimestamp,
  }

  _, err := client.R().
    SetHeader("Content-Type", wyzeContentType).
    SetHeader("Host", wyzeApiHost).
    SetBody(&payload).
    SetResult(&deviceListResponse).
    Post(wyzeGetDeviceListEndpoint)

  if err != nil {
    fmt.Println(err)
    return []WyzeDevice{}
  } else {
    return deviceListResponse.Data.DeviceList
  }
}

func GetWyzeBulbList(client *resty.Client,
    accessToken string) []WyzeDevice {
  devices := GetWyzeDeviceList(client, accessToken)
  var bulbs []WyzeDevice

  for _,device := range devices {
     if (device.ProductType == "MeshLight") {
       device.DeviceMac = device.MAC
       bulbs = append(bulbs, device)
     }
  }

  return bulbs
}

func GetWyzeGroupList(client *resty.Client,
    accessToken string) []WyzeDeviceGroup {
  var deviceListResponse WyzeDeviceListResponse

  payload := WyzeDeviceListRequest{
    AppVer: wyzeDeveloperApi,
    PhoneId: wyzeDeveloperApi,
    AccessToken: accessToken,
    SC: wyzeDeveloperApi,
    SV: wyzeDeveloperApi,
    TS: wyzeRequestTimestamp,
  }

  _, err := client.R().
    SetHeader("Content-Type", wyzeContentType).
    SetHeader("Host", wyzeApiHost).
    SetBody(&payload).
    SetResult(&deviceListResponse).
    Post(wyzeGetDeviceListEndpoint)

  deviceListResponse.Data.DeviceList = IntegrateDeviceProperties(client, accessToken, deviceListResponse.Data.DeviceList)

  devicesMap := make(map[string]WyzeDevice)
  var deviceMacs []string

  for _,device := range deviceListResponse.Data.DeviceList {
    devicesMap[device.MAC] = device
    device.DeviceMac = device.MAC
    deviceMacs = append(deviceMacs, device.MAC)
  }

  var newGroups []WyzeDeviceGroup
  
  for _,group := range deviceListResponse.Data.GroupList {
    group.PoweredOn = false
    var newDeviceList []WyzeDevice
  
    for _,device := range group.DeviceList {
      propDevice := devicesMap[device.DeviceMac]
      group.PoweredOn = (group.PoweredOn || (propDevice.Properties["power_state"] == "1"))
      newDeviceList = append(newDeviceList, propDevice)
    }

    group.DeviceList = newDeviceList
    newGroups = append(newGroups, group)
  }

  deviceListResponse.Data.GroupList = newGroups

  if err != nil {
    fmt.Println(err)
    return []WyzeDeviceGroup{}
  } else {
    return deviceListResponse.Data.GroupList
  }
}

func BuildGroupNameMap(groups []WyzeDeviceGroup) map[string]WyzeDeviceGroup {
  groupMap := make(map[string]WyzeDeviceGroup)
  
  for _,group := range groups {
    groupMap[group.Name] = group
  }
  
  return groupMap
}

func BuildDeviceNameMap(devices []WyzeDevice) map[string]WyzeDevice {
  deviceMap := make(map[string]WyzeDevice)

  for _,device := range devices {
    deviceMap[device.Nickname] = device
  }

  return deviceMap
}

func MakeGroupDeviceMap(groupName string, groupMap map[string]WyzeDeviceGroup) map[string]string {
  deviceMap := make(map[string]string)
  
  group,ok := groupMap[groupName]
  
  if (ok) {
    for _,device := range group.DeviceList {
      deviceMap[device.MAC] = device.Model
    }
  }
  
  return deviceMap
}

func MakeDeviceMap(devices []WyzeDevice) map[string]string {
  deviceMap := make(map[string]string)

  for _,device := range devices {
    deviceMap[device.MAC] = device.Model
  }
  
  return deviceMap
}

func SetWyzeProperties(client *resty.Client,
    accessToken string,
    devices map[string]string,
    properties map[string]string) {

  var propList []WyzeActionProperty

  namesToCodes := getPropertyNamesToCodesMap()
  
  for key, value := range properties {
    wyzeKey, ok := namesToCodes[key]
    var prop WyzeActionProperty

    if (ok) {
      prop = WyzeActionProperty {
        Pid: wyzeKey,
        Pvalue: value,
      }
    } else {
      prop = WyzeActionProperty {
        Pid: key,
        Pvalue: value,
      }
    }
    propList = append(propList, prop)
  }
  
  var actionList []WyzeActionList
  
  for device,model := range devices {
    var paramEntries []WyzeActionParamEntry
    paramEntry := WyzeActionParamEntry{
      MAC: device,
      PList: propList,
    }
    
    paramEntries = append(paramEntries, paramEntry)

    param := WyzeActionParams{
      List: paramEntries,
    }

    action := WyzeActionList{
      ActionKey: "set_mesh_property",
      InstanceId: device,
      ProviderKey: model,
      Params: param,
    }
    
    actionList = append(actionList, action)
  }

  payload := WyzeRunActionListRequest{
    AppVer: wyzeDeveloperApi,
    PhoneId: wyzeDeveloperApi,
    AccessToken: accessToken,
    SC: wyzeDeveloperApi,
    SV: wyzeSvActionValue,
    TS: wyzeRequestTimestamp,
    ActionList: actionList,
  }

  _, err := client.R().
    SetHeader("Content-Type", wyzeContentType).
    SetHeader("Host", wyzeApiHost).
    SetBody(&payload).
    Post(wyzeRunActionEndpoint)
   
  if (err != nil) {
    log.Println(err)
  }
}

func GetWyzeDeviceProperties(client *resty.Client,
    accessToken string, devices []string, properties []string) WyzeDevicePropertyResponse {
  var devicePropertyResponse WyzeDevicePropertyResponse

  payload := WyzeDevicePropertyRequest{
    AppVer: wyzeDeveloperApi,
    PhoneId: wyzeDeveloperApi,
    AccessToken: accessToken,
    SC: wyzeDeveloperApi,
    SV: wyzeDeveloperApi,
    TS: wyzeRequestTimestamp,
    DeviceList: devices,
    TargetPropertyList: properties,
  }

  _, err := client.R().
    SetHeader("Content-Type", wyzeContentType).
    SetHeader("Host", wyzeApiHost).
    SetBody(&payload).
    SetResult(&devicePropertyResponse).
    Post(wyzeGetDevicePropertiesEndpoint)
   
  if (err != nil) {
    log.Println(err)
  }
  
  var newDeviceList []WyzeDeviceProperties

  codeToName := getPropertyCodesToNamesMap()
  
  for _,device := range devicePropertyResponse.Data.DeviceList {
    device.PropertyMap = make(map[string]string)
    for _,props := range device.Properties {
      propName, propOk := codeToName[props.Pid]

      if (propOk) {
        device.PropertyMap[propName] = props.Value
      } else {
        device.PropertyMap[props.Pid] = props.Value
      }
    }
    newDeviceList = append(newDeviceList, device)
  }
  
  devicePropertyResponse.Data.DeviceList = newDeviceList

  return devicePropertyResponse
}

func IntegrateDeviceProperties(client *resty.Client,
    accessToken string, devices []WyzeDevice) []WyzeDevice {
  var deviceMacs []string
  var propDevices []WyzeDevice

  for _,device := range devices {
    deviceMacs = append(deviceMacs, device.MAC)
  }

  properties := GetWyzeDeviceProperties(client, accessToken, deviceMacs, []string{})

  propMap := make(map[string]map[string]string)

  for _,prop := range properties.Data.DeviceList {
    propMap[prop.DeviceMac] = prop.PropertyMap
  }

  for _,device := range devices {
    device.Properties = propMap[device.MAC]
    propDevices = append(propDevices, device)
  }

  return propDevices;
}

func readRefreshToken(filePath string) string {
  data, err := os.ReadFile(filePath)
  if err != nil {
    log.Println(err)
    return ""
  }

  return string(data[:])
}

func writeRefreshToken(refreshToken string, filePath string) {
  f, err := os.Create(filePath)
  if err != nil {
    log.Println(err)
    return
  }

  defer f.Close()

  fmt.Fprintf(f, "%s", refreshToken)
}


func saveFile(url string, fileName string) {
  response, err := http.Get(url)
  if err != nil {
    return
  }
  defer response.Body.Close()

  file, err := os.Create(fileName)
  if err != nil {
    return
  }
  defer file.Close()

  _, err = io.Copy(file, response.Body)
  if err != nil {
    return
  }
}

