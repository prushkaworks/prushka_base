# prushka_base

Для запуска приложения необходимо создать в корне файл .env по шаблону .env.template, 
указать в нем все необходимые коннекты к базам (postgres и redis), а также
данные майлера (лично я использую для этого gmail).


<br>

Запуск основного приложения:
```
$ go run ./cmd/prushka/main.go
```

Запуск контейнера (дампы баз в директории на уровень выше рабочей):
```
$ docker-compose up -d
```

Запуск тестов: 
```
$ go test -v ./internal/server
```

Генерация документации swagger (на всякий):
```
$ ~/go/bin/swag init -d ./internal/ --parseDependency --parseInternal --parseDepth 1 -g ./server/server.go
```