```bash
docker build -t stresstest .
```

```bash
docker run stresstest --url=http://google.com --requests=1000 --concurrency=10
```