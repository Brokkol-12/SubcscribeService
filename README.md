Susbsription Service API  
REST сервис для агреации данных об онлайн-подписках пользователей  

--Установка  
Клонируйте репозиторий командами в cmd:  

git clone https://github.com/Brokkol-12/SubscribeService.git  
cd SubscribeService  

--Настройка окружения  
Прежде чем запускать программу, создайте .env файл в корневой директории.  
Вы можете использовать этот файл для примера: 
cp .env.example .env  

Затем настройте следующие параметры:  

DB_USER=your_db_user  
DB_PASSWORD=your_db_password  

Убедитесь что данные авторзиации совпадают с конфигурацией PostgreSQL.  

--Запуск при помощи Docker (Рекомендовано)  

Соберите и запустите контейнер:  

docker-compose up --build  

После запуска, серис будет доступен по адресу:  
http://localhost:8081  

--Запуск без Docker  
Вам надо установить на свой ПК:  
GO 1.21+  
PostgreSQL  

После, введите в своем CMD в директории API:  
go mod tidy  

И запустите программу:  
go run cmd/main.go  

--Swagger  
После старта у вас будет доступ к документации SWAGGER  
Откройте Ваш браузер и вбейте в поисковую строку:   
http://localhost:8081/swagger/index.html  
Тут вы получите параметры по которым работает API.  
Так же, в самом Swagger возможно проводить тестирование, после раскрытия вкладки с интересующим Вас сервисом, есть подпункт - "Parameters Try it out",  
где система сама подскажет какие и куда данные вводить, а так же выдаст результат

--Тестирование API  
Примеры запросов:  
Create Subscription  
curl -X 'POST' \  
  'http://localhost:8081/subs/create' \  
  -H 'accept: application/json' \  
  -H 'Content-Type: application/json' \  
  -d '{  
  "end_date": "YYYY-MM-DD",  
  "price": 0,  
  "service_name": "string",  
  "start_date": "YYYY-MM-DD",  
  "user_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"  
}  

Важно! UUID должен быть в строгом формате: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx (8-4-4-4-12 символов)  


GetById Subscription  
curl -X 'GET' \  
  'http://localhost:8081/subs/get?id=subscribeID' \  
  -H 'accept: application/json'  

Update Subscription  

curl -X 'PUT' \  
  'http://localhost:8081/subs/update?id=subscribeID' \  
  -H 'accept: application/json' \  
  -H 'Content-Type: application/json' \  
  -d '{  
  "end_date": "YYYY-MM-DD",  
  "price": 0,  
  "service_name": "service",  
  "start_date": "YYYY-MM-DD",  
  "user_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"  
}'  


Delete Subscripion  

curl -X 'DELETE' \  
  'http://localhost:8081/subs/delete?id=subscribeID' \  
  -H 'accept: application/json'  

List Subscripton  

curl -X 'GET' \  
  'http://localhost:8081/subs/list?user_id=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx' \  
  -H 'accept: application/json'  

TotalCalculate Subscription  

curl -X 'GET' \  
  'http://localhost:8081/subs/total?user_id=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx&start=YYYY-MM-DD&end=YYYY-MM-DD' \  
  -H 'accept: application/json'  
 


--API документация  

Swagger UI содержит:  
Create subscription  
Get subscription by ID  
Get all subscriptions  
Update subscription  
Delete subscription
List subscription
Get a total sum on all subscription  

--Tech Stack  
Go 1.21+  
PostgreSQL  
Docker & Docker Compose  
Swagger  

--Замечания  
Убедитесь что PostgreSQL запущен перед стартом программы в Docker.  
.env файл необходим для подключения к базе данных.  
Не привязывайте файл .env к системе управления версиями.  
