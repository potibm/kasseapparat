FROM alpine
COPY dist /app
WORKDIR /app
VOLUME [ "/app/data" ]
RUN chmod +x /app/kasseapparat
CMD ["/app/kasseapparat", "80"] 
EXPOSE 80