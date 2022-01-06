#To Build: docker build -t stocksdash-webapp:v0.1 .
#To Push : docker tag stocksdash-webapp:v0.1 ajayedap/cloudlifter-images:stocksdash-webapp;docker push ajayedap/cloudlifter-images:stocksdash-webapp
#To Run on local : docker run -p 8090:80 stocksdash-webapp:v0.1
FROM nginx:latest
#stocksdash-webapp here acts as app's root context for routing in ingress
RUN mkdir /usr/share/nginx/html/stocksdash-webapp/
COPY src/ /usr/share/nginx/html/stocksdash-webapp/
RUN ls -la /usr/share/nginx/html/stocksdash-webapp/*

