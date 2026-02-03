# Go HTTP Track

Run both services from `go-http/`.

## Run

**Terminal 1 (Provider — Echo API)**  
```bash
go run ./service-a
```

**Terminal 2 (Consumer — Client)**  
```bash
go run ./service-b
```

## Test

**Success**  
```bash
curl "http://127.0.0.1:8081/call-echo?msg=hello"
```  
Expect: `200` and combined JSON: `{"consumer":{"msg":"hello"},"provider":{"echo":"hello"}}`.

**Failure**  
Stop the provider (Ctrl+C in Terminal 1), then:  
```bash
curl "http://127.0.0.1:8081/call-echo?msg=hello"
```  
Expect: `503`, error JSON (`{"error":"provider unavailable","details":"..."}`), and the consumer logs an error line with `service=client endpoint=/call-echo status=503 latency_ms=... error="..."`.
