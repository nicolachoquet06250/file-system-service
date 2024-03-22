# file-system-service

Création d'un service en Go pour pouvoir m'y connecter via mon portfolio-apple et interagire avec le système de fichier via l'IHM.

## Quick start
### Installation

#### Linux & MacOSX
##### Get Last Swagger
```shell
wget https://github.com/nicolachoquet06250/file-system-service/releases/download/$(curl https://api.github.com/repos/nicolachoquet06250/file-system-service/releases | jq .[0].name | sed 's/.\(.*\)/\1/' | sed 's/\(.*\)./\1/')/file-system-service.swagger.yml
```

#### Linux
##### Get Last binary
```shell
last_version="$(curl https://api.github.com/repos/nicolachoquet06250/file-system-service/releases | jq .[0].name | sed 's/.\(.*\)/\1/' | sed 's/\(.*\)./\1/')" && \
  wget -c https://github.com/nicolachoquet06250/file-system-service/releases/download/${last_version}/file-system-service-linux-${last_version}-linux-amd64.tar.gz && \
  tar -xf file-system-service-linux-${last_version}-linux-amd64.tar.gz && \
  rm file-system-service-linux-${last_version}-linux-amd64.tar.gz
```

#### MacOSX amd64
##### Get Last binary
```shell
last_version="$(curl https://api.github.com/repos/nicolachoquet06250/file-system-service/releases | jq .[0].name | sed 's/.\(.*\)/\1/' | sed 's/\(.*\)./\1/')" && \
  wget -c https://github.com/nicolachoquet06250/file-system-service/releases/download/${last_version}/file-system-service-darwin-${last_version}-darwin-amd64.tar.gz && \
  tar -xf file-system-service-darwin-${last_version}-darwin-amd64.tar.gz && \
  rm file-system-service-darwin-${last_version}-darwin-amd64.tar.gz
```

#### MacOSX arm64
##### Get Last binary
```shell
last_version="$(curl https://api.github.com/repos/nicolachoquet06250/file-system-service/releases | jq .[0].name | sed 's/.\(.*\)/\1/' | sed 's/\(.*\)./\1/')" && \
  wget -c https://github.com/nicolachoquet06250/file-system-service/releases/download/${last_version}/file-system-service-darwin-${last_version}-darwin-arm64.tar.gz && \
  tar -xf file-system-service-darwin-${last_version}-darwin-arm64.tar.gz && \
  rm file-system-service-darwin-${last_version}-darwin-arm64.tar.gz
```

#### Windows
https://github.com/nicolachoquet06250/file-system-service/releases/latest

##### Get Last Swagger
Click on `file-system-service.swagger.yml`

##### Get Last binary
Click on `file-system-service-windows-{version}-windows-amd64.zip`

#### Generate Signature token
```shell
file-system-service --generate-signature
```

