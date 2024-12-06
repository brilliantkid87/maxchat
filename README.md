
## Robot Management API

## Instalation
1. Clone repository
2. Jalankan `go mod tidy` untuk mengunduh dependensi

## Running
```bash
go run main.go
```

## Endpoint API

### Get All Robots
`GET /robots`
- Filter opsional:
  - `?model=car` atau `?model=humanoid`
  - `?tech=AI,car` (multi-tech filter)

### Get Robot by Code
`GET /robots/{code}`

### Create Robot
`POST /robots`
Body:
```json
{
    "code": "unique_code",
    "name": "Robot Name",
    "description": "...",
    "model": "car",
    "tech": ["AI", "car"],
    "status": "progress/active/inactive"
}
```

### Update Robot
`PUT /robots/{code}`

### Delete Robot
`DELETE /robots/{code}`

### Get References
`GET /references` - Mendapatkan daftar referensi yang valid untuk model, tech, dan status

## POST References
`POST /references/update`
Body:
```json
{
  "models": ["UpdatedModelA", "UpdatedModelB"],
  "techs": ["AI", "Robotics", "Automation"]
}
```
