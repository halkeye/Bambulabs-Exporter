
# BAMBULABS EXPORTER

`bambulabs-exporter` is a Prometheus exporter for bambulabs (only tested with x1c)

## Usage

`bambulabs-exporter` can be run as a Docker container (preferably) or built from source and run as a static binary.

### docker-compose

```yaml
version: "3"
services:
  plausible-exporter:
    image: ghcr.io/halkeye/bambulabs-exporter:latest
    environment:
      - BAMBULABS_IP='192.168.1.2'
      - BAMBULABS_TOPIC='device/<serialnumber>/report'
      - BAMBULABS_PASSWORD='<password<password>
      - BAMBULABS_USERNAME='bblp'
    ports:
      - 9101:9101
```

### Binary

To run the exporter as binary, clone this repo, build the Go binary and run it:

```sh
git clone https://github.com/halkeye/bambulabs-exporter
go get ./...
go build -o ./bambulabs-exporter .
./bambulabs-exporter
```

### Prometheus Metrics Available
- `*annotates recent changes or additions`

[Sample Metrics Here](sample.md)
| Metric   | Description | Examples |
| ------------- | ------------- |  ------------- |
| ams_humidity  | Humdity of the Enclosure, includes the AMS Number 0-many  | |
| ams_temp  | *Temperature of the AMS, includes the AMS Number 0-many | |
| ams_tray_color | *Filament color in the AMS, includes the AMS Number 0-many & Tray Numbers 0-4 | |
| ams_tray_type | *Filament type in the AMS, includes the AMS Number 0-many & Tray Numbers 0-4 | |
| big_fan1_speed | Big1 Fan Speed  | |
| big_fan2_speed | Big2 Fan Speed  | |
| chamber_temper | Temperature of the Bambu Enclosure  | |
| cooling_fan_speed | Print Head Cooling Fan Speed  | |
| fail_reason | Failure Print Reason Code  | |
| fan_gear | Fan Gear   | |
| layer_number | GCode Layer Number of the Print  | |
| mc_percent | Print Progress in Percentage  | |
| mc_print_error_code | Print Progress Error Code | |
| mc_print_stage | Print Progress Stage | |
| mc_print_sub_stage | Print Progress Sub Stage | |
| mc_remaining_time | Print Progress Remaining Time in minutes  | |
| nozzle_target_temper |Nozzle Target Temperature Metric | |
| nozzle_temper | Nozzle Temperature Metric | |
| print_error | Print Error reported by the Control board | |
| wifi_signal | Wifi Signal Strength in dBm | |

### Grafana

You can use the exported metrics just like you'd use any other metric scraped by Prometheus.

The `examples` folder contains a small [demo dashboard](./examples/grafana-dashboard.json). You can use this as a starting point for integrating the metrics into your own dashboards.

### Credit

Original implementation by [aetrius](https://github.com/aetrius/bambulabs-exporter)
