# FROM alpine:latest

# ADD app /usr/bin/entrypoint
# RUN chmod +x /usr/bin/entrypoint

# ENTRYPOINT ["/usr/bin/entrypoint"]
# CMD []


FROM buildpack-deps:jessie-scm

COPY gobin/benchtest /usr/bin/
RUN chmod +x /usr/bin/benchtest

ENV PARALLEL=10
ENV ITERATIONS=10

CMD ["/usr/bin/benchtest"]