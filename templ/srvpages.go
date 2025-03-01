package templ

import (
	"html/template"
)

const p404str = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>404</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/semantic-ui@2.4.2/dist/semantic.min.css">
</head>
<body>
    <center><h1>404 Page Not found</h1></center>
</body>
</html>
`
const pInternalstr = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>InternalServerError</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/semantic-ui@2.4.2/dist/semantic.min.css">
</head>
<body>
    <center><h1>Internal Server Error</h1></center>
</body>
</html>
`

const pFobiddenstr = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Forbidden</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/semantic-ui@2.4.2/dist/semantic.min.css">
</head>
<body>
    <center><h1>Forbidden</h1></center>
</body>
</html>
`

var Page404 *template.Template = template.Must(template.New("404").Parse(p404str))
var PageInternalServerError *template.Template = template.Must(template.New("InternalServerError").Parse(pInternalstr))
var PageForbidden *template.Template = template.Must(template.New("Forbidden").Parse(pFobiddenstr))
