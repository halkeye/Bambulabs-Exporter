package exporter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestNewExporter(t *testing.T) {
	// Set up test environment variables
	envVars := map[string]string{
		"BAMBULABS_IP":       "192.168.1.100",
		"BAMBULABS_USERNAME": "testuser",
		"BAMBULABS_PASSWORD": "testpass",
		"BAMBULABS_TOPIC":    "device/test123/report",
		"BAMBULABS_DEBUG":    "true",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	exporter := NewExporter()
	if exporter == nil {
		t.Fatal("Expected exporter, got nil")
	}

	config := exporter.GetConfig()
	if config.IP != "192.168.1.100" {
		t.Errorf("Expected IP '192.168.1.100', got '%s'", config.IP)
	}
	if config.Username != "testuser" {
		t.Errorf("Expected Username 'testuser', got '%s'", config.Username)
	}
	if config.Password != "testpass" {
		t.Errorf("Expected Password 'testpass', got '%s'", config.Password)
	}
	if config.Topic != "device/test123/report" {
		t.Errorf("Expected Topic 'device/test123/report', got '%s'", config.Topic)
	}
	if !config.Debug {
		t.Errorf("Expected Debug true, got %v", config.Debug)
	}
}

func TestExporterHTTPEndpoints(t *testing.T) {
	// Reset the default registry to avoid duplicate metric registration
	oldRegistry := prometheus.DefaultRegisterer
	defer func() {
		prometheus.DefaultRegisterer = oldRegistry
	}()

	// Set up test environment variables
	envVars := map[string]string{
		"BAMBULABS_IP":       "192.168.1.100",
		"BAMBULABS_USERNAME": "testuser",
		"BAMBULABS_PASSWORD": "testpass",
		"BAMBULABS_TOPIC":    "device/test123/report",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	// Use a new registry for this test
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	exporter := NewExporter()

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "home endpoint",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "BambuLabs Exporter",
		},
		{
			name:           "healthz endpoint",
			path:           "/healthz",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "metrics endpoint",
			path:           "/metrics",
			expectedStatus: http.StatusOK,
			expectedBody:   "# HELP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			
			switch tt.path {
			case "/":
				exporter.home(rr, req)
			case "/healthz":
				exporter.healthz(rr, req)
			case "/metrics":
				promhttp.Handler().ServeHTTP(rr, req)
			}

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if !strings.Contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("Expected body to contain %s, got %s", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestExporterMessageHandler(t *testing.T) {
	// Reset the default registry to avoid duplicate metric registration
	oldRegistry := prometheus.DefaultRegisterer
	defer func() {
		prometheus.DefaultRegisterer = oldRegistry
	}()

	// Set up test environment variables
	envVars := map[string]string{
		"BAMBULABS_IP":       "192.168.1.100",
		"BAMBULABS_USERNAME": "testuser",
		"BAMBULABS_PASSWORD": "testpass",
		"BAMBULABS_TOPIC":    "device/test123/report",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	// Use a new registry for this test
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	exporter := NewExporter()

	// Test valid push_status message
	validJSON := `{
		"print": {
			"command": "push_status",
			"layer_num": 10,
			"print_error": 0,
			"wifi_signal": "-50dBm",
			"big_fan1_speed": "150",
			"big_fan2_speed": "200",
			"chamber_temper": 30.0,
			"cooling_fan_speed": "100",
			"fail_reason": "0",
			"fan_gear": 3,
			"mc_percent": 50,
			"mc_print_error_code": "0",
			"mc_print_stage": "2",
			"mc_print_sub_stage": 2,
			"mc_remaining_time": 60,
			"nozzle_temper": 220.0,
			"nozzle_target_temper": 230.0,
			"ams": {
				"ams": [
					{
						"id": "0",
						"humidity": "50.0",
						"temp": "25.0",
						"tray": [
							{
								"id": "0",
								"tray_type": "ABS",
								"tray_color": "Blue"
							}
						]
					}
				]
			}
		}
	}`

	// Create a mock MQTT message
	mockMsg := &mockMessage{payload: []byte(validJSON)}
	mockClient := &mockClient{}

	// Call the message handler
	exporter.messagePubHandler(mockClient, mockMsg)

	// Verify metrics were set correctly
	if testutil.ToFloat64(exporter.layerNumberMetric) != 10.0 {
		t.Errorf("Expected layer number 10.0, got %f", testutil.ToFloat64(exporter.layerNumberMetric))
	}
	if testutil.ToFloat64(exporter.printErrorMetric) != 0.0 {
		t.Errorf("Expected print error 0.0, got %f", testutil.ToFloat64(exporter.printErrorMetric))
	}
	if testutil.ToFloat64(exporter.wifiSignalMetric) != -50.0 {
		t.Errorf("Expected wifi signal -50.0, got %f", testutil.ToFloat64(exporter.wifiSignalMetric))
	}
	if testutil.ToFloat64(exporter.chamberTemperMetric) != 30.0 {
		t.Errorf("Expected chamber temperature 30.0, got %f", testutil.ToFloat64(exporter.chamberTemperMetric))
	}
	if testutil.ToFloat64(exporter.fanGearMetric) != 3.0 {
		t.Errorf("Expected fan gear 3.0, got %f", testutil.ToFloat64(exporter.fanGearMetric))
	}
	if testutil.ToFloat64(exporter.mcPercentMetric) != 50.0 {
		t.Errorf("Expected MC percent 50.0, got %f", testutil.ToFloat64(exporter.mcPercentMetric))
	}
	if testutil.ToFloat64(exporter.nozzleTemperMetric) != 220.0 {
		t.Errorf("Expected nozzle temperature 220.0, got %f", testutil.ToFloat64(exporter.nozzleTemperMetric))
	}
	if testutil.ToFloat64(exporter.nozzleTargetTemperMetric) != 230.0 {
		t.Errorf("Expected nozzle target temperature 230.0, got %f", testutil.ToFloat64(exporter.nozzleTargetTemperMetric))
	}

	// Verify AMS metrics
	amsHumidityValue := testutil.ToFloat64(exporter.amsHumidityMetric.With(prometheus.Labels{"ams_number": "0"}))
	if amsHumidityValue != 50.0 {
		t.Errorf("Expected AMS humidity 50.0, got %f", amsHumidityValue)
	}

	amsTempValue := testutil.ToFloat64(exporter.amsTempMetric.With(prometheus.Labels{"ams_number": "0"}))
	if amsTempValue != 25.0 {
		t.Errorf("Expected AMS temperature 25.0, got %f", amsTempValue)
	}

	// Verify tray metrics
	trayColorValue := testutil.ToFloat64(exporter.amsColorMetric.With(prometheus.Labels{
		"ams_number":  "0",
		"tray_number": "0",
		"tray_color":  "Blue",
	}))
	if trayColorValue != 1.0 {
		t.Errorf("Expected tray color metric 1.0, got %f", trayColorValue)
	}

	trayTypeValue := testutil.ToFloat64(exporter.amsTypeMetric.With(prometheus.Labels{
		"ams_number":  "0",
		"tray_number": "0",
		"tray_type":   "ABS",
	}))
	if trayTypeValue != 1.0 {
		t.Errorf("Expected tray type metric 1.0, got %f", trayTypeValue)
	}
}

func TestExporterMessageHandlerInvalidJSON(t *testing.T) {
	// Reset the default registry to avoid duplicate metric registration
	oldRegistry := prometheus.DefaultRegisterer
	defer func() {
		prometheus.DefaultRegisterer = oldRegistry
	}()

	// Set up test environment variables
	envVars := map[string]string{
		"BAMBULABS_IP":       "192.168.1.100",
		"BAMBULABS_USERNAME": "testuser",
		"BAMBULABS_PASSWORD": "testpass",
		"BAMBULABS_TOPIC":    "device/test123/report",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	// Use a new registry for this test
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	exporter := NewExporter()

	invalidJSON := `{"invalid": json}`

	mockMsg := &mockMessage{payload: []byte(invalidJSON)}
	mockClient := &mockClient{}

	// This should not panic, just log an error
	exporter.messagePubHandler(mockClient, mockMsg)
}

func TestExporterMessageHandlerWrongCommand(t *testing.T) {
	// Reset the default registry to avoid duplicate metric registration
	oldRegistry := prometheus.DefaultRegisterer
	defer func() {
		prometheus.DefaultRegisterer = oldRegistry
	}()

	// Set up test environment variables
	envVars := map[string]string{
		"BAMBULABS_IP":       "192.168.1.100",
		"BAMBULABS_USERNAME": "testuser",
		"BAMBULABS_PASSWORD": "testpass",
		"BAMBULABS_TOPIC":    "device/test123/report",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	// Use a new registry for this test
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	exporter := NewExporter()

	wrongCommandJSON := `{
		"print": {
			"command": "wrong_command",
			"layer_num": 10
		}
	}`

	mockMsg := &mockMessage{payload: []byte(wrongCommandJSON)}
	mockClient := &mockClient{}

	// This should not panic, just ignore the message
	exporter.messagePubHandler(mockClient, mockMsg)
}

func TestBambuLabsX1CStruct(t *testing.T) {
	sampleJSON := `{
		"print": {
			"command": "push_status",
			"layer_num": 5,
			"print_error": 0,
			"wifi_signal": "-45dBm",
			"big_fan1_speed": "100",
			"big_fan2_speed": "200",
			"chamber_temper": 25.5,
			"cooling_fan_speed": "150",
			"fail_reason": "0",
			"fan_gear": 2,
			"mc_percent": 25,
			"mc_print_error_code": "0",
			"mc_print_stage": "1",
			"mc_print_sub_stage": 1,
			"mc_remaining_time": 120,
			"nozzle_temper": 200.0,
			"nozzle_target_temper": 210.0,
			"ams": {
				"ams": [
					{
						"id": "0",
						"humidity": "45.2",
						"temp": "23.1",
						"tray": [
							{
								"id": "0",
								"tray_type": "PLA",
								"tray_color": "Red"
							}
						]
					}
				]
			}
		}
	}`

	var data BambuLabsX1C
	err := json.Unmarshal([]byte(sampleJSON), &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Test basic fields
	if data.Print.Command != "push_status" {
		t.Errorf("Expected command 'push_status', got '%s'", data.Print.Command)
	}
	if data.Print.LayerNum != 5 {
		t.Errorf("Expected layer_num 5, got %d", data.Print.LayerNum)
	}
	if data.Print.PrintError != 0 {
		t.Errorf("Expected print_error 0, got %d", data.Print.PrintError)
	}
	if data.Print.WifiSignal != "-45dBm" {
		t.Errorf("Expected wifi_signal '-45dBm', got '%s'", data.Print.WifiSignal)
	}
	if data.Print.ChamberTemper != 25.5 {
		t.Errorf("Expected chamber_temper 25.5, got %f", data.Print.ChamberTemper)
	}
	if data.Print.FanGear != 2 {
		t.Errorf("Expected fan_gear 2, got %d", data.Print.FanGear)
	}
	if data.Print.McPercent != 25 {
		t.Errorf("Expected mc_percent 25, got %d", data.Print.McPercent)
	}
	if data.Print.NozzleTemper != 200.0 {
		t.Errorf("Expected nozzle_temper 200.0, got %f", data.Print.NozzleTemper)
	}
	if data.Print.NozzleTargetTemper != 210.0 {
		t.Errorf("Expected nozzle_target_temper 210.0, got %f", data.Print.NozzleTargetTemper)
	}

	// Test AMS data
	if len(data.Print.Ams.Ams) != 1 {
		t.Errorf("Expected 1 AMS, got %d", len(data.Print.Ams.Ams))
	}
	
	ams := data.Print.Ams.Ams[0]
	if ams.ID != "0" {
		t.Errorf("Expected AMS ID '0', got '%s'", ams.ID)
	}
	if ams.Humidity != "45.2" {
		t.Errorf("Expected humidity '45.2', got '%s'", ams.Humidity)
	}
	if ams.Temp != "23.1" {
		t.Errorf("Expected temp '23.1', got '%s'", ams.Temp)
	}
	
	if len(ams.Tray) != 1 {
		t.Errorf("Expected 1 tray, got %d", len(ams.Tray))
	}
	
	tray := ams.Tray[0]
	if tray.ID != "0" {
		t.Errorf("Expected tray ID '0', got '%s'", tray.ID)
	}
	if tray.TrayType != "PLA" {
		t.Errorf("Expected tray_type 'PLA', got '%s'", tray.TrayType)
	}
	if tray.TrayColor != "Red" {
		t.Errorf("Expected tray_color 'Red', got '%s'", tray.TrayColor)
	}
}

// Mock implementations for testing
type mockMessage struct {
	payload []byte
}

func (m *mockMessage) Duplicate() bool {
	return false
}

func (m *mockMessage) Qos() byte {
	return 0
}

func (m *mockMessage) Retained() bool {
	return false
}

func (m *mockMessage) Topic() string {
	return "test/topic"
}

func (m *mockMessage) MessageID() uint16 {
	return 1
}

func (m *mockMessage) Payload() []byte {
	return m.payload
}

func (m *mockMessage) Ack() {
}

type mockClient struct{}

func (m *mockClient) IsConnected() bool {
	return true
}

func (m *mockClient) IsConnectionOpen() bool {
	return true
}

func (m *mockClient) Connect() mqtt.Token {
	return &mockToken{}
}

func (m *mockClient) Disconnect(quiesce uint) {
}

func (m *mockClient) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	return &mockToken{}
}

func (m *mockClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	return &mockToken{}
}

func (m *mockClient) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	return &mockToken{}
}

func (m *mockClient) Unsubscribe(topics ...string) mqtt.Token {
	return &mockToken{}
}

func (m *mockClient) AddRoute(topic string, callback mqtt.MessageHandler) {
}

func (m *mockClient) RemoveRoute(topic string) {
}

func (m *mockClient) OptionsReader() mqtt.ClientOptionsReader {
	return mqtt.ClientOptionsReader{}
}

type mockToken struct{}

func (m *mockToken) Wait() bool {
	return true
}

func (m *mockToken) WaitTimeout(timeout time.Duration) bool {
	return true
}

func (m *mockToken) Done() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}

func (m *mockToken) Error() error {
	return nil
}