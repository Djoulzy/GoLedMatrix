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