## Swagger
- [Fichiers de définitions json](./swagger/swagger.json)
- [Fichiers de définitions yaml](./swagger/swagger.yaml)
- [Swagger UI accessible ici](http://localhost:3000/swagger)

## API Reference

#### Check Validity

```http request
GET /check-validity
```

##### Response 200

```json
{
  "isValid": true
}
```

#### Authentification

##### Get first token

```http request
POST /auth/get-token
Accept: application/json
Content-Type: application/json
Signature-Token: {generated-signature}
```
###### Response 200

```json
{
  "access_token": "string",
  "refresh_token": "string",
  "expires_in": "int",
  "created_at": "int"
}
```
###### Response 400

```json
{
  "code": 400,
  "message": "string"
}
```
###### Response 500

```json
{
  "code": 500,
  "message": "string"
}
```

##### Refresh token

```http request
PUT /auth/get-token
Accept: application/json
Content-Type: application/json
Signature-Token: {generated-signature}
Refresh-Token: {getted-refresh-token}
```
###### Response 200

```json
{
  "access_token": "string",
  "refresh_token": "string",
  "expires_in": "int",
  "created_at": "int"
}
```
###### Response 400

```json
{
  "code": 400,
  "message": "string"
}
```
###### Response 500

```json
{
  "code": 500,
  "message": "string"
}
```

#### Get directory list content items

```http request
GET /file-system/${path...}
```

| Parameter | Type     | Description                | Default value |
| :-------- | :------- | :------------------------- | :------------ |
| `path`    | `string` | **Optional**. The path of the directory you would like open | / |

##### Response 200

```json
[
    {
        "type": "directory",
        "path": "/home",
        "name": "nchoquet"
    },
    {
        "type": "file",
        "path": "/home",
        "name": ".bashrc"
    },
    "..."
]
```

##### Response 400

```json
{
    "code": 400,
    "message": "an error message"
}
```

##### Response 404

```json
{
  "code": 404,
  "message": "an error message"
}
```

##### Response 403

```json
{
  "code": 403,
  "message": "an error message"
}
```

#### Create a directory

```http request
POST /directory
Accept: application/json
Content-Type: application/json

{
    "type": "directory",
    "path": "string",
    "name": "string"
}
```
##### Response 200

```json
{
  "path": "string",
  "name": "string"
}
```

##### Response 400

```json
{
  "code": 400,
  "message": "an error message"
}
```

##### Response 404

```json
{
  "code": 404,
  "message": "an error message"
}
```

##### Response 403

```json
{
  "code": 403,
  "message": "an error message"
}
```

#### Rename directory

```http request
PATCH /directory/{path...}
Accept: application/json
Content-Type: application/json

{
  "type": "directory",
  "path": "string",
  "name": "string"
}
```

| Parameter | Type     | Description                       | Default value |
| :-------- | :------- | :-------------------------------- | :------------ |
| `path`    | `string` | **Optional**. The path of the directory you would like open | / |


##### Response 200

```json
{
  "path": "string",
  "name": "string"
}
```

##### Response 400

```json
{
  "code": 400,
  "message": "an error message"
}
```

##### Response 404

```json
{
  "code": 404,
  "message": "an error message"
}
```

##### Response 403

```json
{
  "code": 403,
  "message": "an error message"
}
```

#### Delete directory

```http request
DELETE /directory/{path...}
Accept: application/json
```

| Parameter | Type     | Description                       | Default value |
| :-------- | :------- | :-------------------------------- | :------------ |
| `path`    | `string` | **Optional**. The path of the directory you would like open | / |


##### Response 200

```json
{
  "status": "success"
}
```

##### Response 400

```json
{
  "code": 400,
  "message": "an error message"
}
```

##### Response 404

```json
{
  "code": 404,
  "message": "an error message"
}
```

##### Response 403

```json
{
  "code": 403,
  "message": "an error message"
}
```

#### Create file with content

```http request
POST http://localhost:3000/file
Accept: application/json
Content-Type: multipart/form-data; boundary=boundary

--boundary
Content-Disposition: form-data; name="file"
Content-Type: application/json

{
    "type": "file"
    "path": "/path/to/create",
    "name": "file-to-create",
    "extension": "your-extension"
}

--boundary
Content-Disposition: form-data; name="content"
Content-Type: text/plain

Ceci est un test
--boundary--
```

##### Response 200

```json
{
    "path": "/path/to/create",
    "name": "file-to-create",
    "extension": "your-extension"
}
```

##### Response 400

```json
{
    "code": 400,
    "message": "an error message"
}
```

| Parameter | Type     | Description                       | Default value |
| :-------- | :------- | :-------------------------------- | :------------ |
| `path`    | `string` | **Optional**. The path of the directory you would like open | / |

##### Response 403

```json
{
  "code": 403,
  "message": "an error message"
}
```

#### Create file without content

```http request
POST http://localhost:3000/file
Accept: application/json
Content-Type: application/json

{
    "type": "file"
    "path": "/path/to/create",
    "name": "file-to-create",
    "extension": "your-extension"
}
```

##### Response 200

```json
{
    "path": "/path/to/create",
    "name": "file-to-create",
    "extension": "your-extension"
}
```

##### Response 400

```json
{
    "code": 400,
    "message": "an error message"
}
```

| Parameter | Type     | Description                       | Default value |
| :-------- | :------- | :-------------------------------- | :------------ |
| `path`    | `string` | **Optional**. The path of the directory you would like open | / |

##### Response 403

```json
{
  "code": 403,
  "message": "an error message"
}
```

#### Get file content

```http request
GET /file/${path...}
Accept: application/json
```

| Parameter | Type     | Description                       | Default value |
| :-------- | :------- | :-------------------------------- | :------------ |
| `path`    | `string` | **Optional**. The path of the directory you would like open | / |


##### Response 200

###### ***application/json***
###### ***text/plain***
###### ***text/xml***
###### ***application/pdf***

##### Response 400

```json
{
  "code": 400,
  "message": "an error message"
}
```

##### Response 404

```json
{
  "code": 404,
  "message": "an error message"
}
```

##### Response 403

```json
{
  "code": 403,
  "message": "an error message"
}
```

#### Rename selected file

```http request
PATCH /file/{path...}
Accept: application/json
Content-Type: application/json

{
  "type": "file",
  "path": "/path/to/create",
  "name": "file-to-create",
  "extension": "your-extension"
}
```

| Parameter | Type     | Description                       | Default value |
| :-------- | :------- | :-------------------------------- | :------------ |
| `path`    | `string` | **Optional**. The path of the directory you would like open | / |


##### Response 200

```json
{
  "path": "/path/to/create",
  "name": "file-to-create",
  "extension": "your-extension"
}
```

##### Response 400

```json
{
  "code": 400,
  "message": "an error message"
}
```

##### Response 404

```json
{
  "code": 404,
  "message": "an error message"
}
```

##### Response 403

```json
{
  "code": 403,
  "message": "an error message"
}
```

#### Update selected file content

```http request
PUT /file/{path...}
Accept: application/json
Content-Type: text/plain

This is a text for update
the fichier
```

##### Response 200

```json
{
  "status": "success"
}
```

##### Response 400

```json
{
  "code": 400,
  "message": "an error message"
}
```

##### Response 404

```json
{
  "code": 404,
  "message": "an error message"
}
```

##### Response 403

```json
{
  "code": 403,
  "message": "an error message"
}
```

#### Delete selected file

```http request
DELETE /file/{path...}
Accept: application/json
```

##### Response 200

```json
{
  "status": "success"
}
```

##### Response 400

```json
{
  "code": 400,
  "message": "an error message"
}
```

##### Response 404

```json
{
  "code": 404,
  "message": "an error message"
}
```


##### Response 403

```json
{
  "code": 403,
  "message": "an error message"
}
```