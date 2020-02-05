
### Local files uploader to Google Drive

**Official** Google source & docs are available [here](https://developers.google.com/drive/api/v3/quickstart/go).  

**Hardcoded** values*:  
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
$ go build

# upload local files to Drive/'driveFolder'
$ ./g_drive_api upload

# list existing files from Drive/'driveFolder'
$ ./g_drive_api list

# delete 'driveFolder' with all files stored there
$ ./g_drive_api delete
```