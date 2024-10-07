# 使用 Go 1.23.2 作為基底映像檔
FROM golang:1.23.2

# 設定工作目錄
WORKDIR /app

# 將程式碼複製到容器中
COPY . .

# 編譯 Go 程式
RUN go build -o weather-api .

# 設定容器啟動時要執行的命令
CMD ["./weather-api"]