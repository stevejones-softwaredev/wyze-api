package wyze

const wyzeDeveloperApi string = "wyze_developer_api"
const wyzeRequestTimestamp string = "4070908800000" // January 1, 2099 12:00:00 AM UTC for unspecified reasons
const wyzeApiHost string = "api.wyzecam.com"
const wyzeAuthHost string = "auth-prod.api.wyze.com"
const wyzeContentType string = "application/json"
const wyzeSvActionValue string = "5e02224ae0c64d328154737602d28833"

const wyzeAuthEndpoint string = "https://auth-prod.api.wyze.com/api/user/login"
const wyzeAccessTokenEndpoint string = "https://api.wyzecam.com/app/user/refresh_token"
const wyzeGetDeviceListEndpoint string = "https://api.wyzecam.com/app/v2/home_page/get_object_list"
const wyzeGetEventListEndpoint string = "https://api.wyzecam.com/app/v2/device/get_event_list"
const wyzeGetDevicePropertiesEndpoint string = "https://api.wyzecam.com/app/v2/device_list/get_property_list"
const wyzeRunActionEndpoint string = "https://api.wyzecam.com/app/v2/auto/run_action_list"

const wyzeLightPropPowerKey string = "P3"
const wyzeLightPropBrightnessKey string = "P1501"
const wyzeLightPropColorKey string = "P1507"

func getPropertyCodesToNamesMap() map[string]string {
  return map[string]string {
    "P1501": "brightness",
    "P1502": "color_temp",
    "P1505": "remaining_time",
    "P1506": "away_mode",
    "P1507": "color",
    "P1508": "control_light",
    "P1509": "power_loss_recovery",
    "P1510": "delay_off",
    "P1528": "sun_match",
    "P1529": "has_location",
    "P1530": "supports_sun_match",
    "P1531": "supports_timer",
    "P1": "notifications_enabled",
    "P3": "power_state",
    "P5": "online_state",
  };
}

func getPropertyNamesToCodesMap() map[string]string {
  return map[string]string {
    "brightness": "P1501",
    "color_temp": "P1502",
    "remaining_time" : "P1505",
    "away_mode" : "P1506",
    "color" : "P1507",
    "control_light" : "P1508",
    "power_loss_recovery": "P1509",
    "delay_off": "P1510",
    "sun_match": "P1528",
    "has_location": "P1529",
    "supports_sun_match": "P1530",
    "supports_timer": "P1531",
    "notifications_enabled": "P1",
    "power_state": "P3",
    "online_state": "P5",
  };
}
