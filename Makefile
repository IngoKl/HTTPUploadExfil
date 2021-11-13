build:
	go build

cert:
	openssl req -new -newkey rsa:2048 -nodes -keyout HTTPUploadExfil.key -out HTTPUploadExfil.csr
	openssl x509 -req -days 365 -in HTTPUploadExfil.csr -signkey HTTPUploadExfil.key -out HTTPUploadExfil.csr