```
# HELP ams_humidity humidity of the ams
# TYPE ams_humidity gauge
ams_humidity{ams_number="0"} 5
# HELP ams_temp temperature of the ams
# TYPE ams_temp gauge
ams_temp{ams_number="0"} -41.5
# HELP ams_tray_color color of material in ams tray
# TYPE ams_tray_color gauge
ams_tray_color{ams_number="0",tray_color="",tray_number="1"} 1
ams_tray_color{ams_number="0",tray_color="161616FF",tray_number="3"} 1
ams_tray_color{ams_number="0",tray_color="7C4B00FF",tray_number="0"} 1
ams_tray_color{ams_number="0",tray_color="F98C36FF",tray_number="2"} 1
# HELP ams_tray_type type of material in ams tray
# TYPE ams_tray_type gauge
ams_tray_type{ams_number="0",tray_number="0",tray_type="PLA"} 1
ams_tray_type{ams_number="0",tray_number="1",tray_type=""} 1
ams_tray_type{ams_number="0",tray_number="2",tray_type="PLA"} 1
ams_tray_type{ams_number="0",tray_number="3",tray_type="PLA"} 1
# HELP big_fan1_speed Big Fan 1 Speed
# TYPE big_fan1_speed gauge
big_fan1_speed 11
# HELP big_fan2_speed Big Fan 2 Speed
# TYPE big_fan2_speed gauge
big_fan2_speed 0
# HELP chamber_temper Chamber Temperature of Printer
# TYPE chamber_temper gauge
chamber_temper 34
# HELP cooling_fan_speed Cooling Fan Speed
# TYPE cooling_fan_speed gauge
cooling_fan_speed 15
# HELP fail_reason Print Failure Reason
# TYPE fail_reason gauge
fail_reason 0
# HELP fan_gear Fan Gear
# TYPE fan_gear gauge
fan_gear 45823
# HELP layer_number layer number of the print head in gcode
# TYPE layer_number gauge
layer_number 14
# HELP mc_percent Percentage of Progress of print
# TYPE mc_percent gauge
mc_percent 21
# HELP mc_print_error_code Print Progress Error Code
# TYPE mc_print_error_code gauge
mc_print_error_code 0
# HELP mc_print_stage Print Progress Stage
# TYPE mc_print_stage gauge
mc_print_stage 0
# HELP mc_print_sub_stage Print Progress Sub Stage
# TYPE mc_print_sub_stage gauge
mc_print_sub_stage 0
# HELP mc_remaining_time Print Progress Remaining Time in minutes
# TYPE mc_remaining_time gauge
mc_remaining_time 360
# HELP nozzle_target_temper Nozzle Target Temperature Metric
# TYPE nozzle_target_temper gauge
nozzle_target_temper 220
# HELP nozzle_temper Nozzle Temperature Metric
# TYPE nozzle_temper gauge
nozzle_temper 220
# HELP print_error Print error int
# TYPE print_error gauge
print_error 0
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 0
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
# HELP wifi_signal Wifi signal in dBm
# TYPE wifi_signal gauge
wifi_signal -54
```
