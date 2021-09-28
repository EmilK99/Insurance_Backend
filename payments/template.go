package payments

import "html/template"

var (
	SuccessTemplate, _ = template.New("success").Parse(
		`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>success</title>
</head>
<body>
<h1>Transaction successful!</h1>
</body>
</html>`)
	CancelTemplate, _ = template.New("cancel").Parse(
		`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>cancel</title>
</head>
<body>
    <h1>Transaction cancelled!</h1>
</body>
</html>`)
)
