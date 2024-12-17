FROM gomicro/goose

COPY . /migrations/

# Convert line endings for both files and make script executable
RUN sed -i 's/\r$//' /migrations/goose.sh && \
    sed -i 's/\r$//' /migrations/.env && \
    chmod +x /migrations/goose.sh

ENTRYPOINT ["/migrations/goose.sh", "up"]