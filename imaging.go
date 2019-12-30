package onvif

import (
	"fmt"

	"github.com/apex/log"
)

const imageingNameSpace = "http://www.onvif.org/ver20/imaging/wsdl"

var imagingXMLNs = []string{
	`xmlns:timg="http://www.onvif.org/ver20/imaging/wsdl"`,
	`xmlns:tt="http://www.onvif.org/ver10/schema"`,
}

// GetImagingSettings fetch the ImagingConfiguration for the requested VideoSource.
func (device *Device) GetImagingSettings(videoSourceToken string) (ImagingSettings, error) {
	// Create SOAP
	soap := SOAP{
		Body: fmt.Sprintf(`<GetImagingSettings xmlns="%s">
					<VideoSourceToken>%s</VideoSourceToken>
		</GetImagingSettings>`, imageingNameSpace, videoSourceToken),
		// XMLNs:    imagingXMLNs,
		User:     device.User,
		Password: device.Password,
	}

	// Send SOAP request
	response, err := soap.SendRequest(device.Services[imageingNameSpace].XAddr)
	if err != nil {
		return ImagingSettings{}, err
	}

	// Get and parse list of profile to interface
	imagingSettingsInfo, err := response.ValueForPath("Envelope.Body.GetImagingSettingsResponse.ImagingSettings")
	if err != nil {
		return ImagingSettings{}, err
	}

	log.Debugf("VideoSourceToken[%s], imaging settings: %s", videoSourceToken, prettyJSON(imagingSettingsInfo))

	// Parse interface to struct
	imagingSettings := ImagingSettings{}
	if mapSettings, ok := imagingSettingsInfo.(map[string]interface{}); ok {
		if mapSettings["BacklightCompensation"] != nil {
			imagingSettings.BacklightCompensation.Mode = interfaceToString(mapSettings["BacklightCompensation"].(map[string]interface{})["Mode"])
			imagingSettings.BacklightCompensation.Level = interfaceToFloat64(mapSettings["BacklightCompensation"].(map[string]interface{})["Level"])
		}
		imagingSettings.Brightness = interfaceToFloat64(mapSettings["Brightness"])
		imagingSettings.ColorSaturation = interfaceToFloat64(mapSettings["ColorSaturation"])
		imagingSettings.Contrast = interfaceToFloat64(mapSettings["Contrast"])
		if mapSettings["Exposure"] != nil {
			imagingSettings.Exposure.Mode = interfaceToString(mapSettings["Exposure"].(map[string]interface{})["Mode"])
			imagingSettings.Exposure.Priority = interfaceToString(mapSettings["Exposure"].(map[string]interface{})["Priority"])
			imagingSettings.Exposure.MinExposureTime = interfaceToFloat64(mapSettings["Exposure"].(map[string]interface{})["MinExposureTime"])
			imagingSettings.Exposure.MaxExposureTime = interfaceToFloat64(mapSettings["Exposure"].(map[string]interface{})["MaxExposureTime"])
			imagingSettings.Exposure.MinGain = interfaceToFloat64(mapSettings["Exposure"].(map[string]interface{})["MinGain"])
			imagingSettings.Exposure.MaxGain = interfaceToFloat64(mapSettings["Exposure"].(map[string]interface{})["MaxGain"])
			imagingSettings.Exposure.MinIris = interfaceToFloat64(mapSettings["Exposure"].(map[string]interface{})["MinIris"])
			imagingSettings.Exposure.MaxIris = interfaceToFloat64(mapSettings["Exposure"].(map[string]interface{})["MaxIris"])
			imagingSettings.Exposure.ExposureTime = interfaceToFloat64(mapSettings["Exposure"].(map[string]interface{})["ExposureTime"])
			imagingSettings.Exposure.Gain = interfaceToFloat64(mapSettings["Exposure"].(map[string]interface{})["Gain"])
			imagingSettings.Exposure.Iris = interfaceToFloat64(mapSettings["Exposure"].(map[string]interface{})["Iris"])
		}
		if mapSettings["Focus"] != nil {
			imagingSettings.Focus.AutoFocusMode = interfaceToString(mapSettings["Focus"].(map[string]interface{})["AutoFocusMode"])
			imagingSettings.Focus.DefaultSpeed = interfaceToFloat64(mapSettings["Focus"].(map[string]interface{})["DefaultSpeed"])
			imagingSettings.Focus.NearLimit = interfaceToFloat64(mapSettings["Focus"].(map[string]interface{})["NearLimit"])
			imagingSettings.Focus.FarLimit = interfaceToFloat64(mapSettings["Focus"].(map[string]interface{})["FarLimit"])
		}
		imagingSettings.IrCutFilter = interfaceToString(mapSettings["IrCutFilter"])
		imagingSettings.Sharpness = interfaceToFloat64(mapSettings["Sharpness"])
		if mapSettings["WideDynamicRange"] != nil {
			imagingSettings.WideDynamicRange.Mode = interfaceToString(mapSettings["WideDynamicRange"].(map[string]interface{})["Mode"])
			imagingSettings.WideDynamicRange.Level = interfaceToFloat64(mapSettings["WideDynamicRange"].(map[string]interface{})["Level"])
		}
		if mapSettings["WhiteBalance"] != nil {
			imagingSettings.WhiteBalance.Mode = interfaceToString(mapSettings["WhiteBalance"].(map[string]interface{})["Mode"])
			imagingSettings.WhiteBalance.CbGain = interfaceToFloat64(mapSettings["WhiteBalance"].(map[string]interface{})["CbGain"])
			imagingSettings.WhiteBalance.CrGain = interfaceToFloat64(mapSettings["WhiteBalance"].(map[string]interface{})["CrGain"])
		}
	}
	return imagingSettings, err

}

func (device *Device) SetImagingSettings(exposureTime string) error {
	var soap SOAP
	// Create SOAP
	soap = SOAP{
		XMLNs:    []string{`xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"`, `xmlns:xsd="http://www.w3.org/2001/XMLSchema"`},
		User:     device.User,
		Password: device.Password,
	}

	soap.Body = `<SetImagingSettings xmlns="http://www.onvif.org/ver20/imaging/wsdl">
	<VideoSourceToken>VideoSource_1</VideoSourceToken>
    <Exposure>
    <Mode>AUTO</Mode>
    <MinExposureTime>10</MinExposureTime>
    <MaxExposureTime>` + exposureTime + `</MaxExposureTime>
	</Exposure>
	<WideDynamicRange>
	<Mode>ON</Mode>
	<Level>50</Level>
	</WideDynamicRange>
	</SetImagingSettings>`

	fmt.Println(soap.Body)

	// Send SOAP request
	rsp, err := soap.SendRequest(device.XAddr)
	fmt.Println(rsp)
	if err != nil {
		return err
	}
	return nil
}
