# Luzifer / doc-render

`doc-render` is a small webservice serving as a UI to the [`tex-api`](https://github.com/luzifer/tex-api). While the `tex-api` does the heavy-lifting of rendering the LaTeX document into a PDF this provides an UI for the user to interact with LaTeX templates: templates are defined as a LaTeX file in Go-Templating format, a JSON-Schema defining the fields of the template for the frontend to generate a form and optionally some additional files like logo, wallpapers or font files.

![](docs/screenshot.png)

## Template structure

```
templates
└── letter
    ├── background-1.pdf
    ├── background-2.pdf
    ├── main.tex.tpl
    ├── schema.json
    └── signature.png
```

- Folder name must be a slug (i.e. "letter") and is used to identify this template
- `main.tex.tpl` contains the TeX file to render. It contains Go templating and has access to `.Values` (defined in the schema) and `.Recipients` containing addresses passed in through the generator frontend.
- `schema.json` contains a JSON-Schema definition of the template and its `.Values`
  - The `description` is used as a display name
  - `properties` must be flat (no `"type": "object"`) and describe the fields. For example the property `"subject": {"description": "Betreff", "type": "string"}` will yield a text-input field named "Betreff" and its value will be available as `.Values.subject` to the template.
  - `required` properties must have non-empty values
  - Properties having a `default` will display that default in the frontend.
- Additional files can be provided and will be available during rendering

## Server-side storage of pre-filled values

When enabled during deployment `doc-render` allows to store the values filled inside the templates on the server and generate a link to retrieve those values again. The following backends are available:

- `k8s` - Store the values as ConfigMap objects inside the Kubernetes cluster
  - Set `PERSIST_NAMESPACE` to the namespace the ConfigMap objects should be created in
- `mem` - Store the values in an in-memory map (restarting the server will wipe the storage)
- `redis` - Store the values in a Redis instance
  - Set `PERSIST_REDIS` to a redis connection URL ([]`redis://<user>:<password>@<host>:<port>/<db_number>`](https://pkg.go.dev/github.com/redis/go-redis/v9@v9.7.3#ParseURL))
  - Optionally set `PERSIST_REDIS_PREFIX` to a prefix to prepend the object keys
