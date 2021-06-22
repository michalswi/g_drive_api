
## **Local files uploader to Google Drive**


**Official** Google source & docs are available [here](https://developers.google.com/drive/api/v3/quickstart/go).  
**Enable** a Google Workspace API guidance located [here](https://developers.google.com/workspace/guides/create-project).  

**Hardcoded** values*:  
```
const (
	driveFolder    = "audycje"
	sourceDir      = "/tmp/workspace"
	APIcredentials = "credentials.json"
	APItoken       = "token.json"
)
```
are related to another repo where we are downloading files to `sourceDir` directory. Repo available [here](https://github.com/michalswi/broadcast_downloader).  

\* - to be changed in the future

#### **# usage**  
```sh
# upload local files to Drive/'driveFolder'
$ ./g_drive_api upload

# list existing files from Drive/'driveFolder'
$ ./g_drive_api list

# delete 'driveFolder' with all files stored there
$ ./g_drive_api delete
```

#### **# example**

Precondition to move forward is to enable API for Google Drive and have credentials from your [Google Cloud Console](https://console.cloud.google.com/) .

```
$ ./g_drive_api upload
ID: 1Vn4X6ct38uIvbFKF7-k8CQJ7CevboZdh  Name: audycje
Uploading file: /tmp/workspace/28-04-2021-info.mp3
File '28-04-2021-info.mp3' successfully uploaded in 'audycje' directory

$ ./g_drive_api list
Printing all files.
28-04-2021-info.mp3 (15eg_naVks9z-qIxD7OItm_CS4Xn_ATqL)
audycje (1Vn4X6ct38uIvbFKF7-k8CQJ7CevboZdh)

$ ./g_drive_api delete
Removing folder: 'audycje'.
Folder ID: '1Vn4X6ct38uIvbFKF7-k8CQJ7CevboZdh'
Folder with ID: '1Vn4X6ct38uIvbFKF7-k8CQJ7CevboZdh' deleted forever.
```