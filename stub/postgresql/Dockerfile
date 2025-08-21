FROM public_ecr.aws/docker/library/postgres:16

RUN apt-get update && \
    apt-get install -y \
    build-essential \
    postgresql-server-dev-16 \
    postgresql-contrib-16 \
    postgresql-16-partman && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN echo "shared_preload_libraries = 'pg_partman_bgw'" >> /usr/share/postgresql/postgresql.cong.sample