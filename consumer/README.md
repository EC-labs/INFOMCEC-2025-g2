This service will produce and send to the AVG_calc_service

1. a grouped object containing:
- Experiment ID
- Timestamp of measurement
- Array of avg temps

ex: 
```json
{
    "experiment_id": "f55aee0e-3ee9-4de2-b14f-a61fcc4dc258",
    "timestamp": 1231232121,
    "started": true,
    "measurement_count": 3,
    "measuments": [
        46.7,
        43.6
        45.2
    ]
}
```

protocol used for now: gRPC both to avg_calc_service and postgres_service

Consumer service: forwards to postgres_service if not measurement
Consumer service: forwards to avg_calc_service if measurement