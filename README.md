# check-gitlab-patch-status
This is a small tool that checks if a self-hosted GitLab installation is up to date and if there are critical vulnerabilities using the public GitLab API.  
It returns a short text concerning the status of your GitLab installation alongside a status code defined in the following table:
| Code | Meaning |
| ---- | --------|
| 0    | OK      |
| 1    | WARN    |
| 2    | CRIT    |
| 3    | ERR     |

## Usage:
**Linux**
```bash
./check-gitlab-status -H my.gitlab.installation -t myToken
```
<br>

**Windows**
```PowerShell
.\check-gitlab-status.exe -H my.gitlab.installation -t myToken
```

## Parameters
`--host` or `-H`:  
The host or path of your GitLab installation  
<br>
`--token` or `-t`:  
A personal access token with *read_api* privileges
