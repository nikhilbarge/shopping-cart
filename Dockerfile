
FROM ubuntu:18.04

COPY install.sh /install.sh
RUN chmod +x /install.sh
COPY shopping-cart /shopping-cart

ENTRYPOINT ["/install.sh"]
