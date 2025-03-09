[1mdiff --git a/docker-compose.yml b/docker-compose.yml[m
[1mindex 640e92b..ca92cd5 100644[m
[1m--- a/docker-compose.yml[m
[1m+++ b/docker-compose.yml[m
[36m@@ -1,9 +1,37 @@[m
 [m
 services:[m
[32m+[m[32m  postgres:[m
[32m+[m[32m    image: postgres:15[m
[32m+[m[32m    volumes:[m
[32m+[m[32m      - postgresdata:/var/lib/postgresql/data[m
[32m+[m[32m    ports:[m
[32m+[m[32m      - 5432:5432[m
[32m+[m[32m    networks:[m
[32m+[m[32m      - local[m
[32m+[m[32m    healthcheck:[m
[32m+[m[32m      test: ["CMD-SHELL", "pg_isready -U postgres"][m
[32m+[m[32m      interval: 5s[m
[32m+[m[32m      timeout: 30s[m
[32m+[m[32m      retries: 3[m
[32m+[m[32m    environment:[m
[32m+[m[32m      - POSTGRES_PASSWORD=CenalomLoh[m
   back:[m
     build:[m
       dockerfile: ./dockerfile[m
       context: .[m
     ports:[m
       - 3000:3000[m
[31m-    restart: always[m
\ No newline at end of file[m
[32m+[m[32m    restart: always[m
[32m+[m
[32m+[m[32m    depends_on:[m
[32m+[m[32m      postgres:[m
[32m+[m[32m        condition: service_started[m
[32m+[m[32m    networks:[m
[32m+[m[32m      - local[m
[32m+[m
[32m+[m[32mvolumes:[m
[32m+[m[32m  postgresdata:[m
[32m+[m
[32m+[m[32mnetworks:[m
[32m+[m[32m  local:[m
[32m+[m[32m    driver: bridge[m[41m [m
\ No newline at end of file[m
