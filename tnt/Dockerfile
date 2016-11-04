FROM tarantool/tarantool:1.7.3


ADD ./func.so /usr/lib/tarantool/
COPY console /usr/bin/console


RUN chmod +x /usr/bin/console
