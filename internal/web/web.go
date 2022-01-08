package web

import (
	"html/template"
	"io"
)

const askTpl = `
<!doctype html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@4.5.3/dist/css/bootstrap.min.css" integrity="sha384-TX8t27EcRE3e/ihU7zmQxVncDAy5uIKz4rEkgIXeMed4M0jlfIDPvg6uqKI2xXr2" crossorigin="anonymous">
    <base href="http://autorecord.local" />

    <title>Auto Record</title>
  </head>
  <body>
    <header>
        <div class="bg-dark collapse" id="navbarHeader" style="">
            <div class="container">
                <div class="row">
                    <div class="col-sm-8 col-md-7 py-4">
                    <h4 class="text-white">About</h4>
                    <p class="text-muted">The Auto Record program usually machine learning visual recognition to take a picture and play a given record on a Spotify speaker. It requires some configuration below.</p>
                    </div>
                </div>
            </div>
        </div>
        <div class="navbar navbar-dark bg-dark shadow-sm">
            <div class="container d-flex justify-content-between">
            <a href="#" class="navbar-brand d-flex align-items-center">
                <strong>Auto Record</strong>
            </a>
            </div>
        </div>
    </header>
    <main role="main">
        <section class="jumbotron text-center">
            <div class="container">
               <h1 class="jumbotron-heading">{{.Title}}</h1>

                <div class="row">
                    {{range $index, $item := .Items}}
                    <div class="col-sm-6">
                        <div class="card">
                            <div class="card-body">
                                <p class="card-text">{{ $item.Text }}</p>
                                <a href="{{ $item.Path }}" class="btn btn-primary">Let's do it!</a>
                            </div>
                        </div>
                    </div>
                    {{end}}
                </div>
            </div>
        </section>
    </main> 

    <script src="https://code.jquery.com/jquery-3.5.1.slim.min.js" integrity="sha384-DfXdz2htPH0lsSSs5nCTpuj/zy4C+OGpamoFVy38MVBnE+IbbVYUew+OrCXaRkfj" crossorigin="anonymous"></script>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@4.5.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-ho+j7jyWK8fNQe+A12Hb8AhRq26LrZ/JpcUGGOn+Y7RsweNrtN/tE3MoK7ZeZDyx" crossorigin="anonymous"></script>
  </body>
</html>
`

type Item struct {
	Text string
	Path string
}

// Ask writes a webpage to io.Writer asking user to pick from a list of items
func Ask(w io.Writer, title string, items []Item) {
	t, err := template.New("webpage").Parse(askTpl)
	if err != nil {
		panic(err)
	}

	data := struct {
		Title string
		Items []Item
	}{
		Title: title,
		Items: items,
	}

	err = t.Execute(w, data)
	if err != nil {
		panic(err)
	}
}
