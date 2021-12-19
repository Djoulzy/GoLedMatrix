package server

import (
	"html/template"
)

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="user-scalable=no, width=device-width, initial-scale=1.0, maximum-scale=1.0"/>
        <meta name="apple-mobile-web-app-capable" content="yes" />
        <meta name="apple-mobile-web-app-status-bar-style" content="black" />
        <style>
            body {
                font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica,
                    Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
                height: 97vh;
                max-height: 600px;
                background-color: #222222;
            }
        </style>

        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/jquery-modal/0.9.1/jquery.modal.min.css" />

        <script src="https://code.jquery.com/jquery-3.5.1.min.js" integrity="sha256-9/aliU8dGd2tb6OSsuzixeV4y/faTqgFtohetphbbj0=" crossorigin="anonymous"></script>
        <script src="https://cdnjs.cloudflare.com/ajax/libs/jquery-modal/0.9.1/jquery.modal.min.js"></script>
        <script>
            function ajaxDataLoader(url, dataType, method, postVal) {
                var host = document.location.protocol+"//"+document.location.hostname
                if (document.location.port != 0) host += ":"+document.location.port
                host += "/" + url
            
                return $.ajax({
                    url : host,
                    type : method,
                    data: JSON.stringify(postVal),
                    async: true,
                    dataType : dataType,
                })
                .fail(function(data) { console.log("-- Error -- url: ", url) })
            }
        </script>
    </head>
    <body>
        <button type="button" class="modeSelect" data-mode="1">slideShow</button>
        <button type="button" class="modeSelect" data-mode="2">displayGif</button>
        <button type="button" class="modeSelect" data-mode="3">HorloLed</button>
        <button type="button" class="modeSelect" data-mode="4">BouncingBall</button>

        <script>
            $(document).ready(function() {

                $(".modeSelect").on("click", function(e) {
                    var mode = $(e.target).data("mode")
                    const params = {
                        mode: mode,
                    }
                    $.when(ajaxDataLoader("test", 'html', 'POST', params)).done(function(data) {

                    })
                })

            });
        </script>
    </body>
</html>
`))
