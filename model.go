package onvif

// Device contains data of ONVIF camera
type Device struct {
	ID       string
	Name     string
	MACAddr  string
	XAddr    string
	User     string
	Password string
	Services map[string]Service
}

// DeviceInformation contains information of ONVIF camera
type DeviceInformation struct {
	FirmwareVersion string
	HardwareID      string
	Manufacturer    string
	Model           string
	SerialNumber    string
}

type NTPInformation struct {
	FromDHCP    bool
	NTPFromDHCP NetworkHost
	NTPManual   NetworkHost
}

type NetworkHost struct {
	Type        string
	IPv4Address string
	IPv6Address string
	DNSname     string
}

type SystemDateAndTime struct {
	DateTimeType    string
	DaylightSavings bool
	TimeZone        TimeZone_
	UTCDateTime     DateTime
	LocalDateTime   DateTime
}

type TimeZone_ struct {
	TZ string
}

type DateTime struct {
	Time Time_
	Date Date_
}

type Time_ struct {
	Hour   int
	Minute int
	Second int
}

type Date_ struct {
	Year  int
	Month int
	Day   int
}

// NetworkCapabilities contains networking capabilities of ONVIF camera
type NetworkCapabilities struct {
	DynDNS     bool
	IPFilter   bool
	IPVersion6 bool
	ZeroConfig bool
}

// DeviceCapabilities contains capabilities of an ONVIF camera
type DeviceCapabilities struct {
	Network   NetworkCapabilities
	Events    map[string]bool
	Streaming map[string]bool
}

// HostnameInformation contains hostname info of an ONVIF camera
type HostnameInformation struct {
	Name     string
	FromDHCP bool
}

// MediaBounds contains resolution of a video media
type MediaBounds struct {
	Height int
	Width  int
}

// MediaSourceConfig contains configuration of a media source
type MediaSourceConfig struct {
	Name        string
	Token       string
	SourceToken string
	Bounds      MediaBounds
}

// VideoRateControl contains rate control of a video
type VideoRateControl struct {
	BitrateLimit     int
	EncodingInterval int
	FrameRateLimit   int
	ConstantBitRate  bool
}

// VideoEncoderConfig contains configuration of a video encoder
type VideoEncoderConfig struct {
	Name           string
	Token          string
	Encoding       string
	Quality        int
	GovLength      int
	RateControl    VideoRateControl
	Resolution     MediaBounds
	SessionTimeout string
}

// AudioEncoderConfig contains configuration of an audio encoder
type AudioEncoderConfig struct {
	Name           string
	Token          string
	Encoding       string
	Bitrate        int
	SampleRate     int
	SessionTimeout string
}

// PTZConfig contains configuration of a PTZ control in camera
type PTZConfig struct {
	Name      string
	Token     string
	NodeToken string
}

// MediaProfile contains media profile of an ONVIF camera
type MediaProfile struct {
	Name               string
	Token              string
	VideoSourceConfig  MediaSourceConfig
	VideoEncoderConfig VideoEncoderConfig
	AudioSourceConfig  MediaSourceConfig
	AudioEncoderConfig AudioEncoderConfig
	PTZConfig          PTZConfig
}

// MediaURI contains streaming URI of an ONVIF camera
type MediaURI struct {
	URI                 string
	Timeout             string
	InvalidAfterConnect bool
	InvalidAfterReboot  bool
}

// MediaOSDS contains streaming URI of an ONVIF camera
type PosXY struct {
	x float64
	y float64
}

type Position struct {
	Type string
	Pos  PosXY
}

type OSDColor struct {
	Transparent int
	// Color       string
}

type TextString struct {
	IsPersistentText bool
	Type             string
	DateFormat       string
	TimeFormat       string
	FontSize         int
	FontColor        OSDColor
	BackgroundColor  OSDColor
	PlainText        string
}

type OSD struct {
	Token            string
	VideoSourceToken string
	Type             string
	Pos              Position
	Text             TextString
}

type NetworkInterfaces struct {
	Enabled string
	IPv4    IPv4
	// IPv6    IPv6
	Info NetworkInfo
}

type IPv4 struct {
	Enabled    string
	IPv4Config IPv4Config `json:"Config"`
}

type IPv4Config struct {
	DHCP     string
	FromDHCP FromDHCP `json:"FromDHCP"`
}

type FromDHCP struct {
	Address      string
	PrefixLength string
}

type NetworkInfo struct {
	HwAddress string
}

type Service struct {
	NameSpace string
	XAddr     string
	Version   Version
}

type Version struct {
	Major int
	Minor int
}

type ModeAndLevel struct {
	Mode  string  `json:"Mode"`
	Level float64 `json:"Level"`
}

type Exposure20 struct {
	Mode            string  `json:"Mode"`
	Priority        string  `json:"Priority"`
	MinExposureTime float64 `json:"MinExposureTime"`
	MaxExposureTime float64 `json:"MaxExposureTime"`
	MinGain         float64 `json:"MinGain"`
	MaxGain         float64 `json:"MaxGain"`
	MinIris         float64 `json:"MinIris"`
	MaxIris         float64 `json:"MaxIris"`
	ExposureTime    float64 `json:"ExposureTime"`
	Gain            float64 `json:"Gain"`
	Iris            float64 `json:"Iris"`
}

type FocusConfiguration20 struct {
	AutoFocusMode string  `json:"AutoFocusMode"`
	DefaultSpeed  float64 `json:"DefaultSpeed"`
	NearLimit     float64 `json:"NearLimit"`
	FarLimit      float64 `json:"FarLimit"`
}

type WhiteBalance20 struct {
	Mode   string  `json:"Mode"`
	CrGain float64 `json:"CrGain"`
	CbGain float64 `json:"CbGain"`
}

type ImagingSettings struct {
	BacklightCompensation ModeAndLevel
	Brightness            float64 `json:"Brightness"`
	ColorSaturation       float64 `json:"ColorSaturation"`
	Contrast              float64 `json:"Contrast"`
	Exposure              Exposure20
	Focus                 FocusConfiguration20
	IrCutFilter           string  `json:"IrCutFilter"`
	Sharpness             float64 `json:"Sharpness"`
	WideDynamicRange      ModeAndLevel
	WhiteBalance          WhiteBalance20
}
