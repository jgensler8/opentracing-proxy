FROM scratch
COPY opentracing-proxy /opentracing-proxy
CMD ["/opentracing-proxy"]
