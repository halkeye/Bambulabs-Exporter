package exporter

import (
	"cmp"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	Debug    bool
	Username string
	Password string
	IP       string
	Topic    string
}

type Exporter struct {
	config Config
	client mqtt.Client

	// Metrics
	amsHumidityMetric      *prometheus.GaugeVec
	amsTempMetric          *prometheus.GaugeVec
	amsColorMetric         *prometheus.GaugeVec
	amsTypeMetric          *prometheus.GaugeVec
	layerNumberMetric      prometheus.Gauge
	printErrorMetric       prometheus.Gauge
	wifiSignalMetric       prometheus.Gauge
	bigFan1SpeedMetric     prometheus.Gauge
	bigFan2SpeedMetric     prometheus.Gauge
	chamberTemperMetric    prometheus.Gauge
	coolingFanSpeedMetric  prometheus.Gauge
	failReasonMetric       prometheus.Gauge
	fanGearMetric          prometheus.Gauge
	mcPercentMetric        prometheus.Gauge
	mcPrintErrorCodeMetric prometheus.Gauge
	mcPrintStageMetric     prometheus.Gauge
	mcPrintSubStageMetric  prometheus.Gauge
	mcRemainingTimeMetric  prometheus.Gauge
	nozzleTargetTemperMetric prometheus.Gauge
	nozzleTemperMetric     prometheus.Gauge
	bedTargetTemperMetric    prometheus.Gauge
	bedTemperMetric    prometheus.Gauge
}

func NewExporter() *Exporter {
	cfg := Config{}
	err := envconfig.Process("BAMBULABS", &cfg)
	if err != nil {
		panic(err)
	}

	exporter := &Exporter{
		config: cfg,
	}

	exporter.initMetrics()
	return exporter
}

func (e *Exporter) initMetrics() {
	e.amsHumidityMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ams_humidity",
		Help: "humidity of the ams",
	}, []string{"ams_number"})
	e.amsTempMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ams_temp",
		Help: "temperature of the ams",
	}, []string{"ams_number"})
	e.amsColorMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ams_tray_color",
		Help: "color of material in ams tray",
	}, []string{"ams_number", "tray_number", "tray_color"})
	e.amsTypeMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ams_tray_type",
		Help: "type of material in ams tray",
	}, []string{"ams_number", "tray_number", "tray_type"})
	e.layerNumberMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "layer_number",
		Help: "layer number of the print head in gcode",
	})
	e.printErrorMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "print_error",
		Help: "Print error int",
	})
	e.wifiSignalMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wifi_signal",
		Help: "Wifi signal in dBm",
	})
	e.bigFan1SpeedMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "big_fan1_speed",
		Help: "Big Fan 1 Speed",
	})
	e.bigFan2SpeedMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "big_fan2_speed",
		Help: "Big Fan 2 Speed",
	})
	e.chamberTemperMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "chamber_temper",
		Help: "Chamber Temperature of Printer",
	})
	e.coolingFanSpeedMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "cooling_fan_speed",
		Help: "Cooling Fan Speed",
	})
	e.failReasonMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "fail_reason",
		Help: "Print Failure Reason",
	})
	e.fanGearMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "fan_gear",
		Help: "Fan Gear",
	})
	e.mcPercentMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mc_percent",
		Help: "Percentage of Progress of print",
	})
	e.mcPrintErrorCodeMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mc_print_error_code",
		Help: "Print Progress Error Code",
	})
	e.mcPrintStageMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mc_print_stage",
		Help: "Print Progress Stage",
	})
	e.mcPrintSubStageMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mc_print_sub_stage",
		Help: "Print Progress Sub Stage",
	})
	e.mcRemainingTimeMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mc_remaining_time",
		Help: "Print Progress Remaining Time in minutes",
	})
	e.nozzleTargetTemperMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "nozzle_target_temper",
		Help: "Nozzle Target Temperature Metric",
	})
	e.nozzleTemperMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "nozzle_temper",
		Help: "Nozzle Temperature Metric",
	})
	e.bedTargetTemperMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bed_target_temper",
		Help: "Bed target temperature metric",
	})
	e.bedTemperMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bed_temper",
		Help: "Bed temperature metric",
	})
}

