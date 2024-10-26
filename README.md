# open_telemetry
Sistema que exibe a temperatura em Graus Celsius, Farenheit e Kelvin da localidade de um CEP informado. Além disso será implementado tracing e span para rastrear a saúde do sistema.

Execução:
1 - Executar "docker-compose up" na pasta open_telemetry para subir os serviços do ZIPKIN e o collector do open telemetry.
2 - Executar "go run main.go" na pasta "open_telemetry/TemperaturaService/cmd/api" para subir o serviço de temperatura.
3 - Executar "go run main.go na pasta "open_telemetry/CepService/cmd/api" para subir o serviço de CEP.
4 - Acessar a seguinte URL no browser "http://localhost:8080/" e informar um CEP com 8 digitos.
5 - Acessar a seguinte URL no browser "http://localhost:9411/zipkin/" para ter acesso ao ZIPKIN.
