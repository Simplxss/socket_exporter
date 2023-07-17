# socket_exporter
exporter for prometheus that transform tcp/udp protocol to metrics

# Usage
```yaml
listening:
  address: 0.0.0.0:8999               # Address to be listen on
  path:    /metrics                   # Path that prometheus access to
endpoint:
  - address:  192.168.107.245:2000    # Address of the PLC
    type:     tcp                     # Type of connection (tcp|udp)
    length:   158                     # Length of data to read from PLC (NO MORE THAN 1024)
    label:    target=phaseII          # OPTIONAL label for all metrics on this target
    protocol:                         # Protocol to use for this target
      - name:         emergency_stop  # name of metric as it appears to Prometheus
        help:         Emergency Stop  # help text as it appears to Prometheus
        label:        region=roof     # OPTIONAL label of metric as it appears to Prometheus
        datatype:     bool            # Datatype of value in PLC (bool|byte|real|word|dword|int|dint)
        trueValue:    1               # Value that pass to prometheus when the bool value is TRUE
        metricType:   gauge           # Type of metric as it appears to Prometheus (currently only gauge is supported)
        offset:       0               # offset of value in PLC
```
