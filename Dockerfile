FROM busybox
ADD ["go-event-store", "/"]
EXPOSE 8080
ENTRYPOINT ["/go-event-store","-addr","0.0.0.0:8080"]