func (e *Exporter) ConnectToBroker() {
	port := 8883
	clientID := cmp.Or(os.Getenv("OVERRIDE_CLIENT_ID"), "bambulabs-prometheus-exporter")

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ssl://%s:%d", e.config.IP, port))
	opts.SetClientID(clientID)
	opts.SetUsername(e.config.Username)
	opts.SetPassword(e.config.Password)
	opts.SetDefaultPublishHandler(e.messagePubHandler)
	opts.SetAutoReconnect(true)
	opts.OnConnect = e.buildConnectHandler()
	opts.OnConnectionLost = e.buildConnectLostHandler()

	opts.SetTLSConfig(e.newTLSConfig())
	e.client = mqtt.NewClient(opts)
	token := e.client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (e *Exporter) messagePubHandler(client mqtt.Client, msg mqtt.Message) {
	s := msg.Payload()
	data := BambuLabsX1C{}
	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %s\n", err)
		return
	}

	if data.Print.Command != "push_status" {
		fmt.Printf("Ignoring command: %s\n", data.Print.Command)
		return
	}

	e.layerNumberMetric.Set(float64(data.Print.LayerNum))
	e.printErrorMetric.Set(float64(data.Print.PrintError))

	wifi_signal, _ := strconv.ParseFloat(strings.ReplaceAll(data.Print.WifiSignal, "dBm", ""), 64)
	e.wifiSignalMetric.Set(wifi_signal)

	big_fan1_speed, _ := strconv.ParseFloat(data.Print.BigFan1Speed, 64)
	e.bigFan1SpeedMetric.Set(big_fan1_speed)

	big_fan2_speed, _ := strconv.ParseFloat(data.Print.BigFan2Speed, 64)
	e.bigFan2SpeedMetric.Set(big_fan2_speed)

	e.chamberTemperMetric.Set(data.Print.ChamberTemper)

	cooling_fan_speed, _ := strconv.ParseFloat(data.Print.CoolingFanSpeed, 64)
	e.coolingFanSpeedMetric.Set(cooling_fan_speed)

	fail_reason, _ := strconv.ParseFloat(data.Print.FailReason, 64)
	e.failReasonMetric.Set(fail_reason)

	e.fanGearMetric.Set(float64(data.Print.FanGear))
	e.mcPercentMetric.Set(float64(data.Print.McPercent))

	mc_print_error_code, _ := strconv.ParseFloat(data.Print.McPrintErrorCode, 64)
	e.mcPrintErrorCodeMetric.Set(mc_print_error_code)

	mc_print_stage, _ := strconv.ParseFloat(data.Print.McPrintStage, 64)
	e.mcPrintStageMetric.Set(mc_print_stage)

	e.mcPrintStageMetric.Set(float64(data.Print.McPrintSubStage))
	e.mcRemainingTimeMetric.Set(float64(data.Print.McRemainingTime))
	e.nozzleTemperMetric.Set(float64(data.Print.NozzleTemper))
	e.nozzleTargetTemperMetric.Set(float64(data.Print.NozzleTargetTemper))
	e.bedTargetTemperMetric.Set(data.Print.BedTargetTemper)
	e.bedTemperMetric.Set(data.Print.BedTemper)

	for _, ams := range data.Print.Ams.Ams {
		humidity, _ := strconv.ParseFloat(ams.Humidity, 64)
		e.amsHumidityMetric.With(prometheus.Labels{"ams_number": ams.ID}).Set(humidity)

		temp, _ := strconv.ParseFloat(ams.Temp, 64)
		e.amsTempMetric.With(prometheus.Labels{"ams_number": ams.ID}).Set(temp)
		for _, tray := range ams.Tray {
			baseLabels := prometheus.Labels{
				"ams_number":  ams.ID,
				"tray_number": tray.ID,
			}

			e.amsTypeMetric.DeletePartialMatch(baseLabels)
			e.amsTypeMetric.MustCurryWith(baseLabels).With(prometheus.Labels{"tray_type": tray.TrayType}).Set(1)

			e.amsColorMetric.DeletePartialMatch(baseLabels)
			e.amsColorMetric.MustCurryWith(baseLabels).With(prometheus.Labels{"tray_color": tray.TrayColor}).Set(1)
		}
	}
}

