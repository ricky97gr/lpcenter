FROM docker.io/bitnami/nginx:latest

USER root

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

RUN mkdir -p /app/uploads/plugins \
    && mkdir -p /app/config \
    && chmod 777 /app/uploads/plugins

# 先备份原主配置，然后在原配置基础上修复
RUN cp /opt/bitnami/nginx/conf/nginx.conf /opt/bitnami/nginx/conf/nginx.conf.bak
# 直接替换所有以 user 开头的行
RUN sed -i 's/^[[:space:]]*user[[:space:]].*/user root;/' /opt/bitnami/nginx/conf/nginx.conf
# 删除所有 8080 监听
RUN sed -i 's/listen[[:space:]]*8080;//g' /opt/bitnami/nginx/conf/nginx.conf
RUN sed -i 's/listen[[:space:]]*\[::\]:8080;//g' /opt/bitnami/nginx/conf/nginx.conf

COPY docker/nginx.conf /opt/bitnami/nginx/conf/server_blocks/default.conf

COPY docker/start.sh /app/start.sh
RUN chmod +x /app/start.sh

COPY web/dist /app/web

COPY server/bin/lpcenter_server /app/

COPY server/.env /app/config/.env

COPY private.pem /app/
COPY public.pem /app/

EXPOSE 9090 9091 9092

ENTRYPOINT []
CMD ["/app/start.sh"]