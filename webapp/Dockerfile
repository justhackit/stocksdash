#To Build: docker build -t stocksdash-webapp:v0.1 .
#To Push : docker tag stocksdash-webapp:v0.1 ajayedap/cloudlifter-images:stocksdash-webapp;docker push ajayedap/cloudlifter-images:stocksdash-webapp
#To Run on local : docker run -p 8090:80 stocksdash-webapp:v0.1
FROM nginx:latest
COPY src/ /usr/share/nginx/html/
RUN ls -la /usr/share/nginx/html/*

