# ZADANIE

W ramach zadania należy stworzyć aplikację, która nasłuchuje na porcie 8080 i obsługuje protokół HTTP.

Aplikacja powinna wystawiać endpointy rozwiązujące trzy problemy znajdujące się w katalogach atmservice, onlinegame i transactions. 
Specyfikacje endpointów w standarcie OpenAPI 3.0 znajdują się w plikach .json.
Przykładowe wywołanie i odpowiedź, zawarte są w plikach example_request.json i example_response.json

Powodzenia i udanej zabawy! :)


# Jak uruchomić?

Potrzebne jest Go w wersji 1.20.1.

```
$ go get # tylko raz, żeby pobrać zależności
$ go run .
```