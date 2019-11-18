
### Local files uploader to Google Drive

**Official** Google source & docs are available [here](https://developers.google.com/drive/api/v3/quickstart/go).  

**Hardcoded** in code*:  
```
const (
	driveFolder = "audycje"
	sourceDir   = "/tmp/workspace"
)
```
It's related to another repo where we are downloading files to `sourceDir` directory. Repo available [here](https://github.com/michalswi/broadcast_downloader).  

\* - to be changed in the future

**Usage:**  
```sh
$ go get -u google.golang.org/api/drive/v3
$ go get -u golang.org/x/oauth2/google

# upload local files to Drive (create driveFolder first)
$ ./goapi upload
# list existing files in Drive
$ ./goapi list
# delete driveFolder with all files stored there
$ ./goapi delete
```