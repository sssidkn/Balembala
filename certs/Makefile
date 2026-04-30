NAME=kafka
PASSWORD=123456
CNF_PATH=client.cnf

root-ca:
	openssl req -new -x509 -keyout ca.key -out ca.crt -days 365 -subj "/CN=ca" -nodes
# Корневой сертификат (для брокеров и для клиентов) (!Строго один для всех)
# Генерируется приватный ключ ca.key и сертификат ca.crt, который подписывает сертификаты брокеров и клиентов,
# помещается в truststore и брокеров и клиентов для дальнейшей проверки подлинности сертификатов
# CN - Common Name (по сути ни на что не влияет)
# По хорошему ключ еще должен быть зашифрован, но мы обойдемся

truststore:
	keytool -import -file ca.crt -keystore kafka.truststore.jks -storepass $(PASSWORD)
# truststore для хранения корневого сертификата
# клиент проверяет подписан ли сертификат брокера этим сертификатом и совпадает ли SAN с адресом брокера.

ssl:
	openssl req -new -newkey rsa:2048 -keyout $(NAME).key -out $(NAME).csr -config $(CNF_PATH) -nodes
	openssl x509 -req -in $(NAME).csr -CA ca.crt -CAkey ca.key -CAcreateserial -out $(NAME).crt -days 365 -extfile $(CNF_PATH) -extensions v3_req

# Генерируется приватный ключ (расшифровывает приходящие данные) и сертификат (в сертификате публичный ключ для
# шифрования данных, информация о компании, название адрес, то самое CN, по сути для самоподписного сертификата неважно)
# То что нам важно - использовать конфиг с указанием SAN. SAN — типо расширения для сертификата,
# указывает дополнительные DNS и IP, для которых действителен сертификат. CN считается устаревшим, а TLS жестко требует SAN.
# !!! Для Kafka необходимо указывать все DNS и IP, которые могут использоваться для подключения
# Далее сгенерированный сертификат для брокера подписывается тем самым корневым сертификатом

keystore:
	openssl pkcs12 -export -in $(NAME).crt -inkey $(NAME).key -out $(NAME).p12 -name $(NAME) -password pass:$(PASSWORD)
	keytool -importkeystore -destkeystore $(NAME).keystore.jks -srckeystore $(NAME).p12 -srcstoretype PKCS12 -alias $(NAME) -srcstorepass $(PASSWORD) -deststorepass $(PASSWORD)

# Создается файл-хранилище pkcs12, зашифрованное паролем, в котором содержится сертификат и приватный ключ
# Конвертируется в keystore (jks, который поддерживается Kafka)
# Для клиента все то же самое

truststore: root-ca truststore # Один раз

keystore: ssl keystore # Для каждого клиента или брокера

