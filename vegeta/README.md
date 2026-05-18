# Vegeta load test

This folder contains a small local load-test flow for the gateway endpoint:

```text
POST /api/send/{templateId}
```

The send endpoint needs a valid JWT, an existing template, and contacts linked to that template.
Do not edit `targets.txt` by hand. Generate it before each run.

## Prepare test data only

```powershell
powershell -ExecutionPolicy Bypass -File .\vegeta\prepare.ps1 -Contacts 20
```

The script:

- registers or logs in a test user;
- creates a template;
- creates contacts;
- links contacts to the template;
- writes a ready Vegeta target to `targets.txt`.

## Run the load test

```powershell
powershell -ExecutionPolicy Bypass -File .\vegeta\run.ps1 -Rate 10/s -Duration 30s -Contacts 20
```

The script writes:

- `targets.txt` - generated Vegeta target with a real token and template id;
- `normal.bin` - Vegeta binary result;
- `normal.html` - Vegeta HTML plot.

## Run Vegeta manually

```powershell
powershell -ExecutionPolicy Bypass -File .\vegeta\prepare.ps1 -Contacts 20
vegeta attack -targets targets.txt -rate 10/s -duration 30s -output normal.bin
vegeta report normal.bin
vegeta plot normal.bin > normal.html
```

Install Vegeta if needed:

```powershell
go install github.com/tsenart/vegeta@latest
```
