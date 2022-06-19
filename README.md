# logsrv
Сервер для хранения и анализа логов. В отличие от демонстрационного примера https://github.com/n-r-w/log-server-demo готов к использованию в реальной работе

## Поддерживает следующие операции

Добавить записи в лог:

"service", "source" и т.п. поля, которые могут интерпретироваться в зависмости от использования. Обязательных полей нет.
"parameters" - произвольный набор пар ключ-значение.
"body" - произвольный json.

    curl --location --request POST 'http://localhost:8080/api/add' \
    --header 'X-Authorization: dbda0fba4da680c615340d6faa2868eb5413c3b837640078b87149872257f842' \
    --header 'Content-Type: application/json' \
    --data-raw '[
        {
            "service": "WBA",
            "source": "DEMO",
            "category": "",
            "level": "INFO",
            "session": "",
            "info": "произвольная информация для поиска через регулярные выражения",
            "url": "https://github.com/n-r-w/logsrv/edit/main/README.md",
            "httpType": "POST",
            "httpHeaders": {
                "User-Agent": "PostmanRuntime/7.29.1"            
            },
            "properties": {
                "idUser": "456",            
                "idOrder": "123"
            },
            "body": {
                "id": "4286",
                "bank": {},
                "mainInfo": {
                    "type": "personalDataChangePassportAge",
                    "status": "anketaDraft"
                },
                "documentBase": {},
                "passportMain": {
                    "gender": 1,
                    "lastName": "Иванов",
                    "birthDate": "1999-01-01",
                    "firstName": "Андрей",
                    "birthPlace": "Место рождения",
                    "middleName": "Юрьевич",
                    "registrationDate": "2022-06-01"
                },
                "applicant-flags": {
                    "bankComplete": true,
                    "bankNeedVerify": false,
                    "overallComplete": false,
                    "passportComplete": false,
                    "documentBaseComplete": true,
                    "documentBaseNeedVerify": false,
                    "passport1stPageComplete": false,
                    "passport2ndPageComplete": true,
                    "passportNeedVerify_1stPage": true,
                    "passportNeedVerify_2ndPage": true
                },
                "passportRegistration": {
                    "area": null,
                    "city": "г Москва",
                    "flat": "кв 1",
                    "house": "д 5",
                    "index": "125319",
                    "region": "г Москва",
                    "street": "ул Часовая",
                    "areaFiasGuid": null,
                    "cityFiasGuid": "0c5b2444-70a0-4932-980c-b4dc0d3f02b5",
                    "flatFiasGuid": "9eb55994-ab56-4df0-ba0b-798c27af0e91",
                    "houseFiasGuid": "5e626110-547e-4947-b021-cbf8584658c0",
                    "regionFiasGuid": "0c5b2444-70a0-4932-980c-b4dc0d3f02b5",
                    "streetFiasGuid": "d8a334dd-3d3d-4838-aec1-41a40f616318"
                }
            }
        }
    ]'

Поиск записей:

Все параметры не обязательные.
"and" - применять условие И или условие ИЛИ для соединения критериев.
"bodyValues" - позволяет искать по полям json на любом уровне. 
"body" - поиск по конкретным путям json.

    curl --location --request POST 'http://localhost:8080/api/search' \
    --header 'X-Authorization: dbda0fba4da680c615340d6faa2868eb5413c3b837640078b87149872257f842' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "and": false,
        "criteria": [
            {
                "and": true,
                "from": "2022-06-17T07:30:09.029976Z",
                "to": "2022-06-17T18:30:09.029976Z",
                "service": "WBA",
                "level": "INFO",
                "info": "(информация)*(поиска)",
                "url": "(github.com)*(771)",
                "httpHeaders": {
                    "User-Agent": "PostmanRuntime/7.29.1"
                },
                "properties": {
                    "idOrder": "123"
                },
                "bodyValues": {
                    "lastName": "Иванов",
                    "status": "anketaDraft"
                },
                "body": {               
                    "documentBase": {},
                    "passportMain": {
                        "gender": 1,
                        "lastName": "Иванов",
                        "birthDate": "1999-01-01",
                        "firstName": "Андрей",
                        "birthPlace": "Место рождения",
                        "middleName": "Юрьевич",
                        "registrationDate": "2022-06-01"
                    },
                     "mainInfo": {
                        "type": "personalDataChangePassportAge",
                        "status": "anketaDraft"
                    }
                }
            }
        ]
    }'
 
