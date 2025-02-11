```
docker pull fufelx/url-shortener:latest
```
Запуск с использованием in-memory хранилища:
```
docker run -d -p 3030:3030 -e STORAGE_TYPE=in-memory fufelx/url-shortener:latest
```
Запуск с использованием PostgreSQL:
```
docker run -d -p 3030:3030 -e STORAGE_TYPE=pgsql fufelx/url-shortener:latest
```

Для получения сокращенной ссылки
POST http://localhost:3030/api/addurl
```json
{
  "url": "https://www.google.com"
}
```
Ответ
```json
{
  "shorturl": "http://localhost:3030/hHh33JwTCL"
}
```

Для получения оригинальной ссылки
GET http://localhost:3030/api/geturl?shorturl=hHh33JwTCL

```json
{
  "url": "https://www.google.com"
}
```

Для редиректа на оригинальную ссылку
http://localhost:3030/hHh33JwTCL