func (e *Exporter) buildConnectHandler() mqtt.OnConnectHandler {
	return func(client mqtt.Client) {
		dt := time.Now()
		fmt.Printf("Connected: %s\n", dt.String())
		client.Subscribe(e.config.Topic, 1, nil).Wait()
	}
}

func (e *Exporter) buildConnectLostHandler() mqtt.ConnectionLostHandler {
	return func(client mqtt.Client, err error) {
		fmt.Printf("Connect lost: %+v\n", err)
	}
}

func (e *Exporter) newTLSConfig() *tls.Config {
	return &tls.Config{InsecureSkipVerify: true}
}

func (e *Exporter) StartHTTPServer() {
	http.HandleFunc("/", e.home)
	http.HandleFunc("/healthz", e.healthz)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Printf("Listening http://127.0.0.1:9101\n")
}

func (e *Exporter) home(w http.ResponseWriter, r *http.Request) {
	const body = `<html>
				<head>
					<title>BambuLabs Exporter Metrics</title>
				</head>
				<body>
					<h1>BambuLabs Exporter</h1>
					<p><a href='` + "/metrics" + `'>metrics</a></p>
					<p><a href='` + "/healthz" + `'>healthz</a></p>
				</body>
			  </html>`
	fmt.Fprint(w, body)
}

func (e *Exporter) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func (e *Exporter) GetConfig() Config {
	return e.config
}

