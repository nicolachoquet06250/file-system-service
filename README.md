# file-system-service

Création d'un service en Go pour pouvoir m'y connecter via mon portfolio-apple et interagire avec le système de fichier via l'IHM.


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
