This service will calculate the average temp using the received messagaes 
Then it will produce the avg calculation and send it to the notification_service

1. receive array
2. calcute avg for the array
3. send avg for timestamp x experiment_id to notification_service

ex send definition:
```json
{
    "experiment_id": "f55aee0e-3ee9-4de2-b14f-a61fcc4dc258",
    "timestamp": 1231232121,
    "started": true,
    // "measurement_count": 3,
    "avg_measument": 45.6
}
```

protocol used for now: gRPC to notification_service