// BambuLabsX1C represents the structure of the MQTT message from BambuLabs printer
type BambuLabsX1C struct {
	Print struct {
		Ams struct {
			Ams []struct {
				Humidity string `json:"humidity"`
				ID       string `json:"id"`
				Temp     string `json:"temp"`
				Tray     []struct {
					BedTemp       string `json:"bed_temp"`
					BedTempType   string `json:"bed_temp_type"`
					DryingTemp    string `json:"drying_temp"`
					DryingTime    string `json:"drying_time"`
					ID            string `json:"id"`
					NozzleTempMax string `json:"nozzle_temp_max"`
					NozzleTempMin string `json:"nozzle_temp_min"`
					Remain        int    `json:"remain"`
					TagUID        string `json:"tag_uid"`
					TrayColor     string `json:"tray_color"`
					TrayDiameter  string `json:"tray_diameter"`
					TrayIDName    string `json:"tray_id_name"`
					TrayInfoIdx   string `json:"tray_info_idx"`
					TraySubBrands string `json:"tray_sub_brands"`
					TrayType      string `json:"tray_type"`
					TrayUUID      string `json:"tray_uuid"`
					TrayWeight    string `json:"tray_weight"`
					XcamInfo      string `json:"xcam_info"`
				} `json:"tray"`
			} `json:"ams"`
			AmsExistBits     string `json:"ams_exist_bits"`
			InsertFlag       bool   `json:"insert_flag"`
			PowerOnFlag      bool   `json:"power_on_flag"`
			TrayExistBits    string `json:"tray_exist_bits"`
			TrayIsBblBits    string `json:"tray_is_bbl_bits"`
			TrayNow          string `json:"tray_now"`
			TrayReadDoneBits string `json:"tray_read_done_bits"`
			TrayReadingBits  string `json:"tray_reading_bits"`
			TrayTar          string `json:"tray_tar"`
			Version          int    `json:"version"`
		} `json:"ams"`
		AmsRfidStatus           int     `json:"ams_rfid_status"`
		AmsStatus               int     `json:"ams_status"`
		BedTargetTemper         float64 `json:"bed_target_temper"`
		BedTemper               float64 `json:"bed_temper"`
		BigFan1Speed            string  `json:"big_fan1_speed"`
		BigFan2Speed            string  `json:"big_fan2_speed"`
		ChamberTemper           float64 `json:"chamber_temper"`
		Command                 string  `json:"command"`
		CoolingFanSpeed         string  `json:"cooling_fan_speed"`
		FailReason              string  `json:"fail_reason"`
		FanGear                 int     `json:"fan_gear"`
		ForceUpgrade            bool    `json:"force_upgrade"`
		GcodeFile               string  `json:"gcode_file"`
		GcodeFilePreparePercent string  `json:"gcode_file_prepare_percent"`
		GcodeStartTime          string  `json:"gcode_start_time"`
		GcodeState              string  `json:"gcode_state"`
		HeatbreakFanSpeed       string  `json:"heatbreak_fan_speed"`
		Hms                     []any   `json:"hms"`
		HomeFlag                int     `json:"home_flag"`
		HwSwitchState           int     `json:"hw_switch_state"`
		Ipcam                   struct {
			IpcamDev    string `json:"ipcam_dev"`
			IpcamRecord string `json:"ipcam_record"`
			Resolution  string `json:"resolution"`
			Timelapse   string `json:"timelapse"`
		} `json:"ipcam"`
		LayerNum     int    `json:"layer_num"`
		Lifecycle    string `json:"lifecycle"`
		LightsReport []struct {
			Mode string `json:"mode"`
			Node string `json:"node"`
		} `json:"lights_report"`
		Maintain            int     `json:"maintain"`
		McPercent           int     `json:"mc_percent"`
		McPrintErrorCode    string  `json:"mc_print_error_code"`
		McPrintStage        string  `json:"mc_print_stage"`
		McPrintSubStage     int     `json:"mc_print_sub_stage"`
		McRemainingTime     int     `json:"mc_remaining_time"`
		MessProductionState string  `json:"mess_production_state"`
		NozzleTargetTemper  float64 `json:"nozzle_target_temper"`
		NozzleTemper        float64 `json:"nozzle_temper"`
		Online              struct {
			Ahb  bool `json:"ahb"`
			Rfid bool `json:"rfid"`
		} `json:"online"`
		PrintError       int    `json:"print_error"`
		PrintGcodeAction int    `json:"print_gcode_action"`
		PrintRealAction  int    `json:"print_real_action"`
		PrintType        string `json:"print_type"`
		ProfileID        string `json:"profile_id"`
		ProjectID        string `json:"project_id"`
		Sdcard           bool   `json:"sdcard"`
		SequenceID       string `json:"sequence_id"`
		SpdLvl           int    `json:"spd_lvl"`
		SpdMag           int    `json:"spd_mag"`
		Stg              []int  `json:"stg"`
		StgCur           int    `json:"stg_cur"`
		SubtaskID        string `json:"subtask_id"`
		SubtaskName      string `json:"subtask_name"`
		TaskID           string `json:"task_id"`
		TotalLayerNum    int    `json:"total_layer_num"`
		UpgradeState     struct {
			AhbNewVersionNumber string `json:"ahb_new_version_number"`
			AmsNewVersionNumber string `json:"ams_new_version_number"`
			ConsistencyRequest  bool   `json:"consistency_request"`
			DisState            int    `json:"dis_state"`
			ErrCode             int    `json:"err_code"`
			ForceUpgrade        bool   `json:"force_upgrade"`
			Message             string `json:"message"`
			Module              string `json:"module"`
			NewVersionState     int    `json:"new_version_state"`
			OtaNewVersionNumber string `json:"ota_new_version_number"`
			Progress            string `json:"progress"`
			SequenceID          int    `json:"sequence_id"`
			Status              string `json:"status"`
		} `json:"upgrade_state"`
		Upload struct {
			FileSize      int    `json:"file_size"`
			FinishSize    int    `json:"finish_size"`
			Message       string `json:"message"`
			OssURL        string `json:"oss_url"`
			Progress      int    `json:"progress"`
			SequenceID    string `json:"sequence_id"`
			Speed         int    `json:"speed"`
			Status        string `json:"status"`
			TaskID        string `json:"task_id"`
			TimeRemaining int    `json:"time_remaining"`
			TroubleID     string `json:"trouble_id"`
		} `json:"upload"`
		WifiSignal string `json:"wifi_signal"`
		Xcam       struct {
			AllowSkipParts           bool   `json:"allow_skip_parts"`
			BuildplateMarkerDetector bool   `json:"buildplate_marker_detector"`
			FirstLayerInspector      bool   `json:"first_layer_inspector"`
			HaltPrintSensitivity     string `json:"halt_print_sensitivity"`
			PrintHalt                bool   `json:"print_halt"`
			PrintingMonitor          bool   `json:"printing_monitor"`
			SpaghettiDetector        bool   `json:"spaghetti_detector"`
		} `json:"xcam"`
		XcamStatus string `json:"xcam_status"`
	} `json:"print"`
}