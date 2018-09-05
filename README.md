http_attack
==============

Testing http server loading tool, don't be evil

### Usages

Get
```
$> go run http_attack.go --concurrentNum 1  --uri http://www.google.com
```

Post
```
$> go run http_attack.go --concurrentNum 1  --method post --uri http://localhost --params "token=123"
```

PostJSON
```
$> go run http_attack.go --concurrentNum 1  --method postjson --uri http://localhost --params "{'token':'123'}"
```