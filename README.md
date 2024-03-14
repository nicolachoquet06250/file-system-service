# file-system-service

Création d'un service en Go pour pouvoir m'y connecter via mon portfolio-apple et interagire avec le système de fichier via l'IHM.


## API Reference

#### Get directory list content items

```http
  GET /file-system/${path...}
```

| Parameter | Type     | Description                | Default value |
| :-------- | :------- | :------------------------- | :------------ |
| `path`    | `string` | **Optional**. The path of the directory you would like open | / |

### Response 200

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
    ...
]
```

### Response 400

```json
{
    "code": 400,
    "message": "an error message"
}
```

#### Get file content

```http
  GET /file/${path...}
```

| Parameter | Type     | Description                       | Default value |
| :-------- | :------- | :-------------------------------- | :------------ |
| `path`    | `string` | **Optional**. The path of the directory you would like open | / |

#### Create file

```http
  POST /file
  Content-Type: application/json
  Accept: application/json

  {
    "type": "file"
    "path": "/path/to/create",
    "name": "file-to-create",
    "extension": "your-extension"
  }
```

### Response 200

```json
{
    "path": "/path/to/create",
    "name": "file-to-create",
    "extension": "your-extension"
}
```

### Response 400

```json
{
    "code": 400,
    "message": "an error message"
}
```

| Parameter | Type     | Description                       | Default value |
| :-------- | :------- | :-------------------------------- | :------------ |
| `path`    | `string` | **Optional**. The path of the directory you would like open | / |

