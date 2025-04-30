@echo off
echo Запуск triple-s на порту 9000...

cd triple-s
start "Triple-S" triple-s.exe --port 9000 --dir ..\s3-data

echo Triple-S запущен, теперь запускаем основное приложение...
cd ..

echo Создание директорий для хранения...
mkdir s3-data 2>nul
mkdir s3-data\posts 2>nul
mkdir s3-data\comments 2>nul

echo Запуск Docker Compose...
docker-compose up