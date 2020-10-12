# Copying AWS Secrets

Usage
```
go run main.go --region <region_name> --original <original_name> --new <new_name>
```

where:
1. `region_name` -- name of the AWS region for the session
1. `<original_name>` -- friendly name of the existing secret
1. `<new_name>` -- friendly name of the new secret
