# pinback

quick way to export recent pinned urls from pinboard for archiving

## usage

```bash
# check out the repo, then
go run . PINBOARD_API_TOKEN
```

Note: this uses pinboard's [posts/all](https://pinboard.in/api/#posts_all) API 
endpoint, which is rate-limited to one call every 5 minutes.