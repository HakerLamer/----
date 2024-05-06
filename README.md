Для того, чтобы запустить данный контейнер нужно сделать:

Запулить его с Dockerhub:

docker pull tomhetfrainsiden/sbercont_egorshramov 

Запустить docker-контейнер: (-d в фоновом режиме)

docker run -d -p {ваш порт}:8080 sbercont_egorshramov

После этого зайти на http://localhost:{ваш порт}/GetRandomBreeds.