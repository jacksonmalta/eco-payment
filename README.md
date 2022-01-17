# Eco Payment

Adicione saldo, faça saques ou compras à vista ou parcelado

---

Inicinando o ecosistema de pagamento:

Será necessário Docker, Docker compose e Make

Na raiz do repositório executar o comando

```shell
make run/docker
```

---

Criando uma conta:

```shell
curl -i --location --request POST 'localhost:5002/v1/accounts' \
--header 'Content-Type: application/json' \
--data-raw '{
"document_number": "05662459061",
"external_key": "1"
}'
```

document_number: deve conter apenas números e representa o seu CPF ou CNPJ

external_key: é a chave da sua conta e deverá ser único

---

Consultando uma conta:

```shell
curl -i --location --request GET 'localhost:5002/v1/accounts/1'
```

external_key: chave única da sua conta

---

Adicionando saldo:

````shell
curl -i --location --request POST 'localhost:5004/v1/transactions' \
--header 'Content-Type: application/json' \
--data-raw '{
"account_key": "1",
"external_key": "1",
"amount": 1000
}'
````

account_key: é a chave da conta

external_key: é a chave da transação e deverá ser única por conta

amount: é o valor em centavos que deseja adicionar ao saldo

---

Compra à vista ou compra parcelada ou saque:

```shell
curl -i --location --request POST 'localhost:5005/v1/transactions' \
--header 'Content-Type: application/json' \
--data-raw '{
"account_key": "1",
"external_key": "4",
"operation_type": "Withdraw",
"amount": 1000
}'
```

account_key: é a chave da conta

external_key: é a chave da transação e deverá ser única por conta.

operation_type: é o tipo de operação como compra à vista ou compra parcelada ou saque (abaixo está a tabela com os valores possíveis)

amount: é o valor em centavos da operação.

---

Operações possíveis para compras ou saque (operation_type):

Compra à vista: Buying

Compra parcelada: InstallmentBuying

Saque: Withdraw

---

Consulta de saldo por conta:

** Para realizar a consulta abaixo é necessário ter o AWS CLI configurado

```shell
aws --endpoint-url=http://localhost:4566  dynamodb query \
--table-name balance \
--key-condition-expression "AccountKey = :v1" \
--expression-attribute-values '{":v1": {"S": "1"}}'
--return-consumed-capacity TOTAL
```

---