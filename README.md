## Микросервис для работы с балансом пользователей

**Проблема:**

В нашей компании есть много различных микросервисов. Многие из них так или иначе хотят взаимодействовать с балансом
пользователя. На архитектурном комитете приняли решение централизовать работу с балансом пользователя в отдельный
сервис.

**Задача:**

Необходимо реализовать микросервис для работы с балансом пользователей (зачисление средств, списание средств, перевод
средств от пользователя к пользователю, а также метод получения баланса пользователя). Сервис должен предоставлять HTTP
API и принимать/отдавать запросы/ответы в формате JSON.

**Основное задание:**

1. Метод начисления средств на баланс.

_Input_:

```
{
  id_used: int,     #ID Пользователя
  sum: int,         #Кол-во средств на зачисление
}
```

2. Метод резервирования средств с основного баланса на отдельном счете.

_Input_:

```
{
  id_used: int,     #ID Пользователя
  id_setvice: int,  #ID Услуги
  id_order: int,    #ID Заказа
  sum: int,         #Стоймость
}
```

3. Метод признания выручки – списывает из резерва деньги, добавляет данные в отчет для бухгалтерии. Принимает id
   пользователя, ИД услуги, ИД заказа, сумму.

_Input_:

```
{
  id_used: int,      #ID Пользователя
  id_setvice: int,   #ID Услуги
  id_order: int,     #ID Заказа
  sum: int,          #Стоймость
}
```

4. Метод получения баланса пользователя.

_Input_:

```
{
  id_used: int,      #ID Пользователя
}
```

### Руководство по запуску:

---

Для развертывания приложения необходимы docker-compose и утилита [migrate](https://github.com/golang-migrate/migrate)

Последовательность команд для запуска приложения:

```
make build
make run
make migrate-up
```

### Основное задание:

Методы:
* `GET  /api/getBalance/:id`: Метод проверки баланса пользователя
* `POST /api/deposit`: Метод пополнение баланса
* `POST /api/withdrawal`: Метод снятие средств со счета
* `POST /api/transfer`: Метод перевода между пользователями
* `POST /api/reserveService`: Метод резервирования средств с основного баланса на отдельном счете
* `POST /api/approveService`: Метод признания выручки: списывает из резерва деньги, добавляет данные в отчет для бухгалтерии
* `POST /api/failedService`: Метод разрезервирование средств на счету пользователя при отмене покупки услуги
---

#### Метод проверки баланса пользователя

```curl
GET /api/balance/:id
```

Input:

На вход подается только ID пользователя.

Output:

```
{
  "user-balance": int,
  "user-pending-amount": int
}
```

---

#### Метод пополнение баланса

_Примечание_: Создает нового пользователся при отсутствии.

```curl
POST /api/deposit:
```

Input:

На вход подается ID пользователя и кол-во средств на зачисление. При вызове метода для несуществующего пользователя
создается новый.

```
{
  "user-id": int,
  "update-amount": int
}
```

Output:

```
{
  "account-id": int,
  "created-at": time.Time(),
  "operation-status": string,
  "operation-event": string,
  "sum-deposited": int
}
```

---

#### Метод снятие средств со счета

```curl
POST /api/withdrawal
```

Input:

На вход ID пользователя и кол-во средств на снятие.

```
{
  "user-id": int,
  "update-amount": int
}
```

Output:

```
{
  "account-id": int,
  "created-at": time.Time(),
  "operation-event": string,
  "operation-status": string,
  "sum-withdrawn": int
}
```

---

#### Метод перевода между пользователями

```curl
POST /api/transfer
```

Input:

```
{
  "sender-id": int,
  "receiver-id": int,
  "transfer-amount": int
}
```

Output:

```
{
  "amount": int,
  "created-at": time.Time(),
  "event-type": string,
  "receive-account": int,
  "status": string,
  "transfer-account": int
}
```

---

#### Метод резервирования средств с основного баланса на отдельном счете

```curl
POST /api/reserveService
```

Input:

На вход принимает ID пользователя, ID услуги, ID заказа и стоимость.
Записывает все в отдельную таблицу.

```
{
  "user-id": int,
  "service-id": int,
  "order-id": int,
  "payment": int
}
```

Output:

```
{
  "account-id": int,
  "created-at": time.Time(),
  "invoice": int,
  "order-id": int,
  "service-id": int,
  "status": string,
  "updated-at": time.Time()
}
```

--- 

#### Метод признания выручки: списывает из резерва деньги, добавляет данные в отчет для бухгалтерии

```curl
POST /api/approveService
```

Input:

На вход ID пользователя, ID услуги, ID заказа и сумму.
Проверяет последний статус по услуге с принятыми параметрами на предмет конфликта.
Устанавливает статус записи с принятыми параметрами в "Approved".
Уменьшает кредитный баланс пользователя и пытается списать сумму со счета.
В случае нехватки средств откатывает транзакцию и возвращает ошибку.

```
{
  "user-id": int,
  "service-id": int,
  "order-id": int,
  "payment": int
}
```

Output:

```
{
  "account-id": int,
  "created-at": time.Time(),
  "invoice": int,
  "order-id": int,
  "service-id": int,
  "status": string,
  "updated-at": time.Time()
}
```

---

#### Метод разрезервирование средств на счету пользователя при отмене покупки услуги

```curl
POST /api/failedService
```

Input:

На вход дается ID пользователя, ID услуги, ID заказа и сумму.
Проверяет последний статус по услуге с принятыми параметрами 
на предмет конфликта.
Устанавливает статус по услуге в "Cancelled", уменьшает кредитный баланс на стоимость услуги.

```
{
  "user-id": int,
  "service-id": int,
  "order-id": int,
  "payment": int
}
```

Output:

```
{
  "account-id": int,
  "created-at": time.Time(),
  "invoice": int,
  "order-id": int,
  "service-id": int,
  "status": string,
  "updated-at": time.Time()
}